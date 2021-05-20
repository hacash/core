package transactions

import (
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"testing"
)

func Test_t1(t *testing.T) {

	// create account
	account1 := account.CreateAccountByPassword("123456")
	account2 := account.CreateAccountByPassword("qwerty")

	addrprikeys := make(map[string][]byte)
	addrprikeys[string(account1.Address)] = account1.PrivateKey
	addrprikeys[string(account2.Address)] = account2.PrivateKey

	amount1 := fields.NewAmountNumSmallCoin(12)
	// create action
	action1 := actions.NewAction_1_SimpleToTransfer(account2.Address, amount1)
	// create tx
	tx1, e2 := NewEmptyTransaction_2_Simple(account1.Address)
	if e2 != nil {
		fmt.Println(e2)
		return
	}
	tx1.FillNeedSigns(addrprikeys, nil)
	tx1.Fee = *fields.NewAmount(244, []byte{1})
	tx1.AppendAction(action1)

	fmt.Println(tx1.Size(), tx1.FeePurity())

}
