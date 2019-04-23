package BLC

import (
	"github.com/boltdb/bolt"
	"log"
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
		}

		return nil
	})
	if err != nil{
		log.Panic(err)
	}
}