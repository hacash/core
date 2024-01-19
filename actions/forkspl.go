package actions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/sys"
)

type Action_30_SupportDistinguishForkChainID struct {
	CheckChainID fields.VarUint8 // Fork or Test chain ID

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_30_SupportDistinguishForkChainID) Kind() uint16 {
	return 30
}

// json api
func (elm *Action_30_SupportDistinguishForkChainID) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_30_SupportDistinguishForkChainID) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var h1, _ = elm.CheckChainID.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(h1)
	return buffer.Bytes(), nil
}

func (elm *Action_30_SupportDistinguishForkChainID) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.CheckChainID.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_30_SupportDistinguishForkChainID) Size() uint32 {
	return 2 + elm.CheckChainID.Size()
}

func (*Action_30_SupportDistinguishForkChainID) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_30_SupportDistinguishForkChainID) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// chain ID check
	var cid = sys.TransactionSystemCheckChainID
	var tarid = uint64(act.CheckChainID)
	if tarid > 0 && cid != tarid && !sys.NotCheckBlockDifficultyForMiner {
		return fmt.Errorf("Transaction System Check Chain ID need <%d> but got <%d>",
			cid, tarid)
	}

	// submit time check ok pass
	return nil
}

func (act *Action_30_SupportDistinguishForkChainID) WriteinChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// chain ID check
	if sys.TransactionSystemCheckChainID != uint64(act.CheckChainID) {
		return fmt.Errorf("Transaction System Check Chain ID need <%d> but got <%d>",
			sys.TransactionSystemCheckChainID, uint64(act.CheckChainID))
	}

	// submit time check ok pass
	return nil
}

func (act *Action_30_SupportDistinguishForkChainID) RecoverChainState(state interfacev2.ChainStateOperation) error {
	panic("RecoverChainState() removed !")
}

// Set belongs to long_ trs
func (act *Action_30_SupportDistinguishForkChainID) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs.(interfaces.Transaction)
}

func (act *Action_30_SupportDistinguishForkChainID) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs = trs
}

// is burning 90% fees
func (act *Action_30_SupportDistinguishForkChainID) IsBurning90PersentTxFees() bool {
	return false
}
