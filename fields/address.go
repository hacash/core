package fields

import (
	"fmt"
	base58check "github.com/hacash/core/account"
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
