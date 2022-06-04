package account

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/crypto/btcec"
	"math/big"
)

type Account struct {
	AddressReadable string
	Address         []byte
	PublicKey       []byte
	PrivateKey      []byte
	Private         *btcec.PrivateKey
}

const MaxPrikeyValueHex = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364140"

func checkKeyMax(prikey []byte) error {
	maxByte, _ := hex.DecodeString(MaxPrikeyValueHex)
	maxV := big.NewInt(0).SetBytes(maxByte)
	if maxV.Cmp(big.NewInt(0).SetBytes(prikey)) != 1 {
		return fmt.Errorf("Prikey cannot more than %s", MaxPrikeyValueHex)
	}
	// success
	return nil
}

func GetAccountByPriviteKeyHex(hexstr string) (*Account, error) {
	byte1, e1 := hex.DecodeString(hexstr)
	if e1 != nil {
		return nil, e1
	}
	return GetAccountByPriviteKey(byte1)
}

func GetAccountByPriviteKey(byte []byte) (*Account, error) {
	me := checkKeyMax(byte[:])
	if me != nil {
		return nil, me
	}
	privt, e2 := btcec.ToECDSA(byte)
	if e2 != nil {
		return nil, e2
	}
	private := btcec.PrivateKey(*privt)
	return genAccountByPrivateKey(private), nil
}

func CreateAccountByPassword(password string) *Account {
	digest := sha256.Sum256([]byte(password))
	if checkKeyMax(digest[:]) != nil {
		return nil
	}
	privite, _ := btcec.PrivKeyFromBytes(btcec.S256(), digest[:])
	return genAccountByPrivateKey(*privite)
}

func CreateNewRandomAccount() *Account {
	digest := make([]byte, 32)
	for {
		rand.Read(digest)
		if checkKeyMax(digest) == nil {
			break
		}
	}
	privite, _ := btcec.PrivKeyFromBytes(btcec.S256(), digest)
	return genAccountByPrivateKey(*privite)
}

func GetAccountByPrivateKeyOrPassword(password string) *Account {
	// Private key
	if len(password) == 64 {
		if bts, e1 := hex.DecodeString(password); e1 == nil {
			if acc, e2 := GetAccountByPriviteKey(bts); e2 == nil {
				return acc
			}
		}
	}
	// password
	return CreateAccountByPassword(password)

}

func genAccountByPrivateKey(private btcec.PrivateKey) *Account {
	compressedPublic := private.PubKey().SerializeCompressed()
	addr := NewAddressFromPublicKeyV0(compressedPublic)
	readable := NewAddressReadableFromAddress(addr)
	return &Account{
		AddressReadable: readable,
		Address:         addr,
		PublicKey:       compressedPublic,    // 压缩公钥
		PrivateKey:      private.Serialize(), // 私钥
		Private:         &private,
	}
}
