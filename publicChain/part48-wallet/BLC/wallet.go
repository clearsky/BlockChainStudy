package BLC

import "crypto/ecdsa"

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

// 创建一个钱包
func NewWallet() *Wallet{
	private, public := newKye
}
