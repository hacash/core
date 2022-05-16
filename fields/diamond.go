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

// 判断是否为钻石
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
	// 检查成功
	return true
}
