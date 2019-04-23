package BLC

func(cli *CLI) TestMethod(){
	blockchain := GetBlockChainObject()
	defer blockchain.DB.Close()

	utxoSet := &UTXOSet{
		BlockChain:blockchain,
	}
	utxoSet.ResetUTXOSet()
}
