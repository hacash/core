package fields

import "bytes"

const (
	DiamondNameSize = 6
)

type DiamondName = Bytes6
type DiamondNumber = VarUint3

func (d DiamondName) Name() string {
	return string(d)
}

////////////////////////////

var diamond_hash_base_stuff = []byte("0WTYUIAHXVMEKBSZN")

// Judge whether it is a diamond
func IsDiamondValueString(diamondStr string) bool {
	if len(diamondStr) != 6 {
		return false
	}
	for _, a := range diamondStr {
		// drop 0
		if bytes.IndexByte(diamond_hash_base_stuff[1:], byte(a)) == -1 {
			return false
		}
	}
	// Check successful
	return true
}
