package main

type UTXO struct {
	Amount  float64
	PubHash []byte
	Vout    uint64
	TXid    []byte
}
