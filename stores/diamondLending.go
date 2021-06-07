package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	DiamondLendingIdLength = 14
)

type DiamondLending struct {
	IsRansomed          fields.Bool                 // 是否已经赎回(已经被赎回)
	CreateBlockHeight   fields.VarUint5             // 借贷开启时的区块高度
	MainAddress         fields.Address              // 借贷人地址
	MortgageDiamondList fields.DiamondListMaxLen200 // 抵押钻石列表
	LoanTotalAmountMei  fields.VarUint4             // [枚数]总共借出HAC额度，必须等于总可借额度，不能多也不能少
	BorrowPeriod        fields.VarUint1             // 借款周期，一个周期代表 0.5%利息和10000个区块约35天，最低1最高20
}

func NewDiamondLending(address fields.Address) *DiamondLending {
	addr := address.Copy()
	return &DiamondLending{
		MainAddress: addr,
	}
}

func (elm *DiamondLending) Size() uint32 {
	return elm.IsRansomed.Size() +
		elm.CreateBlockHeight.Size() +
		elm.MainAddress.Size() +
		elm.MortgageDiamondList.Size() +
		elm.LoanTotalAmountMei.Size() +
		elm.BorrowPeriod.Size()
}

func (elm *DiamondLending) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var b0, _ = elm.IsRansomed.Serialize()
	var b1, _ = elm.CreateBlockHeight.Serialize()
	var b2, _ = elm.MainAddress.Serialize()
	var b3, _ = elm.MortgageDiamondList.Serialize()
	var b4, _ = elm.LoanTotalAmountMei.Serialize()
	var b5, _ = elm.BorrowPeriod.Serialize()
	buffer.Write(b0)
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	return buffer.Bytes(), nil
}

func (elm *DiamondLending) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.IsRansomed.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.CreateBlockHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.MainAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.MortgageDiamondList.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LoanTotalAmountMei.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.BorrowPeriod.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}
