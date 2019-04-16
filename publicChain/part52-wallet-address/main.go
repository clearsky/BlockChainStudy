package main

import (
	"blockChainStudy/publicChain/part52-wallet-address/BLC"
	"fmt"
)

func main() {

	wallet := BLC.NewWallet()
	address := wallet.GetAddress()
	isValide := wallet.IsValidAddress(address)
	fmt.Printf("%s\n", address)
	fmt.Println(isValide)

}