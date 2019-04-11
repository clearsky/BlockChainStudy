package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {}

func (cli *CLI) addBlock(data string){
	blockchain := GetBlockChainObject()
	defer blockchain.DB.Close()
	blockchain.AddBlockToBlockchain(data)
}

func (cli *CLI) printBlockChainChain(){
	blockchain := GetBlockChainObject()
	defer blockchain.DB.Close()
	blockchain.Printchain()
}

func printUsage(){
	fmt.Println("Usage:")
	fmt.Println("\tcreateBlockChain -data DATA -- 创建区块链，创建创世区块")
	fmt.Println("\taddBlock -data DATA -- 交易数据")
	fmt.Println("\tprintBlockChain -- 输出区块链")
}

func isValidArgs(){
	if len(os.Args) < 2{
		printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) createBlockChain(data string){
	CreateBlockchainWithGenesisBlock(data)
}

func (cli *CLI) Run(){
	isValidArgs() // 判断是否使用命令
	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printBlockChainCmd := flag.NewFlagSet("printBlockChain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data","", "交易数据")
	flagCreateBlockData := createBlockChainCmd.String("data","Genesis Data...", "创世区块交易数据")
	switch os.Args[1] {
	case "addBlock":
		err := addBlockCmd.Parse(os.Args[2:])
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
	default:
		fmt.Printf("命令\"%s\"不是一个有效的命令\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed(){
		data := *flagAddBlockData
		if  data == ""{
			fmt.Println("交易数据不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.addBlock(data)
	}

	if printBlockChainCmd.Parsed(){
		cli.printBlockChainChain()
	}

	if createBlockChainCmd.Parsed(){
		data := *flagCreateBlockData
		if  data == ""{
			fmt.Println("交易数据不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.createBlockChain(data)
	}
}