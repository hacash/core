package transactions

import (
	"github.com/hacash/core/account"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
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

// 创建一笔普通转账交易
func CreateOneTxOfSimpleTransfer(payacc *account.Account, toaddr fields.Address,
	amount *fields.Amount, fee *fields.Amount, timestamp int64) *Transaction_2_Simple {

	// 创建普通转账交易
	newTrs, _ := NewEmptyTransaction_2_Simple(payacc.Address)
	newTrs.Timestamp = fields.VarUint5(timestamp) // 使用时间戳
	newTrs.Fee = *fee                             // set fee
	tranact := actions.NewAction_1_SimpleTransfer(toaddr, amount)
	e9 := newTrs.AppendAction(tranact)
	// sign 私钥签名
	allPrivateKeyBytes := make(map[string][]byte, 1)
	allPrivateKeyBytes[string(payacc.Address)] = payacc.PrivateKey
	e9 = newTrs.FillNeedSigns(allPrivateKeyBytes, nil)
	if e9 != nil {
		return nil
	}
	return newTrs
}
