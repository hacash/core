package actions

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

// diamond 转账
func DoSimpleDiamondTransferFromChainState(state interfaces.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, dia fields.VarUint3) error {
	if bytes.Compare(addr1, addr2) == 0 {
		return nil // 可以自己转给自己，不改变状态，白费手续费
	}
	bls1 := state.Balance(addr1)
	if bls1 == nil {
		return fmt.Errorf("Diamond not find.")
	}
	dia1 := bls1.Diamond
	// 检查余额
	if uint64(dia1) < uint64(dia) {
		return fmt.Errorf("Address %s diamond %d not enough, need more %d.", addr1.ToReadable(), dia1, dia)
	}
	bls2 := state.Balance(addr2)
	if bls2 == nil {
		bls2 = stores.NewEmptyBalance() // create satoshi store
	}
	dia2 := bls2.Diamond
	bls1.Diamond = fields.VarUint3(uint32(dia1) - uint32(dia)) // 扣除
	bse1 := state.BalanceSet(addr1, bls1)
	if bse1 != nil {
		return bse1
	}
	bls2.Diamond = fields.VarUint3(uint32(dia2) + uint32(dia)) // 增加
	bse2 := state.BalanceSet(addr2, bls2)
	if bse2 != nil {
		return bse2
	}
	// return ok
	return nil
}

// 单纯增加 Diamond 余额
func DoAddDiamondFromChainState(state interfaces.ChainStateOperation, addr fields.Address, dia fields.VarUint3) error {
	blssto := state.Balance(addr)
	if blssto == nil {
		blssto = stores.NewEmptyBalance() // first create account
	}
	basedia := blssto.Diamond
	newdia := uint64(basedia) + uint64(dia)
	// 新余额
	blssto.Diamond = fields.VarUint3(newdia)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

// 单纯扣除 diamond 余额
func DoSubDiamondFromChainState(state interfaces.ChainStateOperation, addr fields.Address, dia fields.VarUint3) error {
	blssto := state.Balance(addr)
	if blssto == nil {
		return fmt.Errorf("address %s diamond %d not enough.")
	}
	basedia := blssto.Diamond
	// 检查余额
	if uint64(basedia) < uint64(dia) {
		return fmt.Errorf("address %s satoshi %d not enough, need more %d.", addr.ToReadable(), basedia, dia)
	}
	newdia := uint64(basedia) - uint64(dia)
	blssto.Diamond = fields.VarUint3(newdia)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}
