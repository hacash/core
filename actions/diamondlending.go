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
)

const (
	// Number of credit period blocks
	DiamondsSystemLendingBorrowPeriodBlockNumber uint64 = 10000
)

/*

创建钻石借贷流程：

. 检查合约ID格式
. 检查合约是否已经存在
. 检查钻石状态是否可以借贷，是否归属抵押者
. 检查借贷HAC额度是否匹配
. 检查借贷时间周期 1~20

. 修改钻石抵押状态
. 扣除用户钻石余额
. 增加用户HAC余额
. 保存借贷合约

. 统计实时借贷数量
. 累计借出额度

*/

// Diamond system lending
type Action_15_DiamondsSystemLendingCreate struct {
	//
	LendingID           fields.DiamondSyslendId     // Loan contract ID
	MortgageDiamondList fields.DiamondListMaxLen200 // Mortgage diamond list
	LoanTotalAmount     fields.Amount               // The total lending HAC quota must be equal to the total borrowing quota, and cannot be more or less
	BorrowPeriod        fields.VarUint1             // Borrowing cycle: one cycle represents 0.5% interest and 10000 blocks for about 35 days, with the lowest one and the highest 20, then the annual interest rate is 5%

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_15_DiamondsSystemLendingCreate) Kind() uint16 {
	return 15
}

// json api
func (elm *Action_15_DiamondsSystemLendingCreate) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_15_DiamondsSystemLendingCreate) Serialize() ([]byte, error) {
	var e error = nil
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	b1, e := elm.LendingID.Serialize()
	if e != nil {
		return nil, e
	}
	b2, e := elm.MortgageDiamondList.Serialize()
	if e != nil {
		return nil, e
	}
	b3, e := elm.LoanTotalAmount.Serialize()
	if e != nil {
		return nil, e
	}
	b4, e := elm.BorrowPeriod.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	return buffer.Bytes(), nil
}

func (elm *Action_15_DiamondsSystemLendingCreate) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.LendingID.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.MortgageDiamondList.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LoanTotalAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.BorrowPeriod.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_15_DiamondsSystemLendingCreate) Size() uint32 {
	return 2 +
		elm.LendingID.Size() +
		elm.MortgageDiamondList.Size() +
		elm.LoanTotalAmount.Size() +
		elm.BorrowPeriod.Size()
}

func (*Action_15_DiamondsSystemLendingCreate) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_15_DiamondsSystemLendingCreate) WriteInChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs_v3.GetAddress()

	// Check ID format
	if len(act.LendingID) != stores.DiamondSyslendIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.DiamondSyslendIdLength-1] == 0 {
		return fmt.Errorf("Diamond Lending Id format error.")
	}

	// Query whether the ID exists
	dmdlendObj, e := state.DiamondSystemLending(act.LendingID)
	if e != nil {
		return e
	}
	if dmdlendObj != nil {
		return fmt.Errorf("Diamond Lending <%s> already exist.", hex.EncodeToString(act.LendingID))
	}

	// Quantity check
	dianum := int(act.MortgageDiamondList.Count)
	if dianum == 0 || dianum != len(act.MortgageDiamondList.Diamonds) {
		return fmt.Errorf("Diamonds quantity error")
	}
	if dianum > 200 {
		return fmt.Errorf("Diamonds quantity cannot over 200")
	}

	// Number of inspection cycles
	if act.BorrowPeriod < 1 || act.BorrowPeriod > 20 {
		return fmt.Errorf("BorrowPeriod must between 1 ~ 20")
	}

	// Loanable HAC
	totalLoanHAC := int64(0)

	// Bulk mortgage diamond
	for i := 0; i < len(act.MortgageDiamondList.Diamonds); i++ {
		diamond := act.MortgageDiamondList.Diamonds[i]

		// Query whether the diamond exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		if diaitem == nil {
			return fmt.Errorf("Diamond <%s> not find.", string(diamond))
		}
		// Check address
		if diaitem.Address.NotEqual(feeAddr) {
			return fmt.Errorf("Diamond <%s> not belong to address '%s'", string(diamond), feeAddr.ToReadable())
		}
		// Check the status of diamonds and whether they can be mortgaged
		if diaitem.Status != stores.DiamondStatusNormal {
			return fmt.Errorf("Diamond <%s> has been mortgaged and cannot be transferred.", string(diamond))
		}
		// Mark mortgage diamond
		diaitem.Status = stores.DiamondStatusLendingSystem // Mortgage to system
		e5 := state.DiamondSet(diamond, diaitem)
		if e5 != nil {
			return e5
		}
		diasmelt, e6 := state.BlockStore().ReadDiamond(diamond)
		if e6 != nil {
			return e5
		}
		if diasmelt == nil {
			return fmt.Errorf("Diamond <%s> not exist.", string(diamond))
		}
		// Count the quantity of lending HAC
		totalLoanHAC += int64(diasmelt.AverageBidBurnPrice)
	}

	// Total HAC pieces lent
	totalAmt := fields.NewAmountByUnit248(totalLoanHAC)
	// Verification quantity
	if totalAmt.NotEqual(&act.LoanTotalAmount) {
		return fmt.Errorf("LoanTotalAmountMei must %s but got %s", totalAmt.ToFinString(), act.LoanTotalAmount.ToFinString())
	}

	// Reduce diamond balance
	e9 := DoSubDiamondFromChainStateV3(state, feeAddr, fields.DiamondNumber(dianum))
	if e9 != nil {
		return e9
	}

	// Mortgage successful, HAC balance issued
	e10 := DoAddBalanceFromChainStateV3(state, feeAddr, act.LoanTotalAmount)
	if e10 != nil {
		return e10
	}

	// Preservation of diamond pledge
	paddingHei := state.GetPendingBlockHeight()
	dlsto := &stores.DiamondSystemLending{
		IsRansomed:          fields.CreateBool(false), // 标记未赎回
		CreateBlockHeight:   fields.BlockHeight(paddingHei),
		MainAddress:         feeAddr,
		MortgageDiamondList: act.MortgageDiamondList,
		LoanTotalAmountMei:  fields.VarUint4(totalLoanHAC),
		BorrowPeriod:        act.BorrowPeriod,
	}
	e11 := state.DiamondLendingCreate(act.LendingID, dlsto)
	if e11 != nil {
		return e11
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Increase mortgage statistics of real-time diamond system
	totalsupply.DoAddUint(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCurrentMortgageCount,
		uint64(dianum),
	)
	// Diamond system mortgage quantity statistics cumulative lending flow
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCumulationLoanHacAmount,
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

func (act *Action_15_DiamondsSystemLendingCreate) WriteinChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// Check ID format
	if len(act.LendingID) != stores.DiamondSyslendIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.DiamondSyslendIdLength-1] == 0 {
		return fmt.Errorf("Diamond Lending Id format error.")
	}

	// Query whether the ID exists
	dmdlendObj, e := state.DiamondSystemLending(act.LendingID)
	if e != nil {
		return e
	}
	if dmdlendObj != nil {
		return fmt.Errorf("Diamond Lending <%s> already exist.", hex.EncodeToString(act.LendingID))
	}

	// Quantity check
	dianum := int(act.MortgageDiamondList.Count)
	if dianum == 0 || dianum != len(act.MortgageDiamondList.Diamonds) {
		return fmt.Errorf("Diamonds quantity error")
	}
	if dianum > 200 {
		return fmt.Errorf("Diamonds quantity cannot over 200")
	}

	// Number of inspection cycles
	if act.BorrowPeriod < 1 || act.BorrowPeriod > 20 {
		return fmt.Errorf("BorrowPeriod must between 1 ~ 20")
	}

	// Loanable HAC
	totalLoanHAC := int64(0)

	// Bulk mortgage diamond
	for i := 0; i < len(act.MortgageDiamondList.Diamonds); i++ {
		diamond := act.MortgageDiamondList.Diamonds[i]

		// Query whether the diamond exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		if diaitem == nil {
			return fmt.Errorf("Diamond <%s> not find.", string(diamond))
		}
		// Check address
		if diaitem.Address.NotEqual(feeAddr) {
			return fmt.Errorf("Diamond <%s> not belong to address '%s'", string(diamond), feeAddr.ToReadable())
		}
		// Check the status of diamonds and whether they can be mortgaged
		if diaitem.Status != stores.DiamondStatusNormal {
			return fmt.Errorf("Diamond <%s> has been mortgaged and cannot be transferred.", string(diamond))
		}
		// Mark mortgage diamond
		diaitem.Status = stores.DiamondStatusLendingSystem // Mortgage to system
		e5 := state.DiamondSet(diamond, diaitem)
		if e5 != nil {
			return e5
		}
		diasmelt, e6 := state.BlockStore().ReadDiamond(diamond)
		if e6 != nil {
			return e5
		}
		if diasmelt == nil {
			return fmt.Errorf("Diamond <%s> not exist.", string(diamond))
		}
		// Count the quantity of lending HAC
		totalLoanHAC += int64(diasmelt.AverageBidBurnPrice)
	}

	// Total HAC pieces lent
	totalAmt := fields.NewAmountByUnit248(totalLoanHAC)
	// Verification quantity
	if totalAmt.NotEqual(&act.LoanTotalAmount) {
		return fmt.Errorf("LoanTotalAmountMei must %s but got %s", totalAmt.ToFinString(), act.LoanTotalAmount.ToFinString())
	}

	// Reduce diamond balance
	e9 := DoSubDiamondFromChainState(state, feeAddr, fields.DiamondNumber(dianum))
	if e9 != nil {
		return e9
	}

	// Mortgage successful, HAC balance issued
	e10 := DoAddBalanceFromChainState(state, feeAddr, act.LoanTotalAmount)
	if e10 != nil {
		return e10
	}

	// Preservation of diamond pledge
	paddingHei := state.GetPendingBlockHeight()
	dlsto := &stores.DiamondSystemLending{
		IsRansomed:          fields.CreateBool(false), // 标记未赎回
		CreateBlockHeight:   fields.BlockHeight(paddingHei),
		MainAddress:         feeAddr,
		MortgageDiamondList: act.MortgageDiamondList,
		LoanTotalAmountMei:  fields.VarUint4(totalLoanHAC),
		BorrowPeriod:        act.BorrowPeriod,
	}
	e11 := state.DiamondLendingCreate(act.LendingID, dlsto)
	if e11 != nil {
		return e11
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Increase mortgage statistics of real-time diamond system
	totalsupply.DoAddUint(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCurrentMortgageCount,
		uint64(dianum),
	)
	// Diamond system mortgage quantity statistics cumulative lending flow
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCumulationLoanHacAmount,
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

func (act *Action_15_DiamondsSystemLendingCreate) RecoverChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Roll back all mortgages
	feeAddr := act.belong_trs.GetAddress()

	// Bulk mortgage diamond
	for i := 0; i < len(act.MortgageDiamondList.Diamonds); i++ {
		diamond := act.MortgageDiamondList.Diamonds[i]
		// Query whether the diamond exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		// Mark mortgage diamond
		diaitem.Status = stores.DiamondStatusNormal // State recovery
		e5 := state.DiamondSet(diamond, diaitem)
		if e5 != nil {
			return e5
		}
	}

	// Return diamond balance
	dianum := act.MortgageDiamondList.Count
	e9 := DoAddDiamondFromChainState(state, feeAddr, fields.DiamondNumber(dianum))
	if e9 != nil {
		return e9
	}

	// Retrieve HAC balance
	e10 := DoSubBalanceFromChainState(state, feeAddr, act.LoanTotalAmount)
	if e10 != nil {
		return e10
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Reduce the statistical fallback of mortgage amount of real-time diamond system
	totalsupply.DoSubUint(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCurrentMortgageCount,
		uint64(dianum),
	)
	// Diamond system mortgage quantity statistics cumulative lending daily return
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCumulationLoanHacAmount,
		act.LoanTotalAmount.ToMei(),
	)
	// Update statistics
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	return nil
}

// Set belongs to long_ trs
func (act *Action_15_DiamondsSystemLendingCreate) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}
func (act *Action_15_DiamondsSystemLendingCreate) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_15_DiamondsSystemLendingCreate) IsBurning90PersentTxFees() bool {
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
type Action_16_DiamondsSystemLendingRansom struct {
	//
	LendingID    fields.DiamondSyslendId // Loan contract ID
	RansomAmount fields.Amount           // Redemption amount
	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_16_DiamondsSystemLendingRansom) Kind() uint16 {
	return 16
}

// json api
func (elm *Action_16_DiamondsSystemLendingRansom) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_16_DiamondsSystemLendingRansom) Serialize() ([]byte, error) {
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

func (elm *Action_16_DiamondsSystemLendingRansom) Parse(buf []byte, seek uint32) (uint32, error) {
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

func (elm *Action_16_DiamondsSystemLendingRansom) Size() uint32 {
	return 2 +
		elm.LendingID.Size() +
		elm.RansomAmount.Size()
}

func (*Action_16_DiamondsSystemLendingRansom) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_16_DiamondsSystemLendingRansom) WriteInChainState(state interfaces.ChainStateOperation) error {

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	// Lending period 10000 blocks about 35 days
	dslbpbn := DiamondsSystemLendingBorrowPeriodBlockNumber

	// Test use
	if sys.TestDebugLocalDevelopmentMark {
		dslbpbn = 10 // The test uses 50 blocks as a cycle
	}

	paddingHeight := state.GetPendingBlockHeight()
	feeAddr := act.belong_trs_v3.GetAddress()

	// Check ID format
	if len(act.LendingID) != stores.DiamondSyslendIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.DiamondSyslendIdLength-1] == 0 {
		return fmt.Errorf("Diamond Lending Id format error.")
	}

	// Query whether the ID exists
	dmdlendObj, e := state.DiamondSystemLending(act.LendingID)
	if e != nil {
		return e
	}
	if dmdlendObj == nil {
		return fmt.Errorf("Diamond Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// Check redemption status
	if dmdlendObj.IsRansomed.Check() {
		// Redeemed. Non redeemable
		return fmt.Errorf("Diamond Lending <%s> has been redeemed.", hex.EncodeToString(act.LendingID))
	}

	// Calculate the redemption period and the required redemption amount (judge whether it can be redeemed publicly)
	_, validRansomAmt, e4 := coinbase.CalculationDiamondSystemLendingRedeemAmount(
		feeAddr, dmdlendObj.MainAddress,
		int64(dmdlendObj.BorrowPeriod), int64(dmdlendObj.CreateBlockHeight),
		int64(dmdlendObj.LoanTotalAmountMei),
		int64(dslbpbn), int64(paddingHeight))
	if e4 != nil {
		return e4
	}

	// Check whether the redemption amount is valid (the redemption amount is really greater than the redeemable amount calculated in real time)
	if act.RansomAmount.LessThan(validRansomAmt) {
		return fmt.Errorf("Valid ransom amount must not less than %s but got %s", validRansomAmt.ToFinString(), act.RansomAmount.ToFinString())
	}

	// Redemption operation, deducting HAC balance (so as to check whether the balance is sufficient first)
	e2 := DoSubBalanceFromChainStateV3(state, feeAddr, act.RansomAmount)
	if e2 != nil {
		return e2
	}

	// Operational redemption
	dianum := dmdlendObj.MortgageDiamondList.Count

	// Bulk redemption of diamonds
	for i := 0; i < len(dmdlendObj.MortgageDiamondList.Diamonds); i++ {
		diamond := dmdlendObj.MortgageDiamondList.Diamonds[i]
		// Query whether the diamond exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		if diaitem == nil {
			return fmt.Errorf("diamond <%s> not find.", string(diamond))
		}
		// Check diamond status
		if diaitem.Status != stores.DiamondStatusLendingSystem {
			return fmt.Errorf("diamond <%s> status is not [stores.DiamondStatusLendingSystem].", string(diamond))
		}
		// Mark redemption diamond
		diaitem.Status = stores.DiamondStatusNormal // Redemption diamond status
		diaitem.Address = feeAddr                   // Diamond attribution modification
		e5 := state.DiamondSet(diamond, diaitem)    // 更新钻石
		if e5 != nil {
			return e5
		}
	}

	// Increase diamond balance
	e9 := DoAddDiamondFromChainStateV3(state, feeAddr, fields.DiamondNumber(dianum))
	if e9 != nil {
		return e9
	}

	// Modify the mortgage contract status, mark it as redeemed, and avoid repeated redemption
	e20 := dmdlendObj.SetRansomedStatus(paddingHeight, &act.RansomAmount, feeAddr)
	if e20 != nil {
		return e20
	}
	e10 := state.DiamondLendingUpdate(act.LendingID, dmdlendObj)
	if e10 != nil {
		return e10
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Reduce real-time diamond system mortgage quantity statistics and real-time deduction
	totalsupply.DoSubUint(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCurrentMortgageCount,
		uint64(dianum),
	)
	// Diamond system mortgage quantity statistics cumulative redemption flow
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCumulationRansomHacAmount,
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

func (act *Action_16_DiamondsSystemLendingRansom) WriteinChainState(state interfacev2.ChainStateOperation) error {

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	// Lending period 10000 blocks about 35 days
	dslbpbn := DiamondsSystemLendingBorrowPeriodBlockNumber

	// Test use
	if sys.TestDebugLocalDevelopmentMark {
		dslbpbn = 10 // The test uses 50 blocks as a cycle
	}

	paddingHeight := state.GetPendingBlockHeight()
	feeAddr := act.belong_trs.GetAddress()

	// Check ID format
	if len(act.LendingID) != stores.DiamondSyslendIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.DiamondSyslendIdLength-1] == 0 {
		return fmt.Errorf("Diamond Lending Id format error.")
	}

	// Query whether the ID exists
	dmdlendObj, e := state.DiamondSystemLending(act.LendingID)
	if e != nil {
		return e
	}
	if dmdlendObj == nil {
		return fmt.Errorf("Diamond Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// Check redemption status
	if dmdlendObj.IsRansomed.Check() {
		// Redeemed. Non redeemable
		return fmt.Errorf("Diamond Lending <%s> has been redeemed.", hex.EncodeToString(act.LendingID))
	}

	// Calculate the redemption period and the required redemption amount (judge whether it can be redeemed publicly)
	_, validRansomAmt, e4 := coinbase.CalculationDiamondSystemLendingRedeemAmount(
		feeAddr, dmdlendObj.MainAddress,
		int64(dmdlendObj.BorrowPeriod), int64(dmdlendObj.CreateBlockHeight),
		int64(dmdlendObj.LoanTotalAmountMei),
		int64(dslbpbn), int64(paddingHeight))
	if e4 != nil {
		return e4
	}

	// Check whether the redemption amount is valid (the redemption amount is really greater than the redeemable amount calculated in real time)
	if act.RansomAmount.LessThan(validRansomAmt) {
		return fmt.Errorf("Valid ransom amount must not less than %s but got %s", validRansomAmt.ToFinString(), act.RansomAmount.ToFinString())
	}

	// Redemption operation, deducting HAC balance (so as to check whether the balance is sufficient first)
	e2 := DoSubBalanceFromChainState(state, feeAddr, act.RansomAmount)
	if e2 != nil {
		return e2
	}

	// Operational redemption
	dianum := dmdlendObj.MortgageDiamondList.Count

	// Bulk redemption of diamonds
	for i := 0; i < len(dmdlendObj.MortgageDiamondList.Diamonds); i++ {
		diamond := dmdlendObj.MortgageDiamondList.Diamonds[i]
		// Query whether the diamond exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		if diaitem == nil {
			return fmt.Errorf("diamond <%s> not find.", string(diamond))
		}
		// Check diamond status
		if diaitem.Status != stores.DiamondStatusLendingSystem {
			return fmt.Errorf("diamond <%s> status is not [stores.DiamondStatusLendingSystem].", string(diamond))
		}
		// Mark redemption diamond
		diaitem.Status = stores.DiamondStatusNormal // Redemption diamond status
		diaitem.Address = feeAddr                   // Diamond attribution modification
		e5 := state.DiamondSet(diamond, diaitem)    // 更新钻石
		if e5 != nil {
			return e5
		}
	}

	// Increase diamond balance
	e9 := DoAddDiamondFromChainState(state, feeAddr, fields.DiamondNumber(dianum))
	if e9 != nil {
		return e9
	}

	// Modify the mortgage contract status, mark it as redeemed, and avoid repeated redemption
	e20 := dmdlendObj.SetRansomedStatus(paddingHeight, &act.RansomAmount, feeAddr)
	if e20 != nil {
		return e20
	}
	e10 := state.DiamondLendingUpdate(act.LendingID, dmdlendObj)
	if e10 != nil {
		return e10
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Reduce real-time diamond system mortgage quantity statistics and real-time deduction
	totalsupply.DoSubUint(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCurrentMortgageCount,
		uint64(dianum),
	)
	// Diamond system mortgage quantity statistics cumulative redemption flow
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCumulationRansomHacAmount,
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

func (act *Action_16_DiamondsSystemLendingRansom) RecoverChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// Roll back all redemptions
	dmdlendObj, e := state.DiamondSystemLending(act.LendingID)
	if e != nil {
		return e
	}
	if dmdlendObj == nil {
		return fmt.Errorf("Diamond Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// Return redemption status
	e1 := dmdlendObj.DropRansomedStatus()
	if e1 != nil {
		return e1
	}

	// Diamond batch recovery mortgage
	for i := 0; i < len(dmdlendObj.MortgageDiamondList.Diamonds); i++ {
		diamond := dmdlendObj.MortgageDiamondList.Diamonds[i]
		// Query whether the diamond exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		// Mark mortgage diamond
		diaitem.Status = stores.DiamondStatusLendingSystem // Status still mortgage
		diaitem.Address = dmdlendObj.MainAddress           // Diamonds still belong to the mortgagor
		e5 := state.DiamondSet(diamond, diaitem)
		if e5 != nil {
			return e5
		}
	}

	// Return diamond balance
	dianum := dmdlendObj.MortgageDiamondList.Count
	e9 := DoSubDiamondFromChainState(state, feeAddr, fields.DiamondNumber(dianum))
	if e9 != nil {
		return e9
	}

	// Withdrawal of HAc balance for redemption
	e10 := DoAddBalanceFromChainState(state, feeAddr, act.RansomAmount)
	if e10 != nil {
		return e10
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Increase the mortgage quantity statistics of real-time diamond system, increase and restore
	totalsupply.DoAddUint(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCurrentMortgageCount,
		uint64(dianum),
	)
	// Statistics of mortgage quantity of diamond system, cumulative redemption flow, decrease and refund
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfSystemLendingDiamondCumulationRansomHacAmount,
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
func (act *Action_16_DiamondsSystemLendingRansom) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}

func (act *Action_16_DiamondsSystemLendingRansom) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_16_DiamondsSystemLendingRansom) IsBurning90PersentTxFees() bool {
	return false
}
