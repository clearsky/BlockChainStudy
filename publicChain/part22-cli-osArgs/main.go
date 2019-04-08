package main

import (
	"fmt"
	"os"
)

func main(){
	args := os.Args
	fmt.Printf("args: %v\n", args[1])
}

