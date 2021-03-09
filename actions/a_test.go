package actions

import (
	"fmt"
	"github.com/hacash/core/fields"
	"testing"
	"time"
)

func Test1(t *testing.T) {

	amt1, _ := fields.NewAmountFromFinString("ㄜ1:248")
	amt2, _ := fields.NewAmountFromFinString("ㄜ1:248")

	amt3, amt4, _ := DoAppendCompoundInterestProportionOfHeightV2(amt1, amt2, 42, 1)
	fmt.Println("DoAppendCompoundInterestProportionOfHeight  : ", amt3.ToFinString(), amt4.ToFinString())

	amt5, amt6 := DoAppendCompoundInterest1Of10000By2500Height(amt1, amt2, 42)
	fmt.Println("DoAppendCompoundInterest1Of10000By2500Height: ", amt5.ToFinString(), amt6.ToFinString())

}

func Test2(t *testing.T) {

	ttt := time.Now()

	var total int64 = 0
	for i := 1; i < 10000; i++ {
		total += moveBtcCoinRewardByIdx(int64(i))
	}

	fmt.Println(total, time.Since(ttt).Seconds())

}
