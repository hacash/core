package fields

import (
	"bytes"
	"github.com/hacash/core/crypto/sha3"
)

type Hash = Bytes32
type HashHalfChecker = Bytes16
type HashNonceChecker = Bytes8

const (
	HashSize             = 32
	HashHalfCheckerSize  = 16
	HashNonceCheckerSize = 8
)

// sha3 计算32位哈希值
func CalculateHash(stuff []byte) Hash {
	hx := sha3.Sum256(stuff)
	return hx[:]
}

func (h Hash) GetHalfChecker() HashHalfChecker {
	half := make([]byte, 16)
	copy(half, h)
	return half
}

func (h Hash) GetNonceChecker() HashNonceChecker {
	nonce := make([]byte, 8)
	copy(nonce, h)
	return nonce
}

func (h Hash) Equal(tar Hash) bool {
	return bytes.Compare(h, tar) == 0
}

func (h Hash) NotZeroBlank() bool {
	for i := len(h) - 1; i >= 0; i-- {
		if h[i] > 0 {
			return true
		}
	}
	// empty
	return false
}

func (h HashHalfChecker) Equal(tar HashHalfChecker) bool {
	return bytes.Compare(h, tar) == 0
}

func (h HashNonceChecker) Equal(tar HashNonceChecker) bool {
	return bytes.Compare(h, tar) == 0
}
