package BLC

import "math/big"

// 265位hash前面至少要有16个零
const targetBit = 16

type ProofOfWork struct {
	Block *Block  // 当前要验证的区块
	target  *big.Int  // 大数存储
}

func (proofOfWork *ProofOfWork) Run()([]byte, int64){

	return nil, 0
}

// 创建新的工作量证明对象
func NewProofOfWork(block *Block)  *ProofOfWork{
	// 1.创建一个初始值为1的target
	target := big.NewInt(1)
	// 2.左移256 - targetBit位
	target.Lsh(target, 256 - targetBit)
	return &ProofOfWork{block, target}
}

