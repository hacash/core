package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

type Diamond struct {
	BlockHeight fields.VarInt5
	Number      fields.VarInt3
	Address     fields.Address
}

func (this *Diamond) Size() uint32 {
	return 5 + 3 + 21
}

func (this *Diamond) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.BlockHeight.Serialize()
	b2, _ := this.Number.Serialize()
	b3, _ := this.Address.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	return buffer.Bytes(), nil
}

func (this *Diamond) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = this.BlockHeight.Parse(buf, seek)
	seek, _ = this.Number.Parse(buf, seek)
	seek, _ = this.Address.Parse(buf, seek)
	return seek, nil
}


