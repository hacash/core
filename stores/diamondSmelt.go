package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
//DiamondSmeltSize = 6 + 3 + 5 + 32 + 32 + fields.AddressSize + 4 + 8 + 32 + 2
)

type DiamondSmelt struct {
	Diamond              fields.DiamondName // WTYUIAHXVMEKBSZN
	Number               fields.DiamondNumber
	ContainBlockHeight   fields.BlockHeight
	ContainBlockHash     fields.Hash // current pending block hash
	PrevContainBlockHash fields.Hash // prev block hash
	MinerAddress         fields.Address
	ApproxFeeOffer       fields.Amount  // Fee Amount
	Nonce                fields.Bytes8  // nonce
	CustomMessage        fields.Bytes32 // msg
	// data statistics
	AverageBidBurnPrice fields.VarUint2 // 平均竞价销毁的HAC枚数，向下取整，最低一枚，最高65535枚
}

func (this *DiamondSmelt) Size() uint32 {
	return this.Diamond.Size() +
		this.Number.Size() +
		this.ContainBlockHeight.Size() +
		this.ContainBlockHash.Size() +
		this.PrevContainBlockHash.Size() +
		this.MinerAddress.Size() +
		this.ApproxFeeOffer.Size() +
		this.Nonce.Size() +
		this.CustomMessage.Size() +
		this.AverageBidBurnPrice.Size()
}

func (this *DiamondSmelt) GetApproxFeeOffer() *fields.Amount {
	return &this.ApproxFeeOffer
}

func (this *DiamondSmelt) ParseApproxFeeOffer(amt *fields.Amount) error {
	// 压缩存储空间
	approxfeeoffer, _, e11 := amt.CompressForMainNumLen(4, true)
	if e11 != nil {
		return e11
	}
	this.ApproxFeeOffer = *approxfeeoffer
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
	var e error
	seek, e = this.Diamond.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.Number.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.ContainBlockHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.ContainBlockHash.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.PrevContainBlockHash.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.MinerAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.ApproxFeeOffer.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.Nonce.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.CustomMessage.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.AverageBidBurnPrice.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}
