package transactions

import (
	"bytes"
	"github.com/hacash/core/account"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfacev2"
)

// Take out the action created by the diamond
func CheckoutAction_4_DiamondCreateFromTx(tx interfacev2.Transaction) *actions.Action_4_DiamondCreate {

	// do add is diamond ?
	for _, act := range tx.GetActions() {
		if dcact, ok := act.(*actions.Action_4_DiamondCreate); ok {
			// is diamond create trs
			return dcact // successfully !
		}
	}

	return nil
}

// Create a general transfer transaction
func CreateOneTxOfSimpleTransfer(payacc *account.Account, toaddr fields.Address,
	amount *fields.Amount, fee *fields.Amount, timestamp int64) *Transaction_2_Simple {

	// Create a general transfer transaction
	newTrs, _ := NewEmptyTransaction_2_Simple(payacc.Address)
	newTrs.Timestamp = fields.BlockTxTimestamp(timestamp) // Use timestamp
	newTrs.Fee = *fee                                     // set fee
	tranact := actions.NewAction_1_SimpleToTransfer(toaddr, amount)
	e9 := newTrs.AppendAction(tranact)
	// Sign private key signature
	allPrivateKeyBytes := make(map[string][]byte, 1)
	allPrivateKeyBytes[string(payacc.Address)] = payacc.PrivateKey
	e9 = newTrs.FillNeedSigns(allPrivateKeyBytes, nil)
	if e9 != nil {
		return nil
	}
	return newTrs
}

// Create a BTC transfer transaction
func CreateOneTxOfBTCTransfer(payacc *account.Account, toaddr fields.Address, amount uint64,
	feeacc *account.Account, fee *fields.Amount, timestamp int64) (*Transaction_2_Simple, error) {

	// Sign private key signature
	allPrivateKeyBytes := make(map[string][]byte)
	allPrivateKeyBytes[string(feeacc.Address)] = feeacc.PrivateKey

	// Create transaction
	newTrs, _ := NewEmptyTransaction_2_Simple(feeacc.Address) // 使用手续费地址为主地址
	newTrs.Timestamp = fields.BlockTxTimestamp(timestamp)     // Use timestamp
	newTrs.Fee = *fee                                         // set fee
	var tranact interfacev2.Action = nil
	if bytes.Compare(payacc.Address, feeacc.Address) == 0 {
		tranact = &actions.Action_8_SimpleSatoshiTransfer{
			ToAddress: toaddr,
			Amount:    fields.Satoshi(amount),
		}
	} else {
		tranact = &actions.Action_11_FromToSatoshiTransfer{
			FromAddress: payacc.Address,
			ToAddress:   toaddr,
			Amount:      fields.Satoshi(amount),
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

// Create a hacd transfer transaction
func CreateOneTxOfOutfeeQuantityHACDTransfer(payacc *account.Account, toaddr fields.Address, hacdlistsplitcomma string,
	feeacc *account.Account, fee *fields.Amount, timestamp int64) (*Transaction_2_Simple, error) {

	// Diamond Watch
	var diamonds = fields.NewEmptyDiamondListMaxLen200()
	e0 := diamonds.ParseHACDlistBySplitCommaFromString(hacdlistsplitcomma)
	if e0 != nil {
		return nil, e0
	}

	// Create transaction
	newTrs, _ := NewEmptyTransaction_2_Simple(feeacc.Address) // 使用手续费地址为主地址
	newTrs.Timestamp = fields.BlockTxTimestamp(timestamp)     // Use timestamp
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
	// Sign private key signature
	allPrivateKeyBytes := make(map[string][]byte, 1)
	allPrivateKeyBytes[string(payacc.Address)] = payacc.PrivateKey
	allPrivateKeyBytes[string(feeacc.Address)] = feeacc.PrivateKey
	e9 = newTrs.FillNeedSigns(allPrivateKeyBytes, nil)
	if e9 != nil {
		return nil, e9
	}
	return newTrs, nil
}
