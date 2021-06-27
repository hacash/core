package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	UserLendingIdLength = 17
)

type UserLending struct {
	IsRansomed           fields.Bool // 是否已经赎回(已经被赎回)
	IsRedemptionOvertime fields.Bool // 是否超期仍可赎回（自动展期）
	IsPublicRedeemable   fields.Bool // 到期后是否公共可赎回

	CreateBlockHeight fields.BlockHeight // 借贷开启时的区块高度
	ExpireBlockHeight fields.BlockHeight // 约定到期的区块高度

	MortgagorAddress fields.Address // 抵押人地址
	LenderAddress    fields.Address // 放款人地址

	MortgageBitcoin     fields.SatoshiVariation     // 抵押比特币数量 单位：SAT
	MortgageDiamondList fields.DiamondListMaxLen200 // 抵押钻石表

	LoanTotalAmount        fields.Amount // 总共借出HAC数量，必须小于等于可借数
	AgreedRedemptionAmount fields.Amount // 约定的赎回金额

	PreBurningInterestAmount fields.Amount // 预先销毁的利息，必须大于等于 借出金额的 1%

	// 如已经赎回则写入数据
	RansomBlockHeight fields.BlockHeight // 赎回时的区块高度
	RansomAmount      fields.Amount      // 赎回金额
	RansomAddress     fields.Address     // 赎回人地址
}

func (elm *UserLending) Size() uint32 {
	sz := elm.IsRansomed.Size() +
		elm.IsRedemptionOvertime.Size() +
		elm.IsPublicRedeemable.Size() +
		elm.CreateBlockHeight.Size() +
		elm.ExpireBlockHeight.Size() +
		elm.MortgagorAddress.Size() +
		elm.LenderAddress.Size() +
		elm.MortgageBitcoin.Size() +
		elm.MortgageDiamondList.Size() +
		elm.LoanTotalAmount.Size() +
		elm.AgreedRedemptionAmount.Size() +
		elm.PreBurningInterestAmount.Size()
	// 已经赎回状态
	if elm.IsRansomed.Check() {
		sz += elm.RansomBlockHeight.Size() +
			elm.RansomAmount.Size() +
			elm.RansomAddress.Size()
	}
	return sz
}

func (elm *UserLending) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var b1, _ = elm.IsRansomed.Serialize()
	var b2, _ = elm.IsRedemptionOvertime.Serialize()
	var b3, _ = elm.IsPublicRedeemable.Serialize()
	var b4, _ = elm.CreateBlockHeight.Serialize()
	var b5, _ = elm.ExpireBlockHeight.Serialize()
	var b6, _ = elm.MortgagorAddress.Serialize()
	var b7, _ = elm.LenderAddress.Serialize()
	var b8, _ = elm.MortgageBitcoin.Serialize()
	var b9, _ = elm.MortgageDiamondList.Serialize()
	var b10, _ = elm.LoanTotalAmount.Serialize()
	var b11, _ = elm.AgreedRedemptionAmount.Serialize()
	var b12, _ = elm.PreBurningInterestAmount.Serialize()
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
	buffer.Write(b11)
	buffer.Write(b12)
	// 已赎回状态
	if elm.IsRansomed.Check() {
		var b0, _ = elm.RansomBlockHeight.Serialize()
		var b1, _ = elm.RansomAmount.Serialize()
		var b2, _ = elm.RansomAddress.Serialize()
		buffer.Write(b0)
		buffer.Write(b1)
		buffer.Write(b2)
	}
	return buffer.Bytes(), nil
}

func (elm *UserLending) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.IsRansomed.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.IsRedemptionOvertime.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.IsPublicRedeemable.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.CreateBlockHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ExpireBlockHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.MortgagorAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LenderAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.MortgageBitcoin.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.MortgageDiamondList.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LoanTotalAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.AgreedRedemptionAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.PreBurningInterestAmount.Parse(buf, seek)
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
		seek, e = elm.RansomAddress.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

// 修改、回退赎回状态

func (elm *UserLending) SetRansomedStatus(height uint64, amount *fields.Amount, address fields.Address) error {
	elm.IsRansomed.Set(true) // 设置赎回状态
	elm.RansomBlockHeight = fields.BlockHeight(height)
	elm.RansomAmount = *amount
	elm.RansomAddress = address
	return nil
}

// 移除赎回状态
func (elm *UserLending) DropRansomedStatus() error {
	elm.IsRansomed.Set(false) // 回退赎回状态
	return nil
}
