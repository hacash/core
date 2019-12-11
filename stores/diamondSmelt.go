package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	DiamondSmeltSize = 6 + 3 + 5 + 32 + 32 + fields.AddressSize + 8
)

type DiamondSmelt struct {
	Diamond              fields.Bytes6 // WTYUIAHXVMEKBSZN
	Number               fields.VarInt3
	ContainBlockHeight   fields.VarInt5
	ContainBlockHash     fields.Hash // current pending block hash
	PrevContainBlockHash fields.Hash // prev block hash
	MinerAddress         fields.Address
	Nonce                fields.Bytes8 // nonce
}

func (this *DiamondSmelt) Size() uint32 {
	return uint32(DiamondSmeltSize)
}

func (this *DiamondSmelt) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.Diamond.Serialize()
	b2, _ := this.Number.Serialize()
	b3, _ := this.ContainBlockHeight.Serialize()
	b4, _ := this.ContainBlockHash.Serialize()
	b5, _ := this.PrevContainBlockHash.Serialize()
	b6, _ := this.MinerAddress.Serialize()
	b7, _ := this.Nonce.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	buffer.Write(b7)
	return buffer.Bytes(), nil
}

func (this *DiamondSmelt) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = this.Diamond.Parse(buf, seek)
	seek, _ = this.Number.Parse(buf, seek)
	seek, _ = this.ContainBlockHeight.Parse(buf, seek)
	seek, _ = this.ContainBlockHash.Parse(buf, seek)
	seek, _ = this.PrevContainBlockHash.Parse(buf, seek)
	seek, _ = this.MinerAddress.Parse(buf, seek)
	seek, _ = this.Nonce.Parse(buf, seek)
	return seek, nil
}
