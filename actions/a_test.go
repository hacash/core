package actions

import (
	"fmt"
	"github.com/hacash/core/fields"
	"testing"
)

func Test1(t *testing.T) {

	amt1, _ := fields.NewAmountFromFinString("ㄜ1:248")
	amt2, _ := fields.NewAmountFromFinString("ㄜ1:248")

	amt3, amt4 := DoAppendCompoundInterestProportionOfHeightV2(amt1, amt2, 42, 1)
	fmt.Println("DoAppendCompoundInterestProportionOfHeight  : ", amt3.ToFinString(), amt4.ToFinString())

	amt5, amt6 := DoAppendCompoundInterest1Of10000By2500Height(amt1, amt2, 42)
	fmt.Println("DoAppendCompoundInterest1Of10000By2500Height: ", amt5.ToFinString(), amt6.ToFinString())

}
