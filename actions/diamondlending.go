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
	"math/big"
)

// 钻石系统借贷
type Action_15_DiamondsSystemLendingCreate struct {
	//
	LendingID           fields.Bytes14              // 借贷合约ID
	MortgageDiamondList fields.DiamondListMaxLen200 // 抵押钻石列表
	LoanTotalAmount     fields.Amount               // 总共借出HAC额度，必须等于总可借额度，不能多也不能少
	BorrowPeriod        fields.VarUint1             // 借款周期，一个周期代表 0.5%利息和10000个区块约35天，最低1最高20

	// data ptr
	belong_trs interfaces.Transaction
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
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	var b1, _ = elm.LendingID.Serialize()
	var b2, _ = elm.MortgageDiamondList.Serialize()
	var b3, _ = elm.LoanTotalAmount.Serialize()
	var b4, _ = elm.BorrowPeriod.Serialize()
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

func (act *Action_15_DiamondsSystemLendingCreate) WriteinChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// 检查id格式
	if len(act.LendingID) != stores.DiamondLendingIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.DiamondLendingIdLength-1] == 0 {
		return fmt.Errorf("Diamond Lending Id format error.")
	}

	// 查询id是否存在
	dmdlendObj := state.DiamondLending(act.LendingID)
	if dmdlendObj != nil {
		return fmt.Errorf("Diamond Lending <%d> already exist.", hex.EncodeToString(act.LendingID))
	}

	// 数量检查
	dianum := int(act.MortgageDiamondList.Count)
	if dianum == 0 || dianum != len(act.MortgageDiamondList.Diamonds) {
		return fmt.Errorf("Diamonds quantity error")
	}
	if dianum > 200 {
		return fmt.Errorf("Diamonds quantity cannot over 200")
	}

	// 检查周期数
	if act.BorrowPeriod < 1 || act.BorrowPeriod > 20 {
		return fmt.Errorf("BorrowPeriod must between 1 ~ 20")
	}

	// 可借出HAC
	totalLoanHAC := int64(0)

	// 批量抵押钻石
	for i := 0; i < len(act.MortgageDiamondList.Diamonds); i++ {
		diamond := act.MortgageDiamondList.Diamonds[i]

		// 查询钻石是否存在
		diaitem := state.Diamond(diamond)
		if diaitem == nil {
			return fmt.Errorf("Diamond <%s> not exist.", string(diamond))
		}
		item := diaitem
		// 检查是否已经抵押，是否可以抵押
		if diaitem.Status != stores.DiamondStatusNormal {
			return fmt.Errorf("Diamond <%s> has been mortgaged and cannot be transferred.", string(diamond))
		}
		// 检查所属
		if bytes.Compare(item.Address, feeAddr) != 0 {
			return fmt.Errorf("Diamond <%s> not belong to address '%s'", string(diamond), feeAddr.ToReadable())
		}
		// 标记抵押钻石
		item.Status = stores.DiamondStatusLendingSystem // 抵押给系统
		e5 := state.DiamondSet(diamond, item)
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

		// 统计可借出HAC数量
		totalLoanHAC += int64(diasmelt.AverageBidBurnPrice)
	}

	// 共借出 HAC 枚
	totalAmt, e8 := fields.NewAmountByBigIntWithUnit(big.NewInt(totalLoanHAC), 248)
	if e8 != nil {
		return e8
	}
	// 验证数量
	if totalAmt.Equal(&act.LoanTotalAmount) == false {
		return fmt.Errorf("LoanTotalAmount <%s> and <%s> not match.", totalAmt.ToFinString(), act.LoanTotalAmount.ToFinString())
	}

	// 减少钻石余额
	e9 := DoSubDiamondFromChainState(state, feeAddr, fields.VarUint3(dianum))
	if e9 != nil {
		return e9
	}

	// 抵押成功，发放余额
	e10 := DoAddBalanceFromChainState(state, feeAddr, act.LoanTotalAmount)
	if e10 != nil {
		return e10
	}

	// 完毕
	return nil
}

func (act *Action_15_DiamondsSystemLendingCreate) RecoverChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 回退所有抵押
	feeAddr := act.belong_trs.GetAddress()

	// 批量抵押钻石
	for i := 0; i < len(act.MortgageDiamondList.Diamonds); i++ {
		diamond := act.MortgageDiamondList.Diamonds[i]
		// 查询钻石是否存在
		diaitem := state.Diamond(diamond)
		// 标记抵押钻石
		diaitem.Status = stores.DiamondStatusNormal // 恢复
		e5 := state.DiamondSet(diamond, diaitem)
		if e5 != nil {
			return e5
		}
	}

	// 回退钻石余额
	dianum := act.MortgageDiamondList.Count
	e9 := DoAddDiamondFromChainState(state, feeAddr, fields.VarUint3(dianum))
	if e9 != nil {
		return e9
	}

	// 取回HAC余额
	e10 := DoSubBalanceFromChainState(state, feeAddr, act.LoanTotalAmount)
	if e10 != nil {
		return e10
	}

	return nil
}

// 设置所属 belong_trs
func (act *Action_15_DiamondsSystemLendingCreate) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_15_DiamondsSystemLendingCreate) IsBurning90PersentTxFees() bool {
	return false
}

/////////////////////////////////////////////////

// 钻石系统借贷，赎回
type Action_16_DiamondsSystemLendingRansom struct {
	//
	LendingID           fields.Bytes14              // 借贷合约ID
	MortgageDiamondList fields.DiamondListMaxLen200 // 抵押钻石列表
	LoanTotalAmount     fields.Amount               // 总共借出HAC额度，必须等于总可借额度，不能多也不能少
	BorrowPeriod        fields.VarUint1             // 借款周期，一个周期代表 0.5%利息和10000个区块约35天，最低1最高20

	// data ptr
	belong_trs interfaces.Transaction
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
	var b2, _ = elm.MortgageDiamondList.Serialize()
	var b3, _ = elm.LoanTotalAmount.Serialize()
	var b4, _ = elm.BorrowPeriod.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	return buffer.Bytes(), nil
}

func (elm *Action_16_DiamondsSystemLendingRansom) Parse(buf []byte, seek uint32) (uint32, error) {
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

func (elm *Action_16_DiamondsSystemLendingRansom) Size() uint32 {
	return 2 +
		elm.LendingID.Size() +
		elm.MortgageDiamondList.Size() +
		elm.LoanTotalAmount.Size() +
		elm.BorrowPeriod.Size()
}

func (*Action_16_DiamondsSystemLendingRansom) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_16_DiamondsSystemLendingRansom) WriteinChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	feeAddr := act.belong_trs.GetAddress()

	// 检查id格式
	if len(act.LendingID) != stores.DiamondLendingIdLength ||
		act.LendingID[0] == 0 ||
		act.LendingID[stores.DiamondLendingIdLength-1] == 0 {
		return fmt.Errorf("Diamond Lending Id format error.")
	}

	// 查询id是否存在
	dmdlendObj := state.DiamondLending(act.LendingID)
	if dmdlendObj != nil {
		return fmt.Errorf("Diamond Lending <%d> already exist.", hex.EncodeToString(act.LendingID))
	}

	// 数量检查
	dianum := int(act.MortgageDiamondList.Count)
	if dianum == 0 || dianum != len(act.MortgageDiamondList.Diamonds) {
		return fmt.Errorf("Diamonds quantity error")
	}
	if dianum > 200 {
		return fmt.Errorf("Diamonds quantity cannot over 200")
	}

	// 检查周期数
	if act.BorrowPeriod < 1 || act.BorrowPeriod > 20 {
		return fmt.Errorf("BorrowPeriod must between 1 ~ 20")
	}

	// 可借出HAC
	totalLoanHAC := int64(0)

	// 批量抵押钻石
	for i := 0; i < len(act.MortgageDiamondList.Diamonds); i++ {
		diamond := act.MortgageDiamondList.Diamonds[i]

		// 查询钻石是否存在
		diaitem := state.Diamond(diamond)
		if diaitem == nil {
			return fmt.Errorf("Diamond <%s> not exist.", string(diamond))
		}
		item := diaitem
		// 检查是否已经抵押，是否可以抵押
		if diaitem.Status != stores.DiamondStatusNormal {
			return fmt.Errorf("Diamond <%s> has been mortgaged and cannot be transferred.", string(diamond))
		}
		// 检查所属
		if bytes.Compare(item.Address, feeAddr) != 0 {
			return fmt.Errorf("Diamond <%s> not belong to address '%s'", string(diamond), feeAddr.ToReadable())
		}
		// 标记抵押钻石
		item.Status = stores.DiamondStatusLendingSystem // 抵押给系统
		e5 := state.DiamondSet(diamond, item)
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

		// 统计可借出HAC数量
		totalLoanHAC += int64(diasmelt.AverageBidBurnPrice)
	}

	// 共借出 HAC 枚
	totalAmt, e8 := fields.NewAmountByBigIntWithUnit(big.NewInt(totalLoanHAC), 248)
	if e8 != nil {
		return e8
	}
	// 验证数量
	if totalAmt.Equal(&act.LoanTotalAmount) == false {
		return fmt.Errorf("LoanTotalAmount <%s> and <%s> not match.", totalAmt.ToFinString(), act.LoanTotalAmount.ToFinString())
	}

	// 减少钻石余额
	e9 := DoSubDiamondFromChainState(state, feeAddr, fields.VarUint3(dianum))
	if e9 != nil {
		return e9
	}

	// 抵押成功，发放余额
	e10 := DoAddBalanceFromChainState(state, feeAddr, act.LoanTotalAmount)
	if e10 != nil {
		return e10
	}

	// 完毕
	return nil
}

func (act *Action_16_DiamondsSystemLendingRansom) RecoverChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 回退所有抵押
	feeAddr := act.belong_trs.GetAddress()

	// 批量抵押钻石
	for i := 0; i < len(act.MortgageDiamondList.Diamonds); i++ {
		diamond := act.MortgageDiamondList.Diamonds[i]
		// 查询钻石是否存在
		diaitem := state.Diamond(diamond)
		// 标记抵押钻石
		diaitem.Status = stores.DiamondStatusNormal // 恢复
		e5 := state.DiamondSet(diamond, diaitem)
		if e5 != nil {
			return e5
		}
	}

	// 回退钻石余额
	dianum := act.MortgageDiamondList.Count
	e9 := DoAddDiamondFromChainState(state, feeAddr, fields.VarUint3(dianum))
	if e9 != nil {
		return e9
	}

	// 取回HAC余额
	e10 := DoSubBalanceFromChainState(state, feeAddr, act.LoanTotalAmount)
	if e10 != nil {
		return e10
	}

	return nil
}

// 设置所属 belong_trs
func (act *Action_16_DiamondsSystemLendingRansom) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_16_DiamondsSystemLendingRansom) IsBurning90PersentTxFees() bool {
	return false
}
