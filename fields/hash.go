package fields

import (
	"bytes"
	"github.com/hacash/core/crypto/sha3"
)

type Hash = Bytes32
type HashHalf = Bytes16
type HashNonce = Bytes8

// sha3 计算32位哈希值
func CalculateHash(stuff []byte) Hash {
	hx := sha3.Sum256(stuff)
	return hx[:]
}

func (h Hash) GetHalf() HashHalf {
	half := make([]byte, 16)
	copy(half, h)
	return half
}

func (h Hash) GetNonce() HashNonce {
	nonce := make([]byte, 8)
	copy(nonce, h)
	return nonce
}

func (h Hash) Equal(tar Hash) bool {
	return bytes.Compare(h, tar) == 0
}

func (h HashHalf) Equal(tar HashHalf) bool {
	return bytes.Compare(h, tar) == 0
}

func (h HashNonce) Equal(tar HashNonce) bool {
	return bytes.Compare(h, tar) == 0
}
