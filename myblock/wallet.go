package main

import (
	"blockchain/lib/base58"
	"blockchain/lib/ripemd160"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

//创建一个新钱包
func NewWallet() *Wallet {
	// 生成公私钥
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	publicKey := key.PublicKey
	// 拼接 rs
	pubKey := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	// 返回对象
	return &Wallet{
		PublicKey:  pubKey,
		PrivateKey: key,
	}
}

//创建一个 地址
func (w *Wallet) NewAddress() string {
	// 对公钥 进行 sha256 ripemd160
	key := HashPublicKey(w.PublicKey)
	// 加上version
	version := byte(00)
	//拼接version
	payload := append([]byte{version}, key...)
	// 对 hash 进行 sha256 取 前4
	sum := CheckSum(payload)
	payload = append(payload, sum...)
	// base58
	return base58.Encode(payload)
}

func HashPublicKey(pk []byte) []byte {
	// 对公钥 进行 sha256 ripemd160
	sum256 := sha256.Sum256(pk)
	hash := ripemd160.New()
	_, err := hash.Write(sum256[:])
	if err != nil {
		log.Panic(err)
	}
	return hash.Sum(nil)
}

// 取前4位
func CheckSum(data []byte) []byte {
	// 2次 sha256
	sum1 := sha256.Sum256(data)
	sum2 := sha256.Sum256(sum1[:])
	return sum2[:4]
}

// 获取地址的公钥hash
func getPubKeyHashByAddress(address string) []byte {
	if len(address) < 4 {
		fmt.Printf("无效的地址：%s/n",address)
		return nil
	}
	// base58
	code := base58.Decode(address)
	//截 除最后4个  和 第一个 version
	return code[1 : len(code)-4]

}
