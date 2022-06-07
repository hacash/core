package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	UserLendingIdLength = 17
)

type UserLending struct {
	IsRansomed           fields.Bool // Whether it has been redeemed (redeemed)
	IsRedemptionOvertime fields.Bool // Whether it can be redeemed after expiration (automatic extension)
	IsPublicRedeemable   fields.Bool // Public redeemable after maturity

	CreateBlockHeight fields.BlockHeight // 借贷开启时的区块高度
	ExpireBlockHeight fields.BlockHeight // Agreed expiration block height

	MortgagorAddress fields.Address // Address of mortgagor
	LenderAddress    fields.Address // Lender address

	MortgageBitcoin     fields.SatoshiVariation     // Mortgage bitcoin quantity unit: SAT
	MortgageDiamondList fields.DiamondListMaxLen200 // Mortgage diamond table

	LoanTotalAmount        fields.Amount // The total lending HAC quantity must be less than or equal to the lendable quantity
	AgreedRedemptionAmount fields.Amount // Agreed redemption amount

	PreBurningInterestAmount fields.Amount // Interest for pre destruction must be greater than or equal to 1% of the lending amount

	// Write data if redeemed
	RansomBlockHeight fields.BlockHeight // Block height at redemption
	RansomAmount      fields.Amount      // Redemption amount
	RansomAddress     fields.Address     // Address of the Redeemer
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
	// Redeemed status
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
	// Redeemed status
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
	// Redeemed status
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

// Modify and return redemption status

func (elm *UserLending) SetRansomedStatus(height uint64, amount *fields.Amount, address fields.Address) error {
	elm.IsRansomed.Set(true) // Set redemption status
	elm.RansomBlockHeight = fields.BlockHeight(height)
	elm.RansomAmount = *amount
	elm.RansomAddress = address
	return nil
}

// Remove redemption status
func (elm *UserLending) DropRansomedStatus() error {
	elm.IsRansomed.Set(false) // Return redemption status
	return nil
}
