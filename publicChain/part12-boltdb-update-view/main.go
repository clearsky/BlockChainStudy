package main

import (
	"github.com/boltdb/bolt"
	"log"
)

func main(){
	// 如果不存在，会重新创建一个数据库,mode为权限

	// 创建或者打开数据库
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil{
		log.Panic(err)
	}
	defer db.Close()

	// 更新表
	err = db.Update(func(tx *bolt.Tx) error {
		//创建表
		//b, err := tx.CreateBucket([]byte("BlockBucket"))
		//if err != nil{
		//	return fmt.Errorf("create bucket:%s", err)
		//}
		//// 往表里面存储数据
		//if b != nil{
		//	err := b.Put([]byte("l"), []byte("Send 100 BTC To Tom"))
		//	if err != nil{
		//		log.Panic("数据存储失败......")
		//	}
		//}

		// 返回nil，以便数据库处理相应的操作
		b := tx.Bucket([]byte("BlockBucket"))

		if b != nil {
			err := b.Put([]byte("ll"), []byte("Send 100 BTC To Bob"))
			if err != nil {
				log.Panic("数据存储失败......")
			}
		}

		return nil
	})
	// 更新失败
	if err != nil{
		log.Panic(err)
	}

}