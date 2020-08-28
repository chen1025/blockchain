package main

import (
	"blockchain/bolt"
	"log"
)

type Iterator struct {
	db       *bolt.DB
	nextHash []byte
}

func (i *Iterator) Next() *Block {
	var block Block
	i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DbBucket))
		if bucket == nil {
			log.Panic("数据库为空！")
		}
		get := bucket.Get(i.nextHash)
		block = Deserialize(get)
		i.nextHash = block.PrevHash
		return nil
	})
	return &block
}
