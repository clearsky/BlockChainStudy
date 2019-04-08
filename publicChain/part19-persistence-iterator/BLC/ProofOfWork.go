package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 265位hash前面至少要有16个零
const targetBit = 16

type ProofOfWork struct {
	Block *Block  // 当前要验证的区块
	target  *big.Int  // 大数存储
}

func (proofOfWork *ProofOfWork) IsValid() bool{
	// 1.使用proofOfWork.Block.Hash
	// 2.proofOfWork.Target
	var hashInt big.Int
	hashInt.SetBytes(proofOfWork.Block.Hash)
	if proofOfWork.target.Cmp(&hashInt)  == 1{
		return true
	}
	return false
}

func (proofOfWork *ProofOfWork) prepareData(nonce int64) []byte{
	data := bytes.Join([][]byte{
		proofOfWork.Block.PrevBlockHash,
		proofOfWork.Block.Data,
		IntToHex(proofOfWork.Block.Timestamp),
		IntToHex(int64(targetBit)),
		IntToHex(int64(nonce)),
		IntToHex(int64(proofOfWork.Block.Height)),
	}, []byte{})
	return data
}

func (proofOfWork *ProofOfWork) Run()([]byte, int64){
	// 1.把Block的属性拼接为字节数组
	// 2.生成hash
	// 3.判断hash的有效性，如果满足条件跳出循环

	var nonce int64 = 0

	var hashInt big.Int  //存储新生成的hash值
	var hash [32]byte
	for{
		//准备数据
		dataBytes := proofOfWork.prepareData(nonce)
		//生成hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x", hash)
		//将hash存储到hashint
		hashInt.SetBytes(hash[:])
		//判断hashint是否小于Block里面的target
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if proofOfWork.target.Cmp(&hashInt) == 1{
			break
		}
		nonce++
	}
	fmt.Println()
	return hash[:], nonce
}

// 创建新的工作量证明对象
func NewProofOfWork(block *Block)  *ProofOfWork{
	// 1.创建一个初始值为1的target
	target := big.NewInt(1)
	// 2.左移256 - targetBit位
	target.Lsh(target, 256 - targetBit)
	return &ProofOfWork{block, target}
}

