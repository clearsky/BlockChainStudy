package BLC

import (
	"bytes"
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

// 遍历输出所有区块的信息
func (blc *Blockchain) PrintChain(){

	var block *Block
	var currenHash []byte = blc.Tip
	var end [32]byte
	err := blc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil{
			for {
				// 获取当前区块的字节数组
				blockBytes := b.Get(currenHash)
				block = DeserializeBlock(blockBytes)
				fmt.Println("height: ", block.Height)
				fmt.Printf("preHash: %x\n", block.PrevBlockHash)
				fmt.Println("nonce: ", block.Nonce)
				fmt.Println("data: ", string(block.Data))
				fmt.Printf("hash: %x\n\n", block.Hash)
				if bytes.Equal(block.PrevBlockHash, end[:]){
					break
				}
				currenHash = block.PrevBlockHash
			}
		}

		return nil
	})
	if err != nil{
		log.Panic(err)
	}
}

// 1.创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain{
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil{
		log.Panic(err)
	}
	var blockHash []byte
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b == nil{
			b, err = tx.CreateBucket([]byte(blockTableName))
			if err != nil{
				log.Panic(err)
			}
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