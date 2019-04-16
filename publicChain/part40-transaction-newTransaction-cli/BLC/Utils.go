package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
)

func IntToHex(num int64) []byte{
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil{
		log.Panic(err)
	}
	return buff.Bytes()
}

// 标准的json字符串转数组
func JSONToArray(jsonString string) []string{
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil{
		fmt.Println(err)
	}
	return sArr
}

