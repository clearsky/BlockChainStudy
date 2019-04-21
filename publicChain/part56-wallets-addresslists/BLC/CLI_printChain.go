package BLC
func (cli *CLI) printBlockChainChain(){
	blockchain := GetBlockChainObject()
	defer blockchain.DB.Close()
	blockchain.Printchain()
}