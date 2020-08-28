package main

import "fmt"

func (cli *CLI) AddBlock(data string) {
	cli.blockChain.AddBlock(data)
	fmt.Println("添加区块成功")
}

func (cli *CLI) PrintBlockChain() {

	bc := cli.blockChain
	iterator := bc.NewIterator()
	var i int
	for {
		block := iterator.Next()
		if len(iterator.nextHash) == 0 {
			break
		}
		fmt.Printf("======== 当前区块高度： %d ========\n", i)
		fmt.Printf("前区块哈希值： %x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值： %x\n", block.Hash)
		fmt.Printf("区块数据 :%s\n", block.Data)
		i++
	}

}
