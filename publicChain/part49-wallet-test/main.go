package main

import (
	"blockChainStudy/publicChain/part45-base58/BLC"
	"fmt"
)

func main() {

	// 1.创建钱包
	// (1) 私钥
	// (2) 公钥

	// 2.先将公钥进行一次256hash，再进行一次160hash
	// 20字节的[]byte

	// version {0} + hash160  ->pubkey
	// 256hash pubkey
	// 256 64
	// 最后的四个字节取出来
	// version {} + hash160 +4个字节 -> 25个字节
	// base58 编码

	var msg [23]byte
	// 还有UrlEncoding

	bytes58 := BLC.Base58Decode(msg[:])
	fmt.Printf("%x\n", bytes58)
}