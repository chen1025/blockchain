package main

import (
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	blockChain *BlockChain
}

const usage = `
	addBlock   --data    DATA     "add block"
	print                         "print all Block"
    getBalance --address ADDRESS  "get balance"
	send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
	createWallet                   "create address"
	listAddress                     "print address"
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
		if len(args) == 4 && args[2] == "--data" {
			data := args[3]
			cli.AddBlock(data)
		} else {
			fmt.Printf(usage)
		}
	case "print":
		cli.PrintBlockChain()
	case "getBalance":
		if len(args) == 4 && args[2] == "--address" {
			data := args[3]
			balance := cli.getBalance(data)
			fmt.Println(balance)
		} else {
			fmt.Printf(usage)
		}
	case "send":
		fmt.Printf("转账开始...\n")
		if len(args) != 7 {
			fmt.Printf("参数个数错误，请检查！\n")
			fmt.Printf(usage)
			return
		}
		//./block send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(args[4], 64) //知识点，请注意
		miner := args[5]
		data := args[6]
		cli.Send(from, to, amount, miner, data)
	case "createWallet":
		cli.createWallet()
	case "listAddress":
		cli.listAddress()
	default:
		fmt.Printf("无效的命令/n")
		fmt.Printf(usage)
	}
}
