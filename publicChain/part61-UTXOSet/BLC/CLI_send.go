package BLC
// 转账
func (cli *CLI)send(from []string, to []string, amount []string){

	blockchain := GetBlockChainObject()
	defer blockchain.DB.Close()  // 最后关闭数据库

	blockchain.MineNewBlock(from, to, amount)

}


