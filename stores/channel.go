package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	ChannelIdLength = 16
)

const (
	ChannelStatusOpening                fields.VarUint1 = 0 // Normal opening
	ChannelStatusChallenging            fields.VarUint1 = 1 // Challenging period
	ChannelStatusAgreementClosed        fields.VarUint1 = 2 // After negotiation is closed, reuse can be enabled again
	ChannelStatusFinalArbitrationClosed fields.VarUint1 = 3 // Final arbitration closed, never reusable

)

//
type Channel struct {
	BelongHeight         fields.BlockHeight // Block height when channel is opened
	ArbitrationLockBlock fields.VarUint2    // Number of blocks to be locked for unilateral end channel
	InterestAttribution  fields.VarUint1    // Interest attribution of 1% annualized: 0 Press end to assign 1 All to left 2 Give it all right
	LeftAddress          fields.Address
	LeftAmount           fields.Amount           // HAC
	LeftSatoshi          fields.SatoshiVariation // SAT
	RightAddress         fields.Address
	RightAmount          fields.Amount           // Mortgage amount 2
	RightSatoshi         fields.SatoshiVariation // SAT
	ReuseVersion         fields.VarUint4         // Reuse version number from 1
	Status               fields.VarUint1         // Closed and settled

	// Status = 1 challenge period save data
	IsHaveChallengeLog         fields.Bool             // Record challenge data log
	ChallengeLaunchHeight      fields.BlockHeight      // Block height at the beginning of the challenge
	AssertBillAutoNumber       fields.VarUint8         // Statement serial number provided by the proposer
	AssertAddressIsLeftOrRight fields.Bool             // Whether the proposer is the left address or the right true left false right
	AssertAmount               fields.Amount           // The amount claimed by the proponent
	AssertSatoshi              fields.SatoshiVariation // The sat that proponents claim they should be assigned

	// Status = 2 or status = 3 fund allocation has been closed
	LeftFinalDistributionAmount  fields.Amount           // Final allocation amount on the left
	LeftFinalDistributionSatoshi fields.SatoshiVariation // Final allocation amount on the left

	// cache data
}

func CreateEmptyChannel() *Channel {
	return &Channel{
		BelongHeight:         0,
		ArbitrationLockBlock: 0,
		InterestAttribution:  0,
		LeftAddress:          nil,
		LeftAmount:           fields.NewEmptyAmountValue(),
		LeftSatoshi:          fields.NewEmptySatoshiVariation(),
		RightAddress:         nil,
		RightAmount:          fields.NewEmptyAmountValue(),
		RightSatoshi:         fields.NewEmptySatoshiVariation(),
		ReuseVersion:         1,
		Status:               0,
	}
}

// State judgment
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

// Status operation
func (this *Channel) SetAgreementClosed(leftEndAmt *fields.Amount, satoshi fields.Satoshi) {
	this.Status = ChannelStatusAgreementClosed
	this.LeftFinalDistributionAmount = *leftEndAmt
	this.LeftFinalDistributionSatoshi = satoshi.GetSatoshiVariation()
}
func (this *Channel) SetFinalArbitrationClosed(leftEndAmt *fields.Amount, satoshi fields.Satoshi) {
	this.Status = ChannelStatusFinalArbitrationClosed
	this.LeftFinalDistributionAmount = *leftEndAmt
	this.LeftFinalDistributionSatoshi = satoshi.GetSatoshiVariation()
}
func (this *Channel) SetOpening() {
	this.Status = ChannelStatusOpening
}
func (this *Channel) SetChallenging(blkhei uint64, isLeftAddr bool, assertAmount *fields.Amount, assertSAT fields.Satoshi, billno uint64) {
	this.Status = ChannelStatusChallenging
	this.IsHaveChallengeLog.Set(true)
	this.ChallengeLaunchHeight = fields.BlockHeight(blkhei)
	this.AssertBillAutoNumber = fields.VarUint8(billno)
	this.AssertAddressIsLeftOrRight.Set(isLeftAddr)
	this.AssertAmount = *assertAmount
	this.AssertSatoshi = assertSAT.GetSatoshiVariation()

}
func (this *Channel) CleanChallengingLog() {
	this.IsHaveChallengeLog.Set(false)
	this.ChallengeLaunchHeight = fields.BlockHeight(0)
	this.AssertBillAutoNumber = fields.VarUint8(0)
	this.AssertAddressIsLeftOrRight.Set(false)
	emt := fields.NewEmptyAmount()
	this.AssertAmount = *emt
	this.AssertSatoshi = fields.NewEmptySatoshiVariation()
}

func (this *Channel) Size() uint32 {
	size := this.BelongHeight.Size() +
		this.ArbitrationLockBlock.Size() +
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
			this.AssertAmount.Size() +
			this.AssertSatoshi.Size()
	}
	if this.IsClosed() {
		size += this.LeftFinalDistributionAmount.Size() +
			this.LeftFinalDistributionSatoshi.Size()
	}
	return size
}

func (this *Channel) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = this.BelongHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.ArbitrationLockBlock.Parse(buf, seek)
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
	seek, e = this.IsHaveChallengeLog.Parse(buf, seek)
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
		seek, e = this.AssertSatoshi.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	if this.IsClosed() {
		seek, e = this.LeftFinalDistributionAmount.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		seek, e = this.LeftFinalDistributionSatoshi.Parse(buf, seek)
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
	bt, e = this.ArbitrationLockBlock.Serialize()
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
	bt, e = this.IsHaveChallengeLog.Serialize()
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
		bt, e = this.AssertSatoshi.Serialize()
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
		bt, e = this.LeftFinalDistributionSatoshi.Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bt)
	}
	// ok return
	return buffer.Bytes(), nil
}
