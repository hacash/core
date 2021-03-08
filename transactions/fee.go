package transactions

import (
	"github.com/hacash/core/fields"
	"math/big"
)

const (
	FeePurityUnit = 232                          // = 手续费含量单位 0.00000001铢 = 1烁
	MaxFeePurity  = uint64(18446744073709550000) // = uint64的最大值 约1844枚

)

// 手续费含量 每 8 byte 的含有多少烁代币
// 为何要按 8 个 byte 计算？
// 因为可以避免手续费更多，但是由于手续费字段占用更多空间的问题，导致结构相同的交易
// （比如创建钻石的交易）手续费更高但却排在了后面
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
	calsize := uint64(txsize) / 8 // 每 8 个 byte 的手续费含量
	if calsize == 0 {
		calsize = 1 // 避免除0
	}
	bigfee = bigfee.Div(bigfee, new(big.Int).SetUint64(calsize))
	if bigfee.Cmp(new(big.Int).SetUint64(MaxFeePurity)) == 1 {
		return MaxFeePurity // max
	} else {
		return bigfee.Uint64()
	}
}
