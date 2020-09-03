package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"strings"
)

const reward = 50.0

type Transaction struct {
	TXID    []byte
	Inputs  []Input
	Outputs []Output
}

type Input struct {
	//引用的交易ID
	TXid []byte
	//引用的output的索引值
	Index int64
	//解锁脚本，我们用地址来模拟
	Sig []byte
	//公钥 公钥中的 X Y 拆分成的公钥
	PubKey []byte
}

type Output struct {
	Amount float64
	// 公钥hash
	PubKayHash []byte
}

// 计算交易的 hash
func (t *Transaction) setHash() {
	// 对交易记录进行 序列化
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(t)
	if err != nil {
		log.Panic(err)
	}
	// sha256
	sum256 := sha256.Sum256(buffer.Bytes())
	t.TXID = sum256[:]
}

// 创建一个 挖矿交易
func NewCoinBaseTX(address, data string) *Transaction {
	//挖矿交易的特点：
	//1. 只有一个input
	//2. 无需引用交易id
	//3. 无需引用index
	//矿工由于挖矿时无需指定签名，所以这个sig字段可以由矿工自由填写数据，一般是填写矿池的名字
	in := Input{
		TXid:   []byte{},
		Index:  -1,
		Sig:    nil,
		PubKey: []byte(data),
	}
	out := NewTXOutput(reward, address)
	trans := Transaction{
		TXID:    []byte{},
		Inputs:  []Input{in},
		Outputs: []Output{*out},
	}
	trans.setHash()
	return &trans
}

//设置 公钥hash
func (out *Output) setPubHash(address string) {
	out.PubKayHash = getPubKeyHashByAddress(address)
}

//给TXOutput提供一个创建的方法，否则无法调用Lock
func NewTXOutput(value float64, address string) *Output {
	output := Output{
		Amount: value,
	}
	output.setPubHash(address)
	return &output
}

func (t *Transaction) IsCoinBase() bool {
	if len(t.Inputs) == 1 && len(t.Inputs[0].TXid) == 0 && t.Inputs[0].Index == -1 {
		return true
	}
	return false
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 获取钱包的公私钥
	wallet := NewWallets()
	w := wallet.WsMap[from]
	if w == nil {
		fmt.Println("钱包文件不存在！")
		return nil
	}
	// 获取公钥hash
	publicKey := w.PublicKey
	privateKey := w.PrivateKey
	key := HashPublicKey(publicKey)

	// 获取utxo 列表
	used := 0.0
	utxo := bc.FindUTXOs(key)
	var inputs []Input
	var output []Output
	for _, ut := range utxo {
		if used >= amount {
			break
		}
		// input
		used += ut.Amount
		inputs = append(inputs, Input{
			TXid:   ut.TXid,
			Index:  int64(ut.Vout),
			Sig:    nil,
			PubKey: publicKey,
		})
	}
	if used < amount {
		fmt.Println("余额不足。。。")
		return nil
	}
	// output
	output = append(output, *NewTXOutput(amount, to))
	// 找零
	if used > amount {
		// output
		output = append(output, *NewTXOutput(used-amount, from))
	}
	tran := Transaction{
		TXID:    []byte{},
		Outputs: output,
		Inputs:  inputs,
	}
	tran.setHash()
	// 对交易进行签名
	bc.signTran(&tran, privateKey)
	return &tran
}

// 对没一个 交易进行签名
func (t *Transaction) sign(prevTXs map[string]Transaction, key *ecdsa.PrivateKey) {

	//1. 创建一个当前交易的副本：txCopy，使用函数： TrimmedCopy：要把Signature和PubKey字段设置为nil
	txCopy := t.TrimmedCopy()
	inputs := txCopy.Inputs
	for i, input := range inputs {
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("引用的交易无效")
		}
		// 对交易的每个input 进行签名
		txCopy.Inputs[i].PubKey = prevTXs[string(input.TXid)].Outputs[i].PubKayHash

		//所需要的三个数据都具备了，开始做哈希处理
		//3. 生成要签名的数据。要签名的数据一定是哈希值
		//a. 我们对每一个input都要签名一次，签名的数据是由当前input引用的output的哈希+当前的outputs（都承载在当前这个txCopy里面）
		//b. 要对这个拼好的txCopy进行哈希处理，SetHash得到TXID，这个TXID就是我们要签名最终数据。
		txCopy.setHash()
		//还原，以免影响后面input的签名
		txCopy.Inputs[i].PubKey = nil
		//signDataHash认为是原始数据
		signDataHash := txCopy.TXID
		r, s, err := ecdsa.Sign(rand.Reader, key, signDataHash)
		if err != nil {
			log.Panic(err)
		}
		// 添加的签名
		t.Inputs[i].Sig = append(r.Bytes(), s.Bytes()...)

	}

}

func (t *Transaction) TrimmedCopy() Transaction {
	var inputs []Input
	var outputs []Output

	for _, input := range t.Inputs {
		inputs = append(inputs, Input{input.TXid, input.Index, nil, nil})
	}

	for _, output := range t.Outputs {
		outputs = append(outputs, output)
	}

	return Transaction{t.TXID, inputs, outputs}
}

// 验证交易是否有效
func (t *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if t.IsCoinBase() {
		return true
	}

	//1. 得到签名的数据
	txCopy := t.TrimmedCopy()

	for i, input := range t.Inputs {
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("引用的交易无效")
		}
		txCopy.Inputs[i].PubKey = prevTX.Outputs[input.Index].PubKayHash
		txCopy.setHash()
		dataHash := txCopy.TXID
		// 拿到公钥
		//2. 得到Signature, 反推会r,s
		signature := input.Sig //拆，r,s
		//3. 拆解PubKey, X, Y 得到原生公钥
		pubKey := input.PubKey //拆，X, Y
		var r, s big.Int
		var x, y big.Int
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])
		x.SetBytes(pubKey[:len(pubKey)/2])
		y.SetBytes(pubKey[len(pubKey)/2:])

		// 验证
		pub := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}
		if !ecdsa.Verify(&pub, dataHash, &r, &s) {
			return false
		}
	}
	return true

}

func (t Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", t.TXID))

	for i, input := range t.Inputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.TXid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Index))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Sig))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range t.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %f", output.Amount))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKayHash))
	}

	return strings.Join(lines, "\n")
}
