package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// UTXO
type Transaction struct {
	// 1. 交易hash
	TxHash []byte
	// 2. 输入
	Vins []*TXInput
	// 3. 输出
	Vouts []*TXOutput
}


// 1. Transactino创建
	// 1. 创世区块创建时的Transactino
	func NewCoinbaseTransaction(address string) *Transaction{
		txInput := &TXInput{
			TxHash: []byte{},
			Vout: -1,
			ScriptSig: "Gensis Data...",
		}
		txOutput := &TXOutput{
			Value:10,
			ScriptPubKey: address,
		}
		txCoinbase := &Transaction{
			TxHash:[]byte{},
			Vins: []*TXInput{txInput},
			Vouts: []*TXOutput{txOutput},
		}
		// 设置hash值
		txCoinbase.SetHashTransaction()

		return txCoinbase

	}

//
	func (tx * Transaction) SetHashTransaction(){
		var result bytes.Buffer

		encoder := gob.NewEncoder(&result)

		err := encoder.Encode(tx)
		if err != nil{
			log.Panic(err)
		}
		hash := sha256.Sum256(result.Bytes())

		tx.TxHash = hash[:]
	}
	// 2. 转账时产生的Transaction
func NewSimpleTransaction(from string, to string, amount int) *Transaction{
	var txInputs []*TXInput
	var txOutputs []*TXOutput

	// 消费
		// 消费hash缺失
	txInput := &TXInput{
		TxHash: []byte{},
		Vout: -1,
		ScriptSig: "Gensis Data...",
	}
	txInputs = append(txInputs, txInput)

	// 转账
	txOutput := &TXOutput{
		Value:10,
		ScriptPubKey: to,
	}
	txOutputs = append(txOutputs, txOutput)

	// 找零
	txOutput = &TXOutput{
		Value:nil,
		ScriptPubKey:from,
	}
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{
		TxHash:[]byte{},
		Vins: txInputs,
		Vouts: txOutputs,
	}
	tx.SetHashTransaction()

	return tx
}
