package main

import (
	"blockChainStudy/publicChain/part14-block-boltdb/BLC"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func main()  {
	//var preHash [32]byte
	//block := BLC.NewBlock("Genenis Block", 1, preHash[:])
	//genesisBlock := BLC.CreateGenesisBlock("GenesisBlock")
	//fmt.Println(genesisBlock)
	//blockchain := BLC.CreateBlockchainWithGenesisBlock()
	//
	//blockchain.AddBlockToBlockchain("seconde1")
	//blockchain.AddBlockToBlockchain("seconde1")
	//blockchain.AddBlockToBlockchain("seconde1")
	//blockchain.AddBlockToBlockchain("seconde1")
	//blockchain.AddBlockToBlockchain("seconde1")
	//
	//fmt.Println(blockchain)
	var preHash [32]byte
	block := BLC.NewBlock("123", 1, preHash[:])
	fmt.Println(block.Nonce, block.Hash)
	// 创建或者打开数据库
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil{
		log.Panic(err)
	}
	defer db.Close()
	// 更新数据库
	err = db.Update(func(tx *bolt.Tx) error {
		// 尝试取表对象
		b := tx.Bucket([]byte("Blocks"))
		if b == nil{
			b, err = tx.CreateBucket([]byte("Blocks"))
			if err != nil{
				log.Panic("创建表失败")
			}
		}
		err := b.Put([]byte("l"), block.Serialize())
		if err != nil{
			log.Panic(err)
		}

		return nil
	})
	if err != nil{
		log.Panic(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Blocks"))
		if b != nil{
			blockData := b.Get([]byte("l"))
			fmt.Println(BLC.DeserializeBlock(blockData))
		}

		return nil
	})
	if err != nil{
		log.Panic(err)
	}
}
