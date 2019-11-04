package account

import (
	"crypto/sha256"
	"github.com/hacash/core/crypto/ripemd160"
)

func NewAddressFromPublicKey( version []byte, pubKey []byte ) ([]byte){

	digest := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	hasher.Write(digest[:])
	hs160 := hasher.Sum(nil)

	// 原始以及编码后的
	return append(version, hs160...)
}

func NewAddressReadableFromAddress( address []byte ) (string){
	addr := Encode(address)
	// 原始以及编码后的
	return addr
}





