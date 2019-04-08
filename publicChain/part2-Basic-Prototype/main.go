package main

import (
	"blockChainStudy/publicChain/part2-Basic-Prototype/BLC"
	"fmt"
)

func main()  {
	//var preHash [32]byte
	//block := BLC.NewBlock("Genenis Block", 1, preHash[:])
	genesisBlock := BLC.CreateGenesisBlock("GenesisBlock")
	fmt.Println(genesisBlock)
}
