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
	var moveseek, _ = elm.Address.Parse(buf, seek)
	var moveseek2, _ = elm.Amount.Parse(buf, moveseek)
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
