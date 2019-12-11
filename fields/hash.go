package fields

import "encoding/hex"

type Hash = Bytes32

func (h *Hash) ToHex() string {
	return hex.EncodeToString(*h)
}
