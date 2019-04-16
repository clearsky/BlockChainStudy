package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
)

type Wallet struct {
	// 1.私钥
	PrivateKey ecdsa.PrivateKey  // 椭圆曲线加密
	// 2.公钥
	PublickKey []byte
}

// 创建钱包
func NewWallet() *Wallet{
	privateKey, publicKey := newKeyPair()
	fmt.Println(publicKey)
	fmt.Println(privateKey)
	return &Wallet{
		PrivateKey:privateKey,
		PublickKey:publicKey,
	}
}

// 通过私钥产生公钥
func newKeyPair()(ecdsa.PrivateKey, []byte){

	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil{
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()[0])
	return *private, pubKey
}