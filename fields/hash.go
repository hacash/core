package fields

import "encoding/hex"

type Hash = Bytes32

func (h *Hash) ToHex() string {
	return hex.EncodeToString(*h)
}

func (h *Hash) Equal(tar Hash) bool {
	return bytes.Compare(h, tar) == 0
}
