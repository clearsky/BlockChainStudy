package BLC

import (
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



// 1.创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain{
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil{
		log.Panic(err)
	}
	var blockHash []byte
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil{
			log.Panic(err)
		}
		if b != nil{
			genesis := CreateGenesisBlock("GenesisBlock Data...")
			// 将创世区块存储到表中
			err := b.Put(genesis.Hash, genesis.Serialize())
			if err != nil{
				log.Panic(err)
			}

			// 存储最新的区块的hash
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil{
				log.Panic(err)
			}
			blockHash = genesis.Hash
		}
		return nil
	})
	// 创建创世区块

	// 返回区块链对象
	return &Blockchain{
		Tip: blockHash,
		DB: db,
	}
}

// 2.增加区块到区块链里面
//func (blc *Blockchain) AddBlockToBlockchain(data string){
//	height := blc.Blocks[len(blc.Blocks)-1].Height + 1
//	preHash := blc.Blocks[len(blc.Blocks)-1].Hash
//	newBlock := NewBlock(data, height, preHash)
//	blc.Blocks = append(blc.Blocks, newBlock)
//}