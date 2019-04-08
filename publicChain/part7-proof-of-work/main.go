package main

import (
	"blockChainStudy/publicChain/part7-proof-of-work/BLC"
	"fmt"
)

func main()  {
	//var preHash [32]byte
	//block := BLC.NewBlock("Genenis Block", 1, preHash[:])
	//genesisBlock := BLC.CreateGenesisBlock("GenesisBlock")
	//fmt.Println(genesisBlock)
	blockchain := BLC.CreateBlockchainWithGenesisBlock()

	blockchain.AddBlockToBlockchain("seconde1")
	blockchain.AddBlockToBlockchain("seconde1")
	blockchain.AddBlockToBlockchain("seconde1")
	blockchain.AddBlockToBlockchain("seconde1")
	blockchain.AddBlockToBlockchain("seconde1")

	fmt.Println(blockchain)
}
