package BLC

// 创建按创世区块
func (cli *CLI) createBlockChain(address string){
	blockchain := CreateBlockchainWithGenesisBlock(address)
	defer blockchain.DB.Close()
}
