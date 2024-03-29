package actions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
)

/*


 */

// Diamond engraved
type Action_32_DiamondsEngraved struct {
	//
	DiamondList     fields.DiamondListMaxLen200
	ProtocolCost    fields.Amount   // HAC amount for burning
	EngravedType    fields.VarUint1 //  0:String  1:CompressedDict  51:MD5  52:SHA256 ....
	EngravedContent fields.StringMax255

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_32_DiamondsEngraved) Kind() uint16 {
	return 32
}

// json api
func (elm *Action_32_DiamondsEngraved) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_32_DiamondsEngraved) Serialize() ([]byte, error) {
	var e error = nil
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	b1, e := elm.DiamondList.Serialize()
	if e != nil {
		return nil, e
	}
	b2, e := elm.ProtocolCost.Serialize()
	if e != nil {
		return nil, e
	}
	b3, e := elm.EngravedType.Serialize()
	if e != nil {
		return nil, e
	}
	b4, e := elm.EngravedContent.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	return buffer.Bytes(), nil
}

func (elm *Action_32_DiamondsEngraved) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.DiamondList.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ProtocolCost.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.EngravedType.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.EngravedContent.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_32_DiamondsEngraved) Size() uint32 {
	return 2 +
		elm.DiamondList.Size() +
		elm.ProtocolCost.Size() +
		elm.EngravedType.Size() +
		elm.EngravedContent.Size()
}

func (*Action_32_DiamondsEngraved) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_32_DiamondsEngraved) WriteInChainState(state interfaces.ChainStateOperation) error {

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 518000 start open
	if false == sys.TestDebugLocalDevelopmentMark && state.GetPendingBlockHeight() < 518000 {
		return fmt.Errorf("HIP-15 is effective starting at block 518000")
	}

	if act.ProtocolCost.Size() > 4 {
		return fmt.Errorf("ProtocolCost amount size cannot over 4 bytes")
	}

	var content_type = uint8(act.EngravedType)
	var content = act.EngravedContent.Str
	var content_len = len(content)
	// check String and Visible
	if content_len > 64 {
		return fmt.Errorf("EngravedContent max size is 64")
	}
	if content_type <= 50 {
		if !fields.IsValidVisibleString(content) {
			return fmt.Errorf("EngravedContent must be visible string")
		}
	}

	mainAddr := act.belong_trs_v3.GetAddress()
	//fmt.Println(mainAddr.ToReadable())

	var dias = act.DiamondList.Diamonds
	var dianum = len(dias)
	if dianum > 200 {
		return fmt.Errorf("Diamonds cannot over 200")
	}
	var ttcost = fields.NewEmptyAmount()
	var e error = nil
	for i := 0; i < dianum; i++ {
		var dia = dias[i]
		cost, e := handleEngravedOneDiamond(&mainAddr, dia, &act.EngravedContent, state)
		if e != nil {
			return fmt.Errorf("handleEngravedOneDiamond %s error: %s", dia.Name(), e)
		}
		ttcost, e = ttcost.Add(cost)
		if e != nil {
			return fmt.Errorf("handleEngravedOneDiamond %s error: %s", dia.Name(), e)
		}
	}

	// check cost
	if act.ProtocolCost.LessThan(ttcost) {
		return fmt.Errorf("Engraved Diamond cost error need %s but got %s",
			ttcost.ToFinString(), act.ProtocolCost.ToFinString())
	}

	// sub main addr balance
	if act.ProtocolCost.IsNotEmpty() {
		e = DoSubBalanceFromChainState(state, mainAddr, act.ProtocolCost)
		if e != nil {
			return fmt.Errorf("Engraved Diamond main address balance need %s but not enough",
				act.ProtocolCost.ToFinString())
		}
	}

	// Total supply statistics
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	var costmei = act.ProtocolCost.ToMei()
	totalsupply.DoAddUint(stores.TotalSupplyStoreTypeOfDiamondEngravedOperateCount, uint64(act.DiamondList.Count)) //
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfDiamondEngravedBurning, costmei)                                // Engraved burning
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBurningTotal, costmei)                                          // Total burning
	// update total supply
	e7 := state.UpdateSetTotalSupply(totalsupply)
	if e7 != nil {
		return e7
	}

	// complete
	return nil
}

func (act *Action_32_DiamondsEngraved) WriteinChainState(state interfacev2.ChainStateOperation) error {

	panic("WriteinChainState be deprecated")

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	mainAddr := act.belong_trs.GetAddress()
	fmt.Println(mainAddr.ToReadable())

	// complete
	return nil
}

func (act *Action_32_DiamondsEngraved) RecoverChainState(state interfacev2.ChainStateOperation) error {

	panic("RecoverChainState be deprecated")

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	return nil
}

// Set belongs to long_ trs
func (act *Action_32_DiamondsEngraved) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}
func (act *Action_32_DiamondsEngraved) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning 90% fees
func (act *Action_32_DiamondsEngraved) IsBurning90PersentTxFees() bool {
	return true
}

///////////////////////////////////////////////////////

// Diamond engraved recovery
type Action_33_DiamondsEngravedRecovery struct {
	//
	DiamondList  fields.DiamondListMaxLen200
	ProtocolCost fields.Amount // HAC number for burning

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_33_DiamondsEngravedRecovery) Kind() uint16 {
	return 33
}

// json api
func (elm *Action_33_DiamondsEngravedRecovery) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_33_DiamondsEngravedRecovery) Serialize() ([]byte, error) {
	var e error = nil
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	b1, e := elm.DiamondList.Serialize()
	if e != nil {
		return nil, e
	}
	b2, e := elm.ProtocolCost.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(b1)
	buffer.Write(b2)
	return buffer.Bytes(), nil
}

func (elm *Action_33_DiamondsEngravedRecovery) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.DiamondList.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ProtocolCost.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_33_DiamondsEngravedRecovery) Size() uint32 {
	return 2 +
		elm.DiamondList.Size() +
		elm.ProtocolCost.Size()
}

func (*Action_33_DiamondsEngravedRecovery) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_33_DiamondsEngravedRecovery) WriteInChainState(state interfaces.ChainStateOperation) error {

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	mainAddr := act.belong_trs_v3.GetAddress()
	//fmt.Println(mainAddr.ToReadable())

	var dias = act.DiamondList.Diamonds
	var dianum = len(dias)
	if dianum > 200 {
		return fmt.Errorf("Diamonds cannot over 200")
	}
	var ttcostmei uint32 = 0
	var e error = nil
	for i := 0; i < dianum; i++ {
		var dia = dias[i]
		costmei, e := handleRecoveryEngravedOneDiamond(&mainAddr, dia, state)
		if e != nil {
			return fmt.Errorf("handleRecoveryEngravedOneDiamond %s error: %s", dia.Name(), e)
		}
		ttcostmei += costmei
		if e != nil {
			return fmt.Errorf("handleRecoveryEngravedOneDiamond %s error: %s", dia.Name(), e)
		}
	}
	ttcost := fields.NewAmountByUnitMei(int64(ttcostmei))

	// check cost
	if act.ProtocolCost.LessThan(ttcost) {
		return fmt.Errorf("Engraved Diamond cost error need %s but got %s",
			ttcost.ToFinString(), act.ProtocolCost.ToFinString())
	}

	// sub main addr balance
	if act.ProtocolCost.IsNotEmpty() {
		e = DoSubBalanceFromChainState(state, mainAddr, act.ProtocolCost)
		if e != nil {
			return fmt.Errorf("Engraved Diamond main address balance need %s but not enough",
				act.ProtocolCost.ToFinString())
		}
	}

	// Total supply statistics
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	var costmei = act.ProtocolCost.ToMei()
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfDiamondEngravedBurning, costmei) // Engraved burning
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBurningTotal, costmei)           // Total burning
	// update total supply
	e7 := state.UpdateSetTotalSupply(totalsupply)
	if e7 != nil {
		return e7
	}

	// complete
	return nil
}

func (act *Action_33_DiamondsEngravedRecovery) WriteinChainState(state interfacev2.ChainStateOperation) error {

	panic("WriteinChainState be deprecated")

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	mainAddr := act.belong_trs.GetAddress()
	fmt.Println(mainAddr.ToReadable())

	// complete
	return nil
}

func (act *Action_33_DiamondsEngravedRecovery) RecoverChainState(state interfacev2.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // Waiting for review is not enabled yet
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	return nil
}

// Set belongs to long_ trs
func (act *Action_33_DiamondsEngravedRecovery) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}
func (act *Action_33_DiamondsEngravedRecovery) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning 90% fees
func (act *Action_33_DiamondsEngravedRecovery) IsBurning90PersentTxFees() bool {
	return true
}

////////////////////////////////////////////

func RequestProtocolCostForDiamondList(state interfaces.ChainStateOperationRead,
	dialist *fields.DiamondListMaxLen200, checkaddr *fields.Address) (*fields.Amount, error) {
	var store = state.BlockStoreRead()
	var resamt = fields.NewEmptyAmount()
	dias := dialist.Diamonds
	for i := 0; i < len(dias); i++ {
		var dia = dias[i]
		diasto, e := state.Diamond(dia)
		if e != nil {
			return nil, e
		}
		if diasto == nil {
			return nil, fmt.Errorf(`not find diamond %s`, dia.Name())
		}
		// if check belong address
		if checkaddr != nil {
			e = CheckDiamondStatusNormalAndBelong(&dia, diasto, checkaddr)
			if e != nil {
				return nil, e
			}
		}
		diaslt, e := store.ReadDiamond(dia)
		if e != nil {
			return nil, e
		}
		// inscription
		insnum := diasto.EngravedContents.Count
		if insnum >= 200 {
			return nil, fmt.Errorf("%s inscription too much", dia.Name())
		}
		if insnum >= 10 {
			// append cost fee of 1/10 AverageBidBurnPrice
			one_cost := fields.NewAmountByUnit(int64(diaslt.AverageBidBurnPrice), 247)
			resamt, e = resamt.Add(one_cost)
		}
		// next
	}
	// ok
	return resamt, nil
}

/*
return total cost
*/
func handleEngravedOneDiamond(mainAddr *fields.Address, diamond fields.DiamondName, content *fields.StringMax255, state interfaces.ChainStateOperation) (*fields.Amount, error) {

	var store = state.BlockStore()

	curblkhei := state.GetPendingBlockHeight()
	cost := fields.NewEmptyAmount()

	// load diamond
	dia, err := state.Diamond(diamond)
	if err != nil {
		return nil, err
	}
	var overcheckblknum uint64 = 1000
	if sys.TestDebugLocalDevelopmentMark {
		overcheckblknum /= 100 // 10 blk for test
	}
	if overcheckblknum+uint64(dia.EngravedPrevBlockHeight) > curblkhei {
		return nil, fmt.Errorf("Only one inscription can be made every %d blocks", overcheckblknum)
	}
	// check belong and status
	err = CheckDiamondStatusNormalAndBelong(&diamond, dia, mainAddr)
	if err != nil {
		return nil, err
	}
	diaslt, err := store.ReadDiamond(diamond)
	if err != nil {
		return nil, err
	}
	engsz := int(dia.EngravedContents.Count)
	if engsz >= 200 {
		return nil, fmt.Errorf("The maximum number of inscriptions is 200")
	}
	if engsz >= 10 {
		// burning cost bid fee 1/10
		cost = fields.NewAmountByUnit(int64(diaslt.AverageBidBurnPrice), 247) // 1/10
	}
	// do engraved
	dia.EngravedPrevBlockHeight = fields.BlockHeight(curblkhei)
	dia.EngravedContents.Append(content)
	// save
	err = state.DiamondSet(diamond, dia)
	if err != nil {
		return nil, err
	}

	// ok
	return cost, nil
}

/*
return total cost
*/
func handleRecoveryEngravedOneDiamond(mainAddr *fields.Address, diamond fields.DiamondName, state interfaces.ChainStateOperation) (uint32, error) {

	var store = state.BlockStore()

	// load diamond
	dia, err := state.Diamond(diamond)
	if err != nil {
		return 0, err
	}
	if dia.EngravedContents.Count <= 0 {
		return 0, fmt.Errorf("cannot find any inscriptions in HACD %s", diamond.Name())
	}
	// check belong and status
	err = CheckDiamondStatusNormalAndBelong(&diamond, dia, mainAddr)
	if err != nil {
		return 0, err
	}
	diaslt, err := store.ReadDiamond(diamond)
	if err != nil {
		return 0, err
	}
	// do recovery engraved
	dia.EngravedPrevBlockHeight = fields.BlockHeight(0)            // reset
	dia.EngravedContents = fields.CreateEmptyStringMax255List255() // set empty
	// save
	err = state.DiamondSet(diamond, dia)
	if err != nil {
		return 0, err
	}
	// ok
	return uint32(diaslt.AverageBidBurnPrice), nil
}
