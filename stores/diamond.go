package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	DiamondSize = fields.AddressSize
)

type Diamond struct {
	Address fields.Address
}

func (this *Diamond) Size() uint32 {
	return uint32(DiamondSize)
}

func (this *Diamond) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.Address.Serialize()
	buffer.Write(b1)
	return buffer.Bytes(), nil
}

func (this *Diamond) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = this.Address.Parse(buf, seek)
	return seek, nil
}
