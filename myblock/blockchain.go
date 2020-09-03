package main

import (
	"blockchain/lib/bolt"
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
)

type BlockChain struct {
	db   *bolt.DB
	last []byte
}

const DbFile = "blockChain.db"
const DbBucket = "dbBucket"
const LastKey = "lastHash"

func NewBlockChain(address string) *BlockChain {
	//生成 创世区块
	var lastHash []byte
	db, err := bolt.Open(DbFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DbBucket))
		if bucket == nil {
			//初始化
			block := GenerateGenesisBlock(address)
			bucket, _ = tx.CreateBucket([]byte(DbBucket))
			bucket.Put([]byte(LastKey), block.Hash)
			bucket.Put(block.Hash, block.Serialize())
			lastHash = block.Hash
		} else {
			lastHash = bucket.Get([]byte(LastKey))
		}
		return nil
	})
	// 创建一个chain
	return &BlockChain{
		db:   db,
		last: lastHash,
	}
}

func GenerateGenesisBlock(address string) *Block {
	tx := NewCoinBaseTX(address, "I AM GOD")
	return NewBlock([]*Transaction{tx}, []byte{})
}

func (chain *BlockChain) AddBlock(t []*Transaction) {
	// 验证

	for _,tx := range t {
		flag := chain.VerifyTransaction(tx)
		if !flag {
			//
			fmt.Printf("存在无效交易:%x/n",tx.TXID)
			return
		}
	}

	//a. 创建新的区块
	block := NewBlock(t, chain.last)
	chain.last = block.Hash
	//添加
	chain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DbBucket))
		if bucket == nil {
			log.Panic("bucket not exist")
		}
		bucket.Put([]byte(LastKey), block.Hash)
		bucket.Put(block.Hash, block.Serialize())
		return nil
	})
}

// 生成一个 迭代器
func (chain *BlockChain) NewIterator() *Iterator {
	return &Iterator{
		db:       chain.db,
		nextHash: chain.last,
	}
}

// 找到可用的 UTXO列表
func (chain *BlockChain) FindUTXOs(pubHash []byte) []UTXO {
	//定义一个 utxo
	var utxo []UTXO
	// 遍历 区块找到 有效的output
	iterator := chain.NewIterator()
	spentOutputs := make(map[string][]int64)
	for {
		block := iterator.Next()
		// 遍历 交易
		transactions := block.Transactions
		for _, tx := range transactions {
		OUTPUT:

			// 获取 和我有关的 output
			for i, out := range tx.Outputs {
				if bytes.Equal(out.PubKayHash,pubHash){
					// 排除使用过的
					if spentOutputs[string(tx.TXID)] != nil {
						for _, id := range spentOutputs[string(tx.TXID)] {
							if int64(i) == id {
								//跳转
								continue OUTPUT
							}
						}
					}
					// 添加到列表中
					ut := UTXO{
						TXid:    tx.TXID,
						PubHash: out.PubKayHash,
						Amount:  out.Amount,
						Vout:    uint64(i),
					}
					utxo = append(utxo, ut)
				}
			}
			// 添加以使用过的集合 旷工打包区块不需要
			if !tx.IsCoinBase() {
				for _, in := range tx.Inputs {
					pubKeyHash := HashPublicKey(in.PubKey)
					if bytes.Equal(pubKeyHash,pubHash){
						spentOutputs[string(in.TXid)] = append(spentOutputs[string(in.TXid)], in.Index)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return utxo
}

func (chain *BlockChain) signTran(t *Transaction, key *ecdsa.PrivateKey) {
	//签名，交易创建的最后进行签名
	prevTXs := make(map[string]Transaction)
	//找到 每个 input 的交易
	inputs := t.Inputs
	for _,in := range inputs{
		xid, err := chain.FindTransactionByTXid(in.TXid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[string(in.TXid)] = xid
	}

	// 签名
	t.sign(prevTXs,key)

}

//根据id查找交易本身，需要遍历整个区块链
func (chain *BlockChain) FindTransactionByTXid(id []byte) (Transaction, error) {

	//4. 如果没找到，返回空Transaction，同时返回错误状态

	fmt.Printf("1111111111 : id%x\n", id)
	it := chain.NewIterator()

	//1. 遍历区块链
	for {
		block := it.Next()
		//2. 遍历交易
		for _, tx := range block.Transactions {
			//3. 比较交易，找到了直接退出
			if bytes.Equal(tx.TXID, id) {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			fmt.Printf("区块链遍历结束!\n")
			break
		}
	}

	return Transaction{}, errors.New("无效的交易id，请检查!")
}

func (chain *BlockChain) VerifyTransaction(tx *Transaction) bool {

	if tx.IsCoinBase() {
		return true
	}

	//签名，交易创建的最后进行签名
	prevTXs := make(map[string]Transaction)

	//找到所有引用的交易
	//1. 根据inputs来找，有多少input, 就遍历多少次
	//2. 找到目标交易，（根据TXid来找）
	//3. 添加到prevTXs里面
	for _, input := range tx.Inputs {
		//根据id查找交易本身，需要遍历整个区块链
		fmt.Printf("2222222 : %x\n", input.TXid)
		tx, err := chain.FindTransactionByTXid(input.TXid)

		if err != nil {
			log.Panic(err)
		}

		prevTXs[string(input.TXid)] = tx

	}

	return tx.Verify(prevTXs)
}
