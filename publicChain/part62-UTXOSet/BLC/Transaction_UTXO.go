package BLC

type UTXO struct {
	TXHash []byte
	Index int64
	Output *TXOutput
}
