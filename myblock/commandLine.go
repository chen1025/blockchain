package main

import "fmt"

func (cli *CLI) AddBlock(data string) {
	//cli.blockChain.AddBlock(data)
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
		fmt.Printf("当前版本： %v\n", block.Version)
		fmt.Printf("前区块哈希值： %x\n", block.PrevHash)
		fmt.Printf("当前MerkleRoot： %x\n", block.MerkleRoot)
		fmt.Printf("区块TimeStamp :  %v\n", block.TimeStamp)
		fmt.Printf("当前区块Difficulty： %x\n", block.Difficulty)
		fmt.Printf("区块Nonce:  %v\n", block.Nonce)
		fmt.Printf("当前区块哈希值： %x\n", block.Hash)
		fmt.Printf("区块数据 :%s\n", block.Transactions[0].Inputs[0].Sig)
		i++
	}

}

func (cli *CLI) getBalance(address string) float64 {
	total := 0.0
	utxo := cli.blockChain.FindUTXOs(address)
	for _, o := range utxo {
		total += o.Amount
	}
	return total
}

// 发送btc
func (cli *CLI) Send(from string, to string, amount float64, miner string, data string) {
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
