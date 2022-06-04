package account

import (
	"crypto/sha256"
	"fmt"
	"github.com/hacash/core/crypto/ripemd160"
)

func NewAddressFromPublicKeyV0(pubKey []byte) []byte {
	return NewAddressFromPublicKey([]byte{0}, pubKey)
}

func NewAddressFromPublicKey(version []byte, pubKey []byte) []byte {

	digest := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	hasher.Write(digest[:])
	hs160 := hasher.Sum(nil)

	// Original and encoded
	return append(version, hs160...)
}

func NewAddressReadableFromAddress(address []byte) string {
	addr := Base58CheckEncode(address)
	// Original and encoded
	return addr
}

// check address is ok ?
func CheckReadableAddress(readable string) ([]byte, error) {
	if len(readable) > 34 {
		return nil, fmt.Errorf("Address format error")
	}
	hashhex, e1 := Base58CheckDecode(readable)
	if e1 != nil {
		return nil, fmt.Errorf("Address format error")
	}
	version := uint8(hashhex[0])
	if version > 2 {
		return nil, fmt.Errorf("Address version error")
	}
	addr := hashhex
	return addr, nil
}
