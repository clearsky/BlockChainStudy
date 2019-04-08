package main

import (
	"blockChainStudy/publicChain/part15-persistence/BLC"
	"fmt"
)

func main()  {
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	defer blockchain.DB.Close()

	//blockchain.AddBlockToBlockchain("seconde1")
	//blockchain.AddBlockToBlockchain("seconde1")
	//blockchain.AddBlockToBlockchain("seconde1")
	//blockchain.AddBlockToBlockchain("seconde1")
	//blockchain.AddBlockToBlockchain("seconde1")

	fmt.Println(blockchain)
}
