package main

import (
	"blockChainStudy/publicChain/part53-wallets-address/BLC"
	"fmt"
)

func main() {

	wallets := BLC.NewWallets()
	fmt.Println(wallets)
	wallets.CreateNewWallet()
	fmt.Println(wallets)
}