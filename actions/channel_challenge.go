package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/sys"
)

// tong通过中间实时对账单方面关闭通道，进入挑战期
type Action_22_UnilateralClosePaymentChannelByRealtimeReconciliation struct {
	// 对账单
	Reconciliation OffChainFormPaymentChannelRealtimeReconciliation

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_22_UnilateralClosePaymentChannelByRealtimeReconciliation) Kind() uint16 {
	return 22
}

func (elm *Action_22_UnilateralClosePaymentChannelByRealtimeReconciliation) Size() uint32 {
	return 2 + elm.Reconciliation.Size()
}

// json api
func (elm *Action_22_UnilateralClosePaymentChannelByRealtimeReconciliation) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_22_UnilateralClosePaymentChannelByRealtimeReconciliation) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.Reconciliation.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	return buffer.Bytes(), nil
}

func (elm *Action_22_UnilateralClosePaymentChannelByRealtimeReconciliation) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.Reconciliation.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_22_UnilateralClosePaymentChannelByRealtimeReconciliation) RequestSignAddresses() []fields.Address {
	// 在执行的时候，查询出数据之后再检查检查签名
	return []fields.Address{}
}

func (act *Action_22_UnilateralClosePaymentChannelByRealtimeReconciliation) WriteinChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 查询通道
	paychan := state.Channel(act.Reconciliation.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.Reconciliation.ChannelId))
	}
	// 检查两个账户的签名，仅仅验证这两个地址
	signok, e0 := act.belong_trs.VerifyTargetSigns([]fields.Address{paychan.LeftAddress, paychan.RightAddress})
	if e0 != nil {
		return e0
	}
	if !signok { // 签名检查失败
		return fmt.Errorf("Payment Channel <%s> address signature verify fail.", hex.EncodeToString(act.Reconciliation.ChannelId))
	}
	// 分配金额可以为零但不能为负
	if act.Reconciliation.LeftAmount.IsNegative() {
		return fmt.Errorf("Payment channel distribution amount cannot be negative.")
	}
	// 检查分配金额
	var totalAmount, e1 = paychan.LeftAmount.Add(&paychan.RightAmount)
	if e1 != nil {
		return e1
	}
	// 分配金额不能超过总金额
	if act.Reconciliation.LeftAmount.MoreThan(totalAmount) {
		return fmt.Errorf("LeftAmount %s cannot more than total amount %s.",
			act.Reconciliation.LeftAmount.ToFinString(), totalAmount.ToFinString())
	}
	// 计算右侧金额
	var closedRightAmount, e2 = totalAmount.Sub(&act.Reconciliation.LeftAmount)
	if e2 != nil {
		return e2
	}
	// 写入状态
	return closePaymentChannelWriteinChainState(state, act.Reconciliation.ChannelId,
		paychan, &act.Reconciliation.LeftAmount, closedRightAmount)
}

func (act *Action_22_UnilateralClosePaymentChannelByRealtimeReconciliation) RecoverChainState(state interfaces.ChainStateOperation) error {

	// 查询通道
	paychan := state.Channel(act.Reconciliation.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.Reconciliation.ChannelId))
	}
	// 检查分配金额
	var totalAmount, _ = paychan.LeftAmount.Add(&paychan.RightAmount)
	// 计算右侧金额
	var closedRightAmount, _ = totalAmount.Sub(&act.Reconciliation.LeftAmount)
	return closePaymentChannelRecoverChainState(state, act.Reconciliation.ChannelId, &act.Reconciliation.LeftAmount, closedRightAmount)
}

func (elm *Action_22_UnilateralClosePaymentChannelByRealtimeReconciliation) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_22_UnilateralClosePaymentChannelByRealtimeReconciliation) IsBurning90PersentTxFees() bool {
	return false
}
