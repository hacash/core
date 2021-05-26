package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	DiamondSize = 1 + fields.AddressSize
)

const (
	DiamondStatusNormal           fields.VarUint1 = 0
	DiamondStatusLendingSystem    fields.VarUint1 = 1
	DiamondStatusLendingOtherUser fields.VarUint1 = 2
)

type Diamond struct {
	Status  fields.VarUint1 // 状态 0. 正常可用可转账  1. 抵押给系统  2. 抵押给其他用户
	Address fields.Address
}

func NewDiamond(address fields.Address) *Diamond {
	addr := address.Copy()
	return &Diamond{
		Status:  DiamondStatusNormal,
		Address: addr,
	}
}

func (this *Diamond) Size() uint32 {
	return uint32(DiamondSize)
}

func (this *Diamond) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.Status.Serialize()
	b2, _ := this.Address.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	return buffer.Bytes(), nil
}

func (this *Diamond) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = this.Status.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.Address.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}
