package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	msg := "http://liyuechun.org"
	// 还有UrlEncoding

	haser := sha256.New()
	haser.Write([]byte(msg))
	hash := haser.Sum(nil)
	fmt.Printf("%x\n", hash)


	//encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	//fmt.Println(encoded)
	//decoded, err := base64.StdEncoding.DecodeString(encoded)
	//if err != nil{
	//	log.Panic(err)
	//}
	//fmt.Printf("%s", decoded)


}