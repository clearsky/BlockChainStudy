package BLC

import "fmt"

// 打印所有的钱包地址
func (cli *CLI)AddressLists(){
	wallets, _ := NewWallets()
	for address := range wallets.Wallets{
		fmt.Println(address)
	}
}