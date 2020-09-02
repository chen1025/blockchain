package main

func main() {

	bc := NewBlockChain("tx1001")
	cli := CLI{blockChain: bc}
	cli.Run()
}
