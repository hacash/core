package transactions

import (
	"fmt"
	"github.com/hacash/core/fields"
	"testing"
)

func Test_coinbaseCopy(t *testing.T) {

	cbtrs := NewTransaction_0_Coinbase()
	addr, _ := fields.CheckReadableAddress("1AVRuFXNFi3rdMrPH4hdqSgFrEBnWisWaS")
	cbtrs.Address = *addr
	reward := fields.NewAmountNumSmallCoin(1)
	cbtrs.Reward = *reward
	cbtrs.Message = "ABC123"

	fmt.Println(cbtrs.Serialize())

	clonetrs := cbtrs.Copy()

	fmt.Println(clonetrs.Serialize())

}
