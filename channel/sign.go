package channel

import (
	"github.com/hacash/core/fields"
	"sort"
)

/**
 * 签名相关
 */

// 排序单元
type SortAddresses []fields.Address

func (n SortAddresses) Len() int {
	return len(n)
}
func (n SortAddresses) Less(i, j int) bool {
	a, b := n[i], n[j]
	for k := 0; k < fields.AddressSize; k++ {
		if a[k] < b[k] {
			return true // 字符排序
		}
	}
	return false
}
func (n SortAddresses) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// 去重并排序地址列表
func CleanSortMustSignAddresses(addrs []fields.Address) (fields.VarUint1, []fields.Address) {
	// 去重
	addrsclear := make([]fields.Address, 0)
	repeats := map[string]bool{}
	for _, v := range addrs {
		if _, hav := repeats[string(v)]; hav == false {
			addrsclear = append(addrsclear, v)
		}
		repeats[string(v)] = true
	}
	// 排序
	sort.Sort(SortAddresses(addrsclear))
	return fields.VarUint1(len(addrsclear)), addrsclear
}
