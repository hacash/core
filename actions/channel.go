package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/coinbase"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
)

/**
 * 支付通道交易类型
 */

// Open payment channel
type Action_2_OpenPaymentChannel struct {
	ChannelId    fields.ChannelId // Channel ID
	LeftAddress  fields.Address   // Account 1
	LeftAmount   fields.Amount    // Locked amount
	RightAddress fields.Address   // Account 2
	RightAmount  fields.Amount    // Locked amount

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_2_OpenPaymentChannel) Kind() uint16 {
	return 2
}

func (elm *Action_2_OpenPaymentChannel) Size() uint32 {
	return 2 + elm.ChannelId.Size() +
		elm.LeftAddress.Size() +
		elm.LeftAmount.Size() +
		elm.RightAddress.Size() +
		elm.RightAmount.Size()
}

// json api
func (elm *Action_2_OpenPaymentChannel) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_2_OpenPaymentChannel) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var idBytes, _ = elm.ChannelId.Serialize()
	var addr1Bytes, _ = elm.LeftAddress.Serialize()
	var amt1Bytes, _ = elm.LeftAmount.Serialize()
	var addr2Bytes, _ = elm.RightAddress.Serialize()
	var amt2Bytes, _ = elm.RightAmount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(idBytes)
	buffer.Write(addr1Bytes)
	buffer.Write(amt1Bytes)
	buffer.Write(addr2Bytes)
	buffer.Write(amt2Bytes)
	return buffer.Bytes(), nil
}

func (elm *Action_2_OpenPaymentChannel) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_2_OpenPaymentChannel) RequestSignAddresses() []fields.Address {
	reqs := make([]fields.Address, 2)
	reqs[0] = elm.LeftAddress
	reqs[1] = elm.RightAddress
	return reqs
}

func (act *Action_2_OpenPaymentChannel) WriteInChainState(state interfaces.ChainStateOperation) error {
	return openPaymentChannelWriteinChainState(state,
		act.ChannelId,
		fields.VarUint2(uint16(5000)), // The lock-in period unilaterally proposed is about 17 days
		fields.VarUint1(0),            // By default, interest distribution is shared by the two parties according to the close amount
		act.LeftAddress,
		act.LeftAmount,
		fields.NewEmptySatoshiVariation(), // 0 SAT
		act.RightAddress,
		act.RightAmount,
		fields.NewEmptySatoshiVariation(), // 0 SAT
	)
	/*
		var e error
		// Query whether the channel exists
		sto, e := state.Channel(act.ChannelId)
		if e != nil {
			return e
		}
		// Channel IDS with the same left and right addresses and closed by consensus can be reused
		var reuseVersion fields.VarUint4 = 1
		var isIdCanUse = (sto == nil) ||
			(sto.IsAgreementClosed() && sto.LeftAddress.Equal(act.LeftAddress) && sto.RightAddress.Equal(act.RightAddress))
		if isIdCanUse == false {
			if sto.IsFinalDistributionClosed() {
				return fmt.Errorf("Payment Channel Id <%s> finally arbitration closed.", hex.EncodeToString(act.ChannelId))
			} else if sto.IsOpening() {
				return fmt.Errorf("Payment Channel Id <%s> is opening.", hex.EncodeToString(act.ChannelId))
			}
			return fmt.Errorf("Payment Channel Id <%s> already exist.", hex.EncodeToString(act.ChannelId))
		}
		if sto != nil {
			reuseVersion = sto.ReuseVersion + 1 // Reuse version number growth
		}
		// Channel ID validity
		if len(act.ChannelId) != stores.ChannelIdLength || act.ChannelId[0] == 0 || act.ChannelId[stores.ChannelIdLength-1] == 0 {
			return fmt.Errorf("Payment Channel Id <%s> format error.", hex.EncodeToString(act.ChannelId))
		}
		// Two addresses cannot be the same
		if act.LeftAddress.Equal(act.RightAddress) {
			return fmt.Errorf("Left address cannot equal with right address.")
		}
		// Check the number of digits stored in the amount
		labt, _ := act.LeftAmount.Serialize()
		rabt, _ := act.RightAmount.Serialize()
		if len(labt) > 6 || len(rabt) > 6 {
			// Avoid locking the storage digits of funds too long, resulting in the value storage digits after compound interest calculation exceeding the maximum range
			return fmt.Errorf("Payment Channel create error: left or right Amount bytes too long.")
		}
		// Cannot be negative, or both channels are zero at the same time (one can be positive and the other zero)
		if (!act.LeftAmount.IsPositive() || !act.RightAmount.IsPositive()) ||
			(act.LeftAmount.IsEmpty() && act.RightAmount.IsEmpty()) {
			return fmt.Errorf("Action_2_OpenPaymentChannel Payment Channel create error: left or right Amount is not positive.")
		}
		// Check whether the balance is sufficient
		bls1, e := state.Balance(act.LeftAddress)
		if e != nil {
			return e
		}
		if bls1 == nil {
			return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance cannot empty.", act.LeftAddress.ToReadable())
		}
		amt1 := bls1.Hacash
		if amt1.LessThan(&act.LeftAmount) {
			return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance is not enough. need %s but got %s", act.LeftAddress.ToReadable(), act.LeftAmount.ToFinString(), amt1.ToFinString())
		}
		bls2, e := state.Balance(act.RightAddress)
		if e != nil {
			return e
		}
		if bls2 == nil {
			return fmt.Errorf("Address %s Balance is not enough.", act.RightAddress.ToReadable())
		}
		amt2 := bls2.Hacash
		if amt2.LessThan(&act.RightAmount) {
			return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance is not enough. need %s but got %s", act.RightAddress.ToReadable(), act.RightAmount.ToFinString(), amt2.ToFinString())
		}

		curheight := state.GetPendingBlockHeight()
		// Create channel
		var storeItem = stores.CreateEmptyChannel()
		storeItem.BelongHeight = fields.BlockHeight(curheight)
		storeItem.ArbitrationLockBlock = fields.VarUint2(uint16(5000)) // The lock-in period unilaterally proposed is about 17 days
		storeItem.InterestAttribution = fields.VarUint1(0)             // By default, interest distribution is shared by the two parties according to the close amount
		storeItem.LeftAddress = act.LeftAddress
		storeItem.LeftAmount = act.LeftAmount
		storeItem.RightAddress = act.RightAddress
		storeItem.RightAmount = act.RightAmount
		storeItem.ReuseVersion = reuseVersion // Reuse version number
		storeItem.SetOpening()                // Open status
		// testing environment
		if sys.TestDebugLocalDevelopmentMark {
			storeItem.ArbitrationLockBlock = fields.VarUint2(uint16(20))
		}
		// Deduction balance
		e = DoSubBalanceFromChainStateV3(state, act.LeftAddress, act.LeftAmount)
		if e != nil {
			return e
		}
		e = DoSubBalanceFromChainStateV3(state, act.RightAddress, act.RightAmount)
		if e != nil {
			return e
		}
		// Storage channel
		e = state.ChannelCreate(act.ChannelId, storeItem)
		if e != nil {
			return e
		}
		// Total supply statistics
		totalsupply, e := state.ReadTotalSupply()
		if e != nil {
			return e
		}
		// Cumulative locked HAC
		addamt := act.LeftAmount.ToMei() + act.RightAmount.ToMei()
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfLocatedHACInChannel, addamt)
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfChannelOfOpening, 1)
		// update total supply
		e3 := state.UpdateSetTotalSupply(totalsupply)
		if e3 != nil {
			return e3
		}
		//
		return nil
	*/
}

func (act *Action_2_OpenPaymentChannel) WriteinChainState(state interfacev2.ChainStateOperation) error {
	var e error
	// Query whether the channel exists
	sto, e := state.Channel(act.ChannelId)
	if e != nil {
		return e
	}
	// Channel IDS with the same left and right addresses and closed by consensus can be reused
	var reuseVersion fields.VarUint4 = 1
	var isIdCanUse = (sto == nil) ||
		(sto.IsAgreementClosed() && sto.LeftAddress.Equal(act.LeftAddress) && sto.RightAddress.Equal(act.RightAddress))
	if isIdCanUse == false {
		if sto.IsFinalDistributionClosed() {
			return fmt.Errorf("Payment Channel Id <%s> finally arbitration closed.", hex.EncodeToString(act.ChannelId))
		} else if sto.IsOpening() {
			return fmt.Errorf("Payment Channel Id <%s> is opening.", hex.EncodeToString(act.ChannelId))
		}
		return fmt.Errorf("Payment Channel Id <%s> already exist.", hex.EncodeToString(act.ChannelId))
	}
	if sto != nil {
		reuseVersion = sto.ReuseVersion + 1 // Reuse version number growth
	}
	// Channel ID validity
	if len(act.ChannelId) != stores.ChannelIdLength || act.ChannelId[0] == 0 || act.ChannelId[stores.ChannelIdLength-1] == 0 {
		return fmt.Errorf("Payment Channel Id <%s> format error.", hex.EncodeToString(act.ChannelId))
	}
	// Two addresses cannot be the same
	if act.LeftAddress.Equal(act.RightAddress) {
		return fmt.Errorf("Left address cannot equal with right address.")
	}
	// Check the number of digits stored in the amount
	labt, _ := act.LeftAmount.Serialize()
	rabt, _ := act.RightAmount.Serialize()
	if len(labt) > 6 || len(rabt) > 6 {
		// Avoid locking the storage digits of funds too long, resulting in the value storage digits after compound interest calculation exceeding the maximum range
		return fmt.Errorf("Payment Channel create error: left or right Amount bytes too long.")
	}
	// Cannot be negative, or both channels are zero at the same time (one can be positive and the other zero)
	if (!act.LeftAmount.IsPositive() || !act.RightAmount.IsPositive()) ||
		(act.LeftAmount.IsEmpty() && act.RightAmount.IsEmpty()) {
		return fmt.Errorf("Action_2_OpenPaymentChannel Payment Channel create error: left or right Amount is not positive.")
	}
	// Check whether the balance is sufficient
	bls1, e := state.Balance(act.LeftAddress)
	if e != nil {
		return e
	}
	if bls1 == nil {
		return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance cannot empty.", act.LeftAddress.ToReadable())
	}
	amt1 := bls1.Hacash
	if amt1.LessThan(&act.LeftAmount) {
		return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance is not enough. need %s but got %s", act.LeftAddress.ToReadable(), act.LeftAmount.ToFinString(), amt1.ToFinString())
	}
	bls2, e := state.Balance(act.RightAddress)
	if e != nil {
		return e
	}
	if bls2 == nil {
		return fmt.Errorf("Address %s Balance is not enough.", act.RightAddress.ToReadable())
	}
	amt2 := bls2.Hacash
	if amt2.LessThan(&act.RightAmount) {
		return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance is not enough. need %s but got %s", act.RightAddress.ToReadable(), act.RightAmount.ToFinString(), amt2.ToFinString())
	}
	curheight := state.GetPendingBlockHeight()
	// Create channel
	var storeItem = stores.CreateEmptyChannel()
	storeItem.BelongHeight = fields.BlockHeight(curheight)
	storeItem.ArbitrationLockBlock = fields.VarUint2(uint16(5000)) // The lock-in period unilaterally proposed is about 17 days
	storeItem.InterestAttribution = fields.VarUint1(0)             // By default, interest distribution is shared by the two parties according to the close amount
	storeItem.LeftAddress = act.LeftAddress
	storeItem.LeftAmount = act.LeftAmount
	storeItem.RightAddress = act.RightAddress
	storeItem.RightAmount = act.RightAmount
	storeItem.ReuseVersion = reuseVersion // Reuse version number
	storeItem.SetOpening()                // Open status
	// testing environment
	if sys.TestDebugLocalDevelopmentMark {
		storeItem.ArbitrationLockBlock = fields.VarUint2(uint16(20))
	}
	// Deduction balance
	e = DoSubBalanceFromChainState(state, act.LeftAddress, act.LeftAmount)
	if e != nil {
		return e
	}
	e = DoSubBalanceFromChainState(state, act.RightAddress, act.RightAmount)
	if e != nil {
		return e
	}
	// Storage channel
	e = state.ChannelCreate(act.ChannelId, storeItem)
	if e != nil {
		return e
	}
	// Total supply statistics
	totalsupply, e := state.ReadTotalSupply()
	if e != nil {
		return e
	}
	// Cumulative locked HAC
	addamt := act.LeftAmount.ToMei() + act.RightAmount.ToMei()
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfLocatedHACInChannel, addamt)
	totalsupply.DoAddUint(stores.TotalSupplyStoreTypeOfChannelOfOpening, 1)
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	//
	return nil
}

func (act *Action_2_OpenPaymentChannel) RecoverChainState(state interfacev2.ChainStateOperation) error {

	panic("RecoverChainState be deprecated")

	sto, e := state.Channel(act.ChannelId)
	if e != nil {
		return e
	}
	if sto.ReuseVersion > 1 {
		sto.ReuseVersion = sto.ReuseVersion - 1 // Reuse version number reduction
	} else {
		// Delete channel
		state.ChannelDelete(act.ChannelId)
	}

	// Restore balance
	DoAddBalanceFromChainState(state, act.LeftAddress, act.LeftAmount)
	DoAddBalanceFromChainState(state, act.RightAddress, act.RightAmount)
	// Total supply statistics
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	// Rollback unlocked HAC
	addamt := act.LeftAmount.ToMei() + act.RightAmount.ToMei()
	totalsupply.DoSub(stores.TotalSupplyStoreTypeOfLocatedHACInChannel, addamt)
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	return nil
}

func (elm *Action_2_OpenPaymentChannel) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_2_OpenPaymentChannel) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_2_OpenPaymentChannel) IsBurning90PersentTxFees() bool {
	return false
}

/////////////////////////////////////////////////////////////////

// Close and settle the payment channel (when the fund allocation remains unchanged)
type Action_3_ClosePaymentChannel struct {
	ChannelId fields.ChannelId // Channel ID

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_3_ClosePaymentChannel) Kind() uint16 {
	return 3
}

func (elm *Action_3_ClosePaymentChannel) Size() uint32 {
	return 2 + elm.ChannelId.Size()
}

// json api
func (elm *Action_3_ClosePaymentChannel) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_3_ClosePaymentChannel) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var idBytes, _ = elm.ChannelId.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(idBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_3_ClosePaymentChannel) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = elm.ChannelId.Parse(buf, seek)
	return seek, nil
}

func (elm *Action_3_ClosePaymentChannel) RequestSignAddresses() []fields.Address {
	// During execution, check the signature after querying the data
	return []fields.Address{}
}

func (act *Action_3_ClosePaymentChannel) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}
	// Query channel
	paychan, e := state.Channel(act.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// Judge that the channel has been closed
	if paychan.IsClosed() {
		return fmt.Errorf("Payment Channel <%s> is be closed.", hex.EncodeToString(act.ChannelId))
	}
	// 检查两个账户的签名 // 仅仅验证这两个地址
	signok, e1 := act.belong_trs_v3.VerifyTargetSigns([]fields.Address{paychan.LeftAddress, paychan.RightAddress})
	if e1 != nil {
		return e1
	}
	if !signok { // signature check failed
		return fmt.Errorf("Payment Channel <%s> address signature verify fail.", hex.EncodeToString(act.ChannelId))
	}

	// Write status
	// Calculate channel interest using the deposited amount
	// Sat raised from channel
	leftSAT := paychan.LeftSatoshi.GetRealSatoshi()
	rightSAT := paychan.RightSatoshi.GetRealSatoshi()
	return closePaymentChannelWriteinChainStateV3(state, act.ChannelId, paychan,
		nil, nil, leftSAT, rightSAT, false)
}

func (act *Action_3_ClosePaymentChannel) WriteinChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// Query channel
	paychan, e := state.Channel(act.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// Judge that the channel has been closed
	if paychan.IsClosed() {
		return fmt.Errorf("Payment Channel <%s> is be closed.", hex.EncodeToString(act.ChannelId))
	}
	// 检查两个账户的签名 // 仅仅验证这两个地址
	signok, e1 := act.belong_trs.VerifyTargetSigns([]fields.Address{paychan.LeftAddress, paychan.RightAddress})
	if e1 != nil {
		return e1
	}
	if !signok { // signature check failed
		return fmt.Errorf("Payment Channel <%s> address signature verify fail.", hex.EncodeToString(act.ChannelId))
	}

	// Write status
	// Calculate channel interest using the deposited amount
	// Sat raised from channel
	leftSAT := paychan.LeftSatoshi.GetRealSatoshi()
	rightSAT := paychan.RightSatoshi.GetRealSatoshi()
	return closePaymentChannelWriteinChainState(state, act.ChannelId, paychan,
		nil, nil, leftSAT, rightSAT, false)
}

func (act *Action_3_ClosePaymentChannel) RecoverChainState(state interfacev2.ChainStateOperation) error {
	return closePaymentChannelRecoverChainState_deprecated(state, act.ChannelId, nil, nil, false)
}

func (elm *Action_3_ClosePaymentChannel) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_3_ClosePaymentChannel) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_3_ClosePaymentChannel) IsBurning90PersentTxFees() bool {
	return false
}

/////////////////////////////////////////////////////////////////

// Close and settle payment channels (fund allocation changes)
type Action_12_ClosePaymentChannelBySetupAmount struct {
	ChannelId    fields.ChannelId        // Channel ID
	LeftAddress  fields.Address          // Left account
	LeftAmount   fields.Amount           // Final allocation amount on the left
	LeftSatoshi  fields.SatoshiVariation // Sat assigned on the left
	RightAddress fields.Address          // Right account
	RightAmount  fields.Amount           // Right final allocation amount
	RightSatoshi fields.SatoshiVariation // Sat assigned on the right

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) Kind() uint16 {
	return 12
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) Size() uint32 {
	return 2 + elm.ChannelId.Size() +
		elm.LeftAddress.Size() +
		elm.LeftAmount.Size() +
		elm.LeftSatoshi.Size() +
		elm.RightAddress.Size() +
		elm.RightAmount.Size() +
		elm.RightSatoshi.Size()
}

// json api
func (elm *Action_12_ClosePaymentChannelBySetupAmount) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.ChannelId.Serialize()
	var bt2, _ = elm.LeftAddress.Serialize()
	var bt3, _ = elm.LeftAmount.Serialize()
	var bt4, _ = elm.LeftSatoshi.Serialize()
	var bt5, _ = elm.RightAddress.Serialize()
	var bt6, _ = elm.RightAmount.Serialize()
	var bt7, _ = elm.RightSatoshi.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	buffer.Write(bt3)
	buffer.Write(bt4)
	buffer.Write(bt5)
	buffer.Write(bt6)
	buffer.Write(bt7)
	return buffer.Bytes(), nil
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftSatoshi.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightSatoshi.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) RequestSignAddresses() []fields.Address {
	// Signature required
	return []fields.Address{
		elm.LeftAddress,
		elm.RightAddress,
	}
}

func (act *Action_12_ClosePaymentChannelBySetupAmount) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	// Query channel
	paychan, e := state.Channel(act.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// Check whether the two accounts match
	if paychan.LeftAddress.NotEqual(act.LeftAddress) ||
		paychan.RightAddress.NotEqual(act.RightAddress) {
		// Address check failed
		return fmt.Errorf("Payment Channel <%s> address not match.", act.RightAddress.ToReadable())
	}
	// Write status
	leftSAT := act.LeftSatoshi.GetRealSatoshi()
	rightSAT := act.RightSatoshi.GetRealSatoshi()
	return closePaymentChannelWriteinChainStateV3(state, act.ChannelId,
		paychan, &act.LeftAmount, &act.RightAmount, leftSAT, rightSAT, false)
}

func (act *Action_12_ClosePaymentChannelBySetupAmount) WriteinChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	// Query channel
	paychan, e := state.Channel(act.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// Check whether the two accounts match
	if paychan.LeftAddress.NotEqual(act.LeftAddress) ||
		paychan.RightAddress.NotEqual(act.RightAddress) {
		// Address check failed
		return fmt.Errorf("Payment Channel <%s> address not match.", act.RightAddress.ToReadable())
	}
	// Write status
	leftSAT := act.LeftSatoshi.GetRealSatoshi()
	rightSAT := act.RightSatoshi.GetRealSatoshi()
	return closePaymentChannelWriteinChainState(state, act.ChannelId,
		paychan, &act.LeftAmount, &act.RightAmount, leftSAT, rightSAT, false)
}

func (act *Action_12_ClosePaymentChannelBySetupAmount) RecoverChainState(state interfacev2.ChainStateOperation) error {
	return closePaymentChannelRecoverChainState_deprecated(state, act.ChannelId, &act.LeftAmount, &act.RightAmount, false)
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_12_ClosePaymentChannelBySetupAmount) IsBurning90PersentTxFees() bool {
	return false
}

//////////////////////////////////////

// Closing and settlement payment channels (fund allocation changes) only provide left balance allocation, and automatically calculate the allocation of right
type Action_21_ClosePaymentChannelBySetupOnlyLeftAmount struct {
	ChannelId   fields.ChannelId        // Channel ID
	LeftAmount  fields.Amount           // Left final distribution HAC
	LeftSatoshi fields.SatoshiVariation // Left final assignment sat

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) Kind() uint16 {
	return 21
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) Size() uint32 {
	return 2 + elm.ChannelId.Size() +
		elm.LeftAmount.Size()
}

// json api
func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.ChannelId.Serialize()
	var bt2, _ = elm.LeftAmount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	return buffer.Bytes(), nil
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) RequestSignAddresses() []fields.Address {
	// During execution, check the signature after querying the data
	return []fields.Address{}
}

func (act *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) WriteInChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}
	// Query channel
	paychan, e := state.Channel(act.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// Check the signatures of the two accounts and verify only the two addresses
	signok, e0 := act.belong_trs_v3.VerifyTargetSigns([]fields.Address{paychan.LeftAddress, paychan.RightAddress})
	if e0 != nil {
		return e0
	}
	if !signok { // signature check failed
		return fmt.Errorf("Payment Channel <%s> address signature verify fail.", hex.EncodeToString(act.ChannelId))
	}
	// Allocation amount can be zero but not negative
	if act.LeftAmount.IsNegative() {
		return fmt.Errorf("Payment channel distribution amount cannot be negative.")
	}
	// Check allocation amount
	var totalAmount, e1 = paychan.LeftAmount.Add(&paychan.RightAmount)
	if e1 != nil {
		return e1
	}
	// The allocated amount cannot exceed the total amount
	if act.LeftAmount.MoreThan(totalAmount) {
		return fmt.Errorf("LeftAmount %s cannot more than total amount %s.",
			act.LeftAmount.ToFinString(), totalAmount.ToFinString())
	}
	// Calculate right amount
	var closedRightAmount, e2 = totalAmount.Sub(&act.LeftAmount)
	if e2 != nil {
		return e2
	}
	// Write status
	leftOldSAT := paychan.LeftSatoshi.GetRealSatoshi()
	rightOldSAT := paychan.RightSatoshi.GetRealSatoshi()
	totalOldSAT := leftOldSAT + rightOldSAT
	leftNewSAT := act.LeftSatoshi.GetRealSatoshi()
	if leftNewSAT > totalOldSAT {
		// One side allocation amount cannot exceed the total amount
		return fmt.Errorf("Left satoshi %d cannot more than total %d.", leftNewSAT, totalOldSAT)
	}
	rightNewSAT := totalOldSAT - leftNewSAT
	return closePaymentChannelWriteinChainStateV3(state, act.ChannelId,
		paychan, &act.LeftAmount, closedRightAmount, leftNewSAT, rightNewSAT, false)
}

func (act *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) WriteinChainState(state interfacev2.ChainStateOperation) error {

	//if !sys.TestDebugLocalDevelopmentMark {
	//	return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	//}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// Query channel
	paychan, e := state.Channel(act.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// Check the signatures of the two accounts and verify only the two addresses
	signok, e0 := act.belong_trs.VerifyTargetSigns([]fields.Address{paychan.LeftAddress, paychan.RightAddress})
	if e0 != nil {
		return e0
	}
	if !signok { // signature check failed
		return fmt.Errorf("Payment Channel <%s> address signature verify fail.", hex.EncodeToString(act.ChannelId))
	}
	// Allocation amount can be zero but not negative
	if act.LeftAmount.IsNegative() {
		return fmt.Errorf("Payment channel distribution amount cannot be negative.")
	}
	// Check allocation amount
	var totalAmount, e1 = paychan.LeftAmount.Add(&paychan.RightAmount)
	if e1 != nil {
		return e1
	}
	// The allocated amount cannot exceed the total amount
	if act.LeftAmount.MoreThan(totalAmount) {
		return fmt.Errorf("LeftAmount %s cannot more than total amount %s.",
			act.LeftAmount.ToFinString(), totalAmount.ToFinString())
	}
	// Calculate right amount
	var closedRightAmount, e2 = totalAmount.Sub(&act.LeftAmount)
	if e2 != nil {
		return e2
	}
	// Write status
	leftOldSAT := paychan.LeftSatoshi.GetRealSatoshi()
	rightOldSAT := paychan.RightSatoshi.GetRealSatoshi()
	totalOldSAT := leftOldSAT + rightOldSAT
	leftNewSAT := act.LeftSatoshi.GetRealSatoshi()
	if leftNewSAT > totalOldSAT {
		// One side allocation amount cannot exceed the total amount
		return fmt.Errorf("Left satoshi %d cannot more than total %d.", leftNewSAT, totalOldSAT)
	}
	rightNewSAT := totalOldSAT - leftNewSAT
	return closePaymentChannelWriteinChainState(state, act.ChannelId,
		paychan, &act.LeftAmount, closedRightAmount, leftNewSAT, rightNewSAT, false)
}

func (act *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) RecoverChainState(state interfacev2.ChainStateOperation) error {

	// Query channel
	paychan, e := state.Channel(act.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// Check allocation amount
	var totalAmount, _ = paychan.LeftAmount.Add(&paychan.RightAmount)
	// Calculate right amount
	var closedRightAmount, _ = totalAmount.Sub(&act.LeftAmount)
	return closePaymentChannelRecoverChainState_deprecated(state, act.ChannelId, &act.LeftAmount, closedRightAmount, false)
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) IsBurning90PersentTxFees() bool {
	return false
}

// Open payment channel
type Action_31_OpenPaymentChannelWithSatoshi struct {
	ChannelId            fields.ChannelId // Channel ID
	ArbitrationLockBlock fields.VarUint2  // Number of blocks to be locked for unilateral end channel
	InterestAttribution  fields.VarUint1  // Interest attribution of 1% annualized: 0 Press end to assign 1 All to left 2 Give it all right
	LeftAddress          fields.Address   // Account 1
	LeftAmount           fields.Amount    // Locked amount
	LeftSatoshi          fields.SatoshiVariation
	RightAddress         fields.Address // Account 2
	RightAmount          fields.Amount  // Locked amount
	RightSatoshi         fields.SatoshiVariation

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_31_OpenPaymentChannelWithSatoshi) Kind() uint16 {
	return 31
}

func (elm *Action_31_OpenPaymentChannelWithSatoshi) Size() uint32 {
	return 2 + elm.ChannelId.Size() +
		elm.ArbitrationLockBlock.Size() +
		elm.InterestAttribution.Size() +
		elm.LeftAddress.Size() +
		elm.LeftAmount.Size() +
		elm.LeftSatoshi.Size() +
		elm.RightAddress.Size() +
		elm.RightAmount.Size() +
		elm.RightSatoshi.Size()
}

// json api
func (elm *Action_31_OpenPaymentChannelWithSatoshi) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_31_OpenPaymentChannelWithSatoshi) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var idBytes, _ = elm.ChannelId.Serialize()
	var arbBytes, _ = elm.ArbitrationLockBlock.Serialize()
	var itrBytes, _ = elm.InterestAttribution.Serialize()
	var addr1Bytes, _ = elm.LeftAddress.Serialize()
	var amt1Bytes, _ = elm.LeftAmount.Serialize()
	var sat1Bytes, _ = elm.LeftSatoshi.Serialize()
	var addr2Bytes, _ = elm.RightAddress.Serialize()
	var amt2Bytes, _ = elm.RightAmount.Serialize()
	var sat2Bytes, _ = elm.RightSatoshi.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(idBytes)
	buffer.Write(arbBytes)
	buffer.Write(itrBytes)
	buffer.Write(addr1Bytes)
	buffer.Write(amt1Bytes)
	buffer.Write(sat1Bytes)
	buffer.Write(addr2Bytes)
	buffer.Write(amt2Bytes)
	buffer.Write(sat2Bytes)
	return buffer.Bytes(), nil
}

func (elm *Action_31_OpenPaymentChannelWithSatoshi) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ArbitrationLockBlock.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.InterestAttribution.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftSatoshi.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightSatoshi.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_31_OpenPaymentChannelWithSatoshi) RequestSignAddresses() []fields.Address {
	reqs := make([]fields.Address, 2)
	reqs[0] = elm.LeftAddress
	reqs[1] = elm.RightAddress
	return reqs
}

func (act *Action_31_OpenPaymentChannelWithSatoshi) WriteInChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	return openPaymentChannelWriteinChainState(state,
		act.ChannelId,
		act.ArbitrationLockBlock,
		act.InterestAttribution,
		act.LeftAddress,
		act.LeftAmount,
		act.LeftSatoshi,
		act.RightAddress,
		act.RightAmount,
		act.RightSatoshi,
	)
}

func (act *Action_31_OpenPaymentChannelWithSatoshi) WriteinChainState(state interfacev2.ChainStateOperation) error {

	panic("WriteinChainState in Action_31_OpenPaymentChannelWithSatoshi be deprecated")

	/*
		var e error
		// Query whether the channel exists
		sto, e := state.Channel(act.ChannelId)
		if e != nil {
			return e
		}
		// Channel IDS with the same left and right addresses and closed by consensus can be reused
		var reuseVersion fields.VarUint4 = 1
		var isIdCanUse = (sto == nil) ||
			(sto.IsAgreementClosed() && sto.LeftAddress.Equal(act.LeftAddress) && sto.RightAddress.Equal(act.RightAddress))
		if isIdCanUse == false {
			if sto.IsFinalDistributionClosed() {
				return fmt.Errorf("Payment Channel Id <%s> finally arbitration closed.", hex.EncodeToString(act.ChannelId))
			} else if sto.IsOpening() {
				return fmt.Errorf("Payment Channel Id <%s> is opening.", hex.EncodeToString(act.ChannelId))
			}
			return fmt.Errorf("Payment Channel Id <%s> already exist.", hex.EncodeToString(act.ChannelId))
		}
		if sto != nil {
			reuseVersion = sto.ReuseVersion + 1 // Reuse version number growth
		}
		// Channel ID validity
		if len(act.ChannelId) != stores.ChannelIdLength || act.ChannelId[0] == 0 || act.ChannelId[stores.ChannelIdLength-1] == 0 {
			return fmt.Errorf("Payment Channel Id <%s> format error.", hex.EncodeToString(act.ChannelId))
		}
		// Two addresses cannot be the same
		if act.LeftAddress.Equal(act.RightAddress) {
			return fmt.Errorf("Left address cannot equal with right address.")
		}
		// Check the number of digits stored in the amount
		labt, _ := act.LeftAmount.Serialize()
		rabt, _ := act.RightAmount.Serialize()
		if len(labt) > 6 || len(rabt) > 6 {
			// Avoid locking the storage digits of funds too long, resulting in the value storage digits after compound interest calculation exceeding the maximum range
			return fmt.Errorf("Payment Channel create error: left or right Amount bytes too long.")
		}
		// Cannot be negative, or both channels are zero at the same time (one can be positive and the other zero)
		if (!act.LeftAmount.IsPositive() || !act.RightAmount.IsPositive()) ||
			(act.LeftAmount.IsEmpty() && act.RightAmount.IsEmpty()) {
			return fmt.Errorf("Action_2_OpenPaymentChannel Payment Channel create error: left or right Amount is not positive.")
		}
		// Check whether the balance is sufficient
		bls1, e := state.Balance(act.LeftAddress)
		if e != nil {
			return e
		}
		if bls1 == nil {
			return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance cannot empty.", act.LeftAddress.ToReadable())
		}
		amt1 := bls1.Hacash
		if amt1.LessThan(&act.LeftAmount) {
			return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance is not enough. need %s but got %s", act.LeftAddress.ToReadable(), act.LeftAmount.ToFinString(), amt1.ToFinString())
		}
		bls2, e := state.Balance(act.RightAddress)
		if e != nil {
			return e
		}
		if bls2 == nil {
			return fmt.Errorf("Address %s Balance is not enough.", act.RightAddress.ToReadable())
		}
		amt2 := bls2.Hacash
		if amt2.LessThan(&act.RightAmount) {
			return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance is not enough. need %s but got %s", act.RightAddress.ToReadable(), act.RightAmount.ToFinString(), amt2.ToFinString())
		}
		curheight := state.GetPendingBlockHeight()
		// Create channel
		var storeItem = stores.CreateEmptyChannel()
		storeItem.BelongHeight = fields.BlockHeight(curheight)
		storeItem.ArbitrationLockBlock = fields.VarUint2(uint16(5000)) // The lock-in period unilaterally proposed is about 17 days
		storeItem.InterestAttribution = fields.VarUint1(0)             // By default, interest distribution is shared by the two parties according to the close amount
		storeItem.LeftAddress = act.LeftAddress
		storeItem.LeftAmount = act.LeftAmount
		storeItem.RightAddress = act.RightAddress
		storeItem.RightAmount = act.RightAmount
		storeItem.ReuseVersion = reuseVersion // Reuse version number
		storeItem.SetOpening()                // Open status
		// testing environment
		if sys.TestDebugLocalDevelopmentMark {
			storeItem.ArbitrationLockBlock = fields.VarUint2(uint16(20))
		}
		// Deduction balance
		e = DoSubBalanceFromChainState(state, act.LeftAddress, act.LeftAmount)
		if e != nil {
			return e
		}
		e = DoSubBalanceFromChainState(state, act.RightAddress, act.RightAmount)
		if e != nil {
			return e
		}
		// Storage channel
		e = state.ChannelCreate(act.ChannelId, storeItem)
		if e != nil {
			return e
		}
		// Total supply statistics
		totalsupply, e := state.ReadTotalSupply()
		if e != nil {
			return e
		}
		// Cumulative locked HAC
		addamt := act.LeftAmount.ToMei() + act.RightAmount.ToMei()
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfLocatedHACInChannel, addamt)
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfChannelOfOpening, 1)
		// update total supply
		e3 := state.UpdateSetTotalSupply(totalsupply)
		if e3 != nil {
			return e3
		}
		//
		return nil
	*/
}

func (act *Action_31_OpenPaymentChannelWithSatoshi) RecoverChainState(state interfacev2.ChainStateOperation) error {

	panic("RecoverChainState be deprecated")

	return nil
}

func (elm *Action_31_OpenPaymentChannelWithSatoshi) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_31_OpenPaymentChannelWithSatoshi) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees
func (act *Action_31_OpenPaymentChannelWithSatoshi) IsBurning90PersentTxFees() bool {
	return false
}

//////////////////////////////////////////////////////////

/**
 * open payment channel
 */
func openPaymentChannelWriteinChainState(state interfaces.ChainStateOperation,
	channelId fields.ChannelId,
	arbitrationLockBlock fields.VarUint2,
	interestAttribution fields.VarUint1,
	leftAddress fields.Address,
	leftAmount fields.Amount,
	leftSatoshi fields.SatoshiVariation,
	rightAddress fields.Address,
	rightAmount fields.Amount,
	rightSatoshi fields.SatoshiVariation,
) error {

	var realLeftSatoshi = leftSatoshi.GetRealSatoshi()
	var realRightSatoshi = rightSatoshi.GetRealSatoshi()

	var e error
	// Query whether the channel exists
	sto, e := state.Channel(channelId)
	if e != nil {
		return e
	}
	// Channel IDS with the same left and right addresses and closed by consensus can be reused
	var reuseVersion fields.VarUint4 = 1
	var isIdCanUse = (sto == nil) ||
		(sto.IsAgreementClosed() && sto.LeftAddress.Equal(leftAddress) && sto.RightAddress.Equal(rightAddress))
	if isIdCanUse == false {
		if sto.IsFinalDistributionClosed() {
			return fmt.Errorf("Payment Channel Id <%s> finally arbitration closed.", hex.EncodeToString(channelId))
		} else if sto.IsOpening() {
			return fmt.Errorf("Payment Channel Id <%s> is opening.", hex.EncodeToString(channelId))
		}
		return fmt.Errorf("Payment Channel Id <%s> already exist.", hex.EncodeToString(channelId))
	}
	if sto != nil {
		reuseVersion = sto.ReuseVersion + 1 // Reuse version number growth
	}
	// Channel ID validity
	if len(channelId) != stores.ChannelIdLength || channelId[0] == 0 || channelId[stores.ChannelIdLength-1] == 0 {
		return fmt.Errorf("Payment Channel Id <%s> format error.", hex.EncodeToString(channelId))
	}
	// Two addresses cannot be the same
	if leftAddress.Equal(rightAddress) {
		return fmt.Errorf("Left address cannot equal with right address.")
	}
	// Cannot be negative, or both channels are zero at the same time (one can be positive and the other zero)
	if leftAmount.IsNegative() || leftAmount.IsNegative() ||
		(leftAmount.IsEmpty() && rightAmount.IsEmpty() && realLeftSatoshi == 0 && realRightSatoshi == 0) {
		return fmt.Errorf("OpenPaymentChannel Payment Channel create error: left or right Amount is not positive.")
	}
	// Check the number of digits stored in the amount
	labt, _ := leftAmount.Serialize()
	rabt, _ := rightAmount.Serialize()
	if len(labt) > 6 || len(rabt) > 6 {
		// Avoid locking the storage digits of funds too long, resulting in the value storage digits after compound interest calculation exceeding the maximum range
		return fmt.Errorf("Payment Channel create error: left or right Amount bytes too long.")
	}

	// Check whether the balance is sufficient
	bls1, e := state.Balance(leftAddress)
	if e != nil {
		return e
	}
	if bls1 == nil {
		return fmt.Errorf("OpenPaymentChannel Address %s Balance cannot empty.", leftAddress.ToReadable())
	}
	amt1 := bls1.Hacash
	if amt1.LessThan(&leftAmount) {
		return fmt.Errorf("OpenPaymentChannel Address %s Balance is not enough. need %s but got %s", leftAddress.ToReadable(), leftAmount.ToFinString(), amt1.ToFinString())
	}
	sat1 := bls1.Satoshi
	if sat1 < realLeftSatoshi {
		return fmt.Errorf("OpenPaymentChannel Address %s Satoshi is not enough. need %d but got %d", leftAddress.ToReadable(), leftSatoshi, sat1)
	}

	bls2, e := state.Balance(rightAddress)
	if e != nil {
		return e
	}
	if bls2 == nil {
		return fmt.Errorf("Address %s Balance is not enough.", rightAddress.ToReadable())
	}
	amt2 := bls2.Hacash
	if amt2.LessThan(&rightAmount) {
		return fmt.Errorf("OpenPaymentChannel Address %s Balance is not enough. need %s but got %s", rightAddress.ToReadable(), rightAmount.ToFinString(), amt2.ToFinString())
	}
	sat2 := bls2.Satoshi
	if sat2 < realRightSatoshi {
		return fmt.Errorf("OpenPaymentChannel Address %s Satoshi is not enough. need %d but got %d", rightAddress.ToReadable(), rightSatoshi, sat2)
	}

	curheight := state.GetPendingBlockHeight()
	// Create channel
	var storeItem = stores.CreateEmptyChannel()
	storeItem.BelongHeight = fields.BlockHeight(curheight)
	//storeItem.ArbitrationLockBlock = fields.VarUint2(uint16(5000)) // The lock-in period unilaterally proposed is about 17 days
	//storeItem.InterestAttribution = fields.VarUint1(0)             // By default, interest distribution is shared by the two parties according to the close amount
	storeItem.ArbitrationLockBlock = arbitrationLockBlock
	storeItem.InterestAttribution = interestAttribution
	storeItem.LeftAddress = leftAddress
	storeItem.LeftAmount = leftAmount
	storeItem.LeftSatoshi = leftSatoshi
	storeItem.RightAddress = rightAddress
	storeItem.RightAmount = rightAmount
	storeItem.RightSatoshi = rightSatoshi
	storeItem.ReuseVersion = reuseVersion // Reuse version number
	storeItem.SetOpening()                // Open status
	// testing environment
	if sys.TestDebugLocalDevelopmentMark {
		storeItem.ArbitrationLockBlock = fields.VarUint2(uint16(20))
	}
	// Deduction balance
	if leftAmount.IsPositive() {
		e = DoSubBalanceFromChainStateV3(state, leftAddress, leftAmount)
		if e != nil {
			return e
		}

	}
	if rightAmount.IsPositive() {
		e = DoSubBalanceFromChainStateV3(state, rightAddress, rightAmount)
		if e != nil {
			return e
		}
	}
	if realLeftSatoshi > 0 {
		e = DoSubSatoshiFromChainStateV3(state, leftAddress, realLeftSatoshi)
		if e != nil {
			return e
		}
	}
	if realRightSatoshi > 0 {
		e = DoSubSatoshiFromChainStateV3(state, rightAddress, realRightSatoshi)
		if e != nil {
			return e
		}
	}
	// Storage channel
	e = state.ChannelCreate(channelId, storeItem)
	if e != nil {
		return e
	}
	// Total supply statistics
	totalsupply, e := state.ReadTotalSupply()
	if e != nil {
		return e
	}
	// Cumulative locked HAC
	addamt := leftAmount.ToMei() + rightAmount.ToMei()
	addsat := realRightSatoshi + realRightSatoshi
	if addamt > 0 {
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfLocatedHACInChannel, addamt)
	}
	if addsat > 0 {
		totalsupply.DoAddUint(stores.TotalSupplyStoreTypeOfLocatedSATInChannel, uint64(addsat))
	}
	totalsupply.DoAddUint(stores.TotalSupplyStoreTypeOfChannelOfOpening, 1)
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	//
	return nil
}

// Close channel status write
// isFinalClosed: 是否为仲裁终局结束，不可重用
func closePaymentChannelWriteinChainState(state interfacev2.ChainStateOperation, channelId []byte, paychan *stores.Channel, newLeftAmt *fields.Amount, newRightAmt *fields.Amount, leftNewSAT fields.Satoshi, rightNewSAT fields.Satoshi, isFinalClosed bool) error {
	var e error
	// Judge that the channel has been closed
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(channelId))
	}
	if paychan.IsClosed() {
		return fmt.Errorf("Payment Channel <%s> is be closed.", hex.EncodeToString(channelId))
	}
	// Calculate interest by time
	if newLeftAmt == nil || newRightAmt == nil {
		// Automatically use the deposited amount to calculate interest
		newLeftAmt = &paychan.LeftAmount
		newRightAmt = &paychan.RightAmount
	}
	// Calculate total
	// Allocation amount can be zero but not negative
	if newLeftAmt.IsNegative() || newRightAmt.IsNegative() {
		return fmt.Errorf("Payment channel distribution amount cannot be negative.")
	}
	// Check whether the allocated amount is equal to the deposited amount
	tt1, e1 := newLeftAmt.Add(newRightAmt)
	if e1 != nil {
		return e1
	}
	tt2, e2 := paychan.LeftAmount.Add(&paychan.RightAmount)
	if e2 != nil {
		return e2
	}
	if tt1.NotEqual(tt2) {
		// Unequal
		return fmt.Errorf("HAC distribution amount must equal with lock in.")
	}
	// Calculate the current block height
	//var curheight uint64 = 1
	curheight := state.GetPendingBlockHeight()
	leftAmount, rightAmount, haveinterest, e11 := calculateChannelInterest(
		curheight, uint64(paychan.BelongHeight), newLeftAmt, newRightAmt, paychan.InterestAttribution)
	if e11 != nil {
		return e11
	}
	// Increase the balance (withdraw the locked amount and interest from the channel)
	// HAC
	e = DoAddBalanceFromChainState(state, paychan.LeftAddress, *leftAmount)
	if e != nil {
		return e
	}
	e = DoAddBalanceFromChainState(state, paychan.RightAddress, *rightAmount)
	if e != nil {
		return e
	}
	// Sat raised from channel
	leftOldSAT := paychan.LeftSatoshi.GetRealSatoshi()
	rightOldSAT := paychan.RightSatoshi.GetRealSatoshi()
	totalOldSAT := leftOldSAT + rightOldSAT
	totalNewSAT := leftNewSAT + rightNewSAT
	// Check whether the total quantity matches
	if totalOldSAT != totalNewSAT {
		return fmt.Errorf("SAT distribution error: need total %d SAT but got %d (left: %d, right: %d).",
			totalOldSAT, totalNewSAT, leftNewSAT, rightNewSAT)
	}
	if leftNewSAT > 0 {
		e = DoAddSatoshiFromChainState(state, paychan.LeftAddress, leftNewSAT)
		if e != nil {
			return e
		}
	}
	if rightNewSAT > 0 {
		e = DoAddSatoshiFromChainState(state, paychan.RightAddress, rightNewSAT)
		if e != nil {
			return e
		}
	}
	// Temporarily reserve channels for data fallback
	// Calculate left final allocation
	if isFinalClosed {
		paychan.SetFinalArbitrationClosed(newLeftAmt, leftNewSAT) // Arbitration permanently closed
	} else {
		paychan.SetAgreementClosed(newLeftAmt, leftNewSAT) // Negotiation closed
	}
	e = state.ChannelUpdate(channelId, paychan)
	if e != nil {
		return e
	}
	//
	// Total supply statistics
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	// Reduce unlocked HAC
	lockamt := paychan.LeftAmount.ToMei() + paychan.RightAmount.ToMei()
	totalsupply.DoSub(stores.TotalSupplyStoreTypeOfLocatedHACInChannel, lockamt) // Reduce locked HAC statistics
	if totalNewSAT > 0 {
		totalsupply.DoSubUint(stores.TotalSupplyStoreTypeOfLocatedSATInChannel, uint64(totalNewSAT)) // Reduce locked sat statistics
	}
	totalsupply.DoSubUint(stores.TotalSupplyStoreTypeOfChannelOfOpening, 1) // Reduce channel count
	// Add channel interest statistics
	if haveinterest {
		releaseamt := leftAmount.ToMei() + rightAmount.ToMei()
		//fmt.Println("(act *Action_3_ClosePaymentChannel) WriteinChainState", releaseamt, lockamt, releaseamt - lockamt, )
		//fmt.Println(paychanptr.LeftAddress.ToReadable(), paychanptr.LeftAmount.ToFinString(), paychanptr.LeftAmount.ToMei())
		//fmt.Println(paychanptr.RightAddress.ToReadable(), paychanptr.RightAmount.ToFinString(), paychanptr.RightAmount.ToMei())
		//fmt.Println(leftAmount.ToFinString(), leftAmount.ToMei(), rightAmount.ToFinString(), rightAmount.ToMei())
		if releaseamt-lockamt < 0 {
			return fmt.Errorf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		}
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfChannelInterest, releaseamt-lockamt)
	}
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	return nil

}

// Close channel status fallback
func closePaymentChannelRecoverChainState_deprecated(state interfacev2.ChainStateOperation, channelId []byte, newLeftAmt *fields.Amount, newRightAmt *fields.Amount, backToChallenging bool) error {

	panic("RecoverChainState be deprecated")

	var e error = nil
	// Query channel
	paychan, e := state.Channel(channelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		// The channel must be saved before it can be rolled back
		panic(fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(channelId)))
	}
	// Judge that the channel must be closed
	if paychan.IsClosed() {
		panic(fmt.Errorf("Payment Channel <%s> is be closed.", hex.EncodeToString(channelId)))
	}
	if newLeftAmt == nil || newRightAmt == nil {
		// Automatically use the deposited amount to calculate interest
		newLeftAmt = &paychan.LeftAmount
		newRightAmt = &paychan.RightAmount
	}
	// Calculate difference
	curheight := state.GetPendingBlockHeight()
	// Calculate interest
	leftAmount, rightAmount, haveinterest, e11 := calculateChannelInterest(curheight, uint64(paychan.BelongHeight), newLeftAmt, newRightAmt, 0)
	if e11 != nil {
		return e11
	}
	// Deduct the balance (put the amount back into the channel)
	e = DoSubBalanceFromChainState(state, paychan.LeftAddress, *leftAmount)
	if e != nil {
		return e
	}
	e = DoSubBalanceFromChainState(state, paychan.RightAddress, *rightAmount)
	if e != nil {
		return e
	}
	// Restore channel status
	if backToChallenging {
		paychan.Status = stores.ChannelStatusChallenging // Back to challenge status
	} else {
		paychan.SetOpening() // Relabel the channel as on
	}
	e = state.ChannelUpdate(channelId, paychan)
	if e != nil {
		return e
	}
	// Total supply statistics
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	// Rollback unlocked HAC
	lockamt := paychan.LeftAmount.ToMei() + paychan.RightAmount.ToMei()
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfLocatedHACInChannel, lockamt)
	// Interest statistics of fallback channel
	if haveinterest {
		releaseamt := leftAmount.ToMei() + rightAmount.ToMei()
		totalsupply.DoSub(stores.TotalSupplyStoreTypeOfChannelInterest, releaseamt-lockamt)
	}
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	return nil
}

// Calculate channel interest
// Whether bool has interest
// Interestgiveto whom interest is allocated
func calculateChannelInterest(curheight uint64, openBelongHeight uint64, leftAmount *fields.Amount, rightAmount *fields.Amount, interestgiveto fields.VarUint1) (*fields.Amount, *fields.Amount, bool, error) {
	// Increase interest calculation, compound interest times: about 2500 blocks will increase compound interest by one ten thousandth every 8.68 days, less than 8 days will be ignored, and the annual compound interest is about 0.42%
	//a1, a2 := DoAppendCompoundInterest1Of10000By2500Height(&leftAmount, &rightAmount, insnum)
	var insnum = (curheight - openBelongHeight) / 2500
	var wfzn uint64 = 1 // 万分之一 1/10000
	// Modify the proportion of one-time additional issuance by opening the block height of the channel
	if openBelongHeight > 200000 {
		// Increase interest calculation, compounding times: about 10000 blocks will be compounded once every 34 days, less than 34 days will be ignored, and the annual compound interest is about 1.06%
		insnum = (curheight - openBelongHeight) / 10000
		wfzn = 10 // 千分之一 10/10000
	}
	if insnum > 0 {
		// Calculate channel interest reward
		a1, a2, e := coinbase.DoAppendCompoundInterestProportionOfHeightV2(leftAmount, rightAmount, insnum, wfzn, interestgiveto)
		if e != nil {
			return nil, nil, false, e
		}
		// Plus interest
		return a1, a2, true, nil
	}
	// No interest
	return leftAmount, rightAmount, false, nil
}

//////////////////////////////////////////////////////////

// Close channel status write
// isFinalClosed: 是否为仲裁终局结束，不可重用
func closePaymentChannelWriteinChainStateV3(state interfaces.ChainStateOperation, channelId []byte, paychan *stores.Channel, newLeftAmt *fields.Amount, newRightAmt *fields.Amount, leftNewSAT fields.Satoshi, rightNewSAT fields.Satoshi, isFinalClosed bool) error {
	var e error
	// Judge that the channel has been closed
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(channelId))
	}
	if paychan.IsClosed() {
		return fmt.Errorf("Payment Channel <%s> is be closed.", hex.EncodeToString(channelId))
	}
	// Calculate interest by time
	if newLeftAmt == nil || newRightAmt == nil {
		// Automatically use the deposited amount to calculate interest
		newLeftAmt = &paychan.LeftAmount
		newRightAmt = &paychan.RightAmount
	}
	// Calculate total
	// Allocation amount can be zero but not negative
	if newLeftAmt.IsNegative() || newRightAmt.IsNegative() {
		return fmt.Errorf("Payment channel distribution amount cannot be negative.")
	}
	// Check whether the allocated amount is equal to the deposited amount
	tt1, e1 := newLeftAmt.Add(newRightAmt)
	if e1 != nil {
		return e1
	}
	tt2, e2 := paychan.LeftAmount.Add(&paychan.RightAmount)
	if e2 != nil {
		return e2
	}
	if tt1.NotEqual(tt2) {
		// Unequal
		return fmt.Errorf("HAC distribution amount must equal with lock in.")
	}
	// Calculate the current block height
	//var curheight uint64 = 1

	curheight := state.GetPendingBlockHeight()
	leftAmount, rightAmount, haveinterest, e11 := calculateChannelInterest(
		curheight, uint64(paychan.BelongHeight), newLeftAmt, newRightAmt, paychan.InterestAttribution)
	if e11 != nil {
		return e11
	}
	// Increase the balance (withdraw the locked amount and interest from the channel)
	// HAC
	e = DoAddBalanceFromChainStateV3(state, paychan.LeftAddress, *leftAmount)
	if e != nil {
		return e
	}
	e = DoAddBalanceFromChainStateV3(state, paychan.RightAddress, *rightAmount)
	if e != nil {
		return e
	}
	// Sat raised from channel
	leftOldSAT := paychan.LeftSatoshi.GetRealSatoshi()
	rightOldSAT := paychan.RightSatoshi.GetRealSatoshi()
	totalOldSAT := leftOldSAT + rightOldSAT
	totalNewSAT := leftNewSAT + rightNewSAT
	// Check whether the total quantity matches
	if totalOldSAT != totalNewSAT {
		return fmt.Errorf("SAT distribution error: need total %d SAT but got %d (left: %d, right: %d).",
			totalOldSAT, totalNewSAT, leftNewSAT, rightNewSAT)
	}
	if leftNewSAT > 0 {
		e = DoAddSatoshiFromChainStateV3(state, paychan.LeftAddress, leftNewSAT)
		if e != nil {
			return e
		}
	}
	if rightNewSAT > 0 {
		e = DoAddSatoshiFromChainStateV3(state, paychan.RightAddress, rightNewSAT)
		if e != nil {
			return e
		}
	}
	// Temporarily reserve channels for data fallback
	// Calculate left final allocation
	if isFinalClosed {
		paychan.SetFinalArbitrationClosed(newLeftAmt, leftNewSAT) // Arbitration permanently closed
	} else {
		paychan.SetAgreementClosed(newLeftAmt, leftNewSAT) // Negotiation closed
	}
	e = state.ChannelUpdate(channelId, paychan)
	if e != nil {
		return e
	}
	//
	// Total supply statistics
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	// Reduce unlocked HAC
	lockamt := paychan.LeftAmount.ToMei() + paychan.RightAmount.ToMei()
	totalsupply.DoSub(stores.TotalSupplyStoreTypeOfLocatedHACInChannel, lockamt) // Reduce locked HAC statistics
	if totalNewSAT > 0 {
		totalsupply.DoSubUint(stores.TotalSupplyStoreTypeOfLocatedSATInChannel, uint64(totalNewSAT)) // Reduce locked sat statistics
	}
	totalsupply.DoSubUint(stores.TotalSupplyStoreTypeOfChannelOfOpening, 1) // Reduce channel count
	// Add channel interest statistics
	if haveinterest {
		releaseamt := leftAmount.ToMei() + rightAmount.ToMei()
		//fmt.Println("(act *Action_3_ClosePaymentChannel) WriteinChainState", releaseamt, lockamt, releaseamt - lockamt, )
		//fmt.Println(paychanptr.LeftAddress.ToReadable(), paychanptr.LeftAmount.ToFinString(), paychanptr.LeftAmount.ToMei())
		//fmt.Println(paychanptr.RightAddress.ToReadable(), paychanptr.RightAmount.ToFinString(), paychanptr.RightAmount.ToMei())
		//fmt.Println(leftAmount.ToFinString(), leftAmount.ToMei(), rightAmount.ToFinString(), rightAmount.ToMei())
		if releaseamt-lockamt < 0 {
			return fmt.Errorf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		}
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfChannelInterest, releaseamt-lockamt)
	}
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	return nil

}
