package actions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
)

////////////////////////////////

// Mutual transfer of channel funds and chain funds, atomic exchange

type ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange struct {
	ChannelTranferProveBodyHashChecker fields.HashHalfChecker

	OnChainTranferToAddress fields.Address // Online transfer collection address
	OnChainTranferAmount    fields.Amount  // Online transfer amount

	AddressCount                            fields.VarUint1 // Signature quantity, can only be 2 or 3
	OnchainTransferFromAndMustSignAddresses []fields.Address
	// Two or three addresses, the first of which must be the from address for online transfer
	// The address list must contain the addresses of both sides of the channel, otherwise the verification will fail when submitting the challenge and Arbitration

	// Signature corresponding to address
	MustSigns []fields.Sign // 顺序与 []address 顺序必须一致
}

func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) Size() uint32 {
	size := elm.ChannelTranferProveBodyHashChecker.Size() +
		elm.OnChainTranferToAddress.Size() +
		elm.OnChainTranferAmount.Size() +
		elm.AddressCount.Size()
	size += uint32(len(elm.OnchainTransferFromAndMustSignAddresses)) * fields.AddressSize
	size += uint32(len(elm.MustSigns)) * fields.SignSize
	return size
}

func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) SerializeNoSign() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.ChannelTranferProveBodyHashChecker.Serialize()
	var bt2, _ = elm.OnChainTranferToAddress.Serialize()
	var bt3, _ = elm.OnChainTranferAmount.Serialize()
	var bt4, _ = elm.AddressCount.Serialize()
	buffer.Write(bt1)
	buffer.Write(bt2)
	buffer.Write(bt3)
	buffer.Write(bt4)
	for _, addr := range elm.OnchainTransferFromAndMustSignAddresses {
		var bt1, _ = addr.Serialize()
		buffer.Write(bt1)
	}
	return buffer.Bytes(), nil
}

func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) SignStuffHash() (fields.Hash, error) {
	var conbt, e = elm.SerializeNoSign() // Data body
	if e != nil {
		return nil, e
	}
	return fields.CalculateHash(conbt), nil
}

func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.SerializeNoSign() // Data body
	buffer.Write(bt1)
	for _, sign := range elm.MustSigns {
		var bt1, _ = sign.Serialize()
		buffer.Write(bt1)
	}
	return buffer.Bytes(), nil
}

func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelTranferProveBodyHashChecker.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.OnChainTranferToAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.OnChainTranferAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	// address
	seek, e = elm.AddressCount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	scn := int(elm.AddressCount)
	elm.OnchainTransferFromAndMustSignAddresses = make([]fields.Address, scn)
	for i := 0; i < scn; i++ {
		seek, e = elm.OnchainTransferFromAndMustSignAddresses[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// autograph
	elm.MustSigns = make([]fields.Sign, scn)
	for i := 0; i < scn; i++ {
		seek, e = elm.MustSigns[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// complete
	return seek, nil
}

// Check all signatures
func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) CheckMustAddressAndSigns() error {
	var e error

	// Compute hash
	conhx, e := elm.SignStuffHash()
	if e != nil {
		return e
	}

	// Minimum number of check addresses
	sgmn := len(elm.OnchainTransferFromAndMustSignAddresses)
	if sgmn < 2 || sgmn > 3 || sgmn != int(elm.AddressCount) || sgmn != len(elm.MustSigns) {
		return fmt.Errorf("Address or Sign length error, need 2~3 but got %d, %d, %d.",
			sgmn, int(elm.AddressCount), len(elm.MustSigns))
	}

	// Signatures are arranged by address. Check whether all addresses and signatures match
	for i := 0; i < sgmn; i++ {
		sign := elm.MustSigns[i]
		addr := elm.OnchainTransferFromAndMustSignAddresses[i]
		sgaddr := account.NewAddressFromPublicKeyV0(sign.PublicKey)
		// Judge address order
		if addr.NotEqual(sgaddr) {
			return fmt.Errorf("Address not match, need %s nut got %s.",
				addr.ToReadable(), fields.Address(sgaddr).ToReadable())
		}
		// Check signature
		ok, _ := account.CheckSignByHash32(conhx, sign.PublicKey, sign.Signature)
		if !ok {
			return fmt.Errorf("Left account %s verify signature fail.", addr.ToReadable())
		}
	}

	// All signatures verified successfully
	return nil
}

////////////////////////////////////////////////////////

// Exchange of channels with atoms in chains
// Submit channel and chain swap transactions
type Action_25_PaymantChannelAndOnchainAtomicExchange struct {

	// Atomic swap transaction receipt
	ExchangeEvidence ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) Kind() uint16 {
	return 25
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) Size() uint32 {
	return 2 + elm.ExchangeEvidence.Size()
}

// json api
func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.ExchangeEvidence.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	return buffer.Bytes(), nil
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ExchangeEvidence.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) RequestSignAddresses() []fields.Address {
	// Action internal judgment signature
	return []fields.Address{}
}

func (act *Action_25_PaymantChannelAndOnchainAtomicExchange) WriteInChainState(state interfaces.ChainStateOperation) error {

	var e error

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Query whether it is a duplicate submission
	swaphx := act.ExchangeEvidence.ChannelTranferProveBodyHashChecker
	chaswap, e := state.Chaswap(swaphx)
	if e != nil {
		return e
	}
	if chaswap != nil {
		// Already exists, cannot submit repeatedly
		// Otherwise, it will cause repeated transfers
		return fmt.Errorf("ChannelTranferProveBodyHashChecker <%s> is existence.",
			swaphx.ToHex())
	}

	if len(act.ExchangeEvidence.OnchainTransferFromAndMustSignAddresses) < 2 {
		return fmt.Errorf("Address lenght error.")
	}

	// The content related to the channel is not verified during submission, and only the online transfer is operated
	e = act.ExchangeEvidence.CheckMustAddressAndSigns()
	if e != nil {
		return e
	}

	// Create, save voucher
	objsto := stores.Chaswap{
		IsBeUsed:                                fields.CreateBool(false), // 未使用过
		AddressCount:                            act.ExchangeEvidence.AddressCount,
		OnchainTransferFromAndMustSignAddresses: act.ExchangeEvidence.OnchainTransferFromAndMustSignAddresses,
	}
	e = state.ChaswapCreate(swaphx, &objsto)
	if e != nil {
		return e
	}

	// transfer accounts
	fromAddr := act.ExchangeEvidence.OnchainTransferFromAndMustSignAddresses[0]
	toAddr := act.ExchangeEvidence.OnChainTranferToAddress
	trsAmt := act.ExchangeEvidence.OnChainTranferAmount
	return DoSimpleTransferFromChainState(state, fromAddr, toAddr, trsAmt)
}

func (act *Action_25_PaymantChannelAndOnchainAtomicExchange) WriteinChainState(state interfacev2.ChainStateOperation) error {

	var e error

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Query whether it is a duplicate submission
	swaphx := act.ExchangeEvidence.ChannelTranferProveBodyHashChecker
	chaswap, e := state.Chaswap(swaphx)
	if e != nil {
		return e
	}
	if chaswap != nil {
		// Already exists, cannot submit repeatedly
		// Otherwise, it will cause repeated transfers
		return fmt.Errorf("ChannelTranferProveBodyHashChecker <%s> is existence.",
			swaphx.ToHex())
	}

	if len(act.ExchangeEvidence.OnchainTransferFromAndMustSignAddresses) < 2 {
		return fmt.Errorf("Address lenght error.")
	}

	// The content related to the channel is not verified during submission, and only the online transfer is operated
	e = act.ExchangeEvidence.CheckMustAddressAndSigns()
	if e != nil {
		return e
	}

	// Create, save voucher
	objsto := stores.Chaswap{
		IsBeUsed:                                fields.CreateBool(false), // 未使用过
		AddressCount:                            act.ExchangeEvidence.AddressCount,
		OnchainTransferFromAndMustSignAddresses: act.ExchangeEvidence.OnchainTransferFromAndMustSignAddresses,
	}
	e = state.ChaswapCreate(swaphx, &objsto)
	if e != nil {
		return e
	}

	// transfer accounts
	fromAddr := act.ExchangeEvidence.OnchainTransferFromAndMustSignAddresses[0]
	toAddr := act.ExchangeEvidence.OnChainTranferToAddress
	trsAmt := act.ExchangeEvidence.OnChainTranferAmount
	return DoSimpleTransferFromChainStateV2(state, fromAddr, toAddr, trsAmt)
}

func (act *Action_25_PaymantChannelAndOnchainAtomicExchange) RecoverChainState(state interfacev2.ChainStateOperation) error {

	// Cancel save
	swaphx := act.ExchangeEvidence.ChannelTranferProveBodyHashChecker
	state.ChaswapDelete(swaphx)

	// Transfer return
	fromAddr := act.ExchangeEvidence.OnchainTransferFromAndMustSignAddresses[0]
	toAddr := act.ExchangeEvidence.OnChainTranferToAddress
	trsAmt := act.ExchangeEvidence.OnChainTranferAmount
	return DoSimpleTransferFromChainStateV2(state, toAddr, fromAddr, trsAmt)
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_25_PaymantChannelAndOnchainAtomicExchange) IsBurning90PersentTxFees() bool {
	return false
}
