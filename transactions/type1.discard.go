package transactions

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/sys"
	"math/big"
	"time"

	"github.com/hacash/core/account"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/crypto/btcec"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfacev2"
)

/////////////////////////////////////////////
//////// 【有签名BUG，已废弃！！！】 ///////////
/////////////////////////////////////////////

type Transaction_1_DO_NOT_USE_WITH_BUG struct {
	Timestamp fields.BlockTxTimestamp
	Address   fields.Address
	Fee       fields.Amount

	ActionCount fields.VarUint2
	Actions     []interfaces.Action

	SignCount fields.VarUint2
	Signs     []fields.Sign

	MultisignCount fields.VarUint2
	Multisigns     []fields.Multisign

	// cache data
	hash      []byte
	hashnofee []byte
}

func NewEmptyTransaction_1_Simple(master fields.Address) (*Transaction_1_DO_NOT_USE_WITH_BUG, error) {
	if !master.IsValid() {
		return nil, fmt.Errorf("Master Address is InValid ")
	}
	timeUnix := time.Now().Unix()
	return &Transaction_1_DO_NOT_USE_WITH_BUG{
		Timestamp:      fields.BlockTxTimestamp(uint64(timeUnix)),
		Address:        master,
		Fee:            *fields.NewEmptyAmount(),
		ActionCount:    fields.VarUint2(0),
		SignCount:      fields.VarUint2(0),
		MultisignCount: fields.VarUint2(0),
	}, nil
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) Copy() interfacev2.Transaction {
	// copy
	bodys, _ := trs.Serialize()
	newtrsbts := make([]byte, len(bodys))
	copy(newtrsbts, bodys)
	// create
	var newtrs = new(Transaction_1_DO_NOT_USE_WITH_BUG)
	newtrs.Parse(newtrsbts, 1) // over type
	return newtrs
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) Clone() interfaces.Transaction {
	// copy
	bodys, _ := trs.Serialize()
	newtrsbts := make([]byte, len(bodys))
	copy(newtrsbts, bodys)
	// create
	var newtrs = new(Transaction_1_DO_NOT_USE_WITH_BUG)
	newtrs.Parse(newtrsbts, 1) // over type
	return newtrs
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) Type() uint8 {
	return 1
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) ClearHash() {
	trs.hash = nil
	trs.hashnofee = nil
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) Serialize() ([]byte, error) {
	body, e0 := trs.SerializeNoSign()
	if e0 != nil {
		return nil, e0
	}
	var buffer = new(bytes.Buffer)
	buffer.Write(body)
	// sign
	b1, e1 := trs.SignCount.Serialize()
	if e1 != nil {
		return nil, e1
	}
	buffer.Write(b1)
	for i := 0; i < int(trs.SignCount); i++ {
		var bi, e = trs.Signs[i].Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bi)
	}
	// muilt sign
	b2, e2 := trs.MultisignCount.Serialize()
	if e2 != nil {
		return nil, e2
	}
	buffer.Write(b2)
	for i := 0; i < int(trs.MultisignCount); i++ {
		var bi, e = trs.Multisigns[i].Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bi)
	}
	// ok
	return buffer.Bytes(), nil
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) SerializeNoSign() ([]byte, error) {
	return trs.SerializeNoSignEx(false)
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) SerializeNoSignEx(nofee bool) ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := trs.Timestamp.Serialize()
	b2, _ := trs.Address.Serialize()
	b3, _ := trs.Fee.Serialize()
	b4, _ := trs.ActionCount.Serialize()
	buffer.Write([]byte{trs.Type()}) // type
	buffer.Write(b1)
	buffer.Write(b2)
	if !nofee {
		buffer.Write(b3) // No fee required for signature
	}
	buffer.Write(b4)
	for i := 0; i < len(trs.Actions); i++ {
		var bi, e = trs.Actions[i].Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bi)
	}
	//if nofee {
	//	fmt.Println( "SerializeNoSignEx: " + hex.EncodeToString(buffer.Bytes()))
	//}
	return buffer.Bytes(), nil
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) Parse(buf []byte, seek uint32) (uint32, error) {
	m1, e := trs.Timestamp.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	m2, e := trs.Address.Parse(buf, m1)
	if e != nil {
		return 0, e
	}
	m3, e := trs.Fee.Parse(buf, m2)
	if e != nil {
		return 0, e
	}
	m4, e := trs.ActionCount.Parse(buf, m3)
	if e != nil {
		return 0, e
	}
	iseek := m4
	for i := 0; i < int(trs.ActionCount); i++ {
		var act, sk, err = actions.ParseAction(buf, iseek)
		if err != nil {
			return 0, err
		}
		trs.Actions = append(trs.Actions, act)
		iseek = sk
	}
	iseek, e = trs.SignCount.Parse(buf, iseek)
	if e != nil {
		return 0, e
	}
	for i := 0; i < int(trs.SignCount); i++ {
		var sign fields.Sign
		iseek, e = sign.Parse(buf, iseek)
		if e != nil {
			return 0, e
		}
		trs.Signs = append(trs.Signs, sign)
	}
	iseek, e = trs.MultisignCount.Parse(buf, iseek)
	if e != nil {
		return 0, e
	}
	for i := 0; i < int(trs.MultisignCount); i++ {
		var multisign fields.Multisign
		iseek, e = multisign.Parse(buf, iseek)
		if e != nil {
			return 0, e
		}
		trs.Multisigns = append(trs.Multisigns, multisign)
	}
	return iseek, nil
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) Size() uint32 {
	totalsize := 1 +
		trs.Timestamp.Size() +
		trs.Address.Size() +
		trs.Fee.Size() +
		trs.ActionCount.Size()
	for i := 0; i < int(trs.ActionCount); i++ {
		totalsize += trs.Actions[i].Size()
	}
	totalsize += trs.SignCount.Size()
	for i := 0; i < int(trs.SignCount); i++ {
		totalsize += trs.Signs[i].Size()
	}
	totalsize += trs.MultisignCount.Size()
	for i := 0; i < int(trs.MultisignCount); i++ {
		totalsize += trs.Multisigns[i].Size()
	}
	return totalsize
}

// Transaction unique hash value
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) HashWithFee() fields.Hash {
	if trs.hash == nil {
		return trs.HashWithFeeFresh()
	}
	return trs.hash
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) HashWithFeeFresh() fields.Hash {
	stuff, _ := trs.SerializeNoSign()
	digest := fields.CalculateHash(stuff)
	trs.hash = digest[:]
	return trs.hash
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) Hash() fields.Hash {
	if trs.hashnofee == nil {
		return trs.HashFresh()
	}
	return trs.hashnofee
}
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) HashFresh() fields.Hash {
	notFee := true
	stuff, _ := trs.SerializeNoSignEx(notFee)
	digest := fields.CalculateHash(stuff)
	trs.hashnofee = digest[:]
	return trs.hashnofee
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) AppendAction(action interfacev2.Action) error {
	if trs.ActionCount >= 65535 {
		return fmt.Errorf("Action too much")
	}
	trs.ActionCount += 1
	trs.Actions = append(trs.Actions, action.(interfaces.Action))
	return nil
}

// Get the address to be signed from actions
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) RequestSignAddresses(appends []fields.Address, dropfeeaddr bool) ([]fields.Address, error) {
	if !trs.Address.IsValid() {
		return nil, fmt.Errorf("Master Address is InValid ")
	}
	requests := make([]fields.Address, 0, 32)
	for i := 0; i < int(trs.ActionCount); i++ {
		actreqs := trs.Actions[i].RequestSignAddresses()
		requests = append(requests, actreqs...)
	}
	// duplicate removal
	results := make([][]byte, len(requests))
	has := make(map[string]bool)
	if !dropfeeaddr {
		// No, add the primary address
		results = append(results, trs.Address)
	}
	// 费用方/主地址  去重
	has[string(trs.Address)] = true
	for i := 0; i < len(requests); i++ {
		strkey := string(requests[i])
		if _, ok := has[strkey]; !ok {
			results = append(results, requests[i])
		}
	}
	// return
	return requests, nil
}

// Clear all signatures
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) CleanSigns() {
	trs.SignCount = 0
	trs.Signs = []fields.Sign{}
}

// Return all signatures
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) GetSigns() []fields.Sign {
	return trs.Signs
}

// Set signature
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) SetSigns(allsigns []fields.Sign) {
	num := len(allsigns)
	if num > 65535 {
		panic("Sign is too much.")
	}
	trs.SignCount = fields.VarUint2(num)
	trs.Signs = make([]fields.Sign, 0)
	trs.Signs = append(trs.Signs, allsigns...) // copy
}

// Fill in all required signatures
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) FillTargetSign(signacc *account.Account) error {
	hashNoFee := trs.Hash()
	addrPrivateKeys := map[string][]byte{}
	addrPrivateKeys[string(signacc.Address)] = signacc.PrivateKey
	// 执行一个签名 // BUG !!! 主地址和其他地址错误的都使用 hashNoFee !!!
	return trs.addOneSign(hashNoFee, addrPrivateKeys, signacc.Address)
}

// Fill in signature
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) FillNeedSigns(addrPrivates map[string][]byte, reqs []fields.Address) error {
	// hash := trs.HashWithFeeFresh()
	hashNoFee := trs.Hash()
	requests, e0 := trs.RequestSignAddresses(nil, true)
	if e0 != nil {
		return e0
	}
	// Master signature
	e1 := trs.addOneSign(hashNoFee, addrPrivates, trs.Address)
	if e1 != nil {
		return e1
	}
	// Other signatures
	for i := 0; i < len(requests); i++ {
		e1 := trs.addOneSign(hashNoFee, addrPrivates, requests[i])
		if e1 != nil {
			return e1
		}
	}
	// Filled successfully
	return nil
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) addOneSign(hash []byte, addrPrivates map[string][]byte, address []byte) error {
	privitebytes, has := addrPrivates[string(address)]
	if !has {
		return fmt.Errorf("Private Key '" + account.Base58CheckEncode(address) + "' necessary")
	}
	privite, e1 := account.GetAccountByPriviteKey(privitebytes)
	if e1 != nil {
		return fmt.Errorf("Private Key '" + account.Base58CheckEncode(address) + "' error")
	}
	signature, e2 := privite.Private.Sign(hash)
	if e2 != nil {
		return fmt.Errorf("Private Key '" + account.Base58CheckEncode(address) + "' do sign error")
	}
	// append
	trs.SignCount += 1
	trs.Signs = append(trs.Signs, fields.Sign{
		PublicKey: privite.PublicKey,
		Signature: signature.Serialize64(),
	})
	return nil
}

// Verify required signatures
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) VerifyAllNeedSigns() (bool, error) {
	//hash := trs.HashWithFeeFresh()
	hashNoFee := trs.Hash()
	// Verify all signatures that need to be verified
	requests, e := trs.RequestSignAddresses(nil, true)
	if e != nil {
		return false, e
	}
	if requests == nil || len(requests) == 0 {
		return true, nil // Pass without verifying anything
	}
	allSigns := make(map[string]fields.Sign)
	for i := 0; i < len(trs.Signs); i++ {
		sig := trs.Signs[i]
		addr := account.NewAddressFromPublicKeyV0(sig.PublicKey)
		allSigns[string(addr)] = sig
	}
	// 验证主签名 /// BUG ///
	ok, e := verifyOneSignature_not_use_with_bug(allSigns, trs.Address, hashNoFee)
	if e != nil || !ok {
		return ok, e
	}
	// Verify all other signatures
	for i := 0; i < len(requests); i++ {
		ok, e := verifyOneSignature_not_use_with_bug(allSigns, requests[i], hashNoFee)
		if e != nil || !ok {
			return ok, e
		}
	}
	// Validation successful
	return true, nil
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) VerifyTargetSigns(reqaddrs []fields.Address) (bool, error) {
	hashNoFee := trs.Hash()
	// Start judging
	// 验证主签名 /// BUG ///
	allSigns := make(map[string]fields.Sign)
	for i := 0; i < len(trs.Signs); i++ {
		sig := trs.Signs[i]
		addrbts := account.NewAddressFromPublicKeyV0(sig.PublicKey)
		addr := fields.Address(addrbts)
		allSigns[string(addr)] = sig
	}
	// To verify
	for _, v := range reqaddrs {
		if _, hav := allSigns[string(v)]; hav {
			// /// BUG ///
			ok, e := verifyOneSignature_not_use_with_bug(allSigns, v, hashNoFee)
			if !ok || e != nil {
				return ok, e // Validation failed
			}
			// next
		} else {
			// Validation failed
			return false, nil
		}
	}
	// Validation successful
	return true, nil

}

func verifyOneSignature_not_use_with_bug(allSigns map[string]fields.Sign, address fields.Address, hash []byte) (bool, error) {

	main, ok := allSigns[string(address)]
	if !ok {
		return false, fmt.Errorf("address %s signature not find!", address.ToReadable())
	}
	sigobj, e3 := btcec.ParseSignatureByte64(main.Signature)
	if e3 != nil {
		return false, e3
	}
	pubKey, e4 := btcec.ParsePubKey(main.PublicKey, btcec.S256())
	if e4 != nil {
		return false, e4
	}
	verok := sigobj.Verify(hash, pubKey)
	if !verok {
		return false, fmt.Errorf("verify address %s signature fail.", address.ToReadable())
	}
	// ok
	return true, nil
}

// Balance check required
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) RequestAddressBalance() ([][]byte, []big.Int, error) {
	return nil, nil, nil
}

// 修改 / 恢复 状态数据库
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) WriteInChainState(state interfaces.ChainStateOperation) error {
	/*********************************************************/
	/******* 在区块37000 以上不能接受 trs_type==1 的交易 ********/
	/******* 从而解决第一种交易类型的签名验证的BUG问题     ********/
	/*********************************************************/
	if sys.TransactionSystemCheckChainID > 0 {
		return fmt.Errorf("Test or Fork chain cannot use transaction Type <!>")
	}
	if state.GetPendingBlockHeight() > 37000 {
		return fmt.Errorf("Transaction type<1> be discard DO_NOT_USE_WITH_BUG")
	}
	// actions
	for i := 0; i < len(trs.Actions); i++ {
		trs.Actions[i].(interfaces.Action).SetBelongTrs(trs)
		e := trs.Actions[i].(interfaces.Action).WriteInChainState(state)
		if e != nil {
			return e
		}
	}
	// Deduct handling charges
	return actions.DoSubBalanceFromChainStateV3(state, trs.Address, trs.Fee)
}

// 修改 / 恢复 状态数据库
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) WriteinChainState(state interfacev2.ChainStateOperation) error {
	/*********************************************************/
	/******* 在区块37000 以上不能接受 trs_type==1 的交易 ********/
	/******* 从而解决第一种交易类型的签名验证的BUG问题     ********/
	/*********************************************************/
	if sys.TransactionSystemCheckChainID > 0 {
		return fmt.Errorf("Test or Fork chain cannot use transaction Type <!>")
	}
	if state.GetPendingBlockHeight() > 37000 {
		return fmt.Errorf("Transaction type<1> be discard DO_NOT_USE_WITH_BUG")
	}
	// actions
	for i := 0; i < len(trs.Actions); i++ {
		trs.Actions[i].(interfacev2.Action).SetBelongTransaction(trs)
		e := trs.Actions[i].(interfacev2.Action).WriteinChainState(state)
		if e != nil {
			return e
		}
	}
	// Deduct handling charges
	return actions.DoSubBalanceFromChainState(state, trs.Address, trs.Fee)
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) RecoverChainState(state interfacev2.ChainStateOperation) error {

	panic("RecoverChainState be deprecated")

	// actions
	for i := len(trs.Actions) - 1; i >= 0; i-- {
		trs.Actions[i].(interfacev2.Action).SetBelongTransaction(trs)
		e := trs.Actions[i].(interfacev2.Action).RecoverChainState(state)
		if e != nil {
			return e
		}
	}
	// Refund service charge
	return actions.DoAddBalanceFromChainState(state, trs.Address, trs.Fee)
}

// Service charge content: how many tokens per byte
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) FeePurity() uint64 {
	return CalculateFeePurity(&trs.Fee, trs.Size())
}

// query
func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) GetAddress() fields.Address {
	return trs.Address
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) SetAddress(addr fields.Address) {
	trs.Address = addr
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) GetFeeOfMinerRealReceived() *fields.Amount {
	return &trs.Fee
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) GetFee() *fields.Amount {
	return &trs.Fee
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) SetFee(fee *fields.Amount) {
	trs.Fee = *fee
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) GetActions() []interfacev2.Action {
	var list = make([]interfacev2.Action, len(trs.Actions))
	for i, v := range trs.Actions {
		list[i] = v.(interfacev2.Action)
	}
	return list
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) GetActionList() []interfaces.Action {
	return trs.Actions
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) GetTimestamp() uint64 { // time stamp
	return uint64(trs.Timestamp)
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) SetMessage(fields.TrimString16) {
}

func (trs *Transaction_1_DO_NOT_USE_WITH_BUG) GetMessage() fields.TrimString16 {
	return fields.TrimString16("")
}
