package main

import (
	"blockchain/bolt"
	"fmt"
)

func main() {

	// 打开数据库
	open, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	// 关闭数据库
	defer open.Close()
	// 新增数据
	open.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("v1"))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("v1"))
			if err != nil {
				panic(err)
			}
		}
		bucket.Put([]byte("god"), []byte("first"))
		bucket.Put([]byte("user"), []byte("second"))
		return nil
	})
	// 获取 数据
	open.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("v1"))
		god := bucket.Get([]byte("god"))
		user := bucket.Get([]byte("user"))
		fmt.Printf("%s/n",god)
		fmt.Printf("%s/n",user)
		return nil
	})

}
