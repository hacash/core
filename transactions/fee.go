package transactions

import (
	"github.com/hacash/core/fields"
)

const (
	//FeePurityUnit = 232                          // =Handling charge content unit: 0.00000001 baht = 1 yuan
	//MaxFeePurity  = uint64(18446744073709550000) // =The maximum value of Uint64 is about 1844HAC
	FeePurityUnit = 240                 // =Handling charge content unit: 0.00000001 baht = 1 zhu
	MaxFeePurity  = uint32(42_94967295) // =The maximum value of Uint64 is about 42HAC
)

// How many tokens does the handling charge contain per 32 bytes
// more handling fees can be avoided, but the handling fee field takes up more space, resulting in transactions with the same structure
// (for example, the transaction of creating diamonds) the handling fee is higher, but it ranks behind
func CalculateFeePurity(trsFee *fields.Amount, txsize uint32) uint32 { // zhu
	if trsFee.IsPositive() != true {
		return 0
	}
	feezhu := trsFee.ToZhu()
	if feezhu < 1 {
		return 0 // Below minimum unit zhu
	}
	segsz := txsize/32 + 1 // fee purity every 32 byte
	if txsize%32 == 0 && segsz > 1 {
		segsz -= 1
	}
	purity := feezhu / float64(segsz)
	if purity > float64(MaxFeePurity) {
		return MaxFeePurity
	} else {
		return uint32(purity)
	}
}
