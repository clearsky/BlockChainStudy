package BLC

import (
	"bytes"
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	CurrentHash []byte
	DB *bolt.DB
}

func (blockchainIterator *BlockchainIterator) Next() *Block{

	var block *Block
	var end [32]byte

	if bytes.Equal(blockchainIterator.CurrentHash, end[:]){
		return nil
	}
	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil{
			currentBlockBytes := b.Get(blockchainIterator.CurrentHash)
			block = DeserializeBlock(currentBlockBytes)
			blockchainIterator.CurrentHash = block.PrevBlockHash
		}
		return nil
	})
	if err != nil{
		log.Panic(err)
	}
	return block
}

