package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutputs struct {
	TxOutputs []*TXOutput
}
//序列化
func (txOutPuts *TXOutputs) Serialize() []byte{
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(txOutPuts)
	if err != nil{
		log.Panic(err)
	}
	return result.Bytes()
}

//反序列化
func  DeserializeTXOutPuts(txOutPutsBytes []byte) *TXOutputs{
	var txOutPuts TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(txOutPutsBytes))
	err := decoder.Decode(&txOutPuts)
	if err != nil{
		log.Panic(err)
	}
	return &txOutPuts
}
