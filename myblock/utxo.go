package main

type UTXO struct {
	Amount  float64
	PubHash string
	Vout    uint64
	TXid    []byte
}
