package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

type Block struct {
	//1.版本号
	Version uint64
	//2. 前区块哈希
	PrevHash []byte
	//3. Merkle根（梅克尔根，这就是一个哈希值，我们先不管，我们后面v4再介绍）
	MerkleRoot []byte
	//4. 时间戳
	TimeStamp uint64
	//5. 难度值
	Difficulty uint64
	//6. 随机数，也就是挖矿要找的数据
	Nonce uint64
	//a. 当前区块哈希,正常比特币区块中没有当前区块的哈希，我们为了是方便做了简化！
	Hash []byte
	//b. 数据
	Data []byte
}

func IntToByte(in uint64) []byte {
	var bu bytes.Buffer
	err := binary.Write(&bu, binary.BigEndian, in)
	if err != nil {
		log.Panic(err)
	}
	return bu.Bytes()
}

func NewBlock(data string, prevHash []byte) *Block {
	block := Block{
		Version:    0,
		PrevHash:   prevHash,
		MerkleRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{},
		Data:       []byte(data),
	}
	setHash(&block)
	return &block
}

func setHash(block *Block) {
	blockInfo := [][]byte{
		IntToByte(block.Version),
		block.PrevHash,
		block.MerkleRoot,
		IntToByte(block.TimeStamp),
		IntToByte(block.Difficulty),
		IntToByte(block.Nonce),
		block.Data,
	}
	join := bytes.Join(blockInfo, []byte{})
	sum256 := sha256.Sum256(join)
	block.Hash = sum256[:]
}
