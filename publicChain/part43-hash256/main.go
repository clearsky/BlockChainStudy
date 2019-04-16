package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	haser := sha256.New()
	haser.Write([]byte("http://liyuechun.org"))
	hashBytes := haser.Sum(nil)
	fmt.Printf("%x", hashBytes)
}