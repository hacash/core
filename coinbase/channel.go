package coinbase

import (
	"fmt"
	"github.com/hacash/core/fields"
	"math"
	"math/big"
)

// 计算通道利息奖励 (amt1, amt2, 1, 0.001)
// uint64 溢出 bug 导致 1BjbnHwh....MGRNS3f 地址余额计算错误
func DoAppendCompoundInterestProportionOfHeightV2(amt1 *fields.Amount, amt2 *fields.Amount, caclnum uint64, wfzn uint64, interestgiveto fields.VarUint1) (*fields.Amount, *fields.Amount, error) {
	if caclnum == 0 {
		//panic("insnum cannot be 0.")
		return amt1, amt2, nil
	}
	if len(amt1.Numeral) > 4 || len(amt2.Numeral) > 4 {
		return nil, nil, fmt.Errorf("amount numeral bytes too long.")
	}

	if interestgiveto == 0 {

		// 利息两方分配
		resamt1, e1 := calculateInterestAndPrincipal(amt1, amt1, caclnum, wfzn)
		if e1 != nil {
			return nil, nil, e1
		}
		resamt2, e2 := calculateInterestAndPrincipal(amt2, amt2, caclnum, wfzn)
		if e2 != nil {
			return nil, nil, e2
		}
		return resamt1, resamt2, nil

	} else if interestgiveto == 1 {

		// 利息全部左方得到 left
		totalamt, e0 := amt1.Add(amt2)
		if e0 != nil {
			return nil, nil, e0
		}
		// 利息全给left方
		resamt1, e1 := calculateInterestAndPrincipal(amt1, totalamt, caclnum, wfzn)
		if e1 != nil {
			return nil, nil, e1
		}
		return resamt1, amt2, nil

	} else if interestgiveto == 2 {

		// 利息全部右方得到 right
		totalamt, e0 := amt1.Add(amt2)
		if e0 != nil {
			return nil, nil, e0
		}
		// 利息全给right方
		resamt2, e2 := calculateInterestAndPrincipal(amt2, totalamt, caclnum, wfzn)
		if e2 != nil {
			return nil, nil, e2
		}
		return amt1, resamt2, nil

	}

	return nil, nil, fmt.Errorf("cannot support interestgiveto value: %d", interestgiveto)

}

// 计算利息
// useramt: 用户资金
// basevalamt: 用于计算利息的本金
func calculateInterestAndPrincipal(useramt *fields.Amount, basevalamt *fields.Amount, caclnum uint64, wfzn uint64) (*fields.Amount, error) {

	var interestAmt = fields.NewEmptyAmount()
	var resultAmt = fields.NewEmptyAmount()

	coinnum := new(big.Int).SetBytes(basevalamt.Numeral)
	coinnum = new(big.Int).Mul(coinnum, big.NewInt(100000000))
	for a := uint64(0); a < caclnum; a++ {
		coinnum = new(big.Int).Mul(coinnum, big.NewInt((10000 + int64(wfzn))))
		coinnum = new(big.Int).Div(coinnum, big.NewInt(10000))
	}
	//fmt.Println(".....")
	mulnum := coinnum.Uint64()
	//fmt.Println(mulnum)
	mulnumint := int64(mulnum)
	//fmt.Println(mulnumint)
	newunit := int(basevalamt.Unit) - 8
	if newunit < 0 {
		resultAmt = useramt // 数额极小， 忽略， 余额不变
		return resultAmt, nil
	}
	for {
		if newunit < 255 && mulnumint%10 == 0 {
			mulnumint /= 10
			newunit++
		} else {
			break
		}
	}
	newNumeral := big.NewInt(int64(mulnumint)).Bytes()
	//fmt.Println(newNumeral)
	if newunit > 0 && newunit <= 255 {
		var e error
		newamt := fields.NewAmount(uint8(newunit), newNumeral)
		interestAmt, e = newamt.Sub(basevalamt)
		if e != nil {
			return nil, e
		}
		resultAmt, e = useramt.Add(interestAmt) // 正常情况
		if e != nil {
			return nil, e
		}

	} else {
		resultAmt = useramt // 计算错误， 余额不变
		// 返回错误
		return nil, fmt.Errorf("DoAppendCompoundInterestProportionOfHeightV2 error")
	}

	// ok
	return resultAmt, nil
}

// 2500个区块万分之一的复利计算
func DoAppendCompoundInterest1Of10000By2500Height(amt1 *fields.Amount, amt2 *fields.Amount, insnum uint64) (*fields.Amount, *fields.Amount) {
	if insnum == 0 {
		//panic("insnum cannot be 0.")
		return amt1, amt2
	}
	if len(amt1.Numeral) > 4 || len(amt2.Numeral) > 4 {
		panic("amount numeral bytes too long.")
	}

	amts := []*fields.Amount{amt1, amt2}
	coinnums := make([]*fields.Amount, 2)

	for i := 0; i < 2; i++ {
		//fmt.Println("----------")
		// amt
		coinnum := new(big.Int).SetBytes(amts[i].Numeral).Uint64()
		//fmt.Println(coinnum)
		mulnum := math.Pow(1.0001, float64(insnum)) * float64(coinnum) * float64(100000000)
		//fmt.Println(mulnum)
		mulnumint := int64(mulnum)
		//fmt.Println(mulnumint)
		newunit := int(amts[i].Unit) - 8
		if newunit < 0 {
			coinnums[i] = amts[i] // 数额极小， 忽略， 余额不变
			continue
		}
		for {
			if newunit < 255 && mulnumint%10 == 0 {
				mulnumint /= 10
				newunit++
			} else {
				break
			}
		}
		newNumeral := big.NewInt(int64(mulnumint)).Bytes()
		//fmt.Println(newNumeral)
		if newunit > 0 && newunit <= 255 {
			newamt := fields.NewAmount(uint8(newunit), newNumeral)
			coinnums[i] = newamt // 正常情况
		} else {
			coinnums[i] = amts[i] // 计算错误， 余额不变
		}
	}

	//fmt.Println("insnum: ", insnum)
	//fmt.Println(amts[0].ToFinString(), " => ", coinnums[0].ToFinString())
	//fmt.Println(amts[1].ToFinString(), " => ", coinnums[1].ToFinString())

	return coinnums[0], coinnums[1]

}
