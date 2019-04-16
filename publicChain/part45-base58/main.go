package main

import (
	"blockChainStudy/publicChain/part45-base58/BLC"
	"fmt"
)

func main() {
	bytes := []byte("http://liyuechun.org")  // 相当于公钥
	bytes58 := BLC.Base58Encode(bytes)
	fmt.Printf("%x\n", bytes58)
	fmt.Printf("%s\n", bytes58)

	bytesStr := BLC.Base58Decode(bytes58)
	fmt.Printf("%s\n", bytesStr)
}