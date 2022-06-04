package actions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
)

type Action_29_SubmitTimeLimit struct {
	StartHeight fields.BlockHeight
	EndHeight   fields.BlockHeight

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_29_SubmitTimeLimit) Kind() uint16 {
	return 29
}

// json api
func (elm *Action_29_SubmitTimeLimit) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_29_SubmitTimeLimit) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var h1, _ = elm.StartHeight.Serialize()
	var h2, _ = elm.EndHeight.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(h1)
	buffer.Write(h2)
	return buffer.Bytes(), nil
}

func (elm *Action_29_SubmitTimeLimit) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.StartHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.EndHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_29_SubmitTimeLimit) Size() uint32 {
	return 2 + elm.StartHeight.Size() + elm.EndHeight.Size()
}

func (*Action_29_SubmitTimeLimit) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_29_SubmitTimeLimit) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	curblockhei := state.GetPendingBlockHeight()

	// start check
	if act.StartHeight > 0 {
		if curblockhei < uint64(act.StartHeight) {
			return fmt.Errorf("The transaction cannot be submitted if the current block height is less than %d", act.StartHeight)
		}
	}

	// end check
	if act.EndHeight > 0 {
		if curblockhei > uint64(act.EndHeight) {
			return fmt.Errorf("The transaction cannot be submitted if the current block height is more than %d", act.EndHeight)
		}
	}

	// submit time check ok pass
	return nil
}

func (act *Action_29_SubmitTimeLimit) WriteinChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	curblockhei := state.GetPendingBlockHeight()

	// start check
	if act.StartHeight > 0 {
		if curblockhei < uint64(act.StartHeight) {
			return fmt.Errorf("The transaction cannot be submitted if the current block height is less than %d", act.StartHeight)
		}
	}

	// end check
	if act.EndHeight > 0 {
		if curblockhei > uint64(act.EndHeight) {
			return fmt.Errorf("The transaction cannot be submitted if the current block height is more than %d", act.EndHeight)
		}
	}

	// submit time check ok pass
	return nil
}

func (act *Action_29_SubmitTimeLimit) RecoverChainState(state interfacev2.ChainStateOperation) error {
	panic("RecoverChainState() removed !")
}

// Set belongs to long_ trs
func (act *Action_29_SubmitTimeLimit) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs.(interfaces.Transaction)
}

func (act *Action_29_SubmitTimeLimit) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_29_SubmitTimeLimit) IsBurning90PersentTxFees() bool {
	return false
}
