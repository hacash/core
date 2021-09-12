package fields

import (
	"bytes"
	"fmt"
	base58check "github.com/hacash/core/account"
	"sort"
)

const (
	AddressSize = 21
)

type AddressReadable = TrimString34
type Address = Bytes21

// check address is ok ?
func CheckReadableAddress(readable string) (*Address, error) {
	if len(readable) > 34 {
		return nil, fmt.Errorf("Address format error")
	}
	hashhex, e1 := base58check.Base58CheckDecode(readable)
	if e1 != nil {
		return nil, fmt.Errorf("Address format error")
	}
	version := uint8(hashhex[0])
	if version > 2 {
		return nil, fmt.Errorf("Address version error")
	}
	addr := Address(hashhex)
	return &addr, nil
}

func (this Address) ToReadable() string {
	return base58check.Base58CheckEncode([]byte(this))
}

// check valid
func (this Address) IsValid() bool {
	if this == nil {
		return false
	}
	if len(this) != AddressSize {
		return false
	}
	// ok
	return true
}

// check equal
func (this Address) Equal(tar Address) bool {
	return bytes.Compare(this, tar) == 0
}
func (this Address) NotEqual(tar Address) bool {
	return bytes.Compare(this, tar) != 0
}

// check equal
func (this Address) Copy() Address {
	addr := make([]byte, AddressSize)
	copy(addr, this)
	return addr
}

////////////////////////////////////////////////////////////

// 字母顺序排序
type AddressListByCharacterSort []Address

func (n AddressListByCharacterSort) Len() int {
	return len(n)
}
func (n AddressListByCharacterSort) Less(i, j int) bool {
	a, b := n[i], n[j]
	for k := 0; k < AddressSize; k++ {
		if a[k] < b[k] {
			return true // 字符排序
		}
	}
	return false
}
func (n AddressListByCharacterSort) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// 去重并排序地址列表
func CleanAddressListByCharacterSort(addrs []Address, addradds []Address) (VarUint1, []Address) {
	// 去重
	addrsclear := make([]Address, 0)
	repeats := map[string]bool{}
	if addrs != nil {
		for _, v := range addrs {
			if _, hav := repeats[string(v)]; hav == false {
				addrsclear = append(addrsclear, v)
			}
			repeats[string(v)] = true
		}
	}
	// 添加
	if addradds != nil {
		for _, v := range addradds {
			if _, hav := repeats[string(v)]; hav == false {
				addrsclear = append(addrsclear, v)
			}
			repeats[string(v)] = true
		}
	}
	// 排序
	sort.Sort(AddressListByCharacterSort(addrsclear))
	return VarUint1(len(addrsclear)), addrsclear
}
