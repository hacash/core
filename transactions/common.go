package transactions

import (
	"bytes"
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
	newTrs.Timestamp = fields.BlockTxTimestamp(timestamp) // 使用时间戳
	newTrs.Fee = *fee                                     // set fee
	tranact := actions.NewAction_1_SimpleToTransfer(toaddr, amount)
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

// 创建一笔 BTC 转账交易
func CreateOneTxOfBTCTransfer(payacc *account.Account, toaddr fields.Address, amount uint64,
	feeacc *account.Account, fee *fields.Amount, timestamp int64) (*Transaction_2_Simple, error) {

	// sign 私钥签名
	allPrivateKeyBytes := make(map[string][]byte)
	allPrivateKeyBytes[string(feeacc.Address)] = feeacc.PrivateKey

	// 创建交易
	newTrs, _ := NewEmptyTransaction_2_Simple(feeacc.Address) // 使用手续费地址为主地址
	newTrs.Timestamp = fields.BlockTxTimestamp(timestamp)     // 使用时间戳
	newTrs.Fee = *fee                                         // set fee
	var tranact interfaces.Action = nil
	if bytes.Compare(payacc.Address, feeacc.Address) == 0 {
		tranact = &actions.Action_8_SimpleSatoshiTransfer{
			ToAddress: toaddr,
			Amount:    fields.VarUint8(amount),
		}
	} else {
		tranact = &actions.Action_11_FromToSatoshiTransfer{
			FromAddress: payacc.Address,
			ToAddress:   toaddr,
			Amount:      fields.VarUint8(amount),
		}
		// sign add
		allPrivateKeyBytes[string(payacc.Address)] = payacc.PrivateKey
	}
	e9 := newTrs.AppendAction(tranact)
	if e9 != nil {
		return nil, e9
	}
	e9 = newTrs.FillNeedSigns(allPrivateKeyBytes, nil)
	if e9 != nil {
		return nil, e9
	}
	return newTrs, nil
}

// 创建一笔 HACD 转账交易
func CreateOneTxOfOutfeeQuantityHACDTransfer(payacc *account.Account, toaddr fields.Address, hacdlistsplitcomma string,
	feeacc *account.Account, fee *fields.Amount, timestamp int64) (*Transaction_2_Simple, error) {

	// 钻石表
	var diamonds = fields.NewEmptyDiamondListMaxLen200()
	e0 := diamonds.ParseHACDlistBySplitCommaFromString(hacdlistsplitcomma)
	if e0 != nil {
		return nil, e0
	}

	// 创建交易
	newTrs, _ := NewEmptyTransaction_2_Simple(feeacc.Address) // 使用手续费地址为主地址
	newTrs.Timestamp = fields.BlockTxTimestamp(timestamp)     // 使用时间戳
	newTrs.Fee = *fee                                         // set fee
	tranact := &actions.Action_6_OutfeeQuantityDiamondTransfer{
		FromAddress: payacc.Address,
		ToAddress:   toaddr,
		DiamondList: *diamonds,
	}
	e9 := newTrs.AppendAction(tranact)
	if e9 != nil {
		return nil, e9
	}
	// sign 私钥签名
	allPrivateKeyBytes := make(map[string][]byte, 1)
	allPrivateKeyBytes[string(payacc.Address)] = payacc.PrivateKey
	allPrivateKeyBytes[string(feeacc.Address)] = feeacc.PrivateKey
	e9 = newTrs.FillNeedSigns(allPrivateKeyBytes, nil)
	if e9 != nil {
		return nil, e9
	}
	return newTrs, nil
}
