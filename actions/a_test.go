package actions

import (
	"fmt"
	"github.com/hacash/core/coinbase"
	"github.com/hacash/core/fields"
	"testing"
	"time"
)

func Test9834759374(t *testing.T) {
	var data = "{p:HACD20,op:mint,tick:meme,amt:120k}"
	// b0cc883502556a4466e9aac0bdeed5ff880acf77d6deed2636a2f94d3d426ee5

	fmt.Println(len(data), 64*256/1024)

}

func Test1(t *testing.T) {

	amt1, _ := fields.NewAmountFromFinString("ㄜ1:248")
	amt2, _ := fields.NewAmountFromFinString("ㄜ1:248")
	amtb1, _ := fields.NewAmountFromFinString("ㄜ12345:248")
	amtb2, _ := fields.NewAmountFromFinString("ㄜ12345:248")

	amt3, amt4, _ := coinbase.DoAppendCompoundInterestProportionOfHeightV2(amt1, amt2, 42, 1, 0)
	fmt.Println("DoAppendCompoundInterestProportionOfHeight  : ", amt3.ToFinString(), amt4.ToFinString())

	amt7, amt8, _ := coinbase.DoAppendCompoundInterestProportionOfHeightV2(amtb1, amtb2, 42, 1, 0)
	fmt.Println("DoAppendCompoundInterestProportionOfHeight  : ", amt7.ToFinString(), amt8.ToFinString())

	amt5, amt6 := coinbase.DoAppendCompoundInterest1Of10000By2500Height(amt1, amt2, 42)
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

func Test3(t *testing.T) {

	fmt.Println(moveBtcLockWeekByIdx(2048))

}
