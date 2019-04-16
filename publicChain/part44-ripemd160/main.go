package main

import (
	"fmt"
	"golang.org/x/crypto/ripemd160"
)

func main() {
	haser := ripemd160.New()
	haser.Write([]byte("http://liyuechun.org"))
	hashBytes := haser.Sum(nil)
	fmt.Printf("%x", hashBytes)
}