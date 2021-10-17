package actions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

type Action_8_SimpleSatoshiTransfer struct {
	ToAddress fields.Address
	Amount    fields.Satoshi

	// data ptr
	belong_trs interfaces.Transaction
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

func (act *Action_8_SimpleSatoshiTransfer) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	if act.Amount <= 0 {
		// 转账不能为 0 或负
		return fmt.Errorf("Amount <%d> error.", act.Amount)
	}
	// 转移
	return DoSimpleSatoshiTransferFromChainState(state, act.belong_trs.GetAddress(), act.ToAddress, act.Amount)
}

func (act *Action_8_SimpleSatoshiTransfer) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 回退
	return DoSimpleSatoshiTransferFromChainState(state, act.ToAddress, act.belong_trs.GetAddress(), act.Amount)
}

// 设置所属 belong_trs
func (act *Action_8_SimpleSatoshiTransfer) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
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
	belong_trs interfaces.Transaction
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
	reqs := make([]fields.Address, 1) // 需from签名
	reqs[0] = elm.FromAddress
	return reqs
}

func (act *Action_11_FromToSatoshiTransfer) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	if act.Amount <= 0 {
		// 转账不能为 0 或负
		return fmt.Errorf("Amount <%d> error.", act.Amount)
	}

	// 转移
	return DoSimpleSatoshiTransferFromChainState(state, act.FromAddress, act.ToAddress, act.Amount)
}

func (act *Action_11_FromToSatoshiTransfer) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 回退
	return DoSimpleSatoshiTransferFromChainState(state, act.ToAddress, act.FromAddress, act.Amount)
}

// 设置所属 belong_trs
func (act *Action_11_FromToSatoshiTransfer) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_11_FromToSatoshiTransfer) IsBurning90PersentTxFees() bool {
	return false
}
