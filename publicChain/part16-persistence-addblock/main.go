package main

import (
	"blockChainStudy/publicChain/part16-persistence-addblock/BLC"
)

func main()  {
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	defer blockchain.DB.Close()

	blockchain.AddBlockToBlockchain("seconde1")
	blockchain.AddBlockToBlockchain("seconde1")
	blockchain.AddBlockToBlockchain("seconde1")
	blockchain.AddBlockToBlockchain("seconde1")
	blockchain.AddBlockToBlockchain("seconde1")

}
