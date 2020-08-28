package main

import (
	"blockchain/bolt"
	"log"
)

type BlockChain struct {
	db   *bolt.DB
	last []byte
}

const DbFile = "blockChain.db"
const DbBucket = "dbBucket"
const LastKey = "lastHash"

func NewBlockChain() *BlockChain {
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
			block := GenerateGenesisBlock()
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

func GenerateGenesisBlock() *Block {
	return NewBlock("I AM GOD", []byte{})
}

func (chain *BlockChain) AddBlock(data string) {

	//a. 创建新的区块
	block := NewBlock(data, chain.last)
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

func (chain *BlockChain) NewIterator() *Iterator {
	return &Iterator{
		db:       chain.db,
		nextHash: chain.last,
	}
}
