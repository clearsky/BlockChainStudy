package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const version = byte(0x00)
const addressChecksumLen = 4


type Wallet struct {
	// 1.私钥
	PrivateKey ecdsa.PrivateKey  // 椭圆曲线加密
	// 2.公钥
	PublickKey []byte
}

func CheckSum(versionRipemd160Hash []byte)[]byte{
	hash1 := sha256.Sum256(versionRipemd160Hash)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChecksumLen]
}

func (wallet *Wallet) GetAddress() []byte{
	// 1.hash160
	ripemd160Hash := wallet.Ripemd160Hash(wallet.PublickKey)
	versionRipemd160Hash := append([]byte{version}, ripemd160Hash...)
	checkSumBytes := CheckSum(versionRipemd160Hash)
	allBytes := append(versionRipemd160Hash, checkSumBytes...)
	addressBytes := Base58Encode(allBytes)
	return addressBytes
}

func (wallet *Wallet) Ripemd160Hash(publicKey []byte) []byte{
	// 1.256hash
	hasher := sha256.New()
	hasher.Write(publicKey)
	hash256 := hasher.Sum(nil)

	// 2.165hash
	ripemd160er := ripemd160.New()
	ripemd160er.Write(hash256)
	encoded := ripemd160er.Sum(nil)

	return encoded
}

// 创建钱包
func NewWallet() *Wallet{
	privateKey, publicKey := newKeyPair()
	return &Wallet{
		PrivateKey:privateKey,
		PublickKey:publicKey,
	}
}



// 通过私钥产生公钥
func newKeyPair()(ecdsa.PrivateKey, []byte){

	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil{
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}