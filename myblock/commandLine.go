package main

import (
	"blockchain/lib/base58"
	"bytes"
	"fmt"
)

func (cli *CLI) AddBlock(data string) {
	//cli.blockChain.AddBlock(data)
	fmt.Println("添加区块成功")
}

func (cli *CLI) PrintBlockChain() {
	bc := cli.blockChain
	iterator := bc.NewIterator()
	for {
		block := iterator.Next()
		for _, tx := range block.Transactions{
			fmt.Println(tx)
		}
		if len(iterator.nextHash) == 0 {
			break
		}
	}
}

func (cli *CLI) getBalance(address string) float64 {
	total := 0.0
	key := getPubKeyHashByAddress(address)
	utxo := cli.blockChain.FindUTXOs(key)
	for _, o := range utxo {
		total += o.Amount
	}
	return total
}

// 发送btc
func (cli *CLI) Send(from string, to string, amount float64, miner string, data string) {

	if !IsValidAddress(from) {
		fmt.Printf("地址无效 from: %s\n", from)
		return
	}
	if !IsValidAddress(to) {
		fmt.Printf("地址无效 to: %s\n", to)
		return
	}
	if !IsValidAddress(miner) {
		fmt.Printf("地址无效 miner: %s\n", miner)
		return
	}
	// 创建一个挖矿交易
	tx := NewCoinBaseTX(miner, data)
	// 创建一个 发送交易
	transaction := NewTransaction(from, to, amount, cli.blockChain)
	if transaction == nil {
		return
	}
	// 创建一个区块
	cli.blockChain.AddBlock([]*Transaction{tx, transaction})
	fmt.Println("转账成功")
}

func (cli *CLI) createWallet() {
	ws := NewWallets()
	address := ws.CreateWallets()
	fmt.Printf("地址：%s\n", address)
}

func (cli *CLI) listAddress() {
	ws := NewWallets()
	addresses := ws.listAddress()
	for _, address := range addresses {
		fmt.Printf("地址：%s\n", address)
	}
}

func IsValidAddress(address string) bool {
	//1. 解码
	addressByte := base58.Decode(address)

	if len(addressByte) < 4 {
		return false
	}

	//2. 取数据
	payload := addressByte[:len(addressByte)-4]
	checksum1 := addressByte[len(addressByte)-4: ]

	//3. 做checksum函数
	checksum2 := CheckSum(payload)

	fmt.Printf("checksum1 : %x\n", checksum1)
	fmt.Printf("checksum2 : %x\n", checksum2)

	//4. 比较
	return bytes.Equal(checksum1, checksum2)
}
