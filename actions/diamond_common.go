package actions

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
)

// diamond 转账
func DoSimpleDiamondTransferFromChainState(state interfacev2.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, dia fields.DiamondNumber) error {
	if bytes.Compare(addr1, addr2) == 0 {
		return nil // 可以自己转给自己，不改变状态，白费手续费
	}
	if dia == 0 {
		return nil // 数量为0，直接成功
	}
	bls1, e := state.Balance(addr1)
	if e != nil {
		return e
	}
	if bls1 == nil {
		return fmt.Errorf("Diamond not find.")
	}
	dia1 := bls1.Diamond
	// 检查余额
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
	bls1.Diamond = fields.DiamondNumber(uint32(dia1) - uint32(dia)) // 扣除
	bse1 := state.BalanceSet(addr1, bls1)
	if bse1 != nil {
		return bse1
	}
	bls2.Diamond = fields.DiamondNumber(uint32(dia2) + uint32(dia)) // 增加
	bse2 := state.BalanceSet(addr2, bls2)
	if bse2 != nil {
		return bse2
	}
	// return ok
	return nil
}

// 单纯增加 Diamond 余额
func DoAddDiamondFromChainState(state interfacev2.ChainStateOperation, addr fields.Address, dia fields.DiamondNumber) error {
	if dia == 0 {
		return nil // 数量为0，直接成功
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
	// 新余额
	blssto.Diamond = fields.DiamondNumber(newdia)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

// 单纯扣除 diamond 余额
func DoSubDiamondFromChainState(state interfacev2.ChainStateOperation, addr fields.Address, dia fields.DiamondNumber) error {
	if dia == 0 {
		return nil // 数量为0，直接成功
	}
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
	if blssto == nil {
		return fmt.Errorf("address %s diamond need %d not enough.", addr.ToReadable(), dia)
	}
	basedia := blssto.Diamond
	// 检查余额
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

// diamond 转账
func DoSimpleDiamondTransferFromChainStateV3(state interfaces.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, dia fields.DiamondNumber) error {
	if bytes.Compare(addr1, addr2) == 0 {
		return nil // 可以自己转给自己，不改变状态，白费手续费
	}
	if dia == 0 {
		return nil // 数量为0，直接成功
	}
	bls1, e := state.Balance(addr1)
	if e != nil {
		return e
	}
	if bls1 == nil {
		return fmt.Errorf("Diamond not find.")
	}
	dia1 := bls1.Diamond
	// 检查余额
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
	bls1.Diamond = fields.DiamondNumber(uint32(dia1) - uint32(dia)) // 扣除
	bse1 := state.BalanceSet(addr1, bls1)
	if bse1 != nil {
		return bse1
	}
	bls2.Diamond = fields.DiamondNumber(uint32(dia2) + uint32(dia)) // 增加
	bse2 := state.BalanceSet(addr2, bls2)
	if bse2 != nil {
		return bse2
	}
	// return ok
	return nil
}

// 单纯增加 Diamond 余额
func DoAddDiamondFromChainStateV3(state interfaces.ChainStateOperation, addr fields.Address, dia fields.DiamondNumber) error {
	if dia == 0 {
		return nil // 数量为0，直接成功
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
	// 新余额
	blssto.Diamond = fields.DiamondNumber(newdia)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

// 单纯扣除 diamond 余额
func DoSubDiamondFromChainStateV3(state interfaces.ChainStateOperation, addr fields.Address, dia fields.DiamondNumber) error {
	if dia == 0 {
		return nil // 数量为0，直接成功
	}
	blssto, e := state.Balance(addr)
	if e != nil {
		return e
	}
	if blssto == nil {
		return fmt.Errorf("address %s diamond need %d not enough.", addr.ToReadable(), dia)
	}
	basedia := blssto.Diamond
	// 检查余额
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
