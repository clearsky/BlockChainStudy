package main

import (
	"flag"
	"fmt"
)

func main(){
	flagString := flag.String("printchain", "haha", "输出所有的区块信息")
	flagIng := flag.Int("number", 6, "输出一个整数")
	flageBool := flag.Bool("open", false, "判断真假")  // 不需要传真假，只要传了就为真
	flag.Parse()
	fmt.Printf("%s\n", *flagString)
	fmt.Printf("%d\n", *flagIng)
	fmt.Printf("%v\n", *flageBool)
}

