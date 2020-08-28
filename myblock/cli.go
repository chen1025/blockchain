package main

import (
	"fmt"
	"os"
)

type CLI struct {
	blockChain *BlockChain
}

const usage = `
	addBlock --data DATA  "add block"
	print                 "print all Block"
`

func (cli *CLI) Run() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(usage)
		return
	}
	c := args[1]
	switch c {
	case "addBlock":
		if len(args) == 4 {
			data := args[3]
			cli.AddBlock(data)
		}else {
			fmt.Printf(usage)
		}
	case "print":
		cli.PrintBlockChain()
	default:
		fmt.Printf("无效的命令/n")
		fmt.Printf(usage)
	}
}
