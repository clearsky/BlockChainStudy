package main

import (
	"blockChainStudy/publicChain/part9-Serialize-DeserializeBlock/BLC"
	"fmt"
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
	proofOfWord := BLC.NewProofOfWork(block)
	fmt.Println(proofOfWord.IsValid())

	fmt.Println("------------")
	bytes := block.Serialize()
	fmt.Println(bytes)
	fmt.Println(BLC.DeserializeBlock(bytes))
}
