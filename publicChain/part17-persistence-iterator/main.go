package main

import (
	"blockChainStudy/publicChain/part17-persistence-iterator/BLC"
)

func main()  {
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	defer blockchain.DB.Close()

	blockchain.AddBlockToBlockchain("from mary send 100c to tom")
	blockchain.AddBlockToBlockchain("from mary send 100c to jim")
	blockchain.AddBlockToBlockchain("from tom send 100c to bob")
	blockchain.AddBlockToBlockchain("from tom send 100c to tina")
	blockchain.AddBlockToBlockchain("from bob send 100c to tina")

	blockchain.PrintChain()

}
