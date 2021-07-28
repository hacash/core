package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	ChannelIdLength = 16
)

const (
	ChannelStatusOpening     fields.VarUint1 = 0 // 正常开启
	ChannelStatusChallenging fields.VarUint1 = 1 // 挑战期
	ChannelStatusClosed      fields.VarUint1 = 2 // 已关闭
)

//
type Channel struct {
	BelongHeight fields.BlockHeight // 通道开启时的区块高度
	LockBlock    fields.VarUint2    // 单方面结束通道要锁定的区块数量
	LeftAddress  fields.Address
	LeftAmount   fields.Amount // 抵押数额1
	RightAddress fields.Address
	RightAmount  fields.Amount   // 抵押数额2
	Status       fields.VarUint1 // 已经关闭并结算等状态

	// cache data
}

func (this *Channel) Size() uint32 {
	return this.BelongHeight.Size() +
		this.LockBlock.Size() +
		this.LeftAddress.Size() +
		this.LeftAmount.Size() +
		this.RightAddress.Size() +
		this.RightAmount.Size() +
		this.Status.Size()
}

func (this *Channel) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = this.BelongHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.LockBlock.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.LeftAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.LeftAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.RightAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.RightAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.Status.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (this *Channel) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.BelongHeight.Serialize()
	b2, _ := this.LockBlock.Serialize()
	b3, _ := this.LeftAddress.Serialize()
	b4, _ := this.LeftAmount.Serialize()
	b5, _ := this.RightAddress.Serialize()
	b6, _ := this.RightAmount.Serialize()
	b7, _ := this.Status.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	buffer.Write(b7)
	// ok return
	return buffer.Bytes(), nil
}
