package fields

import (
	"fmt"
	"math/big"
	"testing"
)

func Test_t1(t *testing.T) {

	bignum, _ := new(big.Int).SetString("12340000123400001234000012340000123400001234000012340000123400001234000012340000", 10)

	amount, _ := NewAmountByBigInt(bignum)

	fmt.Println(len(amount.Numeral))
	fmt.Println(amount.Serialize())
	fmt.Println(amount.ToFinString())

	new_amount, changeok, err := amount.EllipsisDecimalFor20SizeStore()
	fmt.Println(new_amount, changeok, err)

	if new_amount != nil {
		fmt.Println(len(new_amount.Numeral))
		fmt.Println(new_amount.Serialize())
		fmt.Println(new_amount.ToFinString())
	}

}
