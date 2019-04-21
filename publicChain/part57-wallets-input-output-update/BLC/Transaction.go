package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
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

// 判断当前的交易是否是Coinbase交易
func (trasaction *Transaction) IsCoinbaseTransaction() bool{
	return len(trasaction.Vins[0].TxHash) == 0 && trasaction.Vins[0].Vout == -1
}

// 1. Transactino创建
	// 1. 创世区块创建时的Transactino
	func NewCoinbaseTransaction(address string) *Transaction{
		txInput := &TXInput{
			TxHash: []byte{},
			Vout: -1,
			Signature:nil,
			Publickey: []byte{},
		}
		txOutput := NewTXOutput(10, address)
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
func NewSimpleTransaction(from string, to string, amount int64, blockchain *Blockchain, txs []*Transaction) *Transaction{

	//UTXOs := blockchain.GetUTXOs(from)
	money, spendableUTXOs := blockchain.FindSpendableUTXOs(from , amount, txs)

	// 1. 函数，返回from这个人所有的未花费交易输出所对应的Transaction
	// 通过一个函数，返回未花费的余额

	wallets, _ := NewWallets()
	wallet := wallets.Wallets[from]
	var txInputs []*TXInput
	var txOutputs []*TXOutput

	// 消费
	for txHash, indexArray := range spendableUTXOs{
		for _, index := range indexArray{
			txHashBytes, _ := hex.DecodeString(txHash)
			txInput := &TXInput{
				TxHash: txHashBytes,
				Vout: index,
				Signature:nil,
				Publickey: wallet.PublickKey,
			}
			txInputs = append(txInputs, txInput)
		}
	}


	// 转账
	txOutput := NewTXOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)

	// 找零
	txOutput = NewTXOutput(money-amount, from)
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{
		TxHash:[]byte{},
		Vins: txInputs,
		Vouts: txOutputs,
	}
	tx.SetHashTransaction()

	return tx
}
