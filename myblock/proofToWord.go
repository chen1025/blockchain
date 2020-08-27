package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofToWork struct {
	block  *Block
	target *big.Int
}

func NewPOW(block *Block) *ProofToWork {
	pow := ProofToWork{
		block: block,
	}
	// 设定 初始的target
	tar := "0000100000000000000000000000000000000000000000000000000000000000"
	var temp big.Int
	tarStr, _ := temp.SetString(tar, 16)

	pow.target = tarStr
	return &pow
}

func (p *ProofToWork) Run() ([]byte, uint64) {
	// 拼装参数
	var nonce uint64
	block := p.block
	for {
		blockInfo := [][]byte{
			IntToByte(block.Version),
			block.PrevHash,
			block.MerkleRoot,
			IntToByte(block.TimeStamp),
			IntToByte(block.Difficulty),
			IntToByte(nonce),
			block.Data,
		}
		join := bytes.Join(blockInfo, []byte{})
		sum256 := sha256.Sum256(join)
		// 计算 big int 值
		var hashInt big.Int
		hi := hashInt.SetBytes(sum256[:])

		if hi.Cmp(p.target) == -1 {
			fmt.Printf("挖矿成功！hash : %x, nonce : %d\n", sum256, nonce)
			return sum256[:], nonce
		} else {
			nonce++
		}
	}
}
