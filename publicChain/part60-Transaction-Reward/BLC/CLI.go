package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {}

//func (cli *CLI) addBlock(txs []*Transaction){
//	blockchain := GetBlockChainObject()
//	defer blockchain.DB.Close()
//	blockchain.AddBlockToBlockchain(txs)
//}



func printUsage(){
	fmt.Println("Usage:")
	fmt.Println("\tcreateBlockChain -address 地址 -- 创建区块链，创建创世区块")
	fmt.Println("\tsend -from 源地址数组 -to 目的地址数组 -amount 交易额数组 -- 交易明细")
	fmt.Println("\tprintBlockChain -- 输出区块链")
	fmt.Println("\tgetBalance -address 地址 -- 输出区块链")
	fmt.Println("\tcreateWallet --创建钱包")
	fmt.Println("\taddressLists --输出所有钱包地址")
}

func isValidArgs(){
	if len(os.Args) < 2{
		printUsage()
		os.Exit(1)
	}
}





func (cli *CLI) Run(){
	isValidArgs() // 判断是否使用命令
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printBlockChainCmd := flag.NewFlagSet("printBlockChain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)
	addressListsCmd := flag.NewFlagSet("addressLists", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from","", "转账源地址...")
	flagTo := sendBlockCmd.String("to","", "转账目的地地址...")
	flagAmout := sendBlockCmd.String("amount","", "转账金额...")
	flagCreateBlockAddress := createBlockChainCmd.String("address","", "创建创世区块的地址")
	flagGetBalanceWithAddress := getBalanceCmd.String("address", "", "需要查询余额的地址")
	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "printBlockChain":
		err := printBlockChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "createBlockChain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "getBalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "createWallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "addressLists":
		err := addressListsCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	default:
		fmt.Printf("命令\"%s\"不是一个有效的命令\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed(){
		from := *flagFrom
		to := *flagTo
		amount := *flagAmout
		if  from == "" || to == "" || amount == ""{
			fmt.Println("交易明细不能为空...")
			fmt.Println(from, to, amount)
			printUsage()
			os.Exit(1)
		}
		jsonFrom := JSONToArray(from)
		jsonTo := JSONToArray(to)
		jsonAmount := JSONToArray(amount)

		for index, fromAddress := range jsonFrom{
			if !IsValidAddress([]byte(fromAddress)) || !IsValidAddress([]byte(jsonTo[index])){
				fmt.Println("地址无效...")
				printUsage()
				os.Exit(1)
			}
		}

		fmt.Println(jsonFrom)
		fmt.Println(jsonTo)
		fmt.Println(jsonAmount)

		// cli.addBlock([]*Transaction{})
		//fmt.Println(from)
		//fmt.Println(to)
		//fmt.Println(amount)
		//
		//fmt.Println(JSONToArray(from))
		//fmt.Println(JSONToArray(to))
		//fmt.Println(JSONToArray(amount))

		cli.send(jsonFrom, jsonTo, jsonAmount)
	}

	if printBlockChainCmd.Parsed(){
		cli.printBlockChainChain()
	}

	if createBlockChainCmd.Parsed(){
		address := *flagCreateBlockAddress
		if  !IsValidAddress([]byte(address)){
			fmt.Println("地址无效...")
			printUsage()
			os.Exit(1)
		}

		cli.createBlockChain(address)
	}

	if getBalanceCmd.Parsed(){
		address := *flagGetBalanceWithAddress
		if  !IsValidAddress([]byte(address)){
			fmt.Println("地址无效...")
			printUsage()
			os.Exit(1)
		}
		// 调用查询余额的函数
		cli.getBalance(address)
	}

	// 创建钱包
	if createWalletCmd.Parsed(){
		cli.CreateWallet()
	}

	// 输出所有创建的地址
	if addressListsCmd.Parsed(){
		cli.AddressLists()
	}
}