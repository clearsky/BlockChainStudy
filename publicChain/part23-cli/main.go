package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func printUsage(){
	fmt.Println("Usage:")
	fmt.Println("\taddBlock -data DATA -- 交易数据")
	fmt.Println("\tprintChain -- 输出区块链")
}

func isValidArgs(){
	if len(os.Args) < 2{
		printUsage()
		os.Exit(1)
	}
}

func main(){
	isValidArgs()
	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)

	flageAddBlockData := addBlockCmd.String("data","", "交易数据")
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
	default:
		printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed(){
		if *flageAddBlockData == ""{
			printUsage()
			os.Exit(1)
		}
		fmt.Println(*flageAddBlockData)
	}

	if printChainCmd.Parsed(){
		fmt.Println("输出所有区块的数据")
	}
}

