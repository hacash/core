package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

// Exchange transactions between channels and atoms on the chain, and save credentials
type Chaswap struct {
	IsBeUsed fields.Bool // Whether it has been used and cannot be reused
	// Signature address
	AddressCount                            fields.VarUint1 // Signature quantity, can only be 2 or 3
	OnchainTransferFromAndMustSignAddresses []fields.Address
}

func (this *Chaswap) Size() uint32 {
	size := this.IsBeUsed.Size() + this.AddressCount.Size()
	size += uint32(this.AddressCount) * fields.AddressSize
	return size
}

func (elm *Chaswap) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.IsBeUsed.Serialize()     // Data body
	var bt2, _ = elm.AddressCount.Serialize() // Data body
	buffer.Write(bt1)
	buffer.Write(bt2)
	for _, addr := range elm.OnchainTransferFromAndMustSignAddresses {
		var bt1, _ = addr.Serialize()
		buffer.Write(bt1)
	}
	return buffer.Bytes(), nil
}

func (elm *Chaswap) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	// address
	seek, e = elm.IsBeUsed.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.AddressCount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	scn := int(elm.AddressCount)
	elm.OnchainTransferFromAndMustSignAddresses = make([]fields.Address, scn)
	for i := 0; i < scn; i++ {
		seek, e = elm.OnchainTransferFromAndMustSignAddresses[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// complete
	return seek, nil
}
