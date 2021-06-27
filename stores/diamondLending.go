package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	DiamondSyslendIdLength = 14
)

type DiamondSystemLending struct {
	IsRansomed          fields.Bool                 // 是否已经赎回(已经被赎回)
	CreateBlockHeight   fields.BlockHeight          // 借贷开启时的区块高度
	MainAddress         fields.Address              // 借贷人地址
	MortgageDiamondList fields.DiamondListMaxLen200 // 抵押钻石列表
	LoanTotalAmountMei  fields.VarUint4             // [枚数]总共借出HAC额度，必须等于总可借额度，不能多也不能少
	BorrowPeriod        fields.VarUint1             // 借款周期，一个周期代表 0.5%利息和10000个区块约35天，最低1最高20

	// 如已经赎回则写入数据
	RansomBlockHeight              fields.BlockHeight     // 赎回时的区块高度
	RansomAmount                   fields.Amount          // 赎回金额
	RansomAddressIfPublicOperation fields.OptionalAddress // 如果是第三方赎回则记录第三方地址

}

func NewDiamondSystemLending(address fields.Address) *DiamondSystemLending {
	addr := address.Copy()
	return &DiamondSystemLending{
		MainAddress: addr,
	}
}

func (elm *DiamondSystemLending) Size() uint32 {
	sz := elm.IsRansomed.Size() +
		elm.CreateBlockHeight.Size() +
		elm.MainAddress.Size() +
		elm.MortgageDiamondList.Size() +
		elm.LoanTotalAmountMei.Size() +
		elm.BorrowPeriod.Size()
	// 已经赎回状态
	if elm.IsRansomed.Check() {
		sz += elm.RansomBlockHeight.Size() +
			elm.RansomAmount.Size() +
			elm.RansomAddressIfPublicOperation.Size()
	}
	return sz
}

func (elm *DiamondSystemLending) Serialize() ([]byte, error) {
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
	// 已赎回状态
	if elm.IsRansomed.Check() {
		var b0, _ = elm.RansomBlockHeight.Serialize()
		var b1, _ = elm.RansomAmount.Serialize()
		var b2, _ = elm.RansomAddressIfPublicOperation.Serialize()
		buffer.Write(b0)
		buffer.Write(b1)
		buffer.Write(b2)
	}
	return buffer.Bytes(), nil
}

func (elm *DiamondSystemLending) Parse(buf []byte, seek uint32) (uint32, error) {
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
	// 已赎回状态
	if elm.IsRansomed.Check() {
		seek, e = elm.RansomBlockHeight.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		seek, e = elm.RansomAmount.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		seek, e = elm.RansomAddressIfPublicOperation.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

// 修改、回退赎回状态

func (elm *DiamondSystemLending) SetRansomedStatus(height uint64, amount *fields.Amount, address fields.Address) error {
	elm.IsRansomed.Set(true) // 设置赎回状态
	elm.RansomBlockHeight = fields.BlockHeight(height)
	elm.RansomAmount = *amount
	elm.RansomAddressIfPublicOperation = fields.NewEmptyOptionalAddress()
	if address.NotEqual(elm.MainAddress) {
		elm.RansomAddressIfPublicOperation.Exist = fields.CreateBool(true)
		elm.RansomAddressIfPublicOperation.Addr = address
	}
	return nil
}

// 移除赎回状态
func (elm *DiamondSystemLending) DropRansomedStatus() error {
	elm.IsRansomed.Set(false) // 回退赎回状态
	return nil
}
