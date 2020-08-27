package main

type BlockChain struct {
	blocks []*Block
}

func NewBlockChain() *BlockChain {
	//生成 创世区块
	block := GenerateGenesisBlock()
	// 创建一个chain
	return &BlockChain{
		blocks: []*Block{block},
	}
}

func GenerateGenesisBlock() *Block {
	return NewBlock("I AM GOD", []byte{})
}

func (chain *BlockChain) AddBlock(data string) {
	//获取 父区块 hash
	perv := chain.blocks[len(chain.blocks)-1]
	//a. 创建新的区块
	block := NewBlock(data, perv.Hash)
	//添加
	chain.blocks = append(chain.blocks, block)
}
