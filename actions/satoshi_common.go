package actions

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

// btc 转账 （amt 单位 聪）
func DoSimpleSatoshiTransferFromChainState(state interfaces.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, amt fields.VarUint8) error {
	if bytes.Compare(addr1, addr2) == 0 {
		return nil // 可以自己转给自己，不改变状态，白费手续费
	}
	bls1 := state.Satoshi(addr1)
	if bls1 == nil {
		return fmt.Errorf("Satoshi not find.")
	}
	amt1 := bls1.Amount
	// 检查余额
	if uint64(amt1) < uint64(amt) {
		return fmt.Errorf("Address %s satoshi %d not enough, need more %d.", addr1.ToReadable(), amt1, amt)
	}
	bls2 := state.Satoshi(addr2)
	if bls2 == nil {
		bls2 = stores.NewEmptySatoshi() // create satoshi store
	}
	amt2 := bls2.Amount
	bls1.Amount = fields.VarUint8(uint64(amt1) - uint64(amt)) // 扣除
	bse1 := state.SatoshiSet(addr1, bls1)
	if bse1 != nil {
		return bse1
	}
	bls2.Amount = fields.VarUint8(uint64(amt2) + uint64(amt)) // 增加
	bse2 := state.SatoshiSet(addr2, bls2)
	if bse2 != nil {
		return bse2
	}
	// return ok
	return nil
}

// 单纯增加 BTC 余额 （amt 单位 聪）
func DoAddSatoshiFromChainState(state interfaces.ChainStateOperation, addr fields.Address, amt fields.VarUint8) error {
	blssto := state.Satoshi(addr)
	if blssto == nil {
		blssto = stores.NewEmptySatoshi() // first create account
	}
	baseamt := blssto.Amount
	newamt := uint64(baseamt) + uint64(amt)
	// 新余额
	blssto.Amount = fields.VarUint8(newamt)
	bserr := state.SatoshiSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

// 单纯扣除 BTC 余额 （amt 单位 聪）
func DoSubSatoshiFromChainState(state interfaces.ChainStateOperation, addr fields.Address, amt fields.VarUint8) error {
	blssto := state.Satoshi(addr)
	baseamt := blssto.Amount
	// 检查余额
	if uint64(baseamt) < uint64(amt) {
		return fmt.Errorf("address %s satoshi %d not enough, need more %d.", addr.ToReadable(), baseamt, amt)
	}
	newamt := uint64(baseamt) - uint64(amt)
	blssto.Amount = fields.VarUint8(newamt)
	bserr := state.SatoshiSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}
