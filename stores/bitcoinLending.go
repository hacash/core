package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	BitcoinSyslendIdLength = 15
)

type BitcoinSystemLending struct {
	IsRansomed                 fields.Bool        // Whether it has been redeemed (redeemed)
	CreateBlockHeight          fields.BlockHeight // 借贷开启时的区块高度
	MainAddress                fields.Address     // Address of Borrower
	MortgageBitcoinPortion     fields.VarUint2    // Number of mortgage bitcoin copies (each = 0.01btc)
	LoanTotalAmount            fields.Amount      // The total lending HAC quantity must be less than or equal to the lendable quantity
	PreBurningInterestAmount   fields.Amount      // Interest for pre destruction must be greater than or equal to the destroyed quantity
	RealtimeTotalMortgageRatio fields.VarUint2    // Value: 0~10000, unit: 10000

	// Write data if redeemed
	RansomBlockHeight              fields.BlockHeight     // Block height at redemption
	RansomAmount                   fields.Amount          // Redemption amount
	RansomAddressIfPublicOperation fields.OptionalAddress // If it is a third party redemption, record the third party address

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
	// Redeemed status
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
		seek, e = elm.RansomAddressIfPublicOperation.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

// Modify and return redemption status

func (elm *BitcoinSystemLending) SetRansomedStatus(height uint64, amount *fields.Amount, address fields.Address) error {
	elm.IsRansomed.Set(true) // Set redemption status
	elm.RansomBlockHeight = fields.BlockHeight(height)
	elm.RansomAmount = *amount
	elm.RansomAddressIfPublicOperation = fields.NewEmptyOptionalAddress()
	if address.NotEqual(elm.MainAddress) {
		elm.RansomAddressIfPublicOperation.Exist = fields.CreateBool(true)
		elm.RansomAddressIfPublicOperation.Addr = address
	}
	return nil
}

// Remove redemption status
func (elm *BitcoinSystemLending) DropRansomedStatus() error {
	elm.IsRansomed.Set(false) // Return redemption status
	return nil
}
