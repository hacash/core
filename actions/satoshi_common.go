package actions

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
)

// BTC transfer (AMT unit Cong)
func DoSimpleSatoshiTransferFromChainState(state interfacev2.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, sat fields.Satoshi) error {
	if sat == 0 {
		return fmt.Errorf("Satoshi transfer amount is empty") // Transfer quantity 0 is not allowed
	}
	bls1, e := state.Balance(addr1)
	if e != nil {
		return e
	}
	if bls1 == nil {
		return fmt.Errorf("Satoshi need %d but empty.", sat)
	}
	sat1 := bls1.Satoshi
	// Check balance
	if uint64(sat1) < uint64(sat) {
		return fmt.Errorf("Address %s satoshi %d not enough, need at least %d.", addr1.ToReadable(), sat1, sat)
	}
	// Check yourself and transfer to yourself
	if bytes.Compare(addr1, addr2) == 0 {
		return nil // You can transfer it to yourself without changing the status. First, check that the balance is sufficient. It is a waste of service fees
	}
	bls2, e := state.Balance(addr2)
	if e != nil {
		return e
	}
	if bls2 == nil {
		bls2 = stores.NewEmptyBalance() // create satoshi store
	}
	sat2 := bls2.Satoshi
	bls1.Satoshi = fields.Satoshi(uint64(sat1) - uint64(sat)) // deduction
	bse1 := state.BalanceSet(addr1, bls1)
	if bse1 != nil {
		return bse1
	}
	bls2.Satoshi = fields.Satoshi(uint64(sat2) + uint64(sat)) // increase
	bse2 := state.BalanceSet(addr2, bls2)
	if bse2 != nil {
		return bse2
	}
	// return ok
	return nil
}

// Simply increase BTC balance (AMT unit)
func DoAddSatoshiFromChainState(state interfacev2.ChainStateOperation, addr fields.Address, sat fields.Satoshi) error {
	if sat == 0 {
		return nil // Quantity is 0, direct success
	}
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
	if blssto == nil {
		blssto = stores.NewEmptyBalance() // first create account
	}
	basesat := blssto.Satoshi
	newsat := uint64(basesat) + uint64(sat) // 增加
	// New balance
	blssto.Satoshi = fields.Satoshi(newsat)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

// Simply deduct BTC balance (AMT unit)
func DoSubSatoshiFromChainState(state interfacev2.ChainStateOperation, addr fields.Address, sat fields.Satoshi) error {
	if sat == 0 {
		return nil // Quantity is 0, direct success
	}
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
	if blssto == nil {
		return fmt.Errorf("address %s satoshi need %d but empty.", addr.ToReadable(), sat)
	}
	basesat := blssto.Satoshi
	// Check balance
	if uint64(basesat) < uint64(sat) {
		return fmt.Errorf("address %s satoshi %d not enough, need more %d.", addr.ToReadable(), basesat, sat)
	}
	newsat := uint64(basesat) - uint64(sat) // 扣除
	blssto.Satoshi = fields.Satoshi(newsat)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

////////////////////////////////////////

// BTC transfer (AMT unit Cong)
func DoSimpleSatoshiTransferFromChainStateV3(state interfaces.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, sat fields.Satoshi) error {
	if sat == 0 {
		return fmt.Errorf("Satoshi transfer amount is empty") // Transfer quantity 0 is not allowed
	}
	bls1, e := state.Balance(addr1)
	if e != nil {
		return e
	}
	if bls1 == nil {
		return fmt.Errorf("Satoshi need %d but empty.", sat)
	}
	sat1 := bls1.Satoshi
	// Check balance
	if uint64(sat1) < uint64(sat) {
		return fmt.Errorf("Address %s satoshi %d not enough, need at least %d.", addr1.ToReadable(), sat1, sat)
	}
	// Check yourself and transfer to yourself
	if bytes.Compare(addr1, addr2) == 0 {
		return nil // You can transfer it to yourself without changing the status. First, check that the balance is sufficient. It is a waste of service fees
	}
	bls2, e := state.Balance(addr2)
	if e != nil {
		return e
	}
	if bls2 == nil {
		bls2 = stores.NewEmptyBalance() // create satoshi store
	}
	sat2 := bls2.Satoshi
	bls1.Satoshi = fields.Satoshi(uint64(sat1) - uint64(sat)) // deduction
	bse1 := state.BalanceSet(addr1, bls1)
	if bse1 != nil {
		return bse1
	}
	bls2.Satoshi = fields.Satoshi(uint64(sat2) + uint64(sat)) // increase
	bse2 := state.BalanceSet(addr2, bls2)
	if bse2 != nil {
		return bse2
	}
	// return ok
	return nil
}

// Simply increase BTC balance (AMT unit)
func DoAddSatoshiFromChainStateV3(state interfaces.ChainStateOperation, addr fields.Address, sat fields.Satoshi) error {
	if sat == 0 {
		return nil // Quantity is 0, direct success
	}
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
	if blssto == nil {
		blssto = stores.NewEmptyBalance() // first create account
	}
	basesat := blssto.Satoshi
	newsat := uint64(basesat) + uint64(sat) // 增加
	// New balance
	blssto.Satoshi = fields.Satoshi(newsat)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

// Simply deduct BTC balance (AMT unit)
func DoSubSatoshiFromChainStateV3(state interfaces.ChainStateOperation, addr fields.Address, sat fields.Satoshi) error {
	if sat == 0 {
		return nil // Quantity is 0, direct success
	}
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
	if blssto == nil {
		return fmt.Errorf("address %s satoshi need %d but empty.", addr.ToReadable(), sat)
	}
	basesat := blssto.Satoshi
	// Check balance
	if uint64(basesat) < uint64(sat) {
		return fmt.Errorf("address %s satoshi %d not enough, need more %d.", addr.ToReadable(), basesat, sat)
	}
	newsat := uint64(basesat) - uint64(sat) // 扣除
	blssto.Satoshi = fields.Satoshi(newsat)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}
