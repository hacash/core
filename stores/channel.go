package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	ChannelIdLength = 16
)

const (
	ChannelStatusOpening                fields.VarUint1 = 0 // 正常开启
	ChannelStatusChallenging            fields.VarUint1 = 1 // 挑战期
	ChannelStatusAgreementClosed        fields.VarUint1 = 2 // 协商关闭，可再次开启重用
	ChannelStatusFinalArbitrationClosed fields.VarUint1 = 3 // 最终仲裁关闭，不可重用

)

//
type Channel struct {
	BelongHeight        fields.BlockHeight // 通道开启时的区块高度
	LockBlock           fields.VarUint2    // 单方面结束通道要锁定的区块数量
	InterestAttribution fields.VarUint1    // 年化 1% 的利息归属： 0.按结束分配 1.全给left 2.全给right
	LeftAddress         fields.Address
	LeftAmount          fields.Amount           // HAC
	LeftSatoshi         fields.SatoshiVariation // SAT
	RightAddress        fields.Address
	RightAmount         fields.Amount           // 抵押数额2
	RightSatoshi        fields.SatoshiVariation // SAT
	ReuseVersion        fields.VarUint4         // 重用版本号 从 1 开始
	Status              fields.VarUint1         // 已经关闭并结算等状态

	// Status = 1 挑战期保存数据
	IsHaveChallengeLog         fields.Bool        // 记录挑战期数据日志
	ChallengeLaunchHeight      fields.BlockHeight // 挑战开始的区块高度
	AssertBillAutoNumber       fields.VarUint8    // 账单流水编号
	AssertAddressIsLeftOrRight fields.Bool        // 主张者是左侧地址还是右侧 true-左  false-右
	AssertAmount               fields.Amount      // 主张者主张自己应该分配的金额

	// Status = 2 or Status = 3 已经关闭资金分配
	LeftFinalDistributionAmount fields.Amount // 左侧最终分配金额

	// cache data
}

func CreateEmptyChannel() *Channel {
	return &Channel{
		BelongHeight:        0,
		LockBlock:           0,
		InterestAttribution: 0,
		LeftAddress:         nil,
		LeftAmount:          fields.NewEmptyAmountValue(),
		LeftSatoshi:         fields.NewEmptySatoshiVariation(),
		RightAddress:        nil,
		RightAmount:         fields.NewEmptyAmountValue(),
		RightSatoshi:        fields.NewEmptySatoshiVariation(),
		ReuseVersion:        1,
		Status:              0,
	}
}

// 状态判断
func (this *Channel) IsOpening() bool {
	return this.Status == ChannelStatusOpening
}
func (this *Channel) IsChallenging() bool {
	return this.Status == ChannelStatusChallenging
}
func (this *Channel) IsAgreementClosed() bool {
	return this.Status == ChannelStatusAgreementClosed
}
func (this *Channel) IsFinalDistributionClosed() bool {
	return this.Status == ChannelStatusFinalArbitrationClosed
}
func (this *Channel) IsClosed() bool {
	return this.Status == ChannelStatusAgreementClosed ||
		this.Status == ChannelStatusFinalArbitrationClosed
}

// 状态操作
func (this *Channel) SetAgreementClosed(leftEndAmt *fields.Amount) {
	this.Status = ChannelStatusAgreementClosed
	this.LeftFinalDistributionAmount = *leftEndAmt
}
func (this *Channel) SetFinalArbitrationClosed(leftEndAmt *fields.Amount) {
	this.Status = ChannelStatusFinalArbitrationClosed
	this.LeftFinalDistributionAmount = *leftEndAmt
}
func (this *Channel) SetOpening() {
	this.Status = ChannelStatusOpening
}
func (this *Channel) SetChallenging(blkhei uint64, isLeftAddr bool, assertAmount *fields.Amount, billno uint64) {
	this.Status = ChannelStatusChallenging
	this.IsHaveChallengeLog.Set(true)
	this.ChallengeLaunchHeight = fields.BlockHeight(blkhei)
	this.AssertBillAutoNumber = fields.VarUint8(billno)
	this.AssertAddressIsLeftOrRight.Set(isLeftAddr)
	this.AssertAmount = *assertAmount
}
func (this *Channel) CleanChallengingLog() {
	this.IsHaveChallengeLog.Set(false)
	this.ChallengeLaunchHeight = fields.BlockHeight(0)
	this.AssertBillAutoNumber = fields.VarUint8(0)
	this.AssertAddressIsLeftOrRight.Set(false)
	emt := fields.NewEmptyAmount()
	this.AssertAmount = *emt
}

func (this *Channel) Size() uint32 {
	size := this.BelongHeight.Size() +
		this.LockBlock.Size() +
		this.LeftAddress.Size() +
		this.LeftAmount.Size() +
		this.LeftSatoshi.Size() +
		this.RightAddress.Size() +
		this.RightAmount.Size() +
		this.RightSatoshi.Size() +
		this.ReuseVersion.Size() +
		this.Status.Size()
	if this.IsHaveChallengeLog.Check() {
		size += this.ChallengeLaunchHeight.Size() +
			this.AssertBillAutoNumber.Size() +
			this.AssertAddressIsLeftOrRight.Size() +
			this.AssertAmount.Size()
	}
	if this.IsClosed() {
		size += this.LeftFinalDistributionAmount.Size()
	}
	return size
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
	seek, e = this.LeftSatoshi.Parse(buf, seek)
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
	seek, e = this.RightSatoshi.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.ReuseVersion.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.Status.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if this.IsHaveChallengeLog.Check() {
		seek, e = this.ChallengeLaunchHeight.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		seek, e = this.AssertBillAutoNumber.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		seek, e = this.AssertAddressIsLeftOrRight.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		seek, e = this.AssertAmount.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	if this.IsClosed() {
		seek, e = this.LeftFinalDistributionAmount.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

func (this *Channel) Serialize() ([]byte, error) {
	var e error
	var bt []byte
	var buffer = new(bytes.Buffer)
	bt, e = this.BelongHeight.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(bt)
	bt, e = this.LockBlock.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(bt)
	bt, e = this.LeftAddress.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(bt)
	bt, e = this.LeftAmount.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(bt)
	bt, e = this.LeftSatoshi.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(bt)
	bt, e = this.RightAddress.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(bt)
	bt, e = this.RightAmount.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(bt)
	bt, e = this.RightSatoshi.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(bt)
	bt, e = this.ReuseVersion.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(bt)
	bt, e = this.Status.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(bt)
	if this.IsHaveChallengeLog.Check() {
		bt, e = this.ChallengeLaunchHeight.Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bt)
		bt, e = this.AssertBillAutoNumber.Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bt)
		bt, e = this.AssertAddressIsLeftOrRight.Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bt)
		bt, e = this.AssertAmount.Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bt)
	}
	if this.IsClosed() {
		bt, e = this.LeftFinalDistributionAmount.Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bt)
	}
	// ok return
	return buffer.Bytes(), nil
}
