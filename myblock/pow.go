package main

func main() {

	bc := NewBlockChain("14A52QNX7hYD21yFPn1qw8u6xRb12F6Cr6")
	cli := CLI{blockChain: bc}
	cli.Run()
}
