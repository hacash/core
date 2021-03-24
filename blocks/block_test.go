package blocks

import (
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/transactions"
	"testing"
)

func Test1(t *testing.T) {

	addr1, _ := fields.CheckReadableAddress("1BjbnHwhV7VgL4kM3EsEHjyjwF5MGRNS3f")

	trs1, _ := transactions.NewEmptyTransaction_2_Simple(*addr1)
	trs1.Timestamp = 1111
	fmt.Println(trs1.HashWithFee().ToHex())

	trxs := []interfaces.Transaction{trs1, trs1, trs1, trs1, trs1, trs1}

	fmt.Println("------------")

	hashs := PickMrklListForCoinbaseTxModify(trxs)
	for _, v := range hashs {
		fmt.Println(v.ToHex())
	}

}
