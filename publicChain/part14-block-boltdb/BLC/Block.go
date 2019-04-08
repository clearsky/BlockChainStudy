package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	// 1.区块高度，编号
	Height int64
	// 2.上一个区块的hash值
	PrevBlockHash []byte
	// 3.交易数据
	Data []byte
	// 4.时间戳
	Timestamp int64
	// 5.本区块的hash值
	Hash []byte
	// 6.Nonce值
	Nonce int64
}

//序列化
func (block *Block) Serialize() []byte{
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil{
		log.Panic(err)
	}
	return result.Bytes()
}

//反序列化
func  DeserializeBlock(blockBytes []byte) *Block{
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil{
		log.Panic(err)
	}
	return &block
}

//func (block *Block) SetHash(){
//	// 1.height转化为字节数组[]byte
//	heightBytes := IntToHex(block.Height)
//	// 2.将时间戳转化为字节数组[]byte
//		// base:2代表2进制
//	timeString := strconv.FormatInt(block.Timestamp, 2)
//	timeBytes := []byte(timeString)
//	// 3.将所有属性拼接起来
//	blockBytes := bytes.Join(
//		[][]byte{
//			heightBytes,
//			block.PrevBlockHash,
//			block.Data,
//			timeBytes,
//			block.Hash,
//		},[]byte{})
//	// 4.生成hash
//	hash := sha256.Sum256(blockBytes)
//	block.Hash = hash[:]
//}

// 1.创建新的区块
func NewBlock(data string, height int64, prevBlockHash []byte) *Block{
	// 1.创建区块
	block := &Block{
		Height:height,
		PrevBlockHash:prevBlockHash,
		Data:[]byte(data),
		Timestamp:time.Now().Unix(),
		Hash:nil,
		Nonce: 0,
	}

	// 2.调用工作量证明的方法，并且返回有效的hash和Nonce值
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()

	block.Hash = hash
	block.Nonce = nonce
	//// 2.设置本区块hash值
	//block.SetHash()
	return block
}

// 2.单独写一个方法，生成创世区块
func CreateGenesisBlock(data string) *Block{
	var preHash [32]byte
	return NewBlock(data, 1, preHash[:])
}