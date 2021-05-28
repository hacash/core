package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	DiamondSmeltSize = 6 + 3 + 5 + 32 + 32 + fields.AddressSize + 4 + 8 + 32 + 2
)

type DiamondSmelt struct {
	Diamond              fields.Bytes6 // WTYUIAHXVMEKBSZN
	Number               fields.VarUint3
	ContainBlockHeight   fields.VarUint5
	ContainBlockHash     fields.Hash // current pending block hash
	PrevContainBlockHash fields.Hash // prev block hash
	MinerAddress         fields.Address
	ApproxFeeOffer       fields.Bytes4  // Fee Amount
	Nonce                fields.Bytes8  // nonce
	CustomMessage        fields.Bytes32 // msg
	// data statistics
	AverageBidBurnPrice fields.VarUint2 // 平均竞价销毁的HAC枚数，向下取整，最低一枚，最高65535枚
}

func (this *DiamondSmelt) Size() uint32 {
	return uint32(DiamondSmeltSize)
}

func (this *DiamondSmelt) GetApproxFeeOffer() *fields.Amount {
	amt := fields.Amount{}
	amt.Parse(this.ApproxFeeOffer, 0)
	return &amt
}

func (this *DiamondSmelt) ParseApproxFeeOffer(amt *fields.Amount) error {
	// 压缩存储空间
	approxfeeoffer, _, e11 := amt.CompressForMainNumLen(4, true)
	if e11 != nil {
		return e11
	}
	approxfeeofferBytes, e12 := approxfeeoffer.Serialize()
	if e12 != nil {
		return e12
	}
	approxfeeofferBytesStores := make([]byte, 4)
	copy(approxfeeofferBytesStores, approxfeeofferBytes)
	this.ApproxFeeOffer = fields.Bytes4(approxfeeofferBytesStores)
	return nil
}

func (this *DiamondSmelt) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.Diamond.Serialize()
	b2, _ := this.Number.Serialize()
	b3, _ := this.ContainBlockHeight.Serialize()
	b4, _ := this.ContainBlockHash.Serialize()
	b5, _ := this.PrevContainBlockHash.Serialize()
	b6, _ := this.MinerAddress.Serialize()
	b7, _ := this.ApproxFeeOffer.Serialize()
	b8, _ := this.Nonce.Serialize()
	b9, _ := this.CustomMessage.Serialize()
	b10, _ := this.AverageBidBurnPrice.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	buffer.Write(b7)
	buffer.Write(b8)
	buffer.Write(b9)
	buffer.Write(b10)
	return buffer.Bytes(), nil
}

func (this *DiamondSmelt) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = this.Diamond.Parse(buf, seek)
	seek, _ = this.Number.Parse(buf, seek)
	seek, _ = this.ContainBlockHeight.Parse(buf, seek)
	seek, _ = this.ContainBlockHash.Parse(buf, seek)
	seek, _ = this.PrevContainBlockHash.Parse(buf, seek)
	seek, _ = this.MinerAddress.Parse(buf, seek)
	seek, _ = this.ApproxFeeOffer.Parse(buf, seek)
	seek, _ = this.Nonce.Parse(buf, seek)
	seek, _ = this.CustomMessage.Parse(buf, seek)
	seek, _ = this.AverageBidBurnPrice.Parse(buf, seek)
	return seek, nil
}
