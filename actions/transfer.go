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

type Action_1_SimpleToTransfer struct {
	ToAddress fields.Address
	Amount    fields.Amount

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func NewAction_1_SimpleToTransfer(addr fields.Address, amt *fields.Amount) *Action_1_SimpleToTransfer {
	return &Action_1_SimpleToTransfer{
		ToAddress: addr,
		Amount:    *amt,
	}
}

func (elm *Action_1_SimpleToTransfer) Kind() uint16 {
	return 1
}

// json api
func (elm *Action_1_SimpleToTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_1_SimpleToTransfer) Serialize() ([]byte, error) {
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

func (elm *Action_1_SimpleToTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
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

func (elm *Action_1_SimpleToTransfer) Size() uint32 {
	return 2 + elm.ToAddress.Size() + elm.Amount.Size()
}

func (*Action_1_SimpleToTransfer) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_1_SimpleToTransfer) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	// check amount value
	if !act.Amount.IsPositive() {
		return fmt.Errorf("Amount is not positive.")
	}
	// 转移
	return DoSimpleTransferFromChainStateV3(state, act.belong_trs_v3.GetAddress(), act.ToAddress, act.Amount)
}

func (act *Action_1_SimpleToTransfer) WriteinChainState(state interfacev2.ChainStateOperation) error {
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

func (act *Action_1_SimpleToTransfer) RecoverChainState(state interfacev2.ChainStateOperation) error {

	panic("RecoverChainState be deprecated")

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 回退
	return DoSimpleTransferFromChainState(state, act.ToAddress, act.belong_trs.GetAddress(), act.Amount)
}

// 设置所属 belong_trs
func (act *Action_1_SimpleToTransfer) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}

func (act *Action_1_SimpleToTransfer) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_1_SimpleToTransfer) IsBurning90PersentTxFees() bool {
	///////////// TEST CODE START /////////////
	//if act.Amount.ToMeiString() == "2" { // 测试转账为 2 枚时，手续费减半
	//	return true
	//}
	///////////// TEST CODE END   /////////////
	return false
}

//////////////////////////////////////////

type Action_13_FromTransfer struct {
	FromAddress fields.Address
	Amount      fields.Amount

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func NewAction_13_FromTransfer(addr fields.Address, amt *fields.Amount) *Action_13_FromTransfer {
	return &Action_13_FromTransfer{
		FromAddress: addr,
		Amount:      *amt,
	}
}

func (elm *Action_13_FromTransfer) Kind() uint16 {
	return 13
}

// json api
func (elm *Action_13_FromTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_13_FromTransfer) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var addrBytes, _ = elm.FromAddress.Serialize()
	var amtBytes, _ = elm.Amount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(addrBytes)
	buffer.Write(amtBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_13_FromTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	moveseek, e := elm.FromAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	moveseek2, e := elm.Amount.Parse(buf, moveseek)
	if e != nil {
		return 0, e
	}
	return moveseek2, nil
}

func (elm *Action_13_FromTransfer) Size() uint32 {
	return 2 + elm.FromAddress.Size() + elm.Amount.Size()
}

func (elm *Action_13_FromTransfer) RequestSignAddresses() []fields.Address {
	return []fields.Address{
		elm.FromAddress,
	} // from sign
}

func (act *Action_13_FromTransfer) WriteInChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	// check amount value
	if !act.Amount.IsPositive() {
		return fmt.Errorf("Amount is not positive.")
	}
	// 转移
	return DoSimpleTransferFromChainStateV3(state, act.FromAddress, act.belong_trs_v3.GetAddress(), act.Amount)
}

func (act *Action_13_FromTransfer) WriteinChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// check amount value
	if !act.Amount.IsPositive() {
		return fmt.Errorf("Amount is not positive.")
	}
	// 转移
	return DoSimpleTransferFromChainState(state, act.FromAddress, act.belong_trs.GetAddress(), act.Amount)
}

func (act *Action_13_FromTransfer) RecoverChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 回退
	return DoSimpleTransferFromChainState(state, act.belong_trs.GetAddress(), act.FromAddress, act.Amount)
}

// 设置所属 belong_trs
func (act *Action_13_FromTransfer) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}

func (act *Action_13_FromTransfer) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_13_FromTransfer) IsBurning90PersentTxFees() bool {
	return false
}

////////////////////////////////////////

type Action_14_FromToTransfer struct {
	FromAddress fields.Address
	ToAddress   fields.Address
	Amount      fields.Amount

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func NewAction_14_FromToTransfer(fromaddr fields.Address, toaddr fields.Address, amt *fields.Amount) *Action_14_FromToTransfer {
	return &Action_14_FromToTransfer{
		FromAddress: fromaddr,
		ToAddress:   toaddr,
		Amount:      *amt,
	}
}

func (elm *Action_14_FromToTransfer) Kind() uint16 {
	return 14
}

// json api
func (elm *Action_14_FromToTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_14_FromToTransfer) Serialize() ([]byte, error) {
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

func (elm *Action_14_FromToTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.FromAddress.Parse(buf, seek)
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

func (elm *Action_14_FromToTransfer) Size() uint32 {
	return 2 + elm.FromAddress.Size() + elm.ToAddress.Size() + elm.Amount.Size()
}

func (elm *Action_14_FromToTransfer) RequestSignAddresses() []fields.Address {
	return []fields.Address{
		elm.FromAddress,
	} // from sign
}

func (act *Action_14_FromToTransfer) WriteInChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	// check amount value
	if !act.Amount.IsPositive() {
		return fmt.Errorf("Amount is not positive.")
	}
	// 转移
	return DoSimpleTransferFromChainStateV3(state, act.FromAddress, act.ToAddress, act.Amount)
}

func (act *Action_14_FromToTransfer) WriteinChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// check amount value
	if !act.Amount.IsPositive() {
		return fmt.Errorf("Amount is not positive.")
	}
	// 转移
	return DoSimpleTransferFromChainState(state, act.FromAddress, act.ToAddress, act.Amount)
}

func (act *Action_14_FromToTransfer) RecoverChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 回退
	return DoSimpleTransferFromChainState(state, act.ToAddress, act.FromAddress, act.Amount)
}

// 设置所属 belong_trs
func (act *Action_14_FromToTransfer) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}

func (act *Action_14_FromToTransfer) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_14_FromToTransfer) IsBurning90PersentTxFees() bool {
	return false
}
