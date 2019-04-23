package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
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
	return (len(trasaction.Vins[0].TxHash) == 0) && (trasaction.Vins[0].Vout == -1)
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
	// 进行签名
	blockchain.SignTransaction(tx, wallet.PrivateKey,txs)
	return tx
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction){
	if tx.IsCoinbaseTransaction(){
		return
	}

	for _, vin := range tx.Vins{
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash  == nil{
			log.Panic("ERROR:PrevTx is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inId, vin := range txCopy.Vins{
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inId].Signature = nil // 设置签名为空
		// 通过索引取到对应的output的160hash，设置publickKey为160hash
		txCopy.Vins[inId].Publickey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		// 设置vin的txHash
		txCopy.TxHash = txCopy.Hash()
		// 吧publickKey设置为nil
		txCopy.Vins[inId].Publickey = nil

		// 签名代码
		r,s,err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil{
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)  // 生成签名
		tx.Vins[inId].Signature = signature
	}
}

func (tx *Transaction) Hash()[]byte{
	txCopy := tx
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

// 序列化
func (tx *Transaction) Serialize()[]byte{
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil{
		log.Panic(err)
	}
	return encoded.Bytes()
}

// 复制tx
func (tx *Transaction) TrimmedCopy()Transaction{
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Vins{
		inputs = append(inputs, &TXInput{  // 只复制txHash和vout索引
			TxHash:vin.TxHash,
			Vout:vin.Vout,
			Signature:nil,
			Publickey:nil,
		})
	}
	for _, vout := range tx.Vouts{
		outputs = append(outputs, &TXOutput{
			Value:vout.Value,
			Ripemd160Hash:vout.Ripemd160Hash,
		})
	}

	txCopy := Transaction{
		TxHash:tx.TxHash,
		Vins:inputs,
		Vouts:outputs,
	}

	return txCopy
}

// 验证数字签名
func (tx *Transaction) Verify(txMap map[string]Transaction) bool{
	if tx.IsCoinbaseTransaction(){
		return true
	}

	for _, vin := range tx.Vins{
		if txMap[hex.EncodeToString(vin.TxHash)].TxHash == nil{
			log.Panic("ERROR:prevTx is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, vin := range tx.Vins{
		prevTx := txMap[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inId].Signature = nil
		txCopy.Vins[inId].Publickey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inId].Publickey = nil

		// 私钥ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen/2)])
		s.SetBytes(vin.Signature[(sigLen/2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.Publickey)
		x.SetBytes(vin.Publickey[:(keyLen/2)])
		y.SetBytes(vin.Publickey[(keyLen/2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false{
			return false
		}
	}
	return true
}
