package transactions

import (
	"github.com/hacash/core/actions"
	"github.com/hacash/core/interfaces"
)

// 取出钻石创建的action
func CheckoutAction_4_DiamondCreateFromTx(tx interfaces.Transaction) *actions.Action_4_DiamondCreate {

	// do add is diamond ?
	for _, act := range tx.GetActions() {
		if dcact, ok := act.(*actions.Action_4_DiamondCreate); ok {
			// is diamond create trs
			return dcact // successfully !
		}
	}

	return nil
}
