package coinbase

import (
	"fmt"
	"github.com/hacash/core/fields"
	"math"
	"math/big"
)

// 计算通道利息奖励 (amt1, amt2, 1, 0.001)
// uint64 溢出 bug 导致 1BjbnHwh....MGRNS3f 地址余额计算错误
func DoAppendCompoundInterestProportionOfHeightV2(amt1 *fields.Amount, amt2 *fields.Amount, caclnum uint64, wfzn uint64) (*fields.Amount, *fields.Amount, error) {
	if caclnum == 0 {
		//panic("insnum cannot be 0.")
		return amt1, amt2, nil
	}
	if len(amt1.Numeral) > 4 || len(amt2.Numeral) > 4 {
		return nil, nil, fmt.Errorf("amount numeral bytes too long.")
	}

	amts := []*fields.Amount{amt1, amt2}
	coinnums := make([]*fields.Amount, 2)

	for i := 0; i < 2; i++ {
		//fmt.Println("----------")
		// amt
		coinnum := new(big.Int).SetBytes(amts[i].Numeral)
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
			// 返回错误
			return nil, nil, fmt.Errorf("DoAppendCompoundInterestProportionOfHeightV2 error")
		}
	}

	//fmt.Println("caclnum: ", caclnum)
	//fmt.Println(amts[0].ToFinString(), " => ", coinnums[0].ToFinString())
	//fmt.Println(amts[1].ToFinString(), " => ", coinnums[1].ToFinString())

	return coinnums[0], coinnums[1], nil

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
