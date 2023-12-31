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
	Status  fields.VarUint1 // Status 0 Normally available and transferable 1 Mortgage to system 2 Mortgage to other users
	Address fields.Address
	// engraved info
	EngravedContents fields.StringMax255List255
}

func NewDiamond(address fields.Address) *Diamond {
	addr := address.Copy()
	return &Diamond{
		Status:           DiamondStatusNormal,
		Address:          addr,
		EngravedContents: fields.CreateEmptyStringMax255List255(),
	}
}

func (this *Diamond) Size() uint32 {
	return this.Status.Size() +
		this.Address.Size() +
		this.EngravedContents.Size()
}

func (this *Diamond) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.Status.Serialize()
	b2, _ := this.Address.Serialize()
	b3, _ := this.EngravedContents.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
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
	seek, e = this.EngravedContents.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}
