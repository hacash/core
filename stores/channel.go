package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	ChannelSize     = 5 + 2 + (21+6)*2 + 1 + 2 + 16 // 80 = 16 × 5
	ChannelIdLength = 16
)

//
type Channel struct {
	BelongHeight fields.BlockHeight // 通道开启时的区块高度
	LockBlock    fields.VarUint2    // 单方面结束通道要锁定的区块数量
	LeftAddress  fields.Address
	LeftAmount   fields.Amount // 抵押数额1  【6位定宽】
	RightAddress fields.Address
	RightAmount  fields.Amount   // 抵押数额2  【6位定宽】
	IsClosed     fields.Bool     // 已经关闭并结算
	ConfigMark   fields.VarUint2 // 标志位
	Others       fields.Bytes16  // 扩展位

	// cache data
}

func (this *Channel) Size() uint32 {
	return uint32(ChannelSize)
}

func (this *Channel) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = this.BelongHeight.Parse(buf, seek)
	seek, _ = this.LockBlock.Parse(buf, seek)
	seek, _ = this.LeftAddress.Parse(buf, seek)
	this.LeftAmount.Parse(buf, seek)
	seek += 6 // 6位定宽
	seek, _ = this.RightAddress.Parse(buf, seek)
	this.RightAmount.Parse(buf, seek)
	seek += 6 // 6位定宽
	seek, _ = this.IsClosed.Parse(buf, seek)
	seek, _ = this.ConfigMark.Parse(buf, seek)
	seek, _ = this.Others.Parse(buf, seek)
	return seek, nil
}

func (this *Channel) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.BelongHeight.Serialize()
	b2, _ := this.LockBlock.Serialize()
	b3, _ := this.LeftAddress.Serialize()
	b4, _ := this.LeftAmount.Serialize()
	if len(b4) < 6 { // 6位定宽，补全6位
		mb4 := make([]byte, 6)
		copy(mb4, b4)
		b4 = mb4
	}
	b5, _ := this.RightAddress.Serialize()
	b6, _ := this.RightAmount.Serialize()
	if len(b6) < 6 { // 6位定宽，补全6位
		mb6 := make([]byte, 6)
		copy(mb6, b6)
		b6 = mb6
	}
	b7, _ := this.IsClosed.Serialize()
	b8, _ := this.ConfigMark.Serialize()
	b9, _ := this.Others.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	buffer.Write(b7)
	buffer.Write(b8)
	buffer.Write(b9)
	// ok return
	return buffer.Bytes(), nil
}
