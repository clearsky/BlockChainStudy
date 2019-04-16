package BLC

type UTXO struct {
	TXHash []byte
	Index int
	Output *TXOutput
}
