package BLC

import (
	"fmt"
	"os"
)

// 获取余额
func (cli *CLI) getBalance(address string){
	fmt.Println(address)
	blockchain := GetBlockChainObject()
	if blockchain == nil{
		fmt.Println("区块链为空")
		printUsage()
		os.Exit(1)
	}
	amount := blockchain.GetBalance(address)
	fmt.Printf("%一共有%d个Token\n", address, amount)
}