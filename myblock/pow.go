package main

func main() {

	bc := NewBlockChain()
	bc.AddBlock("班长向班花转了50枚比特币！")
	bc.AddBlock("班长又向班花转了50枚比特币！")
	cli := CLI{blockChain: bc}
	cli.Run()
}
