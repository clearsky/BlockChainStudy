package BLC

type Wallets struct {
	Wallets map[string]*Wallet
}

// 创建钱包集合
func NewWallets() *Wallets{
	wallets := &Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	return wallets
}

// 创建一个新钱包
func (wallets *Wallets) CreateNewWallet(){
	wallet := NewWallet()
	wallets.Wallets[string(wallet.GetAddress())] = wallet
}
