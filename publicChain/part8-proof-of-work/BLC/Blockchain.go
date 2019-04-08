package BLC

type Blockchain struct {
	Blocks []*Block //存储有序的区块
}

// 1.创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain{
	// 创建创世区块
	genesis := CreateGenesisBlock("GenesisBlock Data...")
	// 返回区块链对象
	return &Blockchain{
		[]*Block{genesis},
	}
}

// 2.增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(data string){
	height := blc.Blocks[len(blc.Blocks)-1].Height + 1
	preHash := blc.Blocks[len(blc.Blocks)-1].Hash
	newBlock := NewBlock(data, height, preHash)
	blc.Blocks = append(blc.Blocks, newBlock)
}