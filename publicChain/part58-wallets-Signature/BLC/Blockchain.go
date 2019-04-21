package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strconv"
)

// 数据库名字
const dbName = "blockchain.db"

// 表的名字
const blockTableName = "blocks"

type Blockchain struct {
	//Blocks []*Block //存储有序的区块
	Tip []byte // 最新区块的hash
	DB  *bolt.DB
}

// 如果一个地址对应的TxOutput未花费， 那么这个Transaction就应该添加到数组中返回
func (blockchain *Blockchain) GetUTXOs(address string, txs []*Transaction) []*UTXO {
	blockchainIterator := blockchain.Iterator()
	spentTXOutputs := make(map[string][]int)
	var end [32]byte
	var UTXOs []*UTXO

	for _, tx := range txs{  // 多笔交易先便利前面的交易
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				// 是否能够解锁
				publicKyeHash := Base58Decode([]byte(address))
				hash160 := publicKyeHash[1:len(publicKyeHash)-4]
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
			if out.UnlockScriptPubKeyWithAddress(address){
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
					if !isSpend{
						utxo := &UTXO{
							TXHash: tx.TxHash,
							Index:int64(index),
							Output:out,
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

	for {
		block := blockchainIterator.Next()
		for i := len(block.Txs)-1; i >=0; i-- {  // 注意顺序问题,从最后一个Transaction往前遍历
			tx := block.Txs[i]
			// Vins
			// 遍历一个Transaction的指定地址输入，将Transaction的hash和输入的输出索引存入字典
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					// 是否能够解锁
					publicKyeHash := Base58Decode([]byte(address))
					hash160 := publicKyeHash[1:len(publicKyeHash)-4]
					if in.UnlockRipemd160Hash(hash160) { // 如果是当前地址名下的消费
						key := hex.EncodeToString(in.TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], int(in.Vout))
					}
				}
			}

			// Vouts
			// 在同一个Transaction中，遍历输出，如果存在与输入的输出索引对应的输出，则代表花费，否则代表未花费，如果
			// 输入列表为空，也代表未花费
		work:
			for index, out := range tx.Vouts {
				if out.UnlockScriptPubKeyWithAddress(address){
					if len(spentTXOutputs) != 0 {
						var isSpend bool
						for txHash, indexArray := range spentTXOutputs {
							for _, i := range indexArray { // 如果存在花费
								if index == i && txHash == hex.EncodeToString(tx.TxHash) { // 如果花费和输出能够对应，则代表花费，不进行操作
									isSpend = true
									continue work // 如果消费了，就重新开始循环遍历vouts
								}

							}
						}
						// 否则代表未花费
						if !isSpend{
							utxo := &UTXO{
								TXHash: tx.TxHash,
								Index:int64(index),
								Output:out,
							}
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

		if bytes.Equal(block.PrevBlockHash, end[:]) { // 遍历到创世区块，退出循环
			break
		}
	}
	return UTXOs
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		CurrentHash: blockchain.Tip,
		DB:          blockchain.DB,
	}
}

func (blc *Blockchain) Printchain() {
	blockchainIterator := blc.Iterator()
	var block *Block
	for {
		block = blockchainIterator.Next()
		if block != nil {
			block.PrintInfo()
		} else {
			break
		}

	}

}

// 判断区块链是否存在
func blockChainExists(db *bolt.DB) bool {
	// 判断blockTable是否存在
	blockTableNameIsExist := false // 默认blockTable不存在
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			tip := b.Get([]byte("l"))
			if tip != nil {
				blockTableNameIsExist = true // 表存在，最新区块的hash存在，则判断区块链存在
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	if blockTableNameIsExist {
		return true
	} else {
		return false
	}
}

// 获取数据库
func getDB() *bolt.DB {
	//获取数据库，如果不存在，会自动创建数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	return db
}

// 1.创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(address string) *Blockchain {
	//获取数据库
	db := getDB()
	var blockchain Blockchain
	//如果blcokTable存在
	if blockChainExists(db) {
		fmt.Println("创世区块已存在...")
		os.Exit(1)
	} else {
		fmt.Println("正在创建创世区块...")
		// 如果blockTable不存在
		err := db.Update(func(tx *bolt.Tx) error {
			//创建blockTable
			b, err := tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panic(err)
			}
			txCoinbase := NewCoinbaseTransaction(address)
			//创建创世区块
			genesis := CreateGenesisBlock([]*Transaction{txCoinbase})
			// 将创世区块存储到blockChain中
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 存储最新的区块的hash
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			fmt.Println("创世区块创建成功...")
			blockchain = Blockchain{
				Tip: genesis.Hash,
				DB:  db,
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
	}
	return &blockchain
}

// 2.增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {
	// 获取上一个区块的信息
	fmt.Println("正在添加新的区块...")
	err := blc.DB.Update(func(tx *bolt.Tx) error {
		// 1.获取表
		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			// 2.创建新区块
			block := DeserializeBlock(b.Get(blc.Tip))
			newBlock := NewBlock(txs, block.Height+1, block.Hash)
			// 3.将区块序列化并且存储到数据库中
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 4.更新数据库里"l"对应的hash
			err = b.Put([]byte("l"), newBlock.Hash)

			// 5.跟新blockchain的Tip
			blc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("新的区块添加成功...")
}

// 返回BlockChain对象
func GetBlockChainObject() *Blockchain {
	db := getDB()
	if blockChainExists(db) {
		var blockchain Blockchain
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blockTableName))
			tip := b.Get([]byte("l"))
			blockchain = Blockchain{
				Tip: tip,
				DB:  db,
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		return &blockchain
	} else {
		fmt.Println("区块链未创建...")
		db.Close()
		os.Exit(1)
		return nil
	}
}

// 挖掘新的区块
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {

	var txs []*Transaction

	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], int64(value), blockchain, txs)
		txs = append(txs, tx)
	}

	var block *Block
	err := blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("l"))           // 获取最新区块的hash值
			blockBytes := b.Get(hash)            // 获取最新区块的二进制
			block = DeserializeBlock(blockBytes) // 反序列化获得最新的区块
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	// 建立新的区块
	block = NewBlock(txs, block.Height+1, block.Hash)

	// 将新区块存储到数据库
	err = blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			err = b.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panic(err)
			}
			blockchain.Tip = block.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

}

// 查询余额
func (blockchain *Blockchain) GetBalance(address string) int64 {
	utxos := blockchain.GetUTXOs(address, []*Transaction{})
	var amount int64
	for _, utxo := range utxos {
		amount = amount + utxo.Output.Value
	}
	return amount
}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpendableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int64) {
	// 1.获取所有的UTXO
	utxos := blockchain.GetUTXOs(from, txs)
	spendableUTXOs := make(map[string][]int64)
	var value int64
	// 2.遍历utxos
	for _, utxo := range utxos {
		value = value + utxo.Output.Value
		hash := hex.EncodeToString(utxo.TXHash)
		spendableUTXOs[hash] = append(spendableUTXOs[hash], utxo.Index)
		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) { // 余额不足的情况
		fmt.Printf("%s 余额不足\n", from)
		os.Exit(1)
	}
	return value, spendableUTXOs
}

func (blockchain *Blockchain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey){
	if tx.IsCoinbaseTransaction(){
		return
	}
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vins{
		prevTX, err := blockchain.FindTransaction(vin.TxHash)
		if err != nil{
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}
	tx.Sign(privateKey, prevTXs)
}

// 通过txHash查找tx
func (blockchain *Blockchain) FindTransaction(TxHash []byte)(Transaction, error){
	bci := blockchain.Iterator()
	end := [32]byte{}
	for{
		block := bci.Next()
		for _, tx := range block.Txs{
			if bytes.Equal(TxHash, tx.TxHash){
				return *tx, nil
			}
		}
		if bytes.Equal(block.PrevBlockHash, end[:]){
			break
		}
	}
	return Transaction{}, nil
}
