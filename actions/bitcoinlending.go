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
	"math/big"
)

/*

比特币系统借贷

*/

const ()

/*

创建比特币系统借贷流程：

. 检查合约ID格式
. 检查合约是否已经存在
. 检查用户比特币余额是否充足
. 检查预付利息是否充足
. 检查借贷HAC额度是否匹配
. 检查系统抵押曲线

. 扣除用户比特币余额
. 扣除用户预付利息
. 增加用户HAC余额
. 保存借贷合约

. 统计实时借贷数量份数
. 累计预销毁利息
. 累计借出流水

*/

// 比特币系统借贷
type Action_17_BitcoinsSystemLendingCreate struct {
	//
	LendingID                fields.Bytes15  // 借贷合约ID
	MortgageBitcoinPortion   fields.VarUint2 // 抵押比特币份数（每份 = 0.01BTC）最多抵押655枚比特币
	LoanTotalAmount          fields.Amount   // 总共借出HAC数量，必须小于等于可借数
	PreBurningInterestAmount fields.Amount   // 预先销毁的利息，必须大于等于销毁数量

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_17_BitcoinsSystemLendingCreate) Kind() uint16 {
	return 17
}

// json api
func (elm *Action_17_BitcoinsSystemLendingCreate) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_17_BitcoinsSystemLendingCreate) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	var b1, _ = elm.LendingID.Serialize()
	var b2, _ = elm.MortgageBitcoinPortion.Serialize()
	var b3, _ = elm.LoanTotalAmount.Serialize()
	var b4, _ = elm.PreBurningInterestAmount.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	return buffer.Bytes(), nil
}

func (elm *Action_17_BitcoinsSystemLendingCreate) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.LendingID.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.MortgageBitcoinPortion.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LoanTotalAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.PreBurningInterestAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_17_BitcoinsSystemLendingCreate) Size() uint32 {
	return 2 +
		elm.LendingID.Size() +
		elm.MortgageBitcoinPortion.Size() +
		elm.LoanTotalAmount.Size() +
		elm.PreBurningInterestAmount.Size()
}

func (*Action_17_BitcoinsSystemLendingCreate) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_17_BitcoinsSystemLendingCreate) WriteinChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// 检查数量上限
	if act.MortgageBitcoinPortion <= 0 {
		return fmt.Errorf("Bitcoin system lending mortgage bitcoin portion cannot empty.")
	}
	if act.MortgageBitcoinPortion > 10000 {
		return fmt.Errorf("Bitcoin system lending mortgage bitcoin portion max is 10000 (100BTC).")
	}

	// 检查id格式
	if len(act.LendingID) != stores.BitcoinSyslendIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.BitcoinSyslendIdLength-1] == 0 {
		return fmt.Errorf("Bitcoin Lending Id format error.")
	}

	// 查询id是否存在
	btclendObj := state.BitcoinSystemLending(act.LendingID)
	if btclendObj != nil {
		return fmt.Errorf("Bitcoin Lending <%s> already exist.", hex.EncodeToString(act.LendingID))
	}

	// 检查扣除比特币余额
	realSat := uint64(act.MortgageBitcoinPortion) * 100 * 10000 // 一份 = 0.01 BTC
	e0 := DoSubSatoshiFromChainState(state, feeAddr, fields.VarUint8(realSat))
	if e0 != nil {
		return e0 // 比特币余额不足
	}
	// 统计信息
	totalsupply, e1 := state.ReadTotalSupply()
	if e1 != nil {
		return e1
	}
	// 当前实时借出的比特币份数
	btcpartcurnum := totalsupply.Get(stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount)
	// 总的比特币份数
	totalbtcpart := totalsupply.Get(stores.TotalSupplyStoreTypeOfTransferBitcoin) * 100
	// 借贷比例，必须算上本次抵押的BTC份数
	alllendper := (btcpartcurnum + float64(act.MortgageBitcoinPortion)) / totalbtcpart * 100

	// 计算可借数量和预付利息，参数单位： %
	canLoanHacPart, predeshac := coinbase.CalculationOfInterestBitcoinMortgageLoanAmount(alllendper)
	canLoanHacPart *= float64(act.MortgageBitcoinPortion)
	predeshac *= float64(act.MortgageBitcoinPortion) // 真实预付利息份数

	// 真实数额，计算时忽略单位 240 后的小数部分
	realMaxLoanAmt, e3 := fields.NewAmountByBigIntWithUnit(
		big.NewInt(int64(canLoanHacPart*100*10000)),
		240,
	)
	if e3 != nil {
		return e3
	}
	realLowPreDes, e4 := fields.NewAmountByBigIntWithUnit(
		big.NewInt(int64(predeshac*100*10000)),
		240,
	)
	if e4 != nil {
		return e4
	}

	// 判断可借数量是否满足
	if act.LoanTotalAmount.MoreThan(realMaxLoanAmt) {
		return fmt.Errorf("Loan total amount %s can not more than real time effective amount %s.", act.LoanTotalAmount.ToFinString(), realMaxLoanAmt.ToFinString())
	}
	// 判断预付利息是否满足
	if act.PreBurningInterestAmount.LessThan(realLowPreDes) {
		return fmt.Errorf("Pre burning interest amount %s can not less than amount %s.", act.PreBurningInterestAmount.ToFinString(), realLowPreDes.ToFinString())
	}

	// 扣除预付利息
	e5 := DoSubBalanceFromChainState(state, feeAddr, act.PreBurningInterestAmount)
	if e5 != nil {
		return e5 // 比特币余额不足
	}

	// 借出数量增加到余额
	e6 := DoAddBalanceFromChainState(state, feeAddr, act.LoanTotalAmount)
	if e6 != nil {
		return e6
	}

	// 保存比特币抵押
	paddingHei := state.GetPendingBlockHeight()
	dlsto := &stores.BitcoinSystemLending{
		IsRansomed:                 fields.CreateBool(false), // 标记未赎回
		CreateBlockHeight:          fields.BlockHeight(paddingHei),
		MainAddress:                feeAddr,
		MortgageBitcoinPortion:     act.MortgageBitcoinPortion,
		LoanTotalAmount:            act.LoanTotalAmount,
		PreBurningInterestAmount:   act.PreBurningInterestAmount,
		RealtimeTotalMortgageRatio: fields.VarUint2(alllendper * 100), // 系统实时算上自己后的总抵押比例，单位：万分之，取值范围 0 ~ 10000
	}
	e11 := state.BitcoinLendingCreate(act.LendingID, dlsto)
	if e11 != nil {
		return e11
	}

	// 系统统计
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// 增加实时比特币系统抵押份数统计
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount,
		float64(act.MortgageBitcoinPortion),
	)
	// 比特币系统抵押累计预先销毁利息
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionBurningInterestHacAmount,
		act.PreBurningInterestAmount.ToMei(),
	)
	// 比特币系统抵押数量统计 累计借出流水
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationLoanHacAmount,
		act.LoanTotalAmount.ToMei(),
	)
	// 更新统计
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	// 完毕
	return nil
}

func (act *Action_17_BitcoinsSystemLendingCreate) RecoverChainState(state interfaces.ChainStateOperation) error {

	var e error = nil

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// 查询id是否存在
	btclendObj := state.BitcoinSystemLending(act.LendingID)
	if btclendObj == nil {
		return fmt.Errorf("Bitcoin Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// 回退比特币余额
	realSat := uint64(act.MortgageBitcoinPortion) * 100 * 10000 // 一份 = 0.01 BTC
	e = DoAddSatoshiFromChainState(state, feeAddr, fields.VarUint8(realSat))
	if e != nil {
		return e
	}

	// 统计信息
	totalsupply, e := state.ReadTotalSupply()
	if e != nil {
		return e
	}

	// 回退预付利息   增加
	e = DoAddBalanceFromChainState(state, feeAddr, act.PreBurningInterestAmount)
	if e != nil {
		return e
	}

	// 回退余额   减少
	e = DoSubBalanceFromChainState(state, feeAddr, act.LoanTotalAmount)
	if e != nil {
		return e
	}

	// 删除比特币抵押合约
	e = state.BitcoinLendingDelete(act.LendingID)
	if e != nil {
		return e
	}

	// 系统统计
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// 增加实时比特币系统抵押份数统计  回退减少
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount,
		float64(act.MortgageBitcoinPortion),
	)
	// 比特币系统抵押累计预先销毁利息     回退减少
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionBurningInterestHacAmount,
		act.PreBurningInterestAmount.ToMei(),
	)
	// 比特币系统抵押数量统计 累计借出流水   回退减少
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationLoanHacAmount,
		act.LoanTotalAmount.ToMei(),
	)
	// 更新统计
	e = state.UpdateSetTotalSupply(totalsupply)
	if e != nil {
		return e
	}

	return nil
}

// 设置所属 belong_trs
func (act *Action_17_BitcoinsSystemLendingCreate) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_17_BitcoinsSystemLendingCreate) IsBurning90PersentTxFees() bool {
	return false
}

/////////////////////////////////////////////////

/*

赎回钻石流程

. 检查抵押ID格式
. 检查抵押合约存在和状态
. 检查每个钻石状态
. 检查是否在私有赎回期
. 检查公共赎回期，并计算利息拍卖扣除
. 计算真实所需的赎回金额

. 检查赎回者HAC余额，并扣除赎回金额
. 增加用户钻石统计
. 修改每枚钻石状态
. 修改借贷合约状态

. 修改钻石抵押实时统计
. 累计赎回销毁HAC数量

*/

// 钻石系统借贷，赎回
type Action_18_BitcoinsSystemLendingRansom struct {
	//
	LendingID    fields.Bytes15 // 借贷合约ID
	RansomAmount fields.Amount  // 赎回金额
	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_18_BitcoinsSystemLendingRansom) Kind() uint16 {
	return 18
}

// json api
func (elm *Action_18_BitcoinsSystemLendingRansom) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_18_BitcoinsSystemLendingRansom) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	var b1, _ = elm.LendingID.Serialize()
	var b2, _ = elm.RansomAmount.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	return buffer.Bytes(), nil
}

func (elm *Action_18_BitcoinsSystemLendingRansom) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.LendingID.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RansomAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_18_BitcoinsSystemLendingRansom) Size() uint32 {
	return 2 +
		elm.LendingID.Size() +
		elm.RansomAmount.Size()
}

func (*Action_18_BitcoinsSystemLendingRansom) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_18_BitcoinsSystemLendingRansom) WriteinChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	paddingHeight := state.GetPendingBlockHeight()
	feeAddr := act.belong_trs.GetAddress()

	// 检查id格式
	if len(act.LendingID) != stores.BitcoinSyslendIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.BitcoinSyslendIdLength-1] == 0 {
		return fmt.Errorf("Bitcoin Lending Id format error.")
	}

	// 查询id是否存在
	btclendObj := state.BitcoinSystemLending(act.LendingID)
	if btclendObj == nil {
		return fmt.Errorf("Bitcoin Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// 检查是否赎回状态
	if btclendObj.IsRansomed.Check() {
		// 已经赎回。不可再次赎回
		return fmt.Errorf("Bitcoin Lending <%s> has been redeemed.", hex.EncodeToString(act.LendingID))
	}

	// 赎回期阶段区块数
	ransomBlockNumberBase := uint64(100000) // 十万个区块约一年
	if sys.TestDebugLocalDevelopmentMark {
		ransomBlockNumberBase = 10 // 测试环境 10 个区块为周期
	}

	// 计算比特币赎回金额
	_, realRansomAmt, e4 := coinbase.CalculationBitcoinSystemLendingRedeemAmount(
		feeAddr, btclendObj.MainAddress, &btclendObj.LoanTotalAmount,
		ransomBlockNumberBase,
		uint64(btclendObj.CreateBlockHeight), paddingHeight,
	)
	if e4 != nil {
		return e4
	}

	// 检查赎回金额是否有效（赎回金额真的大于实时计算的可赎回金额）检查赎回金额是否满足要求
	if act.RansomAmount.LessThan(realRansomAmt) {
		return fmt.Errorf("Ransom amount %s can not less than real time ransom amount %s.", act.RansomAmount.ToFinString(), realRansomAmt.ToFinString())
	}

	// 赎回操作，扣除HAC余额（以便首先检查余额是否充足）
	e2 := DoSubBalanceFromChainState(state, feeAddr, act.RansomAmount)
	if e2 != nil {
		return e2
	}

	// 增加比特币余额
	addSat := uint64(btclendObj.MortgageBitcoinPortion) * 100 * 10000
	e9 := DoAddSatoshiFromChainState(state, feeAddr, fields.VarUint8(addSat))
	if e9 != nil {
		return e9
	}

	// 修改抵押合约状态
	e20 := btclendObj.SetRansomedStatus(paddingHeight, &act.RansomAmount, feeAddr)
	if e20 != nil {
		return e20
	}
	e10 := state.BitcoinLendingUpdate(act.LendingID, btclendObj)
	if e10 != nil {
		return e10
	}

	// 系统统计
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// 减少实时比特币抵押份数统计
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount,
		float64(btclendObj.MortgageBitcoinPortion),
	)
	// 比特币系统抵押累计赎回销毁  流水
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationRansomHacAmount,
		act.RansomAmount.ToMei(),
	)
	// 更新统计
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	// 完毕
	return nil
}

func (act *Action_18_BitcoinsSystemLendingRansom) RecoverChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// 查询id是否存在
	btclendObj := state.BitcoinSystemLending(act.LendingID)
	if btclendObj == nil {
		return fmt.Errorf("Bitcoin Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// 回退赎回操作，增加HAC余额
	e2 := DoAddBalanceFromChainState(state, feeAddr, act.RansomAmount)
	if e2 != nil {
		return e2
	}

	// 回退减少比特币余额
	addSat := uint64(btclendObj.MortgageBitcoinPortion) * 100 * 10000
	e9 := DoSubSatoshiFromChainState(state, feeAddr, fields.VarUint8(addSat))
	if e9 != nil {
		return e9
	}

	// 修改抵押合约状态
	e20 := btclendObj.DropRansomedStatus() // 移除已赎回状态
	if e20 != nil {
		return e20
	}
	e10 := state.BitcoinLendingUpdate(act.LendingID, btclendObj)
	if e10 != nil {
		return e10
	}

	// 系统统计
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// 回退  增加实时比特币抵押份数统计
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount,
		float64(btclendObj.MortgageBitcoinPortion),
	)
	// 比特币系统抵押累计赎回销毁  流水  回退减少
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationRansomHacAmount,
		act.RansomAmount.ToMei(),
	)
	// 更新统计
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	return nil
}

// 设置所属 belong_trs
func (act *Action_18_BitcoinsSystemLendingRansom) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_18_BitcoinsSystemLendingRansom) IsBurning90PersentTxFees() bool {
	return false
}
