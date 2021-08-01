package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/coinbase"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
)

/**
 * 支付通道交易类型
 */

// 开启支付通道
type Action_2_OpenPaymentChannel struct {
	ChannelId    fields.Bytes16 // 通道id
	LeftAddress  fields.Address // 账户1
	LeftAmount   fields.Amount  // 锁定金额
	RightAddress fields.Address // 账户2
	RightAmount  fields.Amount  // 锁定金额

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_2_OpenPaymentChannel) Kind() uint16 {
	return 2
}

func (elm *Action_2_OpenPaymentChannel) Size() uint32 {
	return 2 + elm.ChannelId.Size() +
		elm.LeftAddress.Size() +
		elm.LeftAmount.Size() +
		elm.RightAddress.Size() +
		elm.RightAmount.Size()
}

// json api
func (elm *Action_2_OpenPaymentChannel) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_2_OpenPaymentChannel) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var idBytes, _ = elm.ChannelId.Serialize()
	var addr1Bytes, _ = elm.LeftAddress.Serialize()
	var amt1Bytes, _ = elm.LeftAmount.Serialize()
	var addr2Bytes, _ = elm.RightAddress.Serialize()
	var amt2Bytes, _ = elm.RightAmount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(idBytes)
	buffer.Write(addr1Bytes)
	buffer.Write(amt1Bytes)
	buffer.Write(addr2Bytes)
	buffer.Write(amt2Bytes)
	return buffer.Bytes(), nil
}

func (elm *Action_2_OpenPaymentChannel) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_2_OpenPaymentChannel) RequestSignAddresses() []fields.Address {
	reqs := make([]fields.Address, 2)
	reqs[0] = elm.LeftAddress
	reqs[1] = elm.RightAddress
	return reqs
}

func (act *Action_2_OpenPaymentChannel) WriteinChainState(state interfaces.ChainStateOperation) error {
	// 查询通道是否存在
	sto := state.Channel(act.ChannelId)
	// 左右地址相同且协商一致关闭的通道ID可以被重用
	var reuseVersion fields.VarUint4 = 1
	var isIdCanUse = (sto == nil) ||
		(sto.IsAgreementClosed() && sto.LeftAddress.Equal(act.LeftAddress) && sto.RightAddress.Equal(act.RightAddress))
	if isIdCanUse == false {
		return fmt.Errorf("Payment Channel Id <%s> already exist.", hex.EncodeToString(act.ChannelId))
	}
	if sto != nil {
		reuseVersion = sto.ReuseVersion + 1 // 重用版本号增长
	}
	// 通道id合法性
	if len(act.ChannelId) != stores.ChannelIdLength || act.ChannelId[0] == 0 || act.ChannelId[stores.ChannelIdLength-1] == 0 {
		return fmt.Errorf("Payment Channel Id <%s> format error.", hex.EncodeToString(act.ChannelId))
	}
	// 两个地址不能相同
	if act.LeftAddress.Equal(act.RightAddress) {
		return fmt.Errorf("Left address cannot equal with right address.")
	}
	// 检查金额储存的位数
	labt, _ := act.LeftAmount.Serialize()
	rabt, _ := act.RightAmount.Serialize()
	if len(labt) > 6 || len(rabt) > 6 {
		// 避免锁定资金的储存位数过长，导致的复利计算后的值存储位数超过最大范围
		return fmt.Errorf("Payment Channel create error: left or right Amount bytes too long.")
	}
	// 不能为负数，或者两个通道同时为零（可以一个为正数另一个为零）
	if (!act.LeftAmount.IsPositive() || !act.RightAmount.IsPositive()) ||
		(act.LeftAmount.IsEmpty() && act.RightAmount.IsEmpty()) {
		return fmt.Errorf("Action_2_OpenPaymentChannel Payment Channel create error: left or right Amount is not positive.")
	}
	// 检查余额是否充足
	bls1 := state.Balance(act.LeftAddress)
	if bls1 == nil {
		return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance cannot empty.", act.LeftAddress.ToReadable())
	}
	amt1 := bls1.Hacash
	if amt1.LessThan(&act.LeftAmount) {
		return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance is not enough. need %s but got %s", act.RightAddress.ToReadable(), act.LeftAmount.ToFinString(), amt1.ToFinString())
	}
	bls2 := state.Balance(act.RightAddress)
	if bls2 == nil {
		return fmt.Errorf("Address %s Balance is not enough.", act.RightAddress.ToReadable())
	}
	amt2 := bls2.Hacash
	if amt2.LessThan(&act.RightAmount) {
		return fmt.Errorf("Action_2_OpenPaymentChannel Address %s Balance is not enough. need %s but got %s", act.RightAddress.ToReadable(), act.RightAmount.ToFinString(), amt2.ToFinString())
	}
	curheight := state.GetPendingBlockHeight()
	// 创建 channel
	var storeItem stores.Channel
	storeItem.BelongHeight = fields.BlockHeight(curheight)
	storeItem.LockBlock = fields.VarUint2(uint16(5000)) // 单方面提出的锁定期约为 17 天
	storeItem.LeftAddress = act.LeftAddress
	storeItem.LeftAmount = act.LeftAmount
	storeItem.RightAddress = act.RightAddress
	storeItem.RightAmount = act.RightAmount
	storeItem.ReuseVersion = reuseVersion // 重用版本号
	storeItem.SetOpening()                // 打开状态
	// 扣除余额
	DoSubBalanceFromChainState(state, act.LeftAddress, act.LeftAmount)
	DoSubBalanceFromChainState(state, act.RightAddress, act.RightAmount)
	// 储存通道
	state.ChannelCreate(act.ChannelId, &storeItem)
	// total supply 统计
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	// 累加锁入的HAC
	addamt := act.LeftAmount.ToMei() + act.RightAmount.ToMei()
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfLocatedInChannel, addamt)
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	//
	return nil
}

func (act *Action_2_OpenPaymentChannel) RecoverChainState(state interfaces.ChainStateOperation) error {
	sto := state.Channel(act.ChannelId)
	if sto.ReuseVersion > 1 {
		sto.ReuseVersion = sto.ReuseVersion - 1 // 重用版本号减少
	} else {
		// 删除通道
		state.ChannelDelete(act.ChannelId)
	}

	// 恢复余额
	DoAddBalanceFromChainState(state, act.LeftAddress, act.LeftAmount)
	DoAddBalanceFromChainState(state, act.RightAddress, act.RightAmount)
	// total supply 统计
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	// 回退解锁的HAC
	addamt := act.LeftAmount.ToMei() + act.RightAmount.ToMei()
	totalsupply.DoSub(stores.TotalSupplyStoreTypeOfLocatedInChannel, addamt)
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	return nil
}

func (elm *Action_2_OpenPaymentChannel) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_2_OpenPaymentChannel) IsBurning90PersentTxFees() bool {
	return false
}

/////////////////////////////////////////////////////////////////

// 关闭、结算 支付通道（资金分配不变的情况）
type Action_3_ClosePaymentChannel struct {
	ChannelId fields.Bytes16 // 通道id

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_3_ClosePaymentChannel) Kind() uint16 {
	return 3
}

func (elm *Action_3_ClosePaymentChannel) Size() uint32 {
	return 2 + elm.ChannelId.Size()
}

// json api
func (elm *Action_3_ClosePaymentChannel) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_3_ClosePaymentChannel) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var idBytes, _ = elm.ChannelId.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(idBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_3_ClosePaymentChannel) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = elm.ChannelId.Parse(buf, seek)
	return seek, nil
}

func (elm *Action_3_ClosePaymentChannel) RequestSignAddresses() []fields.Address {
	// 在执行的时候，查询出数据之后再检查检查签名
	return []fields.Address{}
}

func (act *Action_3_ClosePaymentChannel) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 查询通道
	paychan := state.Channel(act.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// 判断通道已经关闭
	if paychan.IsClosed() {
		return fmt.Errorf("Payment Channel <%s> is be closed.", hex.EncodeToString(act.ChannelId))
	}
	// 检查两个账户的签名 // 仅仅验证这两个地址
	signok, e1 := act.belong_trs.VerifyTargetSigns([]fields.Address{paychan.LeftAddress, paychan.RightAddress})
	if e1 != nil {
		return e1
	}
	if !signok { // 签名检查失败
		return fmt.Errorf("Payment Channel <%s> address signature verify fail.", hex.EncodeToString(act.ChannelId))
	}

	// 写入状态
	// 使用存入的金额计算通道利息
	return closePaymentChannelWriteinChainState(state, act.ChannelId, paychan, nil, nil, false)
}

func (act *Action_3_ClosePaymentChannel) RecoverChainState(state interfaces.ChainStateOperation) error {
	return closePaymentChannelRecoverChainState(state, act.ChannelId, nil, nil, false)
}

func (elm *Action_3_ClosePaymentChannel) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_3_ClosePaymentChannel) IsBurning90PersentTxFees() bool {
	return false
}

/////////////////////////////////////////////////////////////////

// 关闭、结算 支付通道（资金分配改变）
type Action_12_ClosePaymentChannelBySetupAmount struct {
	ChannelId    fields.Bytes16 // 通道id
	LeftAddress  fields.Address // 左侧账户
	LeftAmount   fields.Amount  // 左侧最终分配金额
	RightAddress fields.Address // 右侧账户
	RightAmount  fields.Amount  // 右侧最终分配金额

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) Kind() uint16 {
	return 12
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) Size() uint32 {
	return 2 + elm.ChannelId.Size() +
		elm.LeftAddress.Size() +
		elm.LeftAmount.Size() +
		elm.RightAddress.Size() +
		elm.RightAmount.Size()
}

// json api
func (elm *Action_12_ClosePaymentChannelBySetupAmount) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.ChannelId.Serialize()
	var bt2, _ = elm.LeftAddress.Serialize()
	var bt3, _ = elm.LeftAmount.Serialize()
	var bt4, _ = elm.RightAddress.Serialize()
	var bt5, _ = elm.RightAmount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	buffer.Write(bt3)
	buffer.Write(bt4)
	buffer.Write(bt5)
	return buffer.Bytes(), nil
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) RequestSignAddresses() []fields.Address {
	// 必须签名
	return []fields.Address{
		elm.LeftAddress,
		elm.RightAddress,
	}
}

func (act *Action_12_ClosePaymentChannelBySetupAmount) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 查询通道
	paychan := state.Channel(act.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// 检查两个账户是否匹配
	if paychan.LeftAddress.NotEqual(act.LeftAddress) ||
		paychan.RightAddress.NotEqual(act.RightAddress) {
		// 地址检查失败
		return fmt.Errorf("Payment Channel <%s> address not match.", act.RightAddress.ToReadable())
	}
	// 写入状态
	return closePaymentChannelWriteinChainState(state, act.ChannelId,
		paychan, &act.LeftAmount, &act.RightAmount, false)
}

func (act *Action_12_ClosePaymentChannelBySetupAmount) RecoverChainState(state interfaces.ChainStateOperation) error {
	return closePaymentChannelRecoverChainState(state, act.ChannelId, &act.LeftAmount, &act.RightAmount, false)
}

func (elm *Action_12_ClosePaymentChannelBySetupAmount) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_12_ClosePaymentChannelBySetupAmount) IsBurning90PersentTxFees() bool {
	return false
}

//////////////////////////////////////

// 关闭、结算 支付通道（资金分配改变）仅仅提供 left 的余额分配，自动计算 right 的分配
type Action_21_ClosePaymentChannelBySetupOnlyLeftAmount struct {
	ChannelId  fields.Bytes16 // 通道id
	LeftAmount fields.Amount  // 左侧最终分配金额

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) Kind() uint16 {
	return 21
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) Size() uint32 {
	return 2 + elm.ChannelId.Size() +
		elm.LeftAmount.Size()
}

// json api
func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.ChannelId.Serialize()
	var bt2, _ = elm.LeftAmount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	return buffer.Bytes(), nil
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) RequestSignAddresses() []fields.Address {
	// 在执行的时候，查询出数据之后再检查检查签名
	return []fields.Address{}
}

func (act *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) WriteinChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 查询通道
	paychan := state.Channel(act.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// 检查两个账户的签名，仅仅验证这两个地址
	signok, e0 := act.belong_trs.VerifyTargetSigns([]fields.Address{paychan.LeftAddress, paychan.RightAddress})
	if e0 != nil {
		return e0
	}
	if !signok { // 签名检查失败
		return fmt.Errorf("Payment Channel <%s> address signature verify fail.", hex.EncodeToString(act.ChannelId))
	}
	// 分配金额可以为零但不能为负
	if act.LeftAmount.IsNegative() {
		return fmt.Errorf("Payment channel distribution amount cannot be negative.")
	}
	// 检查分配金额
	var totalAmount, e1 = paychan.LeftAmount.Add(&paychan.RightAmount)
	if e1 != nil {
		return e1
	}
	// 分配金额不能超过总金额
	if act.LeftAmount.MoreThan(totalAmount) {
		return fmt.Errorf("LeftAmount %s cannot more than total amount %s.",
			act.LeftAmount.ToFinString(), totalAmount.ToFinString())
	}
	// 计算右侧金额
	var closedRightAmount, e2 = totalAmount.Sub(&act.LeftAmount)
	if e2 != nil {
		return e2
	}
	// 写入状态
	return closePaymentChannelWriteinChainState(state, act.ChannelId,
		paychan, &act.LeftAmount, closedRightAmount, false)
}

func (act *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) RecoverChainState(state interfaces.ChainStateOperation) error {

	// 查询通道
	paychan := state.Channel(act.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// 检查分配金额
	var totalAmount, _ = paychan.LeftAmount.Add(&paychan.RightAmount)
	// 计算右侧金额
	var closedRightAmount, _ = totalAmount.Sub(&act.LeftAmount)
	return closePaymentChannelRecoverChainState(state, act.ChannelId, &act.LeftAmount, closedRightAmount, false)
}

func (elm *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_21_ClosePaymentChannelBySetupOnlyLeftAmount) IsBurning90PersentTxFees() bool {
	return false
}

//////////////////////////////////////////////////////////

// 关闭通道状态写入
// isFinalClosed : 是否为仲裁终局结束，不可重用
func closePaymentChannelWriteinChainState(state interfaces.ChainStateOperation, channelId []byte, paychan *stores.Channel, newLeftAmt *fields.Amount, newRightAmt *fields.Amount, isFinalClosed bool) error {
	var e error
	// 判断通道已经关闭
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(channelId))
	}
	if paychan.IsClosed() {
		return fmt.Errorf("Payment Channel <%s> is be closed.", hex.EncodeToString(channelId))
	}
	// 通过时间计算利息
	if newLeftAmt == nil || newRightAmt == nil {
		// 自动使用存入的金额计算利息
		newLeftAmt = &paychan.LeftAmount
		newRightAmt = &paychan.RightAmount
	}
	// 计算总数
	// 分配金额可以为零但不能为负
	if newLeftAmt.IsNegative() || newRightAmt.IsNegative() {
		return fmt.Errorf("Payment channel distribution amount cannot be negative.")
	}
	// 检查分配金额是否与存入金额相等
	tt1, e1 := newLeftAmt.Add(newRightAmt)
	if e1 != nil {
		return e1
	}
	tt2, e2 := paychan.LeftAmount.Add(&paychan.RightAmount)
	if e2 != nil {
		return e2
	}
	if tt1.NotEqual(tt2) {
		// 不相等
		return fmt.Errorf("Payment channel distribution amount must equal with lock in.")
	}
	// 计算获得当前的区块高度
	//var curheight uint64 = 1
	curheight := state.GetPendingBlockHeight()
	leftAmount, rightAmount, haveinterest, e11 := calculateChannelInterest(
		curheight, uint64(paychan.BelongHeight), newLeftAmt, newRightAmt)
	if e11 != nil {
		return e11
	}
	// 增加余额（将锁定的金额和利息从通道中提取出来）
	e = DoAddBalanceFromChainState(state, paychan.LeftAddress, *leftAmount)
	if e != nil {
		return e
	}
	e = DoAddBalanceFromChainState(state, paychan.RightAddress, *rightAmount)
	if e != nil {
		return e
	}
	// 暂时保留通道用于数据回退
	// 计算左侧最终分配
	if isFinalClosed {
		paychan.SetFinalArbitrationClosed(newLeftAmt) // 仲裁永久关闭
	} else {
		paychan.SetAgreementClosed(newLeftAmt) // 协商关闭
	}
	e = state.ChannelUpdate(channelId, paychan)
	if e != nil {
		return e
	}
	//
	// total supply 统计
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	// 减少解锁的HAC
	lockamt := paychan.LeftAmount.ToMei() + paychan.RightAmount.ToMei()
	totalsupply.DoSub(stores.TotalSupplyStoreTypeOfLocatedInChannel, lockamt)
	// 增加通道利息统计
	if haveinterest {
		releaseamt := leftAmount.ToMei() + rightAmount.ToMei()
		//fmt.Println("(act *Action_3_ClosePaymentChannel) WriteinChainState", releaseamt, lockamt, releaseamt - lockamt, )
		//fmt.Println(paychanptr.LeftAddress.ToReadable(), paychanptr.LeftAmount.ToFinString(), paychanptr.LeftAmount.ToMei())
		//fmt.Println(paychanptr.RightAddress.ToReadable(), paychanptr.RightAmount.ToFinString(), paychanptr.RightAmount.ToMei())
		//fmt.Println(leftAmount.ToFinString(), leftAmount.ToMei(), rightAmount.ToFinString(), rightAmount.ToMei())
		if releaseamt-lockamt < 0 {
			return fmt.Errorf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		}
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfChannelInterest, releaseamt-lockamt)
	}
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	return nil

}

// 关闭通道状态回退
func closePaymentChannelRecoverChainState(state interfaces.ChainStateOperation, channelId []byte, newLeftAmt *fields.Amount, newRightAmt *fields.Amount, backToChallenging bool) error {

	var e error = nil
	// 查询通道
	paychan := state.Channel(channelId)
	if paychan == nil {
		// 通道必须被保存，才能被回退
		panic(fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(channelId)))
	}
	// 判断通道必须是已经关闭的状态
	if paychan.IsClosed() {
		panic(fmt.Errorf("Payment Channel <%s> is be closed.", hex.EncodeToString(channelId)))
	}
	if newLeftAmt == nil || newRightAmt == nil {
		// 自动使用存入的金额计算利息
		newLeftAmt = &paychan.LeftAmount
		newRightAmt = &paychan.RightAmount
	}
	// 计算差额
	curheight := state.GetPendingBlockHeight()
	// 计算利息
	leftAmount, rightAmount, haveinterest, e11 := calculateChannelInterest(curheight, uint64(paychan.BelongHeight), newLeftAmt, newRightAmt)
	if e11 != nil {
		return e11
	}
	// 减除余额（重新将金额放入通道）
	e = DoSubBalanceFromChainState(state, paychan.LeftAddress, *leftAmount)
	if e != nil {
		return e
	}
	e = DoSubBalanceFromChainState(state, paychan.RightAddress, *rightAmount)
	if e != nil {
		return e
	}
	// 恢复通道状态
	if backToChallenging {
		paychan.Status = stores.ChannelStatusChallenging // 回退到挑战状态
	} else {
		paychan.SetOpening() // 重新标记通道为开启状态
	}
	e = state.ChannelUpdate(channelId, paychan)
	if e != nil {
		return e
	}
	// total supply 统计
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	// 回退解锁的HAC
	lockamt := paychan.LeftAmount.ToMei() + paychan.RightAmount.ToMei()
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfLocatedInChannel, lockamt)
	// 回退通道利息统计
	if haveinterest {
		releaseamt := leftAmount.ToMei() + rightAmount.ToMei()
		totalsupply.DoSub(stores.TotalSupplyStoreTypeOfChannelInterest, releaseamt-lockamt)
	}
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	return nil
}

// 计算通道利息
// bool 是否有利息
func calculateChannelInterest(curheight uint64, openBelongHeight uint64, leftAmount *fields.Amount, rightAmount *fields.Amount) (*fields.Amount, *fields.Amount, bool, error) {
	// 增加利息计算，复利次数：约 2500 个区块 8.68 天增加一次万分之一的复利，少于8天忽略不计，年复合利息约 0.42%
	//a1, a2 := DoAppendCompoundInterest1Of10000By2500Height(&leftAmount, &rightAmount, insnum)
	var insnum = (curheight - openBelongHeight) / 2500
	var wfzn uint64 = 1 // 万分之一 1/10000
	// 通过开启通道的区块高度，修改一次增发比例
	if openBelongHeight > 200000 {
		// 增加利息计算，复利次数：约 10000 个区块 34 天增加一次千分之一的复利，少于34天忽略不计，年复合利息约 1.06%
		insnum = (curheight - openBelongHeight) / 10000
		wfzn = 10 // 千分之一 10/10000
	}
	if insnum > 0 {
		// 计算通道利息奖励
		a1, a2, e := coinbase.DoAppendCompoundInterestProportionOfHeightV2(leftAmount, rightAmount, insnum, wfzn)
		if e != nil {
			return nil, nil, false, e
		}
		// 加上了利息
		return a1, a2, true, nil
	}
	// 没有利息
	return leftAmount, rightAmount, false, nil
}
