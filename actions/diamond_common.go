package actions

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
)

// Diamond transfer
func DoSimpleDiamondTransferFromChainState(state interfacev2.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, dia fields.DiamondNumber) error {
	if bytes.Compare(addr1, addr2) == 0 {
		return nil // You can transfer it to yourself without changing the status, which is a waste of service fees
	}
	if dia == 0 {
		return nil // Quantity is 0, direct success
	}
	bls1, e := state.Balance(addr1)
	if e != nil {
		return e
	}
	if bls1 == nil {
		return fmt.Errorf("Diamond not find.")
	}
	dia1 := bls1.Diamond
	// Check balance
	if uint64(dia1) < uint64(dia) {
		return fmt.Errorf("Address %s diamond %d not enough, need more %d.", addr1.ToReadable(), dia1, dia)
	}
	bls2, e := state.Balance(addr2)
	if e != nil {
		return e
	}
	if bls2 == nil {
		bls2 = stores.NewEmptyBalance() // create satoshi store
	}
	dia2 := bls2.Diamond
	bls1.Diamond = fields.DiamondNumber(uint32(dia1) - uint32(dia)) // deduction
	bse1 := state.BalanceSet(addr1, bls1)
	if bse1 != nil {
		return bse1
	}
	bls2.Diamond = fields.DiamondNumber(uint32(dia2) + uint32(dia)) // increase
	bse2 := state.BalanceSet(addr2, bls2)
	if bse2 != nil {
		return bse2
	}
	// return ok
	return nil
}

// Simply increase diamond balance
func DoAddDiamondFromChainState(state interfacev2.ChainStateOperation, addr fields.Address, dia fields.DiamondNumber) error {
	if dia == 0 {
		return nil // Quantity is 0, direct success
	}
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
	if blssto == nil {
		blssto = stores.NewEmptyBalance() // first create account
	}
	basedia := blssto.Diamond
	newdia := uint64(basedia) + uint64(dia)
	// New balance
	blssto.Diamond = fields.DiamondNumber(newdia)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

// Simply deduct diamond balance
func DoSubDiamondFromChainState(state interfacev2.ChainStateOperation, addr fields.Address, dia fields.DiamondNumber) error {
	if dia == 0 {
		return nil // Quantity is 0, direct success
	}
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
	if blssto == nil {
		return fmt.Errorf("address %s diamond need %d not enough.", addr.ToReadable(), dia)
	}
	basedia := blssto.Diamond
	// Check balance
	if uint64(basedia) < uint64(dia) {
		return fmt.Errorf("address %s satoshi %d not enough, need more %d.", addr.ToReadable(), basedia, dia)
	}
	newdia := uint64(basedia) - uint64(dia)
	blssto.Diamond = fields.DiamondNumber(newdia)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

/////////////////////////////////////////////

// Diamond transfer
func DoSimpleDiamondTransferFromChainStateV3(state interfaces.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, dia fields.DiamondNumber) error {
	if bytes.Compare(addr1, addr2) == 0 {
		return nil // You can transfer it to yourself without changing the status, which is a waste of service fees
	}
	if dia == 0 {
		return nil // Quantity is 0, direct success
	}
	bls1, e := state.Balance(addr1)
	if e != nil {
		return e
	}
	if bls1 == nil {
		return fmt.Errorf("Diamond not find.")
	}
	dia1 := bls1.Diamond
	// Check balance
	if uint64(dia1) < uint64(dia) {
		return fmt.Errorf("Address %s diamond %d not enough, need more %d.", addr1.ToReadable(), dia1, dia)
	}
	bls2, e := state.Balance(addr2)
	if e != nil {
		return e
	}
	if bls2 == nil {
		bls2 = stores.NewEmptyBalance() // create satoshi store
	}
	dia2 := bls2.Diamond
	bls1.Diamond = fields.DiamondNumber(uint32(dia1) - uint32(dia)) // deduction
	bse1 := state.BalanceSet(addr1, bls1)
	if bse1 != nil {
		return bse1
	}
	bls2.Diamond = fields.DiamondNumber(uint32(dia2) + uint32(dia)) // increase
	bse2 := state.BalanceSet(addr2, bls2)
	if bse2 != nil {
		return bse2
	}
	// return ok
	return nil
}

// Simply increase diamond balance
func DoAddDiamondFromChainStateV3(state interfaces.ChainStateOperation, addr fields.Address, dia fields.DiamondNumber) error {
	if dia == 0 {
		return nil // Quantity is 0, direct success
	}
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
	if blssto == nil {
		blssto = stores.NewEmptyBalance() // first create account
	}
	basedia := blssto.Diamond
	newdia := uint64(basedia) + uint64(dia)
	// New balance
	blssto.Diamond = fields.DiamondNumber(newdia)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

// Simply deduct diamond balance
func DoSubDiamondFromChainStateV3(state interfaces.ChainStateOperation, addr fields.Address, dia fields.DiamondNumber) error {
	if dia == 0 {
		return nil // Quantity is 0, direct success
	}
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
	if blssto == nil {
		return fmt.Errorf("address %s diamond need %d not enough.", addr.ToReadable(), dia)
	}
	basedia := blssto.Diamond
	// Check balance
	if uint64(basedia) < uint64(dia) {
		return fmt.Errorf("address %s satoshi %d not enough, need more %d.", addr.ToReadable(), basedia, dia)
	}
	newdia := uint64(basedia) - uint64(dia)
	blssto.Diamond = fields.DiamondNumber(newdia)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}
