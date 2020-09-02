package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
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
	Sig string
}

type Output struct {
	Amount  float64
	PubHash string
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
		TXid:  []byte{},
		Index: -1,
		Sig:   data,
	}
	out := Output{
		Amount:  reward,
		PubHash: address,
	}
	trans := Transaction{
		Inputs:  []Input{in},
		Outputs: []Output{out},
	}
	trans.setHash()
	return &trans
}

func (t *Transaction) IsCoinBase() bool {
	if len(t.Inputs) == 1 && len(t.Inputs[0].TXid) == 0 && t.Inputs[0].Index == -1 {
		return true
	}
	return false
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 获取utxo 列表
	used := 0.0
	utxo := bc.FindUTXOs(from)
	var inputs []Input
	var output []Output
	for _, ut := range utxo {
		if used >= amount {
			break
		}
		// input
		used += ut.Amount
		inputs = append(inputs, Input{
			TXid:  ut.TXid,
			Index: int64(ut.Vout),
			Sig:   from,
		})
	}
	if used < amount {
		fmt.Println("余额不足。。。")
		return nil
	}
	// output
	output = append(output, Output{
		Amount:  amount,
		PubHash: to,
	})
	// 找零
	if used > amount {
		// output
		output = append(output, Output{
			Amount:  used - amount,
			PubHash: from,
		})
	}
	tran := Transaction{
		TXID:    []byte{},
		Outputs: output,
		Inputs:  inputs,
	}
	tran.setHash()
	return &tran
}
