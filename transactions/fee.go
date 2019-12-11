package transactions

import (
	"github.com/hacash/core/fields"
	"math/big"
)

const (
	FeePurityUnit = 232                          // = 0.00000001铢 = 1烁
	MaxFeePurity  = uint64(18446744073709550000) // = 约1844枚

)

// 手续费含量 每byte的含有多少烁代币
func CalculateFeePurity(trsFee *fields.Amount, txsize uint32) uint64 {
	if trsFee.IsPositive() != true {
		return 0
	}
	if int(trsFee.Unit) < FeePurityUnit {
		return 0 // 低于最小单位
	}
	feeobj := trsFee.Copy()
	feeobj.Unit -= FeePurityUnit
	bigfee := feeobj.GetValue()
	bigfee = bigfee.Div(bigfee, new(big.Int).SetUint64(uint64(txsize)))
	if bigfee.Cmp(new(big.Int).SetUint64(MaxFeePurity)) == 1 {
		return MaxFeePurity // max
	} else {
		return bigfee.Uint64()
	}
}
