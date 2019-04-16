package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	msg := "http://liyuechun.org"
	// 还有UrlEncoding
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	fmt.Println(encoded)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil{
		log.Panic(err)
	}
	fmt.Printf("%s", decoded)


}