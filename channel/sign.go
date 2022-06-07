package channel

import (
	"github.com/hacash/core/fields"
	"sort"
)

/**
 * 签名相关
 */

// Sorting unit
type SortAddresses []fields.Address

func (n SortAddresses) Len() int {
	return len(n)
}
func (n SortAddresses) Less(i, j int) bool {
	a, b := n[i], n[j]
	for k := 0; k < fields.AddressSize; k++ {
		if a[k] < b[k] {
			return true // Character sorting
		}
	}
	return false
}
func (n SortAddresses) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// De duplication and sorting address list
func CleanSortMustSignAddresses(addrs []fields.Address) (fields.VarUint1, []fields.Address) {
	// duplicate removal
	addrsclear := make([]fields.Address, 0)
	repeats := map[string]bool{}
	for _, v := range addrs {
		if _, hav := repeats[string(v)]; hav == false {
			addrsclear = append(addrsclear, v)
		}
		repeats[string(v)] = true
	}
	// sort
	sort.Sort(SortAddresses(addrsclear))
	return fields.VarUint1(len(addrsclear)), addrsclear
}
