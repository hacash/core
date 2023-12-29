package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/coinbase"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
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

// Bitcoin system lending
type Action_17_BitcoinsSystemLendingCreate struct {
	//
	LendingID                fields.BitcoinSyslendId // Loan contract ID
	MortgageBitcoinPortion   fields.VarUint2         // Mortgage bitcoin shares (each = 0.01btc) up to 655 bitcoins
	LoanTotalAmount          fields.Amount           // The total lending HAC quantity must be less than or equal to the lendable quantity
	PreBurningInterestAmount fields.Amount           // Interest for pre destruction must be greater than or equal to the destroyed quantity

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
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

func (act *Action_17_BitcoinsSystemLendingCreate) WriteInChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs_v3.GetAddress()

	// Upper limit of inspection quantity
	if act.MortgageBitcoinPortion <= 0 {
		return fmt.Errorf("Bitcoin system lending mortgage bitcoin portion cannot empty.")
	}
	if act.MortgageBitcoinPortion > 10000 {
		return fmt.Errorf("Bitcoin system lending mortgage bitcoin portion max is 10000 (100BTC).")
	}

	// Check ID format
	if len(act.LendingID) != stores.BitcoinSyslendIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.BitcoinSyslendIdLength-1] == 0 {
		return fmt.Errorf("Bitcoin Lending Id format error.")
	}

	// Query whether the ID exists
	btclendObj, e := state.BitcoinSystemLending(act.LendingID)
	if e != nil {
		return e
	}
	if btclendObj != nil {
		return fmt.Errorf("Bitcoin Lending <%s> already exist.", hex.EncodeToString(act.LendingID))
	}

	// Check and deduct bitcoin balance
	realSat := uint64(act.MortgageBitcoinPortion) * 100 * 10000 // 一份 = 0.01 BTC
	e0 := DoSubSatoshiFromChainStateV3(state, feeAddr, fields.Satoshi(realSat))
	if e0 != nil {
		return e0 // Insufficient bitcoin balance
	}
	// statistical information
	totalsupply, e1 := state.ReadTotalSupply()
	if e1 != nil {
		return e1
	}
	// Current number of bitcoin copies lent in real time
	btcpartcurnum := float64(totalsupply.GetUint(stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount))
	// Total bitcoin copies
	totalbtcpart := float64(totalsupply.GetUint(stores.TotalSupplyStoreTypeOfTransferBitcoin) * 100)
	// Loan ratio, which must include the BTC copies of this mortgage
	alllendper := (btcpartcurnum + float64(act.MortgageBitcoinPortion)) / totalbtcpart * 100

	// Calculate the debit quantity and prepaid interest. Parameter unit:%
	canLoanHacPart, predeshac := coinbase.CalculationOfInterestBitcoinMortgageLoanAmount(alllendper)
	canLoanHacPart *= float64(act.MortgageBitcoinPortion)
	predeshac *= float64(act.MortgageBitcoinPortion) // Actual prepaid interest shares

	// Actual amount, the decimal part after unit 240 is ignored in calculation
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

	// Judge whether the quantity that can be borrowed meets
	if act.LoanTotalAmount.MoreThan(realMaxLoanAmt) {
		return fmt.Errorf("Loan total amount %s can not more than real time effective amount %s.", act.LoanTotalAmount.ToFinString(), realMaxLoanAmt.ToFinString())
	}
	// Judge whether the prepaid interest is satisfied
	if act.PreBurningInterestAmount.LessThan(realLowPreDes) {
		return fmt.Errorf("Pre burning interest amount %s can not less than amount %s.", act.PreBurningInterestAmount.ToFinString(), realLowPreDes.ToFinString())
	}

	// Deduct prepaid interest
	e5 := DoSubBalanceFromChainStateV3(state, feeAddr, act.PreBurningInterestAmount)
	if e5 != nil {
		return e5 // Insufficient bitcoin balance
	}

	// Increase lending quantity to balance
	e6 := DoAddBalanceFromChainStateV3(state, feeAddr, act.LoanTotalAmount)
	if e6 != nil {
		return e6
	}

	// Save bitcoin mortgage
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

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Increase the statistics of mortgage copies of real-time bitcoin system
	totalsupply.DoAddUint(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount,
		uint64(act.MortgageBitcoinPortion),
	)
	// Accumulated pre destruction interest of bitcoin system mortgage
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionBurningInterestHacAmount,
		act.PreBurningInterestAmount.ToMei(),
	)
	// Bitcoin system mortgage quantity statistics cumulative lending flow
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationLoanHacAmount,
		act.LoanTotalAmount.ToMei(),
	)
	// Update statistics
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	// complete
	return nil
}

func (act *Action_17_BitcoinsSystemLendingCreate) WriteinChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// Upper limit of inspection quantity
	if act.MortgageBitcoinPortion <= 0 {
		return fmt.Errorf("Bitcoin system lending mortgage bitcoin portion cannot empty.")
	}
	if act.MortgageBitcoinPortion > 10000 {
		return fmt.Errorf("Bitcoin system lending mortgage bitcoin portion max is 10000 (100BTC).")
	}

	// Check ID format
	if len(act.LendingID) != stores.BitcoinSyslendIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.BitcoinSyslendIdLength-1] == 0 {
		return fmt.Errorf("Bitcoin Lending Id format error.")
	}

	// Query whether the ID exists
	btclendObj, e := state.BitcoinSystemLending(act.LendingID)
	if e != nil {
		return e
	}
	if btclendObj != nil {
		return fmt.Errorf("Bitcoin Lending <%s> already exist.", hex.EncodeToString(act.LendingID))
	}

	// Check and deduct bitcoin balance
	realSat := uint64(act.MortgageBitcoinPortion) * 100 * 10000 // 一份 = 0.01 BTC
	e0 := DoSubSatoshiFromChainState(state, feeAddr, fields.Satoshi(realSat))
	if e0 != nil {
		return e0 // Insufficient bitcoin balance
	}
	// statistical information
	totalsupply, e1 := state.ReadTotalSupply()
	if e1 != nil {
		return e1
	}
	// Current number of bitcoin copies lent in real time
	btcpartcurnum := float64(totalsupply.GetUint(stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount))
	// Total bitcoin copies
	totalbtcpart := float64(totalsupply.GetUint(stores.TotalSupplyStoreTypeOfTransferBitcoin) * 100)
	// Loan ratio, which must include the BTC copies of this mortgage
	alllendper := (btcpartcurnum + float64(act.MortgageBitcoinPortion)) / totalbtcpart * 100

	// Calculate the debit quantity and prepaid interest. Parameter unit:%
	canLoanHacPart, predeshac := coinbase.CalculationOfInterestBitcoinMortgageLoanAmount(alllendper)
	canLoanHacPart *= float64(act.MortgageBitcoinPortion)
	predeshac *= float64(act.MortgageBitcoinPortion) // Actual prepaid interest shares

	// Actual amount, the decimal part after unit 240 is ignored in calculation
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

	// Judge whether the quantity that can be borrowed meets
	if act.LoanTotalAmount.MoreThan(realMaxLoanAmt) {
		return fmt.Errorf("Loan total amount %s can not more than real time effective amount %s.", act.LoanTotalAmount.ToFinString(), realMaxLoanAmt.ToFinString())
	}
	// Judge whether the prepaid interest is satisfied
	if act.PreBurningInterestAmount.LessThan(realLowPreDes) {
		return fmt.Errorf("Pre burning interest amount %s can not less than amount %s.", act.PreBurningInterestAmount.ToFinString(), realLowPreDes.ToFinString())
	}

	// Deduct prepaid interest
	e5 := DoSubBalanceFromChainState(state, feeAddr, act.PreBurningInterestAmount)
	if e5 != nil {
		return e5 // Insufficient bitcoin balance
	}

	// Increase lending quantity to balance
	e6 := DoAddBalanceFromChainState(state, feeAddr, act.LoanTotalAmount)
	if e6 != nil {
		return e6
	}

	// Save bitcoin mortgage
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

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Increase the statistics of mortgage copies of real-time bitcoin system
	totalsupply.DoAddUint(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount,
		uint64(act.MortgageBitcoinPortion),
	)
	// Accumulated pre destruction interest of bitcoin system mortgage
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionBurningInterestHacAmount,
		act.PreBurningInterestAmount.ToMei(),
	)
	// Bitcoin system mortgage quantity statistics cumulative lending flow
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationLoanHacAmount,
		act.LoanTotalAmount.ToMei(),
	)
	// Update statistics
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	// complete
	return nil
}

func (act *Action_17_BitcoinsSystemLendingCreate) RecoverChainState(state interfacev2.ChainStateOperation) error {

	var e error = nil

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// Query whether the ID exists
	btclendObj, e := state.BitcoinSystemLending(act.LendingID)
	if e != nil {
		return e
	}
	if btclendObj == nil {
		return fmt.Errorf("Bitcoin Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// Back bitcoin balance
	realSat := uint64(act.MortgageBitcoinPortion) * 100 * 10000 // 一份 = 0.01 BTC
	e = DoAddSatoshiFromChainState(state, feeAddr, fields.Satoshi(realSat))
	if e != nil {
		return e
	}

	// statistical information
	totalsupply, e := state.ReadTotalSupply()
	if e != nil {
		return e
	}

	// Increase in prepaid interest refunded
	e = DoAddBalanceFromChainState(state, feeAddr, act.PreBurningInterestAmount)
	if e != nil {
		return e
	}

	// Decrease in return balance
	e = DoSubBalanceFromChainState(state, feeAddr, act.LoanTotalAmount)
	if e != nil {
		return e
	}

	// Delete bitcoin mortgage contract
	e = state.BitcoinLendingDelete(act.LendingID)
	if e != nil {
		return e
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Increase the statistics of mortgage copies of real-time bitcoin system and reduce the fallback
	totalsupply.DoSubUint(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount,
		uint64(act.MortgageBitcoinPortion),
	)
	// Bitcoin system mortgage cumulative pre destruction interest rebate decrease
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionBurningInterestHacAmount,
		act.PreBurningInterestAmount.ToMei(),
	)
	// Bitcoin system mortgage quantity statistics cumulative lending daily return decrease
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationLoanHacAmount,
		act.LoanTotalAmount.ToMei(),
	)
	// Update statistics
	e = state.UpdateSetTotalSupply(totalsupply)
	if e != nil {
		return e
	}

	return nil
}

// Set belongs to long_ trs
func (act *Action_17_BitcoinsSystemLendingCreate) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}
func (act *Action_17_BitcoinsSystemLendingCreate) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
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

// Diamond system lending, redemption
type Action_18_BitcoinsSystemLendingRansom struct {
	//
	LendingID    fields.BitcoinSyslendId // Loan contract ID
	RansomAmount fields.Amount           // Redemption amount
	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
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

func (act *Action_18_BitcoinsSystemLendingRansom) WriteInChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	paddingHeight := state.GetPendingBlockHeight()
	feeAddr := act.belong_trs_v3.GetAddress()

	// Check ID format
	if len(act.LendingID) != stores.BitcoinSyslendIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.BitcoinSyslendIdLength-1] == 0 {
		return fmt.Errorf("Bitcoin Lending Id format error.")
	}

	// Query whether the ID exists
	btclendObj, e := state.BitcoinSystemLending(act.LendingID)
	if e != nil {
		return e
	}
	if btclendObj == nil {
		return fmt.Errorf("Bitcoin Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// Check redemption status
	if btclendObj.IsRansomed.Check() {
		// Redeemed. Non redeemable
		return fmt.Errorf("Bitcoin Lending <%s> has been redeemed.", hex.EncodeToString(act.LendingID))
	}

	// Number of blocks in redemption period
	ransomBlockNumberBase := uint64(100000) // 十万个区块约一年
	if sys.TestDebugLocalDevelopmentMark {
		ransomBlockNumberBase = 10 // Test environment 10 blocks as cycle
	}

	// Calculate bitcoin redemption amount
	_, realRansomAmt, e4 := coinbase.CalculationBitcoinSystemLendingRedeemAmount(
		feeAddr, btclendObj.MainAddress, &btclendObj.LoanTotalAmount,
		ransomBlockNumberBase,
		uint64(btclendObj.CreateBlockHeight), paddingHeight,
	)
	if e4 != nil {
		return e4
	}

	// Check whether the redemption amount is valid (the redemption amount is really greater than the redeemable amount calculated in real time) and whether the redemption amount meets the requirements
	if act.RansomAmount.LessThan(realRansomAmt) {
		return fmt.Errorf("Ransom amount %s can not less than real time ransom amount %s.", act.RansomAmount.ToFinString(), realRansomAmt.ToFinString())
	}

	// Redemption operation, deducting HAC balance (so as to check whether the balance is sufficient first)
	e2 := DoSubBalanceFromChainStateV3(state, feeAddr, act.RansomAmount)
	if e2 != nil {
		return e2
	}

	// Increase bitcoin balance
	addSat := uint64(btclendObj.MortgageBitcoinPortion) * 100 * 10000
	e9 := DoAddSatoshiFromChainStateV3(state, feeAddr, fields.Satoshi(addSat))
	if e9 != nil {
		return e9
	}

	// Modify mortgage contract status
	e20 := btclendObj.SetRansomedStatus(paddingHeight, &act.RansomAmount, feeAddr)
	if e20 != nil {
		return e20
	}
	e10 := state.BitcoinLendingUpdate(act.LendingID, btclendObj)
	if e10 != nil {
		return e10
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Reduce the statistics of real-time bitcoin mortgage shares
	totalsupply.DoSubUint(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount,
		uint64(btclendObj.MortgageBitcoinPortion),
	)
	// Bitcoin system mortgage cumulative redemption destruction flow
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationRansomHacAmount,
		act.RansomAmount.ToMei(),
	)
	// Update statistics
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	// complete
	return nil
}

func (act *Action_18_BitcoinsSystemLendingRansom) WriteinChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	paddingHeight := state.GetPendingBlockHeight()
	feeAddr := act.belong_trs.GetAddress()

	// Check ID format
	if len(act.LendingID) != stores.BitcoinSyslendIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.BitcoinSyslendIdLength-1] == 0 {
		return fmt.Errorf("Bitcoin Lending Id format error.")
	}

	// Query whether the ID exists
	btclendObj, e := state.BitcoinSystemLending(act.LendingID)
	if e != nil {
		return e
	}
	if btclendObj == nil {
		return fmt.Errorf("Bitcoin Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// Check redemption status
	if btclendObj.IsRansomed.Check() {
		// Redeemed. Non redeemable
		return fmt.Errorf("Bitcoin Lending <%s> has been redeemed.", hex.EncodeToString(act.LendingID))
	}

	// Number of blocks in redemption period
	ransomBlockNumberBase := uint64(100000) // 十万个区块约一年
	if sys.TestDebugLocalDevelopmentMark {
		ransomBlockNumberBase = 10 // Test environment 10 blocks as cycle
	}

	// Calculate bitcoin redemption amount
	_, realRansomAmt, e4 := coinbase.CalculationBitcoinSystemLendingRedeemAmount(
		feeAddr, btclendObj.MainAddress, &btclendObj.LoanTotalAmount,
		ransomBlockNumberBase,
		uint64(btclendObj.CreateBlockHeight), paddingHeight,
	)
	if e4 != nil {
		return e4
	}

	// Check whether the redemption amount is valid (the redemption amount is really greater than the redeemable amount calculated in real time) and whether the redemption amount meets the requirements
	if act.RansomAmount.LessThan(realRansomAmt) {
		return fmt.Errorf("Ransom amount %s can not less than real time ransom amount %s.", act.RansomAmount.ToFinString(), realRansomAmt.ToFinString())
	}

	// Redemption operation, deducting HAC balance (so as to check whether the balance is sufficient first)
	e2 := DoSubBalanceFromChainState(state, feeAddr, act.RansomAmount)
	if e2 != nil {
		return e2
	}

	// Increase bitcoin balance
	addSat := uint64(btclendObj.MortgageBitcoinPortion) * 100 * 10000
	e9 := DoAddSatoshiFromChainState(state, feeAddr, fields.Satoshi(addSat))
	if e9 != nil {
		return e9
	}

	// Modify mortgage contract status
	e20 := btclendObj.SetRansomedStatus(paddingHeight, &act.RansomAmount, feeAddr)
	if e20 != nil {
		return e20
	}
	e10 := state.BitcoinLendingUpdate(act.LendingID, btclendObj)
	if e10 != nil {
		return e10
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Reduce the statistics of real-time bitcoin mortgage shares
	totalsupply.DoSubUint(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount,
		uint64(btclendObj.MortgageBitcoinPortion),
	)
	// Bitcoin system mortgage cumulative redemption destruction flow
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationRansomHacAmount,
		act.RansomAmount.ToMei(),
	)
	// Update statistics
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	// complete
	return nil
}

func (act *Action_18_BitcoinsSystemLendingRansom) RecoverChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// Query whether the ID exists
	btclendObj, e := state.BitcoinSystemLending(act.LendingID)
	if e != nil {
		return e
	}
	if btclendObj == nil {
		return fmt.Errorf("Bitcoin Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// Back out redemption operation and increase HAC balance
	e2 := DoAddBalanceFromChainState(state, feeAddr, act.RansomAmount)
	if e2 != nil {
		return e2
	}

	// Rollback to reduce bitcoin balance
	addSat := uint64(btclendObj.MortgageBitcoinPortion) * 100 * 10000
	e9 := DoSubSatoshiFromChainState(state, feeAddr, fields.Satoshi(addSat))
	if e9 != nil {
		return e9
	}

	// Modify mortgage contract status
	e20 := btclendObj.DropRansomedStatus() // 移除已赎回状态
	if e20 != nil {
		return e20
	}
	e10 := state.BitcoinLendingUpdate(act.LendingID, btclendObj)
	if e10 != nil {
		return e10
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Rollback to increase the statistics of real-time bitcoin mortgage copies
	totalsupply.DoAddUint(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount,
		uint64(btclendObj.MortgageBitcoinPortion),
	)
	// Bitcoin system mortgage cumulative redemption destruction daily return decrease
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationRansomHacAmount,
		act.RansomAmount.ToMei(),
	)
	// Update statistics
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	return nil
}

// Set belongs to long_ trs
func (act *Action_18_BitcoinsSystemLendingRansom) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}
func (act *Action_18_BitcoinsSystemLendingRansom) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_18_BitcoinsSystemLendingRansom) IsBurning90PersentTxFees() bool {
	return false
}
