package main

import (
	"blockChainStudy/publicChain/part51-wallet-address/BLC"
	"fmt"
)

func main() {

	wallet := BLC.NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("%s",address)

}