package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	// 1.区块高度，编号
	Height int64
	// 2.上一个区块的hash值
	PrevBlockHash []byte
	// 3.交易数据
	//Data []byte
	Txs []*Transaction
	// 4.时间戳
	Timestamp int64
	// 5.本区块的hash值
	Hash []byte
	// 6.Nonce值
	Nonce int64
}

// 将Txs转换为字节数组
func (block *Block) HashTransactions() []byte{
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Txs{
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

func (block *Block) PrintInfo(){
	fmt.Println("height: ", block.Height)
	fmt.Printf("preHash: %x\n", block.PrevBlockHash)
	fmt.Println("nonce: ", block.Nonce)
	fmt.Printf("time: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05 AM"))
	fmt.Printf("hash: %x\n\n", block.Hash)
	fmt.Println("Txs:")
	for _, tx := range block.Txs{
		fmt.Printf("TxHash:%x\n", tx.TxHash)
		fmt.Println("Vins:")
		for _, in := range tx.Vins {
			fmt.Printf("%x\n", in.TxHash)
			fmt.Printf("%d\n", in.Vout)
			fmt.Printf("%s\n", in.Publickey)
		}
		fmt.Println("Vouts:")
		for _, out := range tx.Vouts{
			fmt.Println(out.Value)
			fmt.Println(out.Ripemd160Hash)
		}
	}
	fmt.Println("--------------")
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
func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block{
	// 1.创建区块
	fmt.Println("正在生成新的区块...")
	block := &Block{
		Height:height,
		PrevBlockHash:prevBlockHash,
		Txs:txs,
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
	fmt.Println("新的区块生成成功...")
	return block
}

// 2.单独写一个方法，生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block{
	var preHash [32]byte
	return NewBlock(txs, 1, preHash[:])
}