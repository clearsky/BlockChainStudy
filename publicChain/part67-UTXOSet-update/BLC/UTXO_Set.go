package BLC

import (
	"bytes"
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
func (utxoSet *UTXOSet) ResetUTXOSet() {
	err := utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoSetTableName))
		if b != nil { // 如果表存在，就删除表
			err := tx.DeleteBucket([]byte(utxoSetTableName))
			if err != nil {
				log.Panic(err)
			}
		}
		// 创建新表
		b, _ = tx.CreateBucket([]byte(utxoSetTableName))
		if b != nil {
			//[string]*TXOutputs
			txOutputMap := utxoSet.BlockChain.FindUTXOMap()
			for keyHash, outs := range txOutputMap {
				txHash, _ := hex.DecodeString(keyHash)
				err := b.Put(txHash, outs.Serialize())
				if err != nil {
					log.Panic(err)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (utxoSet *UTXOSet) FindUTXOForAddress(address string) []*UTXO {
	var UTXOS []*UTXO
	err := utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoSetTableName))
		if b != nil {
			// 游标
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := DeserializeTXOutPuts(v)
				for _, utxo := range txOutputs.UTXOS {
					// 解锁
					if utxo.Output.UnlockScriptPubKeyWithAddress(address) {
						UTXOS = append(UTXOS, utxo)
					}
				}
			}
		} else {
			fmt.Printf("%s表不存在", utxoSetTableName)
			os.Exit(1)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return UTXOS
}

func (utxoSet *UTXOSet) GetBalance(address string) int64 {
	UTXOS := utxoSet.FindUTXOForAddress(address)
	var amount int64
	for _, utxo := range UTXOS {
		amount += utxo.Output.Value
	}
	return amount
}

func (utxoSet *UTXOSet) FindUnPackageSpendableUTXOS(from string,txs []*Transaction) []*UTXO{
	spentTXOutputs := make(map[string][]int)
	var UTXOs []*UTXO

	// 遍历未打包的Transaction
	for _, tx := range txs { // 多笔交易先便利前面的交易
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				// 是否能够解锁
				publicKyeHash := Base58Decode([]byte(from))
				hash160 := publicKyeHash[1 : len(publicKyeHash)-4]
				if in.UnlockRipemd160Hash(hash160) { // 如果是当前地址名下的消费
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], int(in.Vout))
				}
			}
		}

		// Vouts
		// 在同一个Transaction中，遍历输出，如果存在与输入的输出索引对应的输出，则代表花费，否则代表未花费，如果
		// 输入列表为空，也代表未花费
	work1:
		for index, out := range tx.Vouts {
			if out.UnlockScriptPubKeyWithAddress(from) {
				if len(spentTXOutputs) != 0 {
					var isSpend bool
					for txHash, indexArray := range spentTXOutputs {
						for _, i := range indexArray { // 如果存在花费
							if index == i && txHash == hex.EncodeToString(tx.TxHash) { // 如果花费和输出能够对应，则代表花费，不进行操作
								isSpend = true
								continue work1 // 如果消费了，就重新开始循环遍历vouts
							}

						}
					}
					// 否则代表未花费
					if !isSpend {
						utxo := &UTXO{
							TXHash: tx.TxHash,
							Index:  int64(index),
							Output: out,
						}
						fmt.Println(out.Value)
						UTXOs = append(UTXOs, utxo)
					}
				} else {
					utxo := &UTXO{
						TXHash: tx.TxHash,
						Index:  int64(index),
						Output: out,
					}
					UTXOs = append(UTXOs, utxo)
				}
			}
		}
	}
	return UTXOs
}

func (utxoSet *UTXOSet) FindSpendableUTXOS(from string, amount int64, txs []*Transaction)(int64, map[string][]int){
	unPackageUTXOS := utxoSet.FindUnPackageSpendableUTXOS(from, txs)
	// 计数
	var money int64
	spentdableUTXO := make(map[string][]int)
	for _, utxo := range unPackageUTXOS{
		money += utxo.Output.Value
		spentdableUTXO[hex.EncodeToString(utxo.TXHash)] = append(spentdableUTXO[hex.EncodeToString(utxo.TXHash)], int(utxo.Index))
		if money >= amount{
			return money, spentdableUTXO
		}
	}

	// 如果钱不够
	err := utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoSetTableName))
		if b != nil{
			c := b.Cursor()
			UTXOBREAK:
			for k,v := c.First();k!=nil;k,v = c.Next(){
				txOutputs := DeserializeTXOutPuts(v)
				for _, utxo := range txOutputs.UTXOS{
					money += utxo.Output.Value
					txHash := hex.EncodeToString(utxo.TXHash)
					spentdableUTXO[txHash] = append(spentdableUTXO[txHash], int(utxo.Index))
					if money >= amount{
						break UTXOBREAK
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

	if money < amount{
		fmt.Printf("余额不足")
		os.Exit(1)
	}

	return money, spentdableUTXO
}

// 更新
func (utxoSet *UTXOSet) Update(){
	// 最新的block
	block := utxoSet.BlockChain.Iterator().Next()
	ins := []*TXInput{}
	outsMap := make(map[string]*TXOutputs)
	// 找到所有我要删除的数据
	for _, tx := range block.Txs{
		for _, in := range tx.Vins{
			ins = append(ins, in)
		}

	}

	for _, tx := range block.Txs{
		utxos := []*UTXO{}
		 for index, out := range tx.Vouts{
		 	isSpent := false
		 	for _, in := range ins{
		 		if in.Vout == int64(index) && bytes.Equal(out.Ripemd160Hash, Ripemd160Hash(in.Publickey)) && bytes.Equal(tx.TxHash, in.TxHash){
					isSpent = true
		 			continue
				}
			}
		 	if isSpent == false{
		 		utxo := &UTXO{
		 			TXHash:tx.TxHash,
		 			Index:int64(index),
		 			Output:out,
				}
		 		utxos = append(utxos, utxo)
			}
		 }
		 if len(utxos)>0{
		 	txHash := hex.EncodeToString(tx.TxHash)
		 	outsMap[txHash] = &TXOutputs{UTXOS:utxos}
		 }
	}

	// 拿出来处理了，再覆盖回去
	err := utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoSetTableName))
		if b != nil{
			// 删除
			for _, in := range ins{
				txOutputsBytes := b.Get(in.TxHash)
				if len(txOutputsBytes) == 0{
					continue
				}
				txOutputs := DeserializeTXOutPuts(txOutputsBytes)
				UTXOs := []*UTXO{}
				// 判断是否需要删除
				isNeedDelete := false
				for _, utxo := range txOutputs.UTXOS{
					if in.Vout == utxo.Index && bytes.Equal(utxo.Output.Ripemd160Hash, Ripemd160Hash(in.Publickey)){
						isNeedDelete = true
					}else{
						UTXOs = append(UTXOs, utxo)
					}
				}
				if isNeedDelete{
					err := b.Delete(in.TxHash)
					if err != nil{
						log.Panic(err)
					}
					if len(UTXOs) > 0{
						preTXOutputs := outsMap[hex.EncodeToString(in.TxHash)]
						preTXOutputs.UTXOS = append(preTXOutputs.UTXOS, UTXOs...)
						outsMap[hex.EncodeToString(in.TxHash)] = preTXOutputs
					}
				}
			}
			// 新增
			for keyHash, outPuts := range outsMap{
				keyHashBytes, _ := hex.DecodeString(keyHash)
				err := b.Put(keyHashBytes, outPuts.Serialize())
				if err != nil{
					log.Panic(err)
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
}