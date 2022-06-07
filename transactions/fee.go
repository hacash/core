package transactions

import (
	"github.com/hacash/core/fields"
	"math/big"
)

const (
	FeePurityUnit = 232                          // =Handling charge content unit: 0.00000001 baht = 1 yuan
	MaxFeePurity  = uint64(18446744073709550000) // =The maximum value of Uint64 is about 1844

)

// How many tokens does the handling charge contain per 8 bytes
// Why should I calculate by 8 bytes?
// Because more handling fees can be avoided, but the handling fee field takes up more space, resulting in transactions with the same structure
// (for example, the transaction of creating diamonds) the handling fee is higher, but it ranks behind
func CalculateFeePurity(trsFee *fields.Amount, txsize uint32) uint64 {
	if trsFee.IsPositive() != true {
		return 0
	}
	if int(trsFee.Unit) < FeePurityUnit {
		return 0 // Below minimum unit
	}
	feeobj := trsFee.Copy()
	feeobj.Unit -= FeePurityUnit
	bigfee := feeobj.GetValue()
	calsize := uint64(txsize) / 8 // 每 8 个 byte 的手续费含量
	if calsize == 0 {
		calsize = 1 // Avoid division by 0
	}
	bigfee = bigfee.Div(bigfee, new(big.Int).SetUint64(calsize))
	if bigfee.Cmp(new(big.Int).SetUint64(MaxFeePurity)) == 1 {
		return MaxFeePurity // max
	} else {
		return bigfee.Uint64()
	}
}
