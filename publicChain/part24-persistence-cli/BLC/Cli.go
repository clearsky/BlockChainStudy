package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	Blockchain *Blockchain
}

func (cli *CLI) addBlock(data string){
	cli.Blockchain.AddBlockToBlockchain(data)
}

func (cli *CLI) printChain(){
	cli.Blockchain.Printchain()
}

func printUsage(){
	fmt.Println("Usage:")
	fmt.Println("\tcreateBlockChain -data DATA -- 创建区块链，创建创世区块")
	fmt.Println("\taddBlock -data DATA -- 交易数据")
	fmt.Println("\tprintChain -- 输出区块链")
}

func isValidArgs(){
	if len(os.Args) < 2{
		printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) createGenesisBlockChain(data string){
	fmt.Println(data)
}

func (cli *CLI) Run(){
	isValidArgs() // 判断是否使用命令
	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data","", "交易数据")
	flagCreateBlockData := createBlockChainCmd.String("data","Genesis Data...", "创世区块交易数据")
	switch os.Args[1] {
	case "addBlock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
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
			printUsage()
			os.Exit(1)
		}
		cli.addBlock(data)
	}

	if printChainCmd.Parsed(){
		cli.printChain()
	}

	if createBlockChainCmd.Parsed(){
		data := *flagCreateBlockData
		if  data == ""{
			fmt.Println("交易数据不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockChain(data)
	}
}