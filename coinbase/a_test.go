package coinbase

import (
	"fmt"
	"github.com/hacash/core/fields"
	"testing"
)

func Test1(t *testing.T) {

	mainaddr, _ := fields.CheckReadableAddress("1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9")

	bpnum := int64(4)

	for i := int64(0); i < bpnum*5; i++ {
		pdhei := 11 + i*10
		a, b, e := CalculationDiamondSystemLendingRedeemAmount(
			*mainaddr, *mainaddr,
			bpnum, 1,
			100,
			10, pdhei)
		if e == nil {
			fmt.Println(i+1, pdhei, a, b.ToMei())
		} else {
			fmt.Println(e)
		}

	}

}
