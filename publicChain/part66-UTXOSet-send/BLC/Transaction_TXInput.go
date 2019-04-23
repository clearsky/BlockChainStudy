package BLC

import "bytes"

type TXInput struct {
	// 1. 交易的hash
	TxHash []byte
	// 2. 存储TXOutput在Vout里面的索引
	Vout int64
	// 3. 数字签名
	Signature []byte
	// 4. 公钥
	Publickey []byte
}



// 判断当前消费是否由传入地址消费的
func (txInput *TXInput) UnlockRipemd160Hash(ripemd160Hash []byte) bool{
	encoded := Ripemd160Hash(txInput.Publickey)
	return bytes.Equal(encoded, ripemd160Hash)
}