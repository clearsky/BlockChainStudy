package BLC

type TXInput struct {
	// 1. 交易的hash
	TxHash []byte
	// 2. 存储TXOutput在Vout里面的索引
	Vout int64
	// 3.用户名
	ScriptSig string
}

// 判断当前消费是否由传入地址消费的
func (txInput *TXInput) UnlockScriptSigWithAddress(address string) bool{
	return txInput.ScriptSig == address
}