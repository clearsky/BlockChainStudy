package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

// 数据库名字
const  dbName = "blockchain.db"
// 表的名字
const  blockTableName = "blocks"

type Blockchain struct {
	//Blocks []*Block //存储有序的区块
	Tip []byte  // 最新区块的hash
	DB *bolt.DB
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator{
	return &BlockchainIterator{
		CurrentHash: blockchain.Tip,
		DB: blockchain.DB,
	}
}

// 遍历输出所有区块的信息
//func (blc *Blockchain) PrintChain(){
//
//	var block *Block
//	var currenHash []byte = blc.Tip
//	var end [32]byte
//	err := blc.DB.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(blockTableName))
//		if b != nil{
//			for {
//				// 获取当前区块的字节数组
//				blockBytes := b.Get(currenHash)
//				block = DeserializeBlock(blockBytes)
//				block.PrintInfo()
//				if bytes.Equal(block.PrevBlockHash, end[:]){
//					break
//				}
//				currenHash = block.PrevBlockHash
//			}
//		}
//
//		return nil
//	})
//	if err != nil{
//		log.Panic(err)
//	}
//}
func (blc *Blockchain) Printchain(){
	blockchainIterator := blc.Iterator()
	var block *Block
	for {
		block = blockchainIterator.Next()
		if block != nil{
			block.PrintInfo()
		}else{
			break
		}

	}

}

// 判断区块链是否存在
func blockChainExists(db *bolt.DB) bool{
	// 判断blockTable是否存在
	blockTableNameIsExist := false  // 默认blockTable不存在
	err :=db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b !=nil{
			tip := b.Get([]byte("l"))
			if tip != nil{
				blockTableNameIsExist = true  // 表存在，最新区块的hash存在，则判断区块链存在
			}
		}
		return nil
	})
	if err != nil{
		log.Panic(err)
	}
	if blockTableNameIsExist{
		return true
	}else{
		return false
	}
}

// 获取数据库
func getDB() *bolt.DB{
	//获取数据库，如果不存在，会自动创建数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil{
		log.Panic(err)
	}
	return db
}

// 1.创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(data string) *Blockchain{
	//获取数据库
	db := getDB()

	//如果blcokTable存在
	if blockChainExists(db){
		fmt.Println("创世区块已存在...")
		var blockchain Blockchain
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blockTableName))
			tip := b.Get([]byte("l"))
			blockchain = Blockchain{
				Tip: tip,
				DB:db,
			}
			return nil
		})
		if err != nil{
			log.Panic(err)
		}
		return &blockchain
	}

	// 如果blockTable不存在
	var blockHash []byte
	err := db.Update(func(tx *bolt.Tx) error {
		//创建blockTable
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil{
			log.Panic(err)
		}
		//创建创世区块
		genesis := CreateGenesisBlock(data)
		// 将创世区块存储到blockChain中
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil{
			log.Panic(err)
		}
		// 存储最新的区块的hash
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil{
			log.Panic(err)
		}
		blockHash = genesis.Hash
		return nil
	})
	if err != nil{
		log.Panic(err)
	}
	// 创建创世区块

	// 返回区块链对象
	return &Blockchain{
		Tip: blockHash,
		DB: db,
	}
}

// 2.增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(data string){
	// 获取上一个区块的信息

	err := blc.DB.Update(func(tx *bolt.Tx) error {
		// 1.获取表
		b := tx.Bucket([]byte(blockTableName))

		if b != nil{
			// 2.创建新区块
			block := DeserializeBlock(b.Get(blc.Tip))
			newBlock := NewBlock(data,block.Height + 1, block.Hash)
			// 3.将区块序列化并且存储到数据库中
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil{
				log.Panic(err)
			}
			// 4.更新数据库里"l"对应的hash
			err = b.Put([]byte("l"), newBlock.Hash)

			// 5.跟新blockchain的Tip
			blc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil{
		log.Panic(err)
	}
}