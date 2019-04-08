package main

import (
	"blockChainStudy/publicChain/part1-Basic-Prototype/BLC"
	"fmt"
)

func main()  {
	var preHash [32]byte
	block := BLC.NewBlock("Genenis Block", 1, preHash[:])
	fmt.Println(block)
}
