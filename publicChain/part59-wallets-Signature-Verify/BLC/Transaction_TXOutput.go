package BLC

import "bytes"

type TXOutput struct {
	Value int64
	// 公钥
	Ripemd160Hash []byte // 公钥进行256、160hash后的值
}
func (txOutput *TXOutput) Lock(address string){
	publicKeyHash := Base58Decode([]byte(address))
	txOutput.Ripemd160Hash = publicKeyHash[1: len(publicKeyHash)-4]
}

// 解锁
func (txOutput *TXOutput) UnlockScriptPubKeyWithAddress(address string) bool{
	publicKeyHash := Base58Decode([]byte(address))
	hash160 := publicKeyHash[1: len(publicKeyHash)-4]
	return bytes.Equal(txOutput.Ripemd160Hash, hash160)
}

// 创建新的TXOutput
func NewTXOutput(value int64, address string) *TXOutput{
	txOutput := &TXOutput{Value:value, Ripemd160Hash:nil}

	// 设置Ripemd160Hash
	txOutput.Lock(address)

	return txOutput
}