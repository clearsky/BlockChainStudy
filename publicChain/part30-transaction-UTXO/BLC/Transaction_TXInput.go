package BLC

type TXInput struct {
	// 1. 交易的hash
	TxHash []byte
	// 2. 存储TXOutput在Vout里面的索引
	Vout int
	// 3.用户名
	ScriptSig string
}