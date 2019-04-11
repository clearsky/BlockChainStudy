package BLC

type TXInout struct {
	// 1. 交易的hash
	Txid []byte
	// 2. 存储TXOutput在Vout里面的索引
	Vout int
	// 3.用户名
	ScriptSig string
}