package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/channel"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
)

// Without any documents, the channel is closed unilaterally, entering the challenge period
// Fund allocation is calculated based on initial deposit
type Action_22_UnilateralClosePaymentChannelByNothing struct {
	// Channel ID
	ChannelId          fields.ChannelId // Channel ID
	AssertCloseAddress fields.Address   // Proposed address of unilateral claim for closure

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) Kind() uint16 {
	return 22
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) Size() uint32 {
	return 2 + elm.ChannelId.Size() + elm.AssertCloseAddress.Size()
}

// json api
func (elm *Action_22_UnilateralClosePaymentChannelByNothing) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.ChannelId.Serialize()
	var bt2, _ = elm.AssertCloseAddress.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	return buffer.Bytes(), nil
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.AssertCloseAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) RequestSignAddresses() []fields.Address {
	// Proposer must sign
	return []fields.Address{
		elm.AssertCloseAddress,
	}
}

func (act *Action_22_UnilateralClosePaymentChannelByNothing) WriteInChainState(state interfaces.ChainStateOperation) error {
	var e error

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
	// Check status (must be on)
	if paychan.IsOpening() == false {
		return fmt.Errorf("Payment Channel status is not on opening.")
	}
	// Check the two account addresses to see if they match
	addrIsLeft := paychan.LeftAddress.Equal(act.AssertCloseAddress)
	addrIsRight := paychan.RightAddress.Equal(act.AssertCloseAddress)
	if !addrIsLeft && !addrIsRight {
		return fmt.Errorf("Payment Channel <%s> address signature verify fail.", hex.EncodeToString(act.ChannelId))
	}
	// Challenger status
	clghei := state.GetPendingBlockHeight()
	var clgamt = fields.Amount{}
	var clgsat = fields.Satoshi(0)
	if addrIsLeft {
		clgamt = paychan.LeftAmount
		clgsat = paychan.LeftSatoshi.GetRealSatoshi()
	} else {
		clgamt = paychan.RightAmount
		clgsat = paychan.RightSatoshi.GetRealSatoshi()
	}
	// Update to challenge period, no bill number
	paychan.SetChallenging(clghei, addrIsLeft, &clgamt, clgsat, 0)
	// Write status
	e = state.ChannelUpdate(act.ChannelId, paychan)
	if e != nil {
		return e
	}
	return nil
}

func (act *Action_22_UnilateralClosePaymentChannelByNothing) WriteinChainState(state interfacev2.ChainStateOperation) error {
	var e error

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
	// Check status (must be on)
	if paychan.IsOpening() == false {
		return fmt.Errorf("Payment Channel status is not on opening.")
	}
	// Check the two account addresses to see if they match
	addrIsLeft := paychan.LeftAddress.Equal(act.AssertCloseAddress)
	addrIsRight := paychan.RightAddress.Equal(act.AssertCloseAddress)
	if !addrIsLeft && !addrIsRight {
		return fmt.Errorf("Payment Channel <%s> address signature verify fail.", hex.EncodeToString(act.ChannelId))
	}
	// Challenger status
	clghei := state.GetPendingBlockHeight()
	var clgamt = fields.Amount{}
	var clgsat = fields.Satoshi(0)
	if addrIsLeft {
		clgamt = paychan.LeftAmount
		clgsat = paychan.LeftSatoshi.GetRealSatoshi()
	} else {
		clgamt = paychan.RightAmount
		clgsat = paychan.RightSatoshi.GetRealSatoshi()
	}
	// Update to challenge period, no bill number
	paychan.SetChallenging(clghei, addrIsLeft, &clgamt, clgsat, 0)
	// Write status
	e = state.ChannelUpdate(act.ChannelId, paychan)
	if e != nil {
		return e
	}
	return nil
}

func (act *Action_22_UnilateralClosePaymentChannelByNothing) RecoverChainState(state interfacev2.ChainStateOperation) error {

	// Query channel
	paychan, e := state.Channel(act.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// Fallback status
	paychan.SetOpening()
	state.ChannelUpdate(act.ChannelId, paychan)
	return nil
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_22_UnilateralClosePaymentChannelByNothing) IsBurning90PersentTxFees() bool {
	return false
}

/////////////////////////////////////////////////////////////

// 1. close the channel through the intermediate real-time statement and enter the challenge period
// 2. 提供实时对账单，回应挑战，夺取对方全部金额
type Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation struct {
	// Proposer address
	AssertAddress fields.Address
	// Account Statement
	Reconciliation channel.OnChainArbitrationBasisReconciliation

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) Kind() uint16 {
	return 23
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) Size() uint32 {
	return 2 + elm.AssertAddress.Size() + elm.Reconciliation.Size()
}

// json api
func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.AssertAddress.Serialize()
	var bt2, _ = elm.Reconciliation.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	return buffer.Bytes(), nil
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.AssertAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.Reconciliation.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) RequestSignAddresses() []fields.Address {
	// Check signature
	return []fields.Address{
		elm.AssertAddress,
	}
}

func (act *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) WriteInChainState(state interfaces.ChainStateOperation) error {

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	// cid
	channelId := act.Reconciliation.GetChannelId()

	// Query channel
	paychan, e := state.Channel(channelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel <%s> not find.", hex.EncodeToString(channelId))
	}
	// Check the signatures of both account addresses, and both parties check
	// Enter a challenging period or seize funds
	return checkChannelGotoChallegingOrFinalDistributionWriteinChainStateV3(state, act.AssertAddress, paychan, &act.Reconciliation)
}

func (act *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) WriteinChainState(state interfacev2.ChainStateOperation) error {

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// cid
	channelId := act.Reconciliation.GetChannelId()

	// Query channel
	paychan, e := state.Channel(channelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel <%s> not find.", hex.EncodeToString(channelId))
	}
	// Check the signatures of both account addresses, and both parties check
	// Enter a challenging period or seize funds
	return checkChannelGotoChallegingOrFinalDistributionWriteinChainState(state, act.AssertAddress, paychan, &act.Reconciliation)
}

func (act *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) RecoverChainState(state interfacev2.ChainStateOperation) error {

	channelId := act.Reconciliation.GetChannelId()

	// Query channel
	paychan, e := state.Channel(channelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(channelId))
	}

	// Fallback
	return checkChannelGotoChallegingOrFinalDistributionRecoverChainState(state, act.AssertAddress, paychan, &act.Reconciliation)
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) IsBurning90PersentTxFees() bool {
	return false
}

///////////////////////////////////////////////

// Unilateral termination
// 1. unilaterally close the channel through channel chain payment and enter the challenge period
// 2. 提供通道链支付对账单，回应挑战，夺取对方全部金额
type Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody struct {
	// Proposer address
	AssertAddress fields.Address
	// Channel overall payment data
	ChannelChainTransferData channel.OffChainFormPaymentChannelTransfer
	// Payment entity data of this channel
	ChannelChainTransferTargetProveBody channel.ChannelChainTransferProveBodyInfo

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) Kind() uint16 {
	return 24
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) Size() uint32 {
	return 2 + elm.AssertAddress.Size() +
		elm.ChannelChainTransferData.Size() +
		elm.ChannelChainTransferTargetProveBody.Size()
}

// json api
func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.AssertAddress.Serialize()
	var bt2, _ = elm.ChannelChainTransferData.Serialize()
	var bt3, _ = elm.ChannelChainTransferTargetProveBody.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	buffer.Write(bt3)
	return buffer.Bytes(), nil
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.AssertAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ChannelChainTransferData.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ChannelChainTransferTargetProveBody.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) RequestSignAddresses() []fields.Address {
	// Check signature
	return []fields.Address{
		elm.AssertAddress,
	}
}

func (act *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) WriteInChainState(state interfaces.ChainStateOperation) error {

	var e error

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Query channel
	paychan, e := state.Channel(act.ChannelChainTransferTargetProveBody.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel <%s> not find.", hex.EncodeToString(act.ChannelChainTransferTargetProveBody.ChannelId))
	}

	// Check whether the channel hash is correct
	hxhalf := act.ChannelChainTransferTargetProveBody.GetSignStuffHashHalfChecker()
	// Check whether the hash value is included in the list
	var isHashCheckOk = false
	for _, hxckr := range act.ChannelChainTransferData.ChannelTransferProveHashHalfCheckers {
		if hxhalf.Equal(hxckr) {
			isHashCheckOk = true
			break
		}
	}
	if !isHashCheckOk {
		return fmt.Errorf("ChannelChainTransferTargetProveBody hash <%s> not find.", hxhalf.ToHex())
	}

	// Check whether the channel addresses of both parties are included in the signature list
	lsgok := false
	rsgok := false
	for _, v := range act.ChannelChainTransferData.MustSignAddresses {
		if v.Equal(paychan.LeftAddress) {
			lsgok = true
		} else if v.Equal(paychan.RightAddress) {
			rsgok = true
		}
	}
	if !lsgok || !rsgok {
		return fmt.Errorf("Channel signature address is missing.")
	}

	// Check that all signatures are complete and correct
	e = act.ChannelChainTransferData.CheckMustAddressAndSigns()
	if e != nil {
		return e
	}

	// Check whether it is a challenge or a final capture
	return checkChannelGotoChallegingOrFinalDistributionWriteinChainStateV3(state, act.AssertAddress, paychan, &act.ChannelChainTransferTargetProveBody)
}

func (act *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) WriteinChainState(state interfacev2.ChainStateOperation) error {

	var e error

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Query channel
	paychan, e := state.Channel(act.ChannelChainTransferTargetProveBody.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel <%s> not find.", hex.EncodeToString(act.ChannelChainTransferTargetProveBody.ChannelId))
	}

	// Check whether the channel hash is correct
	hxhalf := act.ChannelChainTransferTargetProveBody.GetSignStuffHashHalfChecker()
	// Check whether the hash value is included in the list
	var isHashCheckOk = false
	for _, hxckr := range act.ChannelChainTransferData.ChannelTransferProveHashHalfCheckers {
		if hxhalf.Equal(hxckr) {
			isHashCheckOk = true
			break
		}
	}
	if !isHashCheckOk {
		return fmt.Errorf("ChannelChainTransferTargetProveBody hash <%s> not find.", hxhalf.ToHex())
	}

	// Check whether the channel addresses of both parties are included in the signature list
	lsgok := false
	rsgok := false
	for _, v := range act.ChannelChainTransferData.MustSignAddresses {
		if v.Equal(paychan.LeftAddress) {
			lsgok = true
		} else if v.Equal(paychan.RightAddress) {
			rsgok = true
		}
	}
	if !lsgok || !rsgok {
		return fmt.Errorf("Channel signature address is missing.")
	}

	// Check that all signatures are complete and correct
	e = act.ChannelChainTransferData.CheckMustAddressAndSigns()
	if e != nil {
		return e
	}

	// Check whether it is a challenge or a final capture
	return checkChannelGotoChallegingOrFinalDistributionWriteinChainState(state, act.AssertAddress, paychan, &act.ChannelChainTransferTargetProveBody)
}

func (act *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) RecoverChainState(state interfacev2.ChainStateOperation) error {

	// Query channel
	paychan, e := state.Channel(act.ChannelChainTransferTargetProveBody.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelChainTransferTargetProveBody.ChannelId))
	}

	// Fallback
	return checkChannelGotoChallegingOrFinalDistributionRecoverChainState(state, act.AssertAddress, paychan, &act.ChannelChainTransferTargetProveBody)
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) IsBurning90PersentTxFees() bool {
	return false
}

///////////////////////////////////////////////

// Unilateral termination
// 1. unilaterally close the channel through channel chain atom exchange and enter the challenge period
// 2. 提供通道链支付对账单，回应挑战，夺取对方全部金额
type Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange struct {
	// Proposer address
	AssertAddress fields.Address
	// credential
	ProveBodyHashChecker fields.HashHalfChecker
	// Reconciliation data
	ChannelChainTransferTargetProveBody channel.ChannelChainTransferProveBodyInfo

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) Kind() uint16 {
	return 26
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) Size() uint32 {
	return 2 + elm.AssertAddress.Size() +
		elm.ProveBodyHashChecker.Size() +
		elm.ChannelChainTransferTargetProveBody.Size()
}

// json api
func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.AssertAddress.Serialize()
	var bt2, _ = elm.ProveBodyHashChecker.Serialize()
	var bt3, _ = elm.ChannelChainTransferTargetProveBody.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	buffer.Write(bt3)
	return buffer.Bytes(), nil
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.AssertAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ProveBodyHashChecker.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ChannelChainTransferTargetProveBody.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) RequestSignAddresses() []fields.Address {
	// Check signature
	return []fields.Address{
		elm.AssertAddress,
	}
}

func (act *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) WriteInChainState(state interfaces.ChainStateOperation) error {

	var e error

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	// The fast channel mode cannot be used to initiate challenges and arbitrations. Only the normal mode can be used

	// Query channel
	paychan, e := state.Channel(act.ChannelChainTransferTargetProveBody.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel <%s> not find.", hex.EncodeToString(act.ChannelChainTransferTargetProveBody.ChannelId))
	}

	// Query swap transaction
	swapex, e := state.Chaswap(act.ProveBodyHashChecker)
	if e != nil {
		return e
	}
	if swapex == nil {
		return fmt.Errorf("Chaswap tranfer <%s> not find.", act.ProveBodyHashChecker.ToHex())
	}
	// Has it been used
	if swapex.IsBeUsed.Check() {
		return fmt.Errorf("Chaswap tranfer <%s> already be used.", act.ProveBodyHashChecker.ToHex())
	}

	// Check that the address that must be signed is complete and correct
	addrsmap := make(map[string]bool)
	for _, addr := range swapex.OnchainTransferFromAndMustSignAddresses {
		addrsmap[string(addr)] = true
	}
	_, hasleft := addrsmap[string(paychan.LeftAddress)]
	_, hasright := addrsmap[string(paychan.RightAddress)]
	if !hasleft || !hasright {
		return fmt.Errorf("Chaswap tranfer signature error.")
	}

	// Mark ticket used
	swapex.IsBeUsed.Set(true)
	e = state.ChaswapUpdate(act.ProveBodyHashChecker, swapex)
	if e != nil {
		return e
	}

	// Check whether it is a challenge or a final capture
	return checkChannelGotoChallegingOrFinalDistributionWriteinChainStateV3(state, act.AssertAddress, paychan, &act.ChannelChainTransferTargetProveBody)
}

func (act *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) WriteinChainState(state interfacev2.ChainStateOperation) error {

	var e error

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// The fast channel mode cannot be used to initiate challenges and arbitrations. Only the normal mode can be used

	// Query channel
	paychan, e := state.Channel(act.ChannelChainTransferTargetProveBody.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel <%s> not find.", hex.EncodeToString(act.ChannelChainTransferTargetProveBody.ChannelId))
	}

	// Query swap transaction
	swapex, e := state.Chaswap(act.ProveBodyHashChecker)
	if e != nil {
		return e
	}
	if swapex == nil {
		return fmt.Errorf("Chaswap tranfer <%s> not find.", act.ProveBodyHashChecker.ToHex())
	}
	// Has it been used
	if swapex.IsBeUsed.Check() {
		return fmt.Errorf("Chaswap tranfer <%s> already be used.", act.ProveBodyHashChecker.ToHex())
	}

	// Check that the address that must be signed is complete and correct
	addrsmap := make(map[string]bool)
	for _, addr := range swapex.OnchainTransferFromAndMustSignAddresses {
		addrsmap[string(addr)] = true
	}
	_, hasleft := addrsmap[string(paychan.LeftAddress)]
	_, hasright := addrsmap[string(paychan.RightAddress)]
	if !hasleft || !hasright {
		return fmt.Errorf("Chaswap tranfer signature error.")
	}

	// Mark ticket used
	swapex.IsBeUsed.Set(true)
	e = state.ChaswapUpdate(act.ProveBodyHashChecker, swapex)
	if e != nil {
		return e
	}

	// Check whether it is a challenge or a final capture
	return checkChannelGotoChallegingOrFinalDistributionWriteinChainState(state, act.AssertAddress, paychan, &act.ChannelChainTransferTargetProveBody)
}

func (act *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) RecoverChainState(state interfacev2.ChainStateOperation) error {

	// Query channel
	paychan, e := state.Channel(act.ChannelChainTransferTargetProveBody.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelChainTransferTargetProveBody.ChannelId))
	}

	// Fallback usage status
	swapex, e := state.Chaswap(act.ProveBodyHashChecker)
	if e != nil {
		return e
	}
	swapex.IsBeUsed.Set(false)
	state.ChaswapCreate(act.ProveBodyHashChecker, swapex)

	// Fallback status
	return checkChannelGotoChallegingOrFinalDistributionRecoverChainState(state, act.AssertAddress, paychan, &act.ChannelChainTransferTargetProveBody)
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) IsBurning90PersentTxFees() bool {
	return false
}

//////////////////////////////////////////////////////////////

// The inspection channel enters the challenge period or the final arbitration
func checkChannelGotoChallegingOrFinalDistributionWriteinChainState(state interfacev2.ChainStateOperation, assertAddress fields.Address, paychan *stores.Channel, obj channel.OnChainChannelPaymentArbitrationReconciliationBasis) error {

	channelId := obj.GetChannelId()

	// Channel cannot be closed
	if paychan.IsClosed() {
		return fmt.Errorf("Payment Channel <%s> is closed.", hex.EncodeToString(channelId))
	}
	// Check that the address matches you
	var assertAddressIsLeft = paychan.LeftAddress.Equal(assertAddress)
	var assertAddressIsRight = paychan.RightAddress.Equal(assertAddress)
	if !assertAddressIsLeft && !assertAddressIsRight {
		return fmt.Errorf("Payment Channel AssertAddress is not match left or right.")
	}
	// Check the signatures of both account addresses, and both parties check
	e20 := obj.CheckAddressAndSign(paychan.LeftAddress, paychan.RightAddress)
	if e20 != nil {
		return e20
	}
	// Check the reconciliation fund amount and reuse version
	channelReuseVersion := obj.GetReuseVersion()
	billAutoNumber := obj.GetAutoNumber()
	if channelReuseVersion != uint32(paychan.ReuseVersion) {
		return fmt.Errorf("Payment Channel ReuseVersion is not match, need <%d> but got <%d>.",
			paychan.ReuseVersion, channelReuseVersion)
	}
	// Check the reconciliation fund amount and reuse version
	objlamt := obj.GetLeftBalance()
	objramt := obj.GetRightBalance()
	billTotalAmt, e21 := objlamt.Add(&objramt)
	if e21 != nil {
		return e21
	}
	paychanTotalAmt, e22 := paychan.LeftAmount.Add(&paychan.RightAmount)
	if e22 != nil {
		return e22
	}
	if billTotalAmt.NotEqual(paychanTotalAmt) {
		return fmt.Errorf("Payment Channel Total Amount is not match, need %s but got %s.",
			paychanTotalAmt.ToFinString(), billTotalAmt.ToFinString())
	}
	// Arbitration amount
	var assertTargetAmount = objlamt // left
	var assertTargetSAT = obj.GetLeftSatoshi()
	if assertAddressIsRight {
		assertTargetAmount = objramt // right
		assertTargetSAT = obj.GetRightSatoshi()
	}

	// Judge the channel status, whether to enter the challenge period or finally seize it
	if paychan.IsOpening() {

		// Enter a challenging period
		blkhei := state.GetPendingBlockHeight()
		// Change state
		paychan.SetChallenging(blkhei, assertAddressIsLeft,
			&assertTargetAmount, assertTargetSAT, uint64(billAutoNumber))
		// Write status
		return state.ChannelUpdate(channelId, paychan)

	} else if paychan.IsChallenging() {

		// 只能夺取对方，不能既自己提出仲裁，然后又自己回应挑战
		if paychan.AssertAddressIsLeftOrRight.Check() == assertAddressIsLeft {
			return fmt.Errorf("The arbitration request and the response cannot be the same address")
		}

		// 判断仲裁，是否夺取对方资金
		if billAutoNumber <= uint64(paychan.AssertBillAutoNumber) {
			// The serial number of the bill is not satisfied (it must be greater than the serial number waiting for the challenge)
			return fmt.Errorf("Payment Channel BillAutoNumber must more than %d.", paychan.AssertBillAutoNumber)
		}
		// Higher serial number
		// Seize all funds and close the channel
		var lamt = fields.NewEmptyAmount()
		var ramt = fields.NewEmptyAmount()
		var lsat = fields.Satoshi(0)
		var rsat = fields.Satoshi(0)
		ttsat := paychan.LeftSatoshi.GetRealSatoshi() + paychan.RightSatoshi.GetRealSatoshi()
		if assertAddressIsLeft {
			lamt = paychanTotalAmt // The left account captures all funds, including HAC and sat
			lsat = ttsat
		} else {
			ramt = paychanTotalAmt // The right account captures all funds, including HAC and sat
			rsat = ttsat
		}
		// Close channels and seize all funds and interests
		isFinalClosed := true // 最终仲裁永久关闭
		return closePaymentChannelWriteinChainState(state, channelId, paychan, lamt, ramt, lsat, rsat, isFinalClosed)

	} else {
		return fmt.Errorf("Payment Channel <%s> status error.", hex.EncodeToString(channelId))
	}
}

// Challenge period or final arbitration fallback
func checkChannelGotoChallegingOrFinalDistributionRecoverChainState(state interfacev2.ChainStateOperation, assertAddress fields.Address, paychan *stores.Channel, obj channel.OnChainChannelPaymentArbitrationReconciliationBasis) error {

	panic("RecoverChainState() func is deleted.")

	channelId := obj.GetChannelId()

	// Judge the channel status, whether to enter the challenge period or finally seize it
	if paychan.IsFinalDistributionClosed() {
		if !paychan.IsHaveChallengeLog.Check() {
			return fmt.Errorf("IsHaveChallengeLog is not find.")
		}
		// Calculate the fallback amount
		paychanTotalAmt, _ := paychan.LeftAmount.Add(&paychan.RightAmount)
		var lamt = fields.NewEmptyAmount()
		var ramt = fields.NewEmptyAmount()
		if paychan.LeftAddress.Equal(assertAddress) {
			lamt = paychanTotalAmt // Account on the left seizes all funds
		} else {
			ramt = paychanTotalAmt // Account on the right seizes all funds
		}
		// Return to challenge period
		isBackToChalleging := true
		// Return account balance
		return closePaymentChannelRecoverChainState_deprecated(state, channelId, lamt, ramt, isBackToChalleging)

	} else if paychan.IsChallenging() {

		// Go back to the open state and clear the challenge period data
		paychan.SetOpening()
		paychan.CleanChallengingLog()
		return state.ChannelUpdate(channelId, paychan)

	} else {
		return fmt.Errorf("Payment Channel <%s> status error.", hex.EncodeToString(channelId))
	}
}

//////////////////////////////////////////////////////////////

// The inspection channel enters the challenge period or the final arbitration
func checkChannelGotoChallegingOrFinalDistributionWriteinChainStateV3(state interfaces.ChainStateOperation, assertAddress fields.Address, paychan *stores.Channel, obj channel.OnChainChannelPaymentArbitrationReconciliationBasis) error {

	channelId := obj.GetChannelId()

	// Channel cannot be closed
	if paychan.IsClosed() {
		return fmt.Errorf("Payment Channel <%s> is closed.", hex.EncodeToString(channelId))
	}
	// Check that the address matches you
	var assertAddressIsLeft = paychan.LeftAddress.Equal(assertAddress)
	var assertAddressIsRight = paychan.RightAddress.Equal(assertAddress)
	if !assertAddressIsLeft && !assertAddressIsRight {
		return fmt.Errorf("Payment Channel AssertAddress is not match left or right.")
	}
	// Check the signatures of both account addresses, and both parties check
	e20 := obj.CheckAddressAndSign(paychan.LeftAddress, paychan.RightAddress)
	if e20 != nil {
		return e20
	}
	// Check the reconciliation fund amount and reuse version
	channelReuseVersion := obj.GetReuseVersion()
	billAutoNumber := obj.GetAutoNumber()
	if channelReuseVersion != uint32(paychan.ReuseVersion) {
		return fmt.Errorf("Payment Channel ReuseVersion is not match, need <%d> but got <%d>.",
			paychan.ReuseVersion, channelReuseVersion)
	}
	// Check the reconciliation fund amount and reuse version
	objlamt := obj.GetLeftBalance()
	objramt := obj.GetRightBalance()
	billTotalAmt, e21 := objlamt.Add(&objramt)
	if e21 != nil {
		return e21
	}
	paychanTotalAmt, e22 := paychan.LeftAmount.Add(&paychan.RightAmount)
	if e22 != nil {
		return e22
	}
	if billTotalAmt.NotEqual(paychanTotalAmt) {
		return fmt.Errorf("Payment Channel Total Amount is not match, need %s but got %s.",
			paychanTotalAmt.ToFinString(), billTotalAmt.ToFinString())
	}
	// Arbitration amount
	var assertTargetAmount = objlamt // left
	var assertTargetSAT = obj.GetLeftSatoshi()
	if assertAddressIsRight {
		assertTargetAmount = objramt // right
		assertTargetSAT = obj.GetRightSatoshi()
	}

	// Judge the channel status, whether to enter the challenge period or finally seize it
	if paychan.IsOpening() {

		// Enter a challenging period
		blkhei := state.GetPendingBlockHeight()
		// Change state
		paychan.SetChallenging(blkhei, assertAddressIsLeft,
			&assertTargetAmount, assertTargetSAT, uint64(billAutoNumber))
		// Write status
		return state.ChannelUpdate(channelId, paychan)

	} else if paychan.IsChallenging() {

		// 只能夺取对方，不能既自己提出仲裁，然后又自己回应挑战
		if paychan.AssertAddressIsLeftOrRight.Check() == assertAddressIsLeft {
			return fmt.Errorf("The arbitration request and the response cannot be the same address")
		}

		// 判断仲裁，是否夺取对方资金
		if billAutoNumber <= uint64(paychan.AssertBillAutoNumber) {
			// The serial number of the bill is not satisfied (it must be greater than the serial number waiting for the challenge)
			return fmt.Errorf("Payment Channel BillAutoNumber must more than %d.", paychan.AssertBillAutoNumber)
		}
		// Higher serial number
		// Seize all funds and close the channel
		var lamt = fields.NewEmptyAmount()
		var ramt = fields.NewEmptyAmount()
		var lsat = fields.Satoshi(0)
		var rsat = fields.Satoshi(0)
		ttsat := paychan.LeftSatoshi.GetRealSatoshi() + paychan.RightSatoshi.GetRealSatoshi()
		if assertAddressIsLeft {
			lamt = paychanTotalAmt // The left account captures all funds, including HAC and sat
			lsat = ttsat
		} else {
			ramt = paychanTotalAmt // The right account captures all funds, including HAC and sat
			rsat = ttsat
		}
		// Close channels and seize all funds and interests
		isFinalClosed := true // 最终仲裁永久关闭
		return closePaymentChannelWriteinChainStateV3(state, channelId, paychan, lamt, ramt, lsat, rsat, isFinalClosed)

	} else {
		return fmt.Errorf("Payment Channel <%s> status error.", hex.EncodeToString(channelId))
	}
}

/////////////////////////////////////////////////////

// At the end of the challenge period, channel funds will be finally allocated as claimed
type Action_27_ClosePaymentChannelByClaimDistribution struct {
	// Channel ID
	ChannelId fields.ChannelId // Channel ID

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) Kind() uint16 {
	return 27
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) Size() uint32 {
	return 2 + elm.ChannelId.Size()
}

// json api
func (elm *Action_27_ClosePaymentChannelByClaimDistribution) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.ChannelId.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	return buffer.Bytes(), nil
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) RequestSignAddresses() []fields.Address {
	// No signature required
	return []fields.Address{}
}

func (act *Action_27_ClosePaymentChannelByClaimDistribution) WriteInChainState(state interfaces.ChainStateOperation) error {

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
	// Check status (must be challenge status)
	if paychan.IsChallenging() == false {
		return fmt.Errorf("Payment Channel status is not on challenging.")
	}
	// Check challenge duration
	clghei := state.GetPendingBlockHeight()
	expireHei := uint64(paychan.ChallengeLaunchHeight) + uint64(paychan.ArbitrationLockBlock)
	if clghei <= expireHei {
		// The challenge period is not over yet
		return fmt.Errorf("Payment Channel Challenging expire is %d.", expireHei)
	}
	// Allocate funds as claimed and close the channel
	var lamt = fields.NewEmptyAmount()
	var ramt = fields.NewEmptyAmount()
	ttamt, e := paychan.LeftAmount.Add(&paychan.RightAmount)
	if e != nil {
		return e
	}
	var lsat = fields.Satoshi(0)
	var rsat = fields.Satoshi(0)
	ttsat := paychan.LeftSatoshi.GetRealSatoshi() + paychan.RightSatoshi.GetRealSatoshi()
	if paychan.AssertAddressIsLeftOrRight.Check() {
		lamt = &paychan.AssertAmount // Left main Zhang
		ramt, _ = ttamt.Sub(lamt)
		lsat = paychan.AssertSatoshi.GetRealSatoshi()
		rsat = ttsat - lsat // The right side automatically gets the remaining funds
	} else {
		ramt = &paychan.AssertAmount // Right main Zhang
		lamt, _ = ttamt.Sub(ramt)
		rsat = paychan.AssertSatoshi.GetRealSatoshi()
		lsat = ttsat - rsat // The left side automatically gets the remaining funds
	}

	// Permanently closed
	isFinnalClosed := true
	return closePaymentChannelWriteinChainStateV3(state, act.ChannelId, paychan, lamt, ramt, lsat, rsat, isFinnalClosed)
}

func (act *Action_27_ClosePaymentChannelByClaimDistribution) WriteinChainState(state interfacev2.ChainStateOperation) error {

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
	// Check status (must be challenge status)
	if paychan.IsChallenging() == false {
		return fmt.Errorf("Payment Channel status is not on challenging.")
	}
	// Check challenge duration
	clghei := state.GetPendingBlockHeight()
	expireHei := uint64(paychan.ChallengeLaunchHeight) + uint64(paychan.ArbitrationLockBlock)
	if clghei <= expireHei {
		// The challenge period is not over yet
		return fmt.Errorf("Payment Channel Challenging expire is %d.", expireHei)
	}
	// Allocate funds as claimed and close the channel
	var lamt = fields.NewEmptyAmount()
	var ramt = fields.NewEmptyAmount()
	ttamt, e := paychan.LeftAmount.Add(&paychan.RightAmount)
	if e != nil {
		return e
	}
	var lsat = fields.Satoshi(0)
	var rsat = fields.Satoshi(0)
	ttsat := paychan.LeftSatoshi.GetRealSatoshi() + paychan.RightSatoshi.GetRealSatoshi()
	if paychan.AssertAddressIsLeftOrRight.Check() {
		lamt = &paychan.AssertAmount // Left main Zhang
		ramt, _ = ttamt.Sub(lamt)
		lsat = paychan.AssertSatoshi.GetRealSatoshi()
		rsat = ttsat - lsat // The right side automatically gets the remaining funds
	} else {
		ramt = &paychan.AssertAmount // Right main Zhang
		lamt, _ = ttamt.Sub(ramt)
		rsat = paychan.AssertSatoshi.GetRealSatoshi()
		lsat = ttsat - rsat // The left side automatically gets the remaining funds
	}

	// Permanently closed
	isFinnalClosed := true
	return closePaymentChannelWriteinChainState(state, act.ChannelId, paychan, lamt, ramt, lsat, rsat, isFinnalClosed)
}

func (act *Action_27_ClosePaymentChannelByClaimDistribution) RecoverChainState(state interfacev2.ChainStateOperation) error {

	// Query channel
	paychan, e := state.Channel(act.ChannelId)
	if e != nil {
		return e
	}
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// Fallback status
	// Allocate funds as claimed and close the channel
	var lamt = fields.NewEmptyAmount()
	var ramt = fields.NewEmptyAmount()
	if paychan.AssertAddressIsLeftOrRight.Check() {
		lamt = &paychan.AssertAmount // Left main Zhang
	} else {
		ramt = &paychan.AssertAmount // Right main Zhang
	}

	// close
	isFinnalClosed := true
	return closePaymentChannelRecoverChainState_deprecated(state, act.ChannelId, lamt, ramt, isFinnalClosed)
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_27_ClosePaymentChannelByClaimDistribution) IsBurning90PersentTxFees() bool {
	return false
}
