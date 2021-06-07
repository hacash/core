package actions

import (
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

//////////////////////////////////////////////////////////

// hac转账
func DoSimpleTransferFromChainState(state interfaces.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, amt fields.Amount) error {

	isTrsToMySelf := addr1.Equal(addr2)

	// 如果是数据库升级模式，则去掉所有的判断，直接修改余额
	if state.IsDatabaseVersionRebuildMode() {
		// 判断是否为自己转给自己
		if isTrsToMySelf {
			return nil // 可以自己转给自己，不改变状态，白费手续费
		}
		bls1 := state.Balance(addr1)
		bls2 := state.Balance(addr2)
		if bls2 == nil {
			bls2 = stores.NewEmptyBalance() // create balance store
		}
		amt1 := bls1.Hacash
		amtsub, _ := amt1.Sub(&amt)
		amt2 := bls2.Hacash
		amtadd, _ := amt2.Add(&amt)
		bls1.Hacash = *amtsub
		state.BalanceSet(addr1, bls1)
		bls2.Hacash = *amtadd
		state.BalanceSet(addr2, bls2)
		return nil
	}

	// 正常开始判断
	//fmt.Println("addr1:", addr1.ToReadable(), "addr2:", addr2.ToReadable(), "amt:", amt.ToFinString())
	if state.GetPendingBlockHeight() < 200000 && isTrsToMySelf {
		// 高度 20万 之后，不允许出现自己转给自己的数额大于可用余额的情况！
		return nil // 可以自己转给自己，不改变状态，白费手续费
	}
	// 判断余额是否充足
	bls1 := state.Balance(addr1)
	if bls1 == nil {
		// test
		//fmt.Println( addr1.ToReadable(), "Balance ", amt.ToFinString(), " not find." )
		return fmt.Errorf("Balance not find.")
	}
	amt1 := bls1.Hacash
	//fmt.Println("amt1: " + amt1.ToFinString())
	if amt1.LessThan(&amt) {
		//x, _ := amt.Sub(&amt1)
		//print_xxxxxxx(addr1, x)
		//fmt.Println("[balance not enough]", "addr1: ", addr1.ToReadable(), "amt: " + amt.ToFinString(), "amt1: " + amt1.ToFinString())
		return fmt.Errorf("address %s balance %s not enough， need %s.", addr1.ToReadable(), amt1.ToFinString(), amt.ToFinString())
	}
	// 判断是否为自己转给自己
	if isTrsToMySelf {
		return nil // 可以自己转给自己，不改变状态，白费手续费
	}
	// 查询收款方余额
	bls2 := state.Balance(addr2)
	if bls2 == nil {
		bls2 = stores.NewEmptyBalance() // create balance store
	}
	amt2 := bls2.Hacash
	//fmt.Println("amt2: " + amt2.ToFinString())
	// add
	amtsub, e1 := amt1.Sub(&amt)
	if e1 != nil {
		//fmt.Println("e1: ", e1)
		return e1
	}
	amtadd, e2 := amt2.Add(&amt)
	if e2 != nil {
		//fmt.Println("e2: ", e2)
		return e2
	}
	//fmt.Println("EllipsisDecimalFor23SizeStore: ")
	amtsub_1, ischg1, ederr1 := amtsub.EllipsisDecimalFor11SizeStore()
	amtadd_1, ischg2, ederr2 := amtadd.EllipsisDecimalFor11SizeStore()
	if ederr1 != nil {
		return ederr1
	}
	if ederr2 != nil {
		return ederr2
	}
	if ischg1 || ischg2 {
		return fmt.Errorf("amount can not to store")
	}
	amtsub = amtsub_1
	amtadd = amtadd_1
	//if amtsub.IsEmpty() {
	//	state.BalanceDel(addr1) // 归零
	//} else {
	//fmt.Println("amtsub: " + amtsub.ToFinString())
	bls1.Hacash = *amtsub
	bse1 := state.BalanceSet(addr1, bls1)
	if bse1 != nil {
		return bse1
	}
	//}
	bls2.Hacash = *amtadd
	bse2 := state.BalanceSet(addr2, bls2)
	if bse2 != nil {
		return bse2
	}
	// return ok
	return nil
}

// 单纯增加余额
func DoAddBalanceFromChainState(state interfaces.ChainStateOperation, addr fields.Address, amt fields.Amount) error {
	blssto := state.Balance(addr)
	if blssto == nil {
		blssto = stores.NewEmptyBalance() // first create account
	}
	baseamt := blssto.Hacash
	//fmt.Println( "baseamt: ", baseamt.ToFinString() )
	amtnew, e1 := baseamt.Add(&amt)
	if e1 != nil {
		return e1
	}
	amtsave, ischg, ec1 := amtnew.EllipsisDecimalFor11SizeStore()
	if ec1 != nil {
		return ec1
	}
	if ischg {
		return fmt.Errorf("amount can not to store")
	}
	//addrrr, _ := base58check.Encode(addr)
	//fmt.Println( "DoAddBalanceFromChainState: ++++++++++ ", addr.ToReadable(), amtsave.ToFinString() )
	blssto.Hacash = *amtsave
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

// 单纯扣除余额
func DoSubBalanceFromChainState(state interfaces.ChainStateOperation, addr fields.Address, amt fields.Amount) error {
	blssto := state.Balance(addr)
	if blssto == nil {
		return fmt.Errorf("address %s amount need %s not enough.", addr.ToReadable(), amt.ToFinString())
	}
	baseamt := blssto.Hacash
	//fmt.Println("baseamt: " + baseamt.ToFinString())
	if baseamt.LessThan(&amt) {
		//x, _ := amt.Sub(&baseamt)
		//print_xxxxxxx(addr, x)
		//fmt.Println("[balance not enough]", "block height: 0", "addr: ", addr.ToReadable(), "baseamt: " + baseamt.ToFinString(), "amt: " + amt.ToFinString())
		return fmt.Errorf("address %s balance %s not enough， need %s.", addr.ToReadable(), baseamt.ToFinString(), amt.ToFinString())
	}
	//fmt.Println("amt fee: " + amt.ToFinString())
	amtnew, e1 := baseamt.Sub(&amt)
	if e1 != nil {
		return e1
	}
	amtnew1, ischg, ec1 := amtnew.EllipsisDecimalFor11SizeStore()
	if ec1 != nil {
		return ec1
	}
	if ischg {
		return fmt.Errorf("amount can not to store")
	}
	//fmt.Println("amtnew1: " + amtnew1.ToFinString())
	blssto.Hacash = *amtnew1
	//fmt.Println("state.BalanceSet: ", addr.ToReadable(), amtnew1.ToFinString())
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

/*
func init() {

	amt1, _ := fields.NewAmountFromFinString("ㄜ1:248")
	amt2, _ := fields.NewAmountFromFinString("ㄜ1:248")

	amt5, amt6 := DoAppendCompoundInterest1Of10000By2500Height(amt1, amt2, 1200)

	fmt.Println("DoAppendCompoundInterest1Of10000By2500Height: ", amt5.ToFinString(), amt6.ToFinString())
}
*/

/*
func init1() {

	amt1, _ := fields.NewAmountFromFinString("ㄜ5:250")
	amt2, _ := fields.NewAmountFromFinString("ㄜ5:250")
	amt3, _ := fields.NewAmountFromFinString("ㄜ55799:246")
	amt4, _ := fields.NewAmountFromFinString("ㄜ5279999:244")

	amt5, amt6 := DoAppendCompoundInterest1Of10000By2500Height(amt1, amt2, 0)
	amt7, amt8 := DoAppendCompoundInterest1Of10000By2500Height(amt3, amt4, 1)

	tt1, _ := amt5.Sub(amt1)
	tt2, _ := amt6.Sub(amt2)
	tt3, _ := amt7.Sub(amt3)
	tt4, _ := amt8.Sub(amt4)

	totalsub, _ := tt1.Add(tt2)
	totalsub, _ = totalsub.Add(tt3)
	totalsub, _ = totalsub.Add(tt4)

	fmt.Println("yangjie: ", totalsub.ToFinString())
}
*/

/*

///////  余额检查测试  ///////

var amtxx *fields.Amount = nil
var amts []*fields.Amount = nil
func print_xxxxxxx(addr fields.Address, amtx *fields.Amount)  {
	if amtxx == nil {
		amtxx = fields.NewEmptyAmount()
		amts = []*fields.Amount{
			fields.NewEmptyAmount(),
			fields.NewEmptyAmount(),
			fields.NewEmptyAmount(),
			fields.NewEmptyAmount(),
		}
	}
	adname := addr.ToReadable()
	amtxx, _ = amtxx.Add(amtx)
	idx := -1
	if strings.Index(adname, "1LsQL") > -1 {
		idx = 0
	}else if strings.Index(adname, "12vi7") > -1 {
		idx = 1
	}else if strings.Index(adname, "1NUgK") > -1 {
		idx = 2
	}else if strings.Index(adname, "1HE2qA") > -1 {
		idx = 3
	}else{
		//panic(adname)
	}
	amts[idx], _ = amts[idx].Add(amtx)

	fmt.Println("addr ", adname, " add ", amtx.ToFinString(), "addr amt", amts[idx].ToFinString(), " total amt:", amtxx.ToFinString())
	for _, v := range amts {
		fmt.Print(v.ToFinString()+", ")
	}
	fmt.Println("")
}

*/
