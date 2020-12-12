package account

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/hacash/core/crypto/btcec"
)

type Account struct {
	AddressReadable string
	Address         []byte
	PublicKey       []byte
	PrivateKey      []byte
	Private         *btcec.PrivateKey
}

func GetAccountByPriviteKeyHex(hexstr string) (*Account, error) {
	byte, e1 := hex.DecodeString(hexstr)
	if e1 != nil {
		return nil, e1
	}
	return GetAccountByPriviteKey(byte)
}

func GetAccountByPriviteKey(byte []byte) (*Account, error) {
	privt, e2 := btcec.ToECDSA(byte)
	if e2 != nil {
		return nil, e2
	}
	private := btcec.PrivateKey(*privt)
	return genAccountByPrivateKey(private), nil
}

func CreateAccountByPassword(password string) *Account {
	digest := sha256.Sum256([]byte(password))
	privite, _ := btcec.PrivKeyFromBytes(btcec.S256(), digest[:])
	return genAccountByPrivateKey(*privite)
}

func CreateNewRandomAccount() *Account {
	digest := make([]byte, 32)
	rand.Read(digest)
	privite, _ := btcec.PrivKeyFromBytes(btcec.S256(), digest)
	return genAccountByPrivateKey(*privite)
}

func GetAccountByPrivateKeyOrPassword(password string) *Account {
	// 私钥
	if len(password) == 64 {
		if bts, e1 := hex.DecodeString(password); e1 == nil {
			if acc, e2 := GetAccountByPriviteKey(bts); e2 == nil {
				return acc
			}
		}
	}
	// 密码
	return CreateAccountByPassword(password)

}

func genAccountByPrivateKey(private btcec.PrivateKey) *Account {
	compressedPublic := private.PubKey().SerializeCompressed()
	addr := NewAddressFromPublicKey([]byte{0}, compressedPublic)
	readable := NewAddressReadableFromAddress(addr)
	return &Account{
		AddressReadable: readable,
		Address:         addr,
		PublicKey:       compressedPublic,    // 压缩公钥
		PrivateKey:      private.Serialize(), // 私钥
		Private:         &private,
	}
}
