package actions

import (
	"bytes"
	"encoding/binary"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

type Action_8_SimpleSatoshiTransfer struct {
	Address fields.Address
	Amount  fields.VarUint8

	// data ptr
	belong_trs interfaces.Transaction
}

func NewAction_8_SimpleSatoshiTransfer(addr fields.Address, amt fields.VarUint8) *Action_8_SimpleSatoshiTransfer {
	return &Action_8_SimpleSatoshiTransfer{
		Address: addr,
		Amount:  amt,
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
	var addrBytes, _ = elm.Address.Serialize()
	var amtBytes, _ = elm.Amount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(addrBytes)
	buffer.Write(amtBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_8_SimpleSatoshiTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	moveseek, e := elm.Address.Parse(buf, seek)
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
	return 2 + elm.Address.Size() + elm.Amount.Size()
}

func (elm *Action_8_SimpleSatoshiTransfer) RequestSignAddresses() []fields.Address {
	return []fields.Address{}
}

func (act *Action_8_SimpleSatoshiTransfer) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 转移
	return DoSimpleSatoshiTransferFromChainState(state, act.belong_trs.GetAddress(), act.Address, act.Amount)
}

func (act *Action_8_SimpleSatoshiTransfer) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 回退
	return DoSimpleSatoshiTransferFromChainState(state, act.Address, act.belong_trs.GetAddress(), act.Amount)
}

// 设置所属 belone_trs
func (act *Action_8_SimpleSatoshiTransfer) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}

///////////////////////////////////////////////////////////////////////////////////////////////

type Action_11_FromToSatoshiTransfer struct {
	FromAddress fields.Address
	ToAddress   fields.Address
	Amount      fields.VarUint8

	// data ptr
	belong_trs interfaces.Transaction
}

func NewAction_11_FromToSatoshiTransfer(fromaddr fields.Address, toaddr fields.Address, amt fields.VarUint8) *Action_11_FromToSatoshiTransfer {
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

// 设置所属 belone_trs
func (act *Action_11_FromToSatoshiTransfer) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}
