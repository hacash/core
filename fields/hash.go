package fields

import (
	"bytes"
)

type Hash = Bytes32

func (h Hash) Equal(tar Hash) bool {
	return bytes.Compare(h, tar) == 0
}
