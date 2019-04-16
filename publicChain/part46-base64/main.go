package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	msg := "Hello, 世界"
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	fmt.Println(encoded)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil{
		log.Panic(err)
	}
	fmt.Printf("%s", decoded)


}