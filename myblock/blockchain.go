package main

import (
	"blockchain/lib/bolt"
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
func (chain *BlockChain) FindUTXOs(address string) []UTXO {
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
				if out.PubHash == address {
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
						PubHash: out.PubHash,
						Amount:  out.Amount,
						Vout:    uint64(i),
					}
					utxo = append(utxo, ut)
				}
			}
			// 添加以使用过的集合 旷工打包区块不需要
			if !tx.IsCoinBase() {
				for _, in := range tx.Inputs {
					if in.Sig == address {
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
