package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletsFile = "wallets.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

// 创建钱包集合
func NewWallets() (*Wallets, error){
	if _, err := os.Stat(walletsFile);os.IsNotExist(err){
		wallets := &Wallets{}
		wallets.Wallets = make(map[string]*Wallet)
		return wallets, err
	}

	// 读取文件
	fileContent, err := ioutil.ReadFile(walletsFile)
	if err != nil{
		log.Panic(err)
	}

	// 反序列化
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil{
		log.Panic(err)
	}

	return &wallets, nil
}

// 创建一个新钱包
func (wallets *Wallets) CreateNewWallet(){
	wallet := NewWallet()
	wallets.Wallets[string(wallet.GetAddress())] = wallet
	wallets.SaveWallets()
	fmt.Printf("Address:%s\n", wallet.GetAddress())
}

// 将钱包信息写入到文件
func (wallets *Wallets) SaveWallets(){
	var content bytes.Buffer

	// 注册的目的，是为了，可以序列化任何类型
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallets)
	if err != nil{
		log.Panic(err)
	}

	// 将序列化以后的数据，写入到文件，原来文件的数据会被覆盖掉
	err = ioutil.WriteFile(walletsFile, content.Bytes(), 0644)
	if err != nil{
		log.Panic(err)
	}
}
