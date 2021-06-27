package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	BitcoinSyslendIdLength = 15
)

type BitcoinSystemLending struct {
	IsRansomed                 fields.Bool        // 是否已经赎回(已经被赎回)
	CreateBlockHeight          fields.BlockHeight // 借贷开启时的区块高度
	MainAddress                fields.Address     // 借贷人地址
	MortgageBitcoinPortion     fields.VarUint2    // 抵押比特币份数（每份 = 0.01BTC）
	LoanTotalAmount            fields.Amount      // 总共借出HAC数量，必须小于等于可借数
	PreBurningInterestAmount   fields.Amount      // 预先销毁的利息，必须大于等于销毁数量
	RealtimeTotalMortgageRatio fields.VarUint2    // 取值 0~10000， 单位 万分之

	// 如已经赎回则写入数据
	RansomBlockHeight              fields.BlockHeight     // 赎回时的区块高度
	RansomAmount                   fields.Amount          // 赎回金额
	RansomAddressIfPublicOperation fields.OptionalAddress // 如果是第三方赎回则记录第三方地址

}

func NewBitcoinSystemLending(address fields.Address) *BitcoinSystemLending {
	addr := address.Copy()
	return &BitcoinSystemLending{
		MainAddress: addr,
	}
}

func (elm *BitcoinSystemLending) Size() uint32 {
	sz := elm.IsRansomed.Size() +
		elm.CreateBlockHeight.Size() +
		elm.MainAddress.Size() +
		elm.MortgageBitcoinPortion.Size() +
		elm.LoanTotalAmount.Size() +
		elm.PreBurningInterestAmount.Size() +
		elm.RealtimeTotalMortgageRatio.Size()
	if elm.IsRansomed.Check() {
		sz += elm.RansomBlockHeight.Size() +
			elm.RansomAmount.Size() +
			elm.RansomAddressIfPublicOperation.Size()
	}
	return sz
}

func (elm *BitcoinSystemLending) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var b0, _ = elm.IsRansomed.Serialize()
	var b1, _ = elm.CreateBlockHeight.Serialize()
	var b2, _ = elm.MainAddress.Serialize()
	var b3, _ = elm.MortgageBitcoinPortion.Serialize()
	var b4, _ = elm.LoanTotalAmount.Serialize()
	var b5, _ = elm.PreBurningInterestAmount.Serialize()
	var b6, _ = elm.RealtimeTotalMortgageRatio.Serialize()
	buffer.Write(b0)
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
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

func (elm *BitcoinSystemLending) Parse(buf []byte, seek uint32) (uint32, error) {
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
	seek, e = elm.MortgageBitcoinPortion.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LoanTotalAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.PreBurningInterestAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RealtimeTotalMortgageRatio.Parse(buf, seek)
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

func (elm *BitcoinSystemLending) SetRansomedStatus(height uint64, amount *fields.Amount, address fields.Address) error {
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
func (elm *BitcoinSystemLending) DropRansomedStatus() error {
	elm.IsRansomed.Set(false) // 回退赎回状态
	return nil
}
