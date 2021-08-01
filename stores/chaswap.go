package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

// 通道与链上原子互换交易，保存凭据
type Chaswap struct {
	IsBeUsed fields.Bool // 是否已经使用过，不可重复使用
	// 签名地址
	AddressCount                            fields.VarUint1 // 签名数量，只能取值 2 或 3
	OnchainTransferFromAndMustSignAddresses []fields.Address
}

func (this *Chaswap) Size() uint32 {
	size := this.IsBeUsed.Size() + this.AddressCount.Size()
	size += uint32(this.AddressCount) * fields.AddressSize
	return size
}

func (elm *Chaswap) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.IsBeUsed.Serialize()     // 数据体
	var bt2, _ = elm.AddressCount.Serialize() // 数据体
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
	// 地址
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
	// 完成
	return seek, nil
}
