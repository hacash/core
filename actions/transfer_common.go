package actions

import (
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
)

//////////////////////////////////////////////////////////

// HAC transfer
func DoSimpleTransferFromChainState(state interfacev2.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, amt fields.Amount) error {

	isTrsToMySelf := addr1.Equal(addr2)

	// In the database upgrade mode, all judgments are removed and the balance is modified directly
	if state.IsDatabaseVersionRebuildMode() {
		// Determine whether to transfer for yourself
		if isTrsToMySelf {
			return nil // You can transfer it to yourself without changing the status, which is a waste of service fees
		}
		bls1, e := state.Balance(addr1)
		if e != nil {
			return e
		}
		bls2, e := state.Balance(addr2)
		if e != nil {
			return e
		}
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

	// Normal start judgment
	//fmt.Println("addr1:", addr1.ToReadable(), "addr2:", addr2.ToReadable(), "amt:", amt.ToFinString())
	if state.GetPendingBlockHeight() < 200000 && isTrsToMySelf {
		// After 200000 yuan, the amount transferred to you is not allowed to be greater than the available balance!
		return nil // You can transfer it to yourself without changing the status, which is a waste of service fees
	}
	// Judge whether the balance is sufficient
	bls1, e := state.Balance(addr1)
	if e != nil {
		return e
	}
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
	// Determine whether to transfer for yourself
	if isTrsToMySelf {
		return nil // You can transfer it to yourself without changing the status, which is a waste of service fees
	}
	// Query payee balance
	bls2, e := state.Balance(addr2)
	if e != nil {
		return e
	}
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

// Simply increase balance
func DoAddBalanceFromChainState(state interfacev2.ChainStateOperation, addr fields.Address, amt fields.Amount) error {
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
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

// Net balance
func DoSubBalanceFromChainState(state interfacev2.ChainStateOperation, addr fields.Address, amt fields.Amount) error {
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
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

//////////////////////////////////////////////////////////

// HAC transfer
func DoSimpleTransferFromChainStateV3(state interfaces.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, amt fields.Amount) error {

	isTrsToMySelf := addr1.Equal(addr2)

	// In the database upgrade mode, all judgments are removed and the balance is modified directly
	if state.IsDatabaseVersionRebuildMode() {
		// Determine whether to transfer for yourself
		if isTrsToMySelf {
			return nil // You can transfer it to yourself without changing the status, which is a waste of service fees
		}
		bls1, e := state.Balance(addr1)
		if e != nil {
			return e
		}
		bls2, e := state.Balance(addr2)
		if e != nil {
			return e
		}
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

	// Normal start judgment
	//fmt.Println("addr1:", addr1.ToReadable(), "addr2:", addr2.ToReadable(), "amt:", amt.ToFinString())
	if state.GetPendingBlockHeight() < 200000 && isTrsToMySelf {
		// After 200000 yuan, the amount transferred to you is not allowed to be greater than the available balance!
		return nil // You can transfer it to yourself without changing the status, which is a waste of service fees
	}

	// First, judge whether the balance is sufficient
	// And then judge whether it is self transfer
	bls1, e := state.Balance(addr1)
	if e != nil {
		return e
	}
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
	// Determine whether to transfer for yourself
	if isTrsToMySelf {
		return nil // You can transfer it to yourself without changing the status, which is a waste of service fees
	}
	// Query payee balance
	bls2, e := state.Balance(addr2)
	if e != nil {
		return e
	}
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

// Simply increase balance
func DoAddBalanceFromChainStateV3(state interfaces.ChainStateOperation, addr fields.Address, amt fields.Amount) error {
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
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

// Net balance
func DoSubBalanceFromChainStateV3(state interfaces.ChainStateOperation, addr fields.Address, amt fields.Amount) error {
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
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
