package actions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
)

type Action_8_SimpleSatoshiTransfer struct {
	ToAddress fields.Address
	Amount    fields.Satoshi

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func NewAction_8_SimpleSatoshiTransfer(addr fields.Address, amt fields.Satoshi) *Action_8_SimpleSatoshiTransfer {
	return &Action_8_SimpleSatoshiTransfer{
		ToAddress: addr,
		Amount:    amt,
	}
}

func (elm *Action_8_SimpleSatoshiTransfer) Kind() uint16 {
	return 8
}

// json api
func (elm *Action_8_SimpleSatoshiTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_8_SimpleSatoshiTransfer) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var addrBytes, _ = elm.ToAddress.Serialize()
	var amtBytes, _ = elm.Amount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(addrBytes)
	buffer.Write(amtBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_8_SimpleSatoshiTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	moveseek, e := elm.ToAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	moveseek2, e := elm.Amount.Parse(buf, moveseek)
	if e != nil {
		return 0, e
	}
	return moveseek2, nil
}

func (elm *Action_8_SimpleSatoshiTransfer) Size() uint32 {
	return 2 + elm.ToAddress.Size() + elm.Amount.Size()
}

func (elm *Action_8_SimpleSatoshiTransfer) RequestSignAddresses() []fields.Address {
	return []fields.Address{}
}

func (act *Action_8_SimpleSatoshiTransfer) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	if act.Amount <= 0 {
		// Transfer cannot be 0 or negative
		return fmt.Errorf("Amount <%d> error.", act.Amount)
	}
	// transfer
	fromAddress := act.belong_trs_v3.GetAddress()
	return DoSimpleSatoshiTransferFromChainStateV3(state, fromAddress, act.ToAddress, act.Amount)
}

func (act *Action_8_SimpleSatoshiTransfer) WriteinChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	if act.Amount <= 0 {
		// Transfer cannot be 0 or negative
		return fmt.Errorf("Amount <%d> error.", act.Amount)
	}
	// transfer
	fromAddress := act.belong_trs.GetAddress()
	return DoSimpleSatoshiTransferFromChainState(state, fromAddress, act.ToAddress, act.Amount)
}

func (act *Action_8_SimpleSatoshiTransfer) RecoverChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// Fallback
	return DoSimpleSatoshiTransferFromChainState(state, act.ToAddress, act.belong_trs.GetAddress(), act.Amount)
}

// Set belongs to long_ trs
func (act *Action_8_SimpleSatoshiTransfer) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}

func (act *Action_8_SimpleSatoshiTransfer) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_8_SimpleSatoshiTransfer) IsBurning90PersentTxFees() bool {
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////

type Action_11_FromToSatoshiTransfer struct {
	FromAddress fields.Address
	ToAddress   fields.Address
	Amount      fields.Satoshi

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func NewAction_11_FromToSatoshiTransfer(fromaddr fields.Address, toaddr fields.Address, amt fields.Satoshi) *Action_11_FromToSatoshiTransfer {
	return &Action_11_FromToSatoshiTransfer{
		FromAddress: fromaddr,
		ToAddress:   toaddr,
		Amount:      amt,
	}
}

func (elm *Action_11_FromToSatoshiTransfer) Kind() uint16 {
	return 11
}

// json api
func (elm *Action_11_FromToSatoshiTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_11_FromToSatoshiTransfer) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var addr1Bytes, _ = elm.FromAddress.Serialize()
	var addr2Bytes, _ = elm.ToAddress.Serialize()
	var amtBytes, _ = elm.Amount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(addr1Bytes)
	buffer.Write(addr2Bytes)
	buffer.Write(amtBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_11_FromToSatoshiTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, e := elm.FromAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ToAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.Amount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_11_FromToSatoshiTransfer) Size() uint32 {
	return 2 + elm.FromAddress.Size() + +elm.ToAddress.Size() + elm.Amount.Size()
}

func (elm *Action_11_FromToSatoshiTransfer) RequestSignAddresses() []fields.Address {
	return []fields.Address{
		elm.FromAddress, // From signature required
	}
}

func (act *Action_11_FromToSatoshiTransfer) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	if act.Amount <= 0 {
		// Transfer cannot be 0 or negative
		return fmt.Errorf("Amount <%d> error.", act.Amount)
	}

	// transfer
	return DoSimpleSatoshiTransferFromChainStateV3(state, act.FromAddress, act.ToAddress, act.Amount)
}

func (act *Action_11_FromToSatoshiTransfer) WriteinChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	if act.Amount <= 0 {
		// Transfer cannot be 0 or negative
		return fmt.Errorf("Amount <%d> error.", act.Amount)
	}

	// transfer
	return DoSimpleSatoshiTransferFromChainState(state, act.FromAddress, act.ToAddress, act.Amount)
}

func (act *Action_11_FromToSatoshiTransfer) RecoverChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// Fallback
	return DoSimpleSatoshiTransferFromChainState(state, act.ToAddress, act.FromAddress, act.Amount)
}

// Set belongs to long_ trs
func (act *Action_11_FromToSatoshiTransfer) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}
func (act *Action_11_FromToSatoshiTransfer) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_11_FromToSatoshiTransfer) IsBurning90PersentTxFees() bool {
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////

type Action_28_FromSatoshiTransfer struct {
	FromAddress fields.Address
	Amount      fields.Satoshi

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func NewAction_28_FromSatoshiTransfer(fromaddr fields.Address, amt fields.Satoshi) *Action_28_FromSatoshiTransfer {
	return &Action_28_FromSatoshiTransfer{
		FromAddress: fromaddr,
		Amount:      amt,
	}
}

func (elm *Action_28_FromSatoshiTransfer) Kind() uint16 {
	return 28
}

// json api
func (elm *Action_28_FromSatoshiTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_28_FromSatoshiTransfer) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var addr1Bytes, _ = elm.FromAddress.Serialize()
	var amtBytes, _ = elm.Amount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(addr1Bytes)
	buffer.Write(amtBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_28_FromSatoshiTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, e := elm.FromAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.Amount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_28_FromSatoshiTransfer) Size() uint32 {
	return 2 + elm.FromAddress.Size() + elm.Amount.Size()
}

func (elm *Action_28_FromSatoshiTransfer) RequestSignAddresses() []fields.Address {
	return []fields.Address{
		elm.FromAddress, // From signature required
	}
}

func (act *Action_28_FromSatoshiTransfer) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	if act.Amount <= 0 {
		// Transfer cannot be 0 or negative
		return fmt.Errorf("Amount <%d> error.", act.Amount)
	}

	// transfer
	toAddress := act.belong_trs_v3.GetAddress()
	return DoSimpleSatoshiTransferFromChainStateV3(state, act.FromAddress, toAddress, act.Amount)
}

func (act *Action_28_FromSatoshiTransfer) WriteinChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	if act.Amount <= 0 {
		// Transfer cannot be 0 or negative
		return fmt.Errorf("Amount <%d> error.", act.Amount)
	}

	// transfer
	toAddress := act.belong_trs.GetAddress()
	return DoSimpleSatoshiTransferFromChainState(state, act.FromAddress, toAddress, act.Amount)
}

func (act *Action_28_FromSatoshiTransfer) RecoverChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// Fallback
	return DoSimpleSatoshiTransferFromChainState(state, act.belong_trs.GetAddress(), act.FromAddress, act.Amount)
}

// Set belongs to long_ trs
func (act *Action_28_FromSatoshiTransfer) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}

func (act *Action_28_FromSatoshiTransfer) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_28_FromSatoshiTransfer) IsBurning90PersentTxFees() bool {
	return false
}
