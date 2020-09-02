package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

const FileName = "wallet.bat"

type Wallets struct {
	WsMap map[string]*Wallet
}

// 创建一个 钱包
func NewWallets() *Wallets {
	var ws Wallets
	ws.WsMap = make(map[string]*Wallet)
	// 加载文件
	ws.loadFile()
	return &ws
}

// 生成一个钱包
func (ws *Wallets) CreateWallets() string {
	wallet := NewWallet()
	address := wallet.NewAddress()
	//保存
	ws.WsMap[address] = wallet
	ws.saveToFile()
	return address
}

//保存进钱包文件
func (ws *Wallets) saveToFile() {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	// 写进文件
	err = ioutil.WriteFile(FileName, buffer.Bytes(), 0600)
	if err != nil {
		log.Panic(err)
	}
}

// 加载之前的钱包文件
func (ws *Wallets) loadFile() {
	// 判断文件是否存在
	_, err := os.Stat(FileName)
	if os.IsNotExist(err) {
		return
	}
	file, err := ioutil.ReadFile(FileName)
	if err != nil {
		log.Panic(err)
	}
	// 解码
	var wls Wallets
	// 注册
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(file))
	err = decoder.Decode(&wls)
	if err != nil {
		log.Panic(err)
	}
	ws.WsMap = wls.WsMap
}

//遍历 地址
func (ws *Wallets) listAddress() []string {
	var adds []string
	for address := range ws.WsMap {
		adds = append(adds, address)
	}
	return adds
}
