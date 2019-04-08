package main

import (
	"github.com/boltdb/bolt"
	"log"
)

func main(){
	// 如果不存在，会重新创建一个数据库,mode为权限
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil{
		log.Panic(err)
	}
	defer db.Close()
}