package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
)

/*

用户之间系统借贷

*/

const ()

/*

创建用户见借贷流程：

. 检查合约ID格式
. 检查合约是否已经存在
. 检查钻石状态是否可以借贷，是否归属抵押者
. 检查比特币余额充足
. 检查贷出者余额
. 检查贷出者销毁利息

. 修改钻石抵押状态
. 扣除用户钻石余额
. 转移用户HAC余额
. 保存借贷合约

. 累计钻石、比特币借出流水
. 累计借出HAC额度
. 统计销毁1%利息

*/

// Inter user system loan
type Action_19_UsersLendingCreate struct {
	//
	LendingID fields.UserLendingId // Loan contract ID

	IsRedemptionOvertime    fields.Bool        // Whether it can be redeemed after expiration (automatic extension)
	IsPublicRedeemable      fields.Bool        // Public redeemable after maturity
	AgreedExpireBlockHeight fields.BlockHeight // Agreed expiration block height

	MortgagorAddress fields.Address // Address of mortgagor
	LenderAddress    fields.Address // Lender address

	MortgageBitcoin     fields.SatoshiVariation     // Mortgage bitcoin quantity unit: SAT
	MortgageDiamondList fields.DiamondListMaxLen200 // Mortgage diamond list

	LoanTotalAmount        fields.Amount // Total lending HAC limit
	AgreedRedemptionAmount fields.Amount // Agreed redemption amount

	PreBurningInterestAmount fields.Amount // Interest for pre destruction must be greater than or equal to 1% of the lending amount

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_19_UsersLendingCreate) Kind() uint16 {
	return 19
}

// json api
func (elm *Action_19_UsersLendingCreate) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_19_UsersLendingCreate) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	var b1, _ = elm.LendingID.Serialize()
	var b2, _ = elm.IsRedemptionOvertime.Serialize()
	var b3, _ = elm.IsPublicRedeemable.Serialize()
	var b4, _ = elm.AgreedExpireBlockHeight.Serialize()
	var b5, _ = elm.MortgagorAddress.Serialize()
	var b6, _ = elm.LenderAddress.Serialize()
	var b7, _ = elm.MortgageBitcoin.Serialize()
	var b8, _ = elm.MortgageDiamondList.Serialize()
	var b9, _ = elm.LoanTotalAmount.Serialize()
	var b10, _ = elm.AgreedRedemptionAmount.Serialize()
	var b11, _ = elm.PreBurningInterestAmount.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	buffer.Write(b7)
	buffer.Write(b8)
	buffer.Write(b9)
	buffer.Write(b10)
	buffer.Write(b11)
	return buffer.Bytes(), nil
}

func (elm *Action_19_UsersLendingCreate) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.LendingID.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.IsRedemptionOvertime.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.IsPublicRedeemable.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.AgreedExpireBlockHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.MortgagorAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LenderAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.MortgageBitcoin.Parse(buf, seek)
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
	seek, e = elm.AgreedRedemptionAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.PreBurningInterestAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_19_UsersLendingCreate) Size() uint32 {
	return 2 + elm.LendingID.Size() +
		elm.IsRedemptionOvertime.Size() +
		elm.IsPublicRedeemable.Size() +
		elm.AgreedExpireBlockHeight.Size() +
		elm.MortgagorAddress.Size() +
		elm.LenderAddress.Size() +
		elm.MortgageBitcoin.Size() +
		elm.MortgageDiamondList.Size() +
		elm.LoanTotalAmount.Size() +
		elm.AgreedRedemptionAmount.Size() +
		elm.PreBurningInterestAmount.Size()
}

func (act *Action_19_UsersLendingCreate) RequestSignAddresses() []fields.Address {
	return []fields.Address{
		act.MortgagorAddress,
		act.LenderAddress,
	} // Both the mortgagor and the lender need to sign
}

func (act *Action_19_UsersLendingCreate) WriteInChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 不能自己借给自己
	if act.MortgagorAddress.Equal(act.LenderAddress) {
		return fmt.Errorf("Cannot lending to myself.")
	}

	// Check amount length
	if len(act.LoanTotalAmount.Numeral) > 4 {
		return fmt.Errorf("Amount <%s> byte length is too long.", act.LoanTotalAmount.ToFinString())
	}
	if len(act.AgreedRedemptionAmount.Numeral) > 4 {
		return fmt.Errorf("Amount <%s> byte length is too long.", act.AgreedRedemptionAmount.ToFinString())
	}
	if len(act.PreBurningInterestAmount.Numeral) > 4 {
		return fmt.Errorf("Amount <%s> byte length is too long.", act.PreBurningInterestAmount.ToFinString())
	}
	// 借贷数额不能为空
	if act.LoanTotalAmount.IsEmpty() || act.AgreedRedemptionAmount.IsEmpty() || act.PreBurningInterestAmount.IsEmpty() {
		return fmt.Errorf("Amount cannot be empty.")
	}

	// block
	paddingHeight := state.GetPendingBlockHeight()

	// Check ID format
	if len(act.LendingID) != stores.UserLendingIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.UserLendingIdLength-1] == 0 {
		return fmt.Errorf("Diamond Lending Id format error.")
	}

	// Query whether the ID exists
	usrlendObj, e := state.UserLending(act.LendingID)
	if e != nil {
		return e
	}
	if usrlendObj != nil {
		return fmt.Errorf("User Lending <%s> already exist.", hex.EncodeToString(act.LendingID))
	}

	// Diamond quantity check
	dianum := int(act.MortgageDiamondList.Count)

	// Check the quantity of collateral
	if act.MortgageBitcoin.NotEmpty.Is(false) && dianum == 0 {
		return fmt.Errorf("Mortgage diamond and bitcoin cannot be empty at the same time")
	}

	// Check redemption period height
	effectiveExpireBlockHeight := paddingHeight + 288
	if sys.TestDebugLocalDevelopmentMark {
		effectiveExpireBlockHeight = paddingHeight + 10 // Test environment 10 blocks
	}
	if uint64(act.AgreedExpireBlockHeight) < effectiveExpireBlockHeight {
		// The agreed redemption period is at least after 288 blocks
		return fmt.Errorf("AgreedExpireBlockHeight %d is too short, must over than %d.", act.AgreedExpireBlockHeight, effectiveExpireBlockHeight)
	}

	if dianum != len(act.MortgageDiamondList.Diamonds) {
		return fmt.Errorf("Diamonds quantity error")
	}
	if dianum > 200 {
		return fmt.Errorf("Diamonds quantity cannot over 200")
	}

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
		// Check the diamond address
		if diaitem.Address.NotEqual(act.MortgagorAddress) {
			return fmt.Errorf("Diamond <%s> not belong to address '%s'", string(diamond), act.MortgagorAddress.ToReadable())
		}
		// Check whether it has been mortgaged and whether it can be mortgaged
		if diaitem.Status != stores.DiamondStatusNormal {
			return fmt.Errorf("Diamond <%s> has been mortgaged.", string(diamond))
		}
		// Mark mortgage diamond
		diaitem.Status = stores.DiamondStatusLendingOtherUser // Mark mortgage to other users
		e5 := state.DiamondSet(diamond, diaitem)
		if e5 != nil {
			return e5
		}
	}

	// 减少抵押人钻石余额
	e9 := DoSubDiamondFromChainStateV3(state, act.MortgagorAddress, fields.DiamondNumber(dianum))
	if e9 != nil {
		return e9
	}

	// 是否抵押比特币  扣除抵押人比特币余额
	if act.MortgageBitcoin.NotEmpty.Check() {
		// Deduct bitcoin
		e := DoSubSatoshiFromChainStateV3(state, act.MortgagorAddress, act.MortgageBitcoin.ValueSAT)
		if e != nil {
			return e // Deduction failed for insufficient bitcoin balance
		}
	}

	// Check the amount of interest destroyed
	mustBurnDesk := act.LoanTotalAmount.Copy()
	if mustBurnDesk.Unit > 2 {
		mustBurnDesk.Unit -= 2 // 1% destroy at least 1% of the loan quantity
	}
	if act.PreBurningInterestAmount.LessThan(mustBurnDesk) {
		// The destruction interest cannot be less than 1% of the loan amount
		return fmt.Errorf("PreBurningInterestAmount <%s> can not less than <%s>", act.PreBurningInterestAmount.ToFinString(), mustBurnDesk.ToFinString())
	}

	// Interest on destruction, paid by the Lender
	e10 := DoSubBalanceFromChainStateV3(state, act.LenderAddress, act.PreBurningInterestAmount)
	if e10 != nil {
		return e10 // Insufficient interest balance destroyed
	}

	// Mortgage succeeded, transfer balance: lender - > mortgage Borrower
	e11 := DoSimpleTransferFromChainStateV3(state, act.LenderAddress, act.MortgagorAddress, act.LoanTotalAmount)
	if e11 != nil {
		return e11 // Insufficient lender balance
	}

	// Save mortgage loan contract
	dlsto := &stores.UserLending{
		IsRansomed:               fields.CreateBool(false), // 标记未赎回
		IsRedemptionOvertime:     act.IsRedemptionOvertime,
		IsPublicRedeemable:       act.IsPublicRedeemable,
		CreateBlockHeight:        fields.BlockHeight(paddingHeight),
		ExpireBlockHeight:        act.AgreedExpireBlockHeight,
		MortgagorAddress:         act.MortgagorAddress,
		LenderAddress:            act.LenderAddress,
		MortgageBitcoin:          act.MortgageBitcoin,
		MortgageDiamondList:      act.MortgageDiamondList,
		LoanTotalAmount:          act.LoanTotalAmount,
		AgreedRedemptionAmount:   act.AgreedRedemptionAmount,
		PreBurningInterestAmount: act.PreBurningInterestAmount,
	}
	e12 := state.UserLendingCreate(act.LendingID, dlsto)
	if e12 != nil {
		return e12
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Increase diamond lending flow
	if dianum > 0 {
		totalsupply.DoAdd(
			stores.TotalSupplyStoreTypeOfUsersLendingCumulationDiamond,
			float64(dianum),
		)
	}
	// Increase bitcoin loan quantity flow
	if act.MortgageBitcoin.NotEmpty.Check() {
		totalsupply.DoAdd(
			stores.TotalSupplyStoreTypeOfUsersLendingCumulationBitcoin,
			float64(act.MortgageBitcoin.ValueSAT),
		)

	}
	// HAC flow of inter user loan amount
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfUsersLendingCumulationHacAmount,
		act.LoanTotalAmount.ToMei(),
	)
	// Pre destruction 1% interest accumulation
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfUsersLendingBurningOnePercentInterestHacAmount,
		act.PreBurningInterestAmount.ToMei(),
	)
	// Update statistics
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	// complete
	return nil
}

func (act *Action_19_UsersLendingCreate) WriteinChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 不能自己借给自己
	if act.MortgagorAddress.Equal(act.LenderAddress) {
		return fmt.Errorf("Cannot lending to myself.")
	}

	// Check amount length
	if len(act.LoanTotalAmount.Numeral) > 4 {
		return fmt.Errorf("Amount <%s> byte length is too long.", act.LoanTotalAmount.ToFinString())
	}
	if len(act.AgreedRedemptionAmount.Numeral) > 4 {
		return fmt.Errorf("Amount <%s> byte length is too long.", act.AgreedRedemptionAmount.ToFinString())
	}
	if len(act.PreBurningInterestAmount.Numeral) > 4 {
		return fmt.Errorf("Amount <%s> byte length is too long.", act.PreBurningInterestAmount.ToFinString())
	}
	// 借贷数额不能为空
	if act.LoanTotalAmount.IsEmpty() || act.AgreedRedemptionAmount.IsEmpty() || act.PreBurningInterestAmount.IsEmpty() {
		return fmt.Errorf("Amount cannot be empty.")
	}

	// block
	paddingHeight := state.GetPendingBlockHeight()

	// Check ID format
	if len(act.LendingID) != stores.UserLendingIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.UserLendingIdLength-1] == 0 {
		return fmt.Errorf("Diamond Lending Id format error.")
	}

	// Query whether the ID exists
	usrlendObj, e := state.UserLending(act.LendingID)
	if e != nil {
		return e
	}
	if usrlendObj != nil {
		return fmt.Errorf("User Lending <%s> already exist.", hex.EncodeToString(act.LendingID))
	}

	// Diamond quantity check
	dianum := int(act.MortgageDiamondList.Count)

	// Check the quantity of collateral
	if act.MortgageBitcoin.NotEmpty.Is(false) && dianum == 0 {
		return fmt.Errorf("Mortgage diamond and bitcoin cannot be empty at the same time")
	}

	// Check redemption period height
	effectiveExpireBlockHeight := paddingHeight + 288
	if sys.TestDebugLocalDevelopmentMark {
		effectiveExpireBlockHeight = paddingHeight + 10 // Test environment 10 blocks
	}
	if uint64(act.AgreedExpireBlockHeight) < effectiveExpireBlockHeight {
		// The agreed redemption period is at least after 288 blocks
		return fmt.Errorf("AgreedExpireBlockHeight %d is too short, must over than %d.", act.AgreedExpireBlockHeight, effectiveExpireBlockHeight)
	}

	if dianum != len(act.MortgageDiamondList.Diamonds) {
		return fmt.Errorf("Diamonds quantity error")
	}
	if dianum > 200 {
		return fmt.Errorf("Diamonds quantity cannot over 200")
	}

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
		// Check the diamond address
		if diaitem.Address.NotEqual(act.MortgagorAddress) {
			return fmt.Errorf("Diamond <%s> not belong to address '%s'", string(diamond), act.MortgagorAddress.ToReadable())
		}
		// Check whether it has been mortgaged and whether it can be mortgaged
		if diaitem.Status != stores.DiamondStatusNormal {
			return fmt.Errorf("Diamond <%s> has been mortgaged.", string(diamond))
		}
		// Mark mortgage diamond
		diaitem.Status = stores.DiamondStatusLendingOtherUser // Mark mortgage to other users
		e5 := state.DiamondSet(diamond, diaitem)
		if e5 != nil {
			return e5
		}
	}

	// 减少抵押人钻石余额
	e9 := DoSubDiamondFromChainState(state, act.MortgagorAddress, fields.DiamondNumber(dianum))
	if e9 != nil {
		return e9
	}

	// 是否抵押比特币  扣除抵押人比特币余额
	if act.MortgageBitcoin.NotEmpty.Check() {
		// Deduct bitcoin
		e := DoSubSatoshiFromChainState(state, act.MortgagorAddress, act.MortgageBitcoin.ValueSAT)
		if e != nil {
			return e // Deduction failed for insufficient bitcoin balance
		}
	}

	// Check the amount of interest destroyed
	mustBurnDesk := act.LoanTotalAmount.Copy()
	if mustBurnDesk.Unit > 2 {
		mustBurnDesk.Unit -= 2 // 1% destroy at least 1% of the loan quantity
	}
	if act.PreBurningInterestAmount.LessThan(mustBurnDesk) {
		// The destruction interest cannot be less than 1% of the loan amount
		return fmt.Errorf("PreBurningInterestAmount <%s> can not less than <%s>", act.PreBurningInterestAmount.ToFinString(), mustBurnDesk.ToFinString())
	}

	// Interest on destruction, paid by the Lender
	e10 := DoSubBalanceFromChainState(state, act.LenderAddress, act.PreBurningInterestAmount)
	if e10 != nil {
		return e10 // Insufficient interest balance destroyed
	}

	// Mortgage succeeded, transfer balance: lender - > mortgage Borrower
	e11 := DoSimpleTransferFromChainState(state, act.LenderAddress, act.MortgagorAddress, act.LoanTotalAmount)
	if e11 != nil {
		return e11 // Insufficient lender balance
	}

	// Save mortgage loan contract
	dlsto := &stores.UserLending{
		IsRansomed:               fields.CreateBool(false), // 标记未赎回
		IsRedemptionOvertime:     act.IsRedemptionOvertime,
		IsPublicRedeemable:       act.IsPublicRedeemable,
		CreateBlockHeight:        fields.BlockHeight(paddingHeight),
		ExpireBlockHeight:        act.AgreedExpireBlockHeight,
		MortgagorAddress:         act.MortgagorAddress,
		LenderAddress:            act.LenderAddress,
		MortgageBitcoin:          act.MortgageBitcoin,
		MortgageDiamondList:      act.MortgageDiamondList,
		LoanTotalAmount:          act.LoanTotalAmount,
		AgreedRedemptionAmount:   act.AgreedRedemptionAmount,
		PreBurningInterestAmount: act.PreBurningInterestAmount,
	}
	e12 := state.UserLendingCreate(act.LendingID, dlsto)
	if e12 != nil {
		return e12
	}

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Increase diamond lending flow
	if dianum > 0 {
		totalsupply.DoAdd(
			stores.TotalSupplyStoreTypeOfUsersLendingCumulationDiamond,
			float64(dianum),
		)
	}
	// Increase bitcoin loan quantity flow
	if act.MortgageBitcoin.NotEmpty.Check() {
		totalsupply.DoAdd(
			stores.TotalSupplyStoreTypeOfUsersLendingCumulationBitcoin,
			float64(act.MortgageBitcoin.ValueSAT),
		)

	}
	// HAC flow of inter user loan amount
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfUsersLendingCumulationHacAmount,
		act.LoanTotalAmount.ToMei(),
	)
	// Pre destruction 1% interest accumulation
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfUsersLendingBurningOnePercentInterestHacAmount,
		act.PreBurningInterestAmount.ToMei(),
	)
	// Update statistics
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	// complete
	return nil
}

func (act *Action_19_UsersLendingCreate) RecoverChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Diamond quantity
	dianum := int(act.MortgageDiamondList.Count)

	// Query whether the ID exists
	usrlendObj, e := state.UserLending(act.LendingID)
	if e != nil {
		return e
	}
	if usrlendObj == nil {
		return fmt.Errorf("User Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// Refund of bulk mortgage diamonds
	for i := 0; i < len(act.MortgageDiamondList.Diamonds); i++ {
		diamond := act.MortgageDiamondList.Diamonds[i]
		// Query whether the diamond exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		// Mark mortgage diamond
		diaitem.Status = stores.DiamondStatusNormal // Fallback status
		state.DiamondSet(diamond, diaitem)
	}

	// Increase in returned diamond balance
	if dianum > 0 {
		DoAddDiamondFromChainState(state, act.MortgagorAddress, fields.DiamondNumber(dianum))
	}

	// Back off deduction of bitcoin increase
	if act.MortgageBitcoin.NotEmpty.Check() {
		DoAddSatoshiFromChainState(state, act.MortgagorAddress, act.MortgageBitcoin.ValueSAT)
	}

	// Increase in interest for rollback destruction
	DoAddBalanceFromChainState(state, act.LenderAddress, act.PreBurningInterestAmount)

	// Mortgage succeeded, transfer balance: refund the lender < - Mortgage Borrower
	DoSimpleTransferFromChainState(state, act.MortgagorAddress, act.LenderAddress, act.LoanTotalAmount)

	// Delete mortgage loan contract
	state.UserLendingDelete(act.LendingID)

	// System statistics
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// Rebate deduction increase diamond lending amount daily deduction
	if dianum > 0 {
		totalsupply.DoSub(
			stores.TotalSupplyStoreTypeOfUsersLendingCumulationDiamond,
			float64(dianum),
		)
	}
	// Fallback deduction to increase bitcoin loan quantity flow
	if act.MortgageBitcoin.NotEmpty.Check() {
		totalsupply.DoSub(
			stores.TotalSupplyStoreTypeOfUsersLendingCumulationBitcoin,
			float64(act.MortgageBitcoin.ValueSAT),
		)

	}
	// Rollback deduction of inter user loan amount daily
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfUsersLendingCumulationHacAmount,
		act.LoanTotalAmount.ToMei(),
	)
	// Refund deducting 1% interest accumulation of pre destruction
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfUsersLendingBurningOnePercentInterestHacAmount,
		act.PreBurningInterestAmount.ToMei(),
	)
	// Update statistics
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	return nil
}

// Set belongs to long_ trs
func (act *Action_19_UsersLendingCreate) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}

func (act *Action_19_UsersLendingCreate) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_19_UsersLendingCreate) IsBurning90PersentTxFees() bool {
	return false
}

/////////////////////////////////////////////////

/*

用户间借贷，赎回流程

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
type Action_20_UsersLendingRansom struct {
	//
	LendingID    fields.UserLendingId // Loan contract ID
	RansomAmount fields.Amount        // Redemption amount

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_20_UsersLendingRansom) Kind() uint16 {
	return 20
}

// json api
func (elm *Action_20_UsersLendingRansom) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_20_UsersLendingRansom) Serialize() ([]byte, error) {
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

func (elm *Action_20_UsersLendingRansom) Parse(buf []byte, seek uint32) (uint32, error) {
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

func (elm *Action_20_UsersLendingRansom) Size() uint32 {
	return 2 +
		elm.LendingID.Size() +
		elm.RansomAmount.Size()
}

func (*Action_20_UsersLendingRansom) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_20_UsersLendingRansom) WriteInChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	paddingHeight := state.GetPendingBlockHeight()
	feeAddr := act.belong_trs_v3.GetAddress()

	// Check ID format
	if len(act.LendingID) != stores.UserLendingIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.UserLendingIdLength-1] == 0 {
		return fmt.Errorf("User Lending Id format error.")
	}

	// Query whether the ID exists
	usrlendObj, e := state.UserLending(act.LendingID)
	if e != nil {
		return e
	}
	if usrlendObj == nil {
		return fmt.Errorf("User Lending <%s> not exist.", act.LendingID.ToHex())
	}

	// Check redemption status
	if usrlendObj.IsRansomed.Check() {
		// Redeemed. Non redeemable
		return fmt.Errorf("User Lending <%s> has been redeemed.", act.LendingID.ToHex())
	}

	// Type of Redeemer
	isMortgagorDoRedeem := feeAddr.Equal(usrlendObj.MortgagorAddress)
	isLendersDoRedeem := feeAddr.Equal(usrlendObj.LenderAddress)
	isPublicDoRedeem := !isMortgagorDoRedeem && !isLendersDoRedeem
	// Whether it is within the mortgage period
	isWithinMortgageTime := paddingHeight <= uint64(usrlendObj.ExpireBlockHeight)

	// If public redemption is not opened, the third party cannot redeem
	if usrlendObj.IsPublicRedeemable.Is(false) && isPublicDoRedeem {
		// Public redemption closure
		return fmt.Errorf("Public redeemable not open.")
	}

	// Non mortgagor cannot redeem during mortgage period
	if isWithinMortgageTime && !isMortgagorDoRedeem {
		return fmt.Errorf("only %s can do redeem before height %d.", usrlendObj.MortgagorAddress.ToReadable(), usrlendObj.ExpireBlockHeight)
	}

	// Beyond the mortgage period

	// If automatic extension or public redemption is not started, the mortgagor cannot redeem after the expiration
	if usrlendObj.IsPublicRedeemable.Is(false) &&
		usrlendObj.IsRedemptionOvertime.Is(false) &&
		isMortgagorDoRedeem {
		// Outside the mortgage period, if there is no agreement on automatic extension and public redemption, the mortgagor cannot redeem
		return fmt.Errorf("only %s can do redeem after height %d.", usrlendObj.LenderAddress.ToReadable(), usrlendObj.ExpireBlockHeight)
	}

	// The lender withholds the collateral, and the redemption amount must be zero
	if isLendersDoRedeem && act.RansomAmount.IsNotEmpty() {
		return fmt.Errorf("Ransom amount must be zore but got %s with lender address %s do redeem", act.RansomAmount.ToFinString(), usrlendObj.LenderAddress.ToReadable())
	}

	// If it is withheld by the non lender and redeemed by the mortgagor or a third party, the redemption amount shall be checked
	if !isLendersDoRedeem && act.RansomAmount.LessThan(&usrlendObj.AgreedRedemptionAmount) {
		return fmt.Errorf("Ransom amount cannot less than %s but got %s.", usrlendObj.AgreedRedemptionAmount.ToFinString(), act.RansomAmount.ToFinString())
	}

	// Payment of ransom
	if isLendersDoRedeem {
		// Lender seizure
		// Without paying any ransom, directly detain
	} else {
		// Redemption by mortgagor or third party
		// Transfer HAC and pay ransom
		e2 := DoSimpleTransferFromChainStateV3(state, feeAddr, usrlendObj.LenderAddress, act.RansomAmount)
		if e2 != nil {
			return e2
		}
	}

	// Operational redemption or seizure
	dianum := usrlendObj.MortgageDiamondList.Count

	// Mass redemption or seizure of diamonds
	for i := 0; i < len(usrlendObj.MortgageDiamondList.Diamonds); i++ {
		diamond := usrlendObj.MortgageDiamondList.Diamonds[i]
		// Query whether the diamond exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		if diaitem == nil {
			return fmt.Errorf("diamond <%s> not find.", string(diamond))
		}
		// Check diamond home address
		if diaitem.Address.NotEqual(usrlendObj.MortgagorAddress) {
			return fmt.Errorf("diamond <%s> not belong to address %s", string(diamond), usrlendObj.MortgagorAddress.ToReadable())
		}
		// Check diamond status
		if diaitem.Status != stores.DiamondStatusLendingOtherUser {
			return fmt.Errorf("diamond <%s> status is not [stores.DiamondStatusLendingOtherUser].", string(diamond))
		}
		// Marker diamond
		diaitem.Status = stores.DiamondStatusNormal // Redemption diamond status
		diaitem.Address = feeAddr                   // Diamond ownership modification redemption, public redemption or seizure
		e5 := state.DiamondSet(diamond, diaitem)    // 更新钻石
		if e5 != nil {
			return e5
		}
	}

	// Increase diamond balance (Redeemer or lender)
	if dianum > 0 {
		e9 := DoAddDiamondFromChainStateV3(state, feeAddr, fields.DiamondNumber(dianum))
		if e9 != nil {
			return e9
		}
	}

	// Increase bitcoin balance
	if usrlendObj.MortgageBitcoin.NotEmpty.Check() {
		e10 := DoAddSatoshiFromChainStateV3(state, feeAddr, usrlendObj.MortgageBitcoin.ValueSAT)
		if e10 != nil {
			return e10
		}
	}

	// Modify mortgage contract status
	e13 := usrlendObj.SetRansomedStatus(paddingHeight, &act.RansomAmount, feeAddr) // 标记已经赎回，避免重复赎回
	if e13 != nil {
		return e13
	}
	e11 := state.UserLendingUpdate(act.LendingID, usrlendObj)
	if e11 != nil {
		return e11
	}

	// complete
	return nil
}

func (act *Action_20_UsersLendingRansom) WriteinChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	paddingHeight := state.GetPendingBlockHeight()
	feeAddr := act.belong_trs.GetAddress()

	// Check ID format
	if len(act.LendingID) != stores.UserLendingIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.UserLendingIdLength-1] == 0 {
		return fmt.Errorf("User Lending Id format error.")
	}

	// Query whether the ID exists
	usrlendObj, e := state.UserLending(act.LendingID)
	if e != nil {
		return e
	}
	if usrlendObj == nil {
		return fmt.Errorf("User Lending <%s> not exist.", act.LendingID.ToHex())
	}

	// Check redemption status
	if usrlendObj.IsRansomed.Check() {
		// Redeemed. Non redeemable
		return fmt.Errorf("User Lending <%s> has been redeemed.", act.LendingID.ToHex())
	}

	// Type of Redeemer
	isMortgagorDoRedeem := feeAddr.Equal(usrlendObj.MortgagorAddress)
	isLendersDoRedeem := feeAddr.Equal(usrlendObj.LenderAddress)
	isPublicDoRedeem := !isMortgagorDoRedeem && !isLendersDoRedeem
	// Whether it is within the mortgage period
	isWithinMortgageTime := paddingHeight <= uint64(usrlendObj.ExpireBlockHeight)

	// If public redemption is not opened, the third party cannot redeem
	if usrlendObj.IsPublicRedeemable.Is(false) && isPublicDoRedeem {
		// Public redemption closure
		return fmt.Errorf("Public redeemable not open.")
	}

	// Non mortgagor cannot redeem during mortgage period
	if isWithinMortgageTime && !isMortgagorDoRedeem {
		return fmt.Errorf("only %s can do redeem before height %d.", usrlendObj.MortgagorAddress.ToReadable(), usrlendObj.ExpireBlockHeight)
	}

	// Beyond the mortgage period

	// If automatic extension or public redemption is not started, the mortgagor cannot redeem after the expiration
	if usrlendObj.IsPublicRedeemable.Is(false) &&
		usrlendObj.IsRedemptionOvertime.Is(false) &&
		isMortgagorDoRedeem {
		// Outside the mortgage period, if there is no agreement on automatic extension and public redemption, the mortgagor cannot redeem
		return fmt.Errorf("only %s can do redeem after height %d.", usrlendObj.LenderAddress.ToReadable(), usrlendObj.ExpireBlockHeight)
	}

	// The lender withholds the collateral, and the redemption amount must be zero
	if isLendersDoRedeem && act.RansomAmount.IsNotEmpty() {
		return fmt.Errorf("Ransom amount must be zore but got %s with lender address %s do redeem", act.RansomAmount.ToFinString(), usrlendObj.LenderAddress.ToReadable())
	}

	// If it is withheld by the non lender and redeemed by the mortgagor or a third party, the redemption amount shall be checked
	if !isLendersDoRedeem && act.RansomAmount.LessThan(&usrlendObj.AgreedRedemptionAmount) {
		return fmt.Errorf("Ransom amount cannot less than %s but got %s.", usrlendObj.AgreedRedemptionAmount.ToFinString(), act.RansomAmount.ToFinString())
	}

	// Payment of ransom
	if isLendersDoRedeem {
		// Lender seizure
		// Without paying any ransom, directly detain
	} else {
		// Redemption by mortgagor or third party
		// Transfer HAC and pay ransom
		e2 := DoSimpleTransferFromChainState(state, feeAddr, usrlendObj.LenderAddress, act.RansomAmount)
		if e2 != nil {
			return e2
		}
	}

	// Operational redemption or seizure
	dianum := usrlendObj.MortgageDiamondList.Count

	// Mass redemption or seizure of diamonds
	for i := 0; i < len(usrlendObj.MortgageDiamondList.Diamonds); i++ {
		diamond := usrlendObj.MortgageDiamondList.Diamonds[i]
		// Query whether the diamond exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		if diaitem == nil {
			return fmt.Errorf("diamond <%s> not find.", string(diamond))
		}
		// Check diamond home address
		if diaitem.Address.NotEqual(usrlendObj.MortgagorAddress) {
			return fmt.Errorf("diamond <%s> not belong to address %s", string(diamond), usrlendObj.MortgagorAddress.ToReadable())
		}
		// Check diamond status
		if diaitem.Status != stores.DiamondStatusLendingOtherUser {
			return fmt.Errorf("diamond <%s> status is not [stores.DiamondStatusLendingOtherUser].", string(diamond))
		}
		// Marker diamond
		diaitem.Status = stores.DiamondStatusNormal // Redemption diamond status
		diaitem.Address = feeAddr                   // Diamond ownership modification redemption, public redemption or seizure
		e5 := state.DiamondSet(diamond, diaitem)    // 更新钻石
		if e5 != nil {
			return e5
		}
	}

	// Increase diamond balance (Redeemer or lender)
	if dianum > 0 {
		e9 := DoAddDiamondFromChainState(state, feeAddr, fields.DiamondNumber(dianum))
		if e9 != nil {
			return e9
		}
	}

	// Increase bitcoin balance
	if usrlendObj.MortgageBitcoin.NotEmpty.Check() {
		e10 := DoAddSatoshiFromChainState(state, feeAddr, usrlendObj.MortgageBitcoin.ValueSAT)
		if e10 != nil {
			return e10
		}
	}

	// Modify mortgage contract status
	e13 := usrlendObj.SetRansomedStatus(paddingHeight, &act.RansomAmount, feeAddr) // 标记已经赎回，避免重复赎回
	if e13 != nil {
		return e13
	}
	e11 := state.UserLendingUpdate(act.LendingID, usrlendObj)
	if e11 != nil {
		return e11
	}

	// complete
	return nil
}

func (act *Action_20_UsersLendingRansom) RecoverChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// Query whether the ID exists
	usrlendObj, e := state.UserLending(act.LendingID)
	if e != nil {
		return e
	}
	if usrlendObj == nil {
		return fmt.Errorf("User Lending <%s> not exist.", hex.EncodeToString(act.LendingID))
	}

	// Redemption type
	isLendersDoRedeem := feeAddr.Equal(usrlendObj.LenderAddress)

	// The borrower redeems as agreed
	if isLendersDoRedeem {
		// The lender is not required to pay any ransom
	} else {
		// Transfer HAC fallback HAC
		DoSimpleTransferFromChainState(state, usrlendObj.LenderAddress, usrlendObj.MortgagorAddress, act.RansomAmount)
	}

	// Operational redemption or seizure
	dianum := usrlendObj.MortgageDiamondList.Count

	// Mass redemption or seizure of diamonds
	for i := 0; i < len(usrlendObj.MortgageDiamondList.Diamonds); i++ {
		diamond := usrlendObj.MortgageDiamondList.Diamonds[i]
		// Query whether the diamond exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		// Marker diamond
		diaitem.Status = stores.DiamondStatusLendingOtherUser // Rollback diamond status
		diaitem.Address = usrlendObj.MortgagorAddress         // Return of diamond ownership to the mortgagor
		state.DiamondSet(diamond, diaitem)                    // Update diamond
	}

	// Decrease of diamond balance by rollback
	if dianum > 0 {
		DoSubDiamondFromChainState(state, feeAddr, fields.DiamondNumber(dianum))
	}

	// Rollback to reduce bitcoin balance
	if usrlendObj.MortgageBitcoin.NotEmpty.Check() {
		DoSubSatoshiFromChainState(state, feeAddr, usrlendObj.MortgageBitcoin.ValueSAT)
	}

	// Modify mortgage contract status
	usrlendObj.DropRansomedStatus() // Mark not redeemed, seized
	state.UserLendingUpdate(act.LendingID, usrlendObj)

	return nil
}

// Set belongs to long_ trs
func (act *Action_20_UsersLendingRansom) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}

func (act *Action_20_UsersLendingRansom) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_20_UsersLendingRansom) IsBurning90PersentTxFees() bool {
	return false
}
