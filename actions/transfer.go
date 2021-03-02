package actions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

type Action_1_SimpleTransfer struct {
	ToAddress fields.Address
	Amount    fields.Amount

	// data ptr
	belong_trs interfaces.Transaction
}

func NewAction_1_SimpleTransfer(addr fields.Address, amt *fields.Amount) *Action_1_SimpleTransfer {
	return &Action_1_SimpleTransfer{
		ToAddress: addr,
		Amount:    *amt,
	}
}

func (elm *Action_1_SimpleTransfer) Kind() uint16 {
	return 1
}

// json api
func (elm *Action_1_SimpleTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_1_SimpleTransfer) Serialize() ([]byte, error) {
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

func (elm *Action_1_SimpleTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
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

func (elm *Action_1_SimpleTransfer) Size() uint32 {
	return 2 + elm.ToAddress.Size() + elm.Amount.Size()
}

func (*Action_1_SimpleTransfer) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_1_SimpleTransfer) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// check amount value
	if !act.Amount.IsPositive() {
		return fmt.Errorf("Amount is not positive.")
	}
	// 转移
	return DoSimpleTransferFromChainState(state, act.belong_trs.GetAddress(), act.ToAddress, act.Amount)
}

func (act *Action_1_SimpleTransfer) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 回退
	return DoSimpleTransferFromChainState(state, act.ToAddress, act.belong_trs.GetAddress(), act.Amount)
}

// 设置所属 belone_trs
func (act *Action_1_SimpleTransfer) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_1_SimpleTransfer) IsBurning90PersentTxFees() bool {
	return false
}
