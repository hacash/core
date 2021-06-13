package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
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

// 用户间系统借贷
type Action_19_UsersLendingCreate struct {
	//
	LendingID fields.Bytes17 // 借贷合约ID

	IsRedemptionOvertime    fields.Bool     // 是否超期仍可赎回（自动展期）
	IsPublicRedeemable      fields.Bool     // 到期后是否公共可赎回
	AgreedExpireBlockHeight fields.VarUint5 // 约定到期的区块高度

	MortgagorAddress fields.Address // 抵押人地址
	LendersAddress   fields.Address // 放款人地址

	MortgageBitcoin     fields.SatoshiVariation     // 抵押比特币数量 单位：SAT
	MortgageDiamondList fields.DiamondListMaxLen200 // 抵押钻石列表

	LoanTotalAmount        fields.Amount // 总共借出HAC额度
	AgreedRedemptionAmount fields.Amount // 约定的赎回金额

	PreBurningInterestAmount fields.Amount // 预先销毁的利息，必须大于等于 借出金额的 1%

	// data ptr
	belong_trs interfaces.Transaction
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
	var b6, _ = elm.LendersAddress.Serialize()
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
	seek, e = elm.LendersAddress.Parse(buf, seek)
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
		elm.LendersAddress.Size() +
		elm.MortgageBitcoin.Size() +
		elm.MortgageDiamondList.Size() +
		elm.LoanTotalAmount.Size() +
		elm.AgreedRedemptionAmount.Size() +
		elm.PreBurningInterestAmount.Size()
}

func (act *Action_19_UsersLendingCreate) RequestSignAddresses() []fields.Address {
	return []fields.Address{
		act.MortgagorAddress,
		act.LendersAddress,
	} // 抵押人和放款人都需要签名
}

func (act *Action_19_UsersLendingCreate) WriteinChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 不能自己借给自己
	if act.MortgagorAddress.Equal(act.LendersAddress) {
		return fmt.Errorf("Cannot lending to myself.")
	}

	// 检查数额长度
	if len(act.LoanTotalAmount.Numeral) > 4 {
		return fmt.Errorf("Amount <%s> byte length is too long.", act.LoanTotalAmount.ToFinString())
	}
	if len(act.AgreedRedemptionAmount.Numeral) > 4 {
		return fmt.Errorf("Amount <%s> byte length is too long.", act.AgreedRedemptionAmount.ToFinString())
	}
	if len(act.PreBurningInterestAmount.Numeral) > 4 {
		return fmt.Errorf("Amount <%s> byte length is too long.", act.PreBurningInterestAmount.ToFinString())
	}

	// 区块
	paddingHeight := state.GetPendingBlockHeight()

	// 检查id格式
	if len(act.LendingID) != stores.UserLendingIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.UserLendingIdLength-1] == 0 {
		return fmt.Errorf("Diamond Lending Id format error.")
	}

	// 查询id是否存在
	usrlendObj := state.UserLending(act.LendingID)
	if usrlendObj != nil {
		return fmt.Errorf("User Lending <%d> already exist.", hex.EncodeToString(act.LendingID))
	}

	// 钻石数量检查
	dianum := int(act.MortgageDiamondList.Count)

	// 检查抵押物数量
	if act.MortgageBitcoin.NotEmpty.Is(false) && dianum == 0 {
		return fmt.Errorf("Mortgage diamond and bitcoin cannot be empty at the same time")
	}

	// 检查赎回期限高度
	effectiveExpireBlockHeight := paddingHeight + 288
	if sys.TestDebugLocalDevelopmentMark {
		effectiveExpireBlockHeight = paddingHeight + 10 // 测试环境10个区块
	}
	if uint64(act.AgreedExpireBlockHeight) < effectiveExpireBlockHeight {
		// 约定赎回期至少在288个区块以后
		return fmt.Errorf("AgreedExpireBlockHeight %d is too short, must over than %d.", act.AgreedExpireBlockHeight, effectiveExpireBlockHeight)
	}

	if dianum != len(act.MortgageDiamondList.Diamonds) {
		return fmt.Errorf("Diamonds quantity error")
	}
	if dianum > 200 {
		return fmt.Errorf("Diamonds quantity cannot over 200")
	}

	// 批量抵押钻石
	for i := 0; i < len(act.MortgageDiamondList.Diamonds); i++ {
		diamond := act.MortgageDiamondList.Diamonds[i]
		// 查询钻石是否存在
		diaitem := state.Diamond(diamond)
		if diaitem == nil {
			return fmt.Errorf("Diamond <%s> not find.", string(diamond))
		}
		// 检查所属
		if diaitem.Address.NotEqual(act.MortgagorAddress) {
			return fmt.Errorf("Diamond <%s> not belong to address '%s'", string(diamond), act.MortgagorAddress.ToReadable())
		}
		// 检查是否已经抵押，是否可以抵押
		if diaitem.Status != stores.DiamondStatusNormal {
			return fmt.Errorf("Diamond <%s> has been mortgaged.", string(diamond))
		}
		// 标记抵押钻石
		diaitem.Status = stores.DiamondStatusLendingOtherUser // 标记抵押给其它用户
		e5 := state.DiamondSet(diamond, diaitem)
		if e5 != nil {
			return e5
		}
	}

	// 减少抵押人钻石余额
	e9 := DoSubDiamondFromChainState(state, act.MortgagorAddress, fields.VarUint3(dianum))
	if e9 != nil {
		return e9
	}

	// 是否抵押比特币  扣除抵押人比特币余额
	if act.MortgageBitcoin.NotEmpty.Check() {
		// 扣除比特币
		e := DoSubSatoshiFromChainState(state, act.MortgagorAddress, act.MortgageBitcoin.ValueSAT)
		if e != nil {
			return e
		}
	}

	// 检查销毁利息数额
	mustBurnDesk := act.LoanTotalAmount.Copy()
	if mustBurnDesk.Unit > 2 {
		mustBurnDesk.Unit -= 2 // 1% 销毁最小为借贷数量的 1%
	}
	if act.PreBurningInterestAmount.LessThan(mustBurnDesk) {
		// 销毁利息不能少于借贷数额的 1%
		return fmt.Errorf("PreBurningInterestAmount <%s> can not less than <%s>", act.PreBurningInterestAmount.ToFinString(), mustBurnDesk.ToFinString())
	}

	// 销毁利息，由放款人支付
	e10 := DoSubBalanceFromChainState(state, act.LendersAddress, act.PreBurningInterestAmount)
	if e10 != nil {
		return e10
	}

	// 抵押成功，转移余额：  放款人 -> 抵押借款者
	e11 := DoSimpleTransferFromChainState(state, act.LendersAddress, act.MortgagorAddress, act.LoanTotalAmount)
	if e11 != nil {
		return e11
	}

	// 保存抵押借贷合约
	dlsto := &stores.UserLending{
		IsRansomed:               fields.CreateBool(false), // 标记未赎回
		IsRedemptionOvertime:     act.IsRedemptionOvertime,
		IsPublicRedeemable:       act.IsPublicRedeemable,
		CreateBlockHeight:        fields.VarUint5(paddingHeight),
		ExpireBlockHeight:        act.AgreedExpireBlockHeight,
		MortgagorAddress:         act.MortgagorAddress,
		LendersAddress:           act.LendersAddress,
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

	// 系统统计
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// 增加钻石借贷数量流水
	if dianum > 0 {
		totalsupply.DoAdd(
			stores.TotalSupplyStoreTypeOfUsersLendingCumulationDiamond,
			float64(dianum),
		)
	}
	// 增加比特币借贷数量流水
	if act.MortgageBitcoin.NotEmpty.Check() {
		totalsupply.DoAdd(
			stores.TotalSupplyStoreTypeOfUsersLendingCumulationBitcoin,
			float64(act.MortgageBitcoin.ValueSAT),
		)

	}
	// 用户间借贷额 HAC 流水
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfUsersLendingCumulationHacAmount,
		act.LoanTotalAmount.ToMei(),
	)
	// 预先销毁 1% 利息流水
	totalsupply.DoAdd(
		stores.TotalSupplyStoreTypeOfUsersLendingBurningOnePercentInterestHacAmount,
		act.PreBurningInterestAmount.ToMei(),
	)
	// 更新统计
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	// 完毕
	return nil
}

func (act *Action_19_UsersLendingCreate) RecoverChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 钻石数量
	dianum := int(act.MortgageDiamondList.Count)

	// 查询id是否存在
	usrlendObj := state.UserLending(act.LendingID)
	if usrlendObj == nil {
		return fmt.Errorf("User Lending <%d> not exist.", hex.EncodeToString(act.LendingID))
	}

	// 回退 批量抵押钻石
	for i := 0; i < len(act.MortgageDiamondList.Diamonds); i++ {
		diamond := act.MortgageDiamondList.Diamonds[i]
		// 查询钻石是否存在
		diaitem := state.Diamond(diamond)
		// 标记抵押钻石
		diaitem.Status = stores.DiamondStatusNormal // 回退状态
		state.DiamondSet(diamond, diaitem)
	}

	// 回退 钻石余额 增加
	DoAddDiamondFromChainState(state, act.MortgagorAddress, fields.VarUint3(dianum))

	// 回退  扣除比特币 增加
	DoAddSatoshiFromChainState(state, act.MortgagorAddress, act.MortgageBitcoin.ValueSAT)

	// 回退 销毁利息 增加
	DoAddBalanceFromChainState(state, act.LendersAddress, act.PreBurningInterestAmount)

	// 抵押成功，转移余额： 回退  放款人 <- 抵押借款者
	DoSimpleTransferFromChainState(state, act.MortgagorAddress, act.LendersAddress, act.LoanTotalAmount)

	// 删除抵押借贷合约
	state.UserLendingDelete(act.LendingID)

	// 系统统计
	totalsupply, e20 := state.ReadTotalSupply()
	if e20 != nil {
		return e20
	}
	// 回退扣除  增加钻石借贷数量流水 扣除
	if dianum > 0 {
		totalsupply.DoSub(
			stores.TotalSupplyStoreTypeOfUsersLendingCumulationDiamond,
			float64(dianum),
		)
	}
	// 回退扣除   增加比特币借贷数量流水
	if act.MortgageBitcoin.NotEmpty.Check() {
		totalsupply.DoSub(
			stores.TotalSupplyStoreTypeOfUsersLendingCumulationBitcoin,
			float64(act.MortgageBitcoin.ValueSAT),
		)

	}
	// 回退扣除   用户间借贷额流水
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfUsersLendingCumulationHacAmount,
		act.LoanTotalAmount.ToMei(),
	)
	// 回退扣除   预先销毁 1% 利息流水
	totalsupply.DoSub(
		stores.TotalSupplyStoreTypeOfUsersLendingBurningOnePercentInterestHacAmount,
		act.PreBurningInterestAmount.ToMei(),
	)
	// 更新统计
	e21 := state.UpdateSetTotalSupply(totalsupply)
	if e21 != nil {
		return e21
	}

	return nil
}

// 设置所属 belong_trs
func (act *Action_19_UsersLendingCreate) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
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

// 钻石系统借贷，赎回
type Action_20_UsersLendingRansom struct {
	//
	LendingID    fields.Bytes17 // 借贷合约ID
	RansomAmount fields.Amount  // 赎回金额

	// data ptr
	belong_trs interfaces.Transaction
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

func (act *Action_20_UsersLendingRansom) WriteinChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	paddingHeight := state.GetPendingBlockHeight()
	feeAddr := act.belong_trs.GetAddress()

	// 检查id格式
	if len(act.LendingID) != stores.UserLendingIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.UserLendingIdLength-1] == 0 {
		return fmt.Errorf("User Lending Id format error.")
	}

	// 查询id是否存在
	usrlendObj := state.UserLending(act.LendingID)
	if usrlendObj == nil {
		return fmt.Errorf("User Lending <%d> not exist.", hex.EncodeToString(act.LendingID))
	}

	// 检查是否赎回状态
	if usrlendObj.IsRansomed.Check() {
		// 已经赎回。不可再次赎回
		return fmt.Errorf("User Lending <%d> has been redeemed.", hex.EncodeToString(act.LendingID))
	}

	// 赎回人类型
	isMortgagorDoRedeem := feeAddr.Equal(usrlendObj.MortgagorAddress)
	isLendersDoRedeem := feeAddr.Equal(usrlendObj.LendersAddress)
	isPublicDoRedeem := !isMortgagorDoRedeem && !isLendersDoRedeem
	// 是否处于抵押期内
	isWithinMortgageTime := paddingHeight <= uint64(usrlendObj.ExpireBlockHeight)

	// 抵押期内非抵押人不能赎回
	if isWithinMortgageTime && !isMortgagorDoRedeem {
		return fmt.Errorf("only %s can do redeem before height %d.", usrlendObj.MortgagorAddress.ToReadable(), usrlendObj.ExpireBlockHeight)
	}

	// 超出抵押期限

	// 未开启公开赎回，则第三方不能赎回
	if usrlendObj.IsPublicRedeemable.Is(false) && isPublicDoRedeem {
		// 公共赎回关闭
		return fmt.Errorf("Public redeemable not open.")
	}

	// 未开启自动展期和公共赎回，则抵押人超期不能赎回
	if usrlendObj.IsPublicRedeemable.Is(false) &&
		usrlendObj.IsRedemptionOvertime.Is(false) &&
		isMortgagorDoRedeem {
		// 抵押期外，没有约定自动展期和公共可赎回，则抵押人不能赎回
		return fmt.Errorf("only %s can do redeem after height %d.", usrlendObj.LendersAddress.ToReadable(), usrlendObj.ExpireBlockHeight)

	}

	// 放款人扣留抵押品，赎回金额必须为零
	if isLendersDoRedeem && act.RansomAmount.IsEmpty() == false {
		return fmt.Errorf("Ransom amount must be zore but got %s with lender address %s do redeem", act.RansomAmount.ToFinString(), usrlendObj.LendersAddress.ToReadable())
	}

	// 非放款人扣留，则检查赎回金额
	if !isLendersDoRedeem && act.RansomAmount.LessThan(&usrlendObj.AgreedRedemptionAmount) {
		return fmt.Errorf("Ransom amount cannot less than %s but got %s.", usrlendObj.AgreedRedemptionAmount.ToFinString(), act.RansomAmount.ToFinString())
	}

	// 借款人或第三方 赎回
	if !isLendersDoRedeem {
		// 转移 HAC，支付赎金
		e2 := DoSimpleTransferFromChainState(state, feeAddr, usrlendObj.LendersAddress, act.RansomAmount)
		if e2 != nil {
			return e2
		}
	}

	// 放款人扣押
	if isLendersDoRedeem {
		// 无需支付任何赎金
	}

	// 操作赎回或扣押
	dianum := usrlendObj.MortgageDiamondList.Count

	// 批量赎回或扣押钻石
	for i := 0; i < len(usrlendObj.MortgageDiamondList.Diamonds); i++ {
		diamond := usrlendObj.MortgageDiamondList.Diamonds[i]
		// 查询钻石是否存在
		diaitem := state.Diamond(diamond)
		if diaitem == nil {
			return fmt.Errorf("diamond <%s> not find.", string(diamond))
		}
		// 检查钻石归属地址
		if diaitem.Address.NotEqual(usrlendObj.MortgagorAddress) {
			return fmt.Errorf("diamond <%s> not belong to address %s", usrlendObj.MortgagorAddress.ToReadable())
		}
		// 检查钻石状态
		if diaitem.Status != stores.DiamondStatusLendingOtherUser {
			return fmt.Errorf("diamond <%s> status is not [stores.DiamondStatusLendingOtherUser].", string(diamond))
		}
		// 标记钻石
		diaitem.Status = stores.DiamondStatusNormal // 赎回钻石状态
		diaitem.Address = feeAddr                   // 钻石 归属 修改 赎回、公共赎回或扣押
		e5 := state.DiamondSet(diamond, diaitem)    // 更新钻石
		if e5 != nil {
			return e5
		}
	}

	// 增加钻石余额（赎回人或放款人）
	if dianum > 0 {
		e9 := DoAddDiamondFromChainState(state, feeAddr, fields.VarUint3(dianum))
		if e9 != nil {
			return e9
		}
	}

	// 增加比特币余额
	if usrlendObj.MortgageBitcoin.NotEmpty.Check() {
		e10 := DoAddSatoshiFromChainState(state, feeAddr, usrlendObj.MortgageBitcoin.ValueSAT)
		if e10 != nil {
			return e10
		}
	}

	// 修改抵押合约状态
	usrlendObj.IsRansomed.Set(true) // 标记已经赎回，避免重复赎回
	e11 := state.UserLendingUpdate(act.LendingID, usrlendObj)
	if e11 != nil {
		return e11
	}

	// 完毕
	return nil
}

func (act *Action_20_UsersLendingRansom) RecoverChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// 查询id是否存在
	usrlendObj := state.UserLending(act.LendingID)
	if usrlendObj == nil {
		return fmt.Errorf("User Lending <%d> not exist.", hex.EncodeToString(act.LendingID))
	}

	// 赎回类型
	isLendersDoRedeem := feeAddr.Equal(usrlendObj.LendersAddress)

	// 借款人按约定赎回
	if !isLendersDoRedeem {
		// 转移 HAC  回退HAC
		DoSimpleTransferFromChainState(state, usrlendObj.LendersAddress, usrlendObj.MortgagorAddress, act.RansomAmount)
	}

	// 放款人扣押
	if isLendersDoRedeem {
		// 无需支付任何赎金
	}

	// 操作赎回或扣押
	dianum := usrlendObj.MortgageDiamondList.Count

	// 批量赎回或扣押钻石
	for i := 0; i < len(usrlendObj.MortgageDiamondList.Diamonds); i++ {
		diamond := usrlendObj.MortgageDiamondList.Diamonds[i]
		// 查询钻石是否存在
		diaitem := state.Diamond(diamond)
		// 标记钻石
		diaitem.Status = stores.DiamondStatusLendingOtherUser // 回退钻石状态
		diaitem.Address = usrlendObj.MortgagorAddress         // 钻石 归属 回退至抵押人
		state.DiamondSet(diamond, diaitem)                    // 更新钻石
	}

	// 回退 减少钻石余额
	if dianum > 0 {
		DoSubDiamondFromChainState(state, feeAddr, fields.VarUint3(dianum))
	}

	// 回退减少比特币余额
	if usrlendObj.MortgageBitcoin.NotEmpty.Check() {
		DoSubSatoshiFromChainState(state, feeAddr, usrlendObj.MortgageBitcoin.ValueSAT)
	}

	// 修改抵押合约状态
	usrlendObj.IsRansomed.Set(false) // 标记未赎回或扣押
	state.UserLendingUpdate(act.LendingID, usrlendObj)

	return nil
}

// 设置所属 belong_trs
func (act *Action_20_UsersLendingRansom) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_20_UsersLendingRansom) IsBurning90PersentTxFees() bool {
	return false
}
