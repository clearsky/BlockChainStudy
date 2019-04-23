package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

// 1.有一个方法，功能：
// 遍历整个数据库，遍历所有的未花费的UTXO，然后将所有的UTXO存储到数据库
// 去遍历数据库时
// [string]*TXOutputs

//txHash, TXOutputs := range txOutputs{
//
//}

const utxoSetTableName = "utxoSetTableName"
type UTXOSet struct {
	BlockChain *Blockchain
}

// 重置数据库表
func (utxoSet *UTXOSet) ResetUTXOSet(){
	err := utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoSetTableName))
		if b!= nil{ // 如果表存在，就删除表
			err := tx.DeleteBucket([]byte(utxoSetTableName))
			if err != nil{
				log.Panic(err)
			}
		}
		// 创建新表
		b, _ =tx.CreateBucket([]byte(utxoSetTableName))
		if b != nil{
			//[string]*TXOutputs
			txOutputMap := utxoSet.BlockChain.FindUTXOMap()
			for keyHash, outs := range txOutputMap{
				txHash,_ := hex.DecodeString(keyHash)
				err := b.Put(txHash, outs.Serialize())
				if err != nil{
					log.Panic(err)
				}
			}
		}

		return nil
	})
	if err != nil{
		log.Panic(err)
	}
}

func (utxoSet *UTXOSet) FindUTXOForAddress(address string)[]*UTXO{
	var UTXOS []*UTXO
	err := utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoSetTableName))
		if b != nil{
			// 游标
			c := b.Cursor()
			for k, v := c.First();k!=nil;k, v = c.Next(){
				txOutputs := DeserializeTXOutPuts(v)
				for _, utxo := range txOutputs.UTXOS{
					// 解锁
					if utxo.Output.UnlockScriptPubKeyWithAddress(address){
						UTXOS = append(UTXOS, utxo)
					}
				}
			}
		}else{
			fmt.Printf("%s表不存在", utxoSetTableName)
			os.Exit(1)
		}
		return nil
	})
	if err != nil{
		log.Panic(err)
	}
	return UTXOS
}

func (utxoSet *UTXOSet) GetBalance(address string) int64{
	UTXOS := utxoSet.FindUTXOForAddress(address)
	var amount int64
	for _, utxo := range UTXOS{
		amount += utxo.Output.Value
	}
	return amount
}