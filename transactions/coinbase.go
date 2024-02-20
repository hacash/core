package transactions

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
	"math/big"
)

type Transaction_0_Coinbase struct {
	// type
	Address fields.Address
	Reward  fields.Amount
	Message fields.TrimString16
	// Version number
	ExtendDataVersion fields.VarUint1 // Equal to 0 before 220000 height

	// When extenddataversion > = 1, it has the following fields:
	MinerNonce   fields.Bytes32
	WitnessCount fields.VarUint1 // Number of voting witnesses
	WitnessSigs  []uint8         // Witness specified hash mantissa
	Witnesses    []fields.Sign   // Signature of prev block hash, voting fork

	/* -------- -------- */

	// cache data
	TotalFeeUserPayed     fields.Amount // Total transaction fee of the block
	TotalFeeMinerReceived fields.Amount // Total transaction fee of the block
}

func NewTransaction_0_CoinbaseV0() *Transaction_0_Coinbase {
	return &Transaction_0_Coinbase{
		ExtendDataVersion: 0,
	}
}

func NewTransaction_0_CoinbaseV1() *Transaction_0_Coinbase {
	return &Transaction_0_Coinbase{
		ExtendDataVersion: 1,
		MinerNonce:        make([]byte, 32),
		WitnessCount:      0,
	}
}

func (trs *Transaction_0_Coinbase) Describe(isUnitMei, isForMining bool) map[string]interface{} {
	cbinfo := make(map[string]interface{})
	cbinfo["address"] = trs.Address.ToReadable()
	cbinfo["reward"] = trs.Reward.ToMeiOrFinString(isUnitMei)
	msg, _ := trs.Message.Serialize()
	if isForMining {
		cbinfo["message_hex"] = hex.EncodeToString(msg)
	}
	cbinfo["message"] = trs.Message
	extendDataVersion := uint8(trs.ExtendDataVersion)
	cbinfo["extend_data_version"] = extendDataVersion
	if extendDataVersion >= 1 {
		if !isForMining {
			cbinfo["miner_nonce"] = trs.MinerNonce.ToHex() // unwanted
		}
		wcnum := int(trs.WitnessCount)
		cbinfo["witness_count"] = wcnum
		wtnum := make([]int, wcnum)
		wtsig := make([]string, wcnum)
		for i := 0; i < wcnum; i++ {
			wtnum[i] = int(trs.WitnessSigs[i])
			sigbts, _ := trs.Witnesses[i].Serialize()
			wtsig[i] = hex.EncodeToString(sigbts)
		}
		cbinfo["witness_sigs"] = wtnum
		cbinfo["witnesses"] = wtsig
	}
	return cbinfo
}

func (trs *Transaction_0_Coinbase) ClearHash() {

}

func (trs *Transaction_0_Coinbase) CopyForMining() *Transaction_0_Coinbase {
	// copy
	bodys, _ := trs.Serialize()
	newtrsbts := make([]byte, len(bodys))
	copy(newtrsbts, bodys)
	// create
	var newtrs = new(Transaction_0_Coinbase)
	newtrs.Parse(newtrsbts, 1) // over type
	return newtrs
}

func (trs *Transaction_0_Coinbase) Copy() interfacev2.Transaction {
	return trs.CopyForMining()
}

func (trs *Transaction_0_Coinbase) Clone() interfaces.Transaction {
	return trs.CopyForMining()
}

func (trs *Transaction_0_Coinbase) GetReward() *fields.Amount {
	return &trs.Reward
}

func (trs *Transaction_0_Coinbase) Type() uint8 {
	return 0
}

func (trs *Transaction_0_Coinbase) Serialize() ([]byte, error) {

	var buffer bytes.Buffer
	b1, _ := trs.Address.Serialize()
	b2, _ := trs.Reward.Serialize()
	b3, _ := trs.Message.Serialize()
	b4, _ := trs.ExtendDataVersion.Serialize()
	// fmt.Println("trs.Message=", trs.Message)
	buffer.Write([]byte{trs.Type()}) // type
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	// Version number
	version := uint8(trs.ExtendDataVersion)
	// ----------- 按版本区分 ----------- //
	if version == 0 {
		// No subsequent fields
	}
	if version >= 1 {
		// Additional nonce values
		buffer.Write(trs.MinerNonce)
		// witness
		witnessCount := uint8(trs.WitnessCount)
		buffer.Write([]byte{witnessCount})
		for i := uint8(0); i < witnessCount; i++ {
			b := trs.WitnessSigs[i]
			buffer.Write([]byte{b})
		}
		for i := uint8(0); i < witnessCount; i++ {
			b := trs.Witnesses[i]
			s1, e := b.Serialize()
			if e != nil {
				return nil, e
			}
			buffer.Write(s1)
		}
	}
	return buffer.Bytes(), nil
}

func (trs *Transaction_0_Coinbase) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = trs.ParseHead(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = trs.ExtendDataVersion.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	// ----------- 按版本区分 ----------- //
	// Judge version number
	version := uint8(trs.ExtendDataVersion)
	if version == 0 {
		// No subsequent fields
	}
	if version >= 1 {
		// Nonce value
		seek, e = trs.MinerNonce.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		// witness
		seek, e = trs.WitnessCount.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		if trs.WitnessCount > 0 {
			lenwc := int(trs.WitnessCount)
			trs.WitnessSigs = make([]uint8, lenwc)
			trs.Witnesses = make([]fields.Sign, lenwc)
			for i := 0; i < lenwc; i++ {
				if seek >= uint32(len(buf)) {
					return 0, fmt.Errorf("seek out of buf len.")
				}
				trs.WitnessSigs[i] = buf[seek]
				seek++
			}
			for i := 0; i < lenwc; i++ {
				var sign fields.Sign
				if seek >= uint32(len(buf)) {
					return 0, fmt.Errorf("seek out of buf len.")
				}
				seek, e = sign.Parse(buf, seek)
				if e != nil {
					return 0, e
				}
				trs.Witnesses[i] = sign
			}
		}
	}
	return seek, nil
}

func (trs *Transaction_0_Coinbase) ParseHead(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = trs.Address.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = trs.Reward.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = trs.Message.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (trs *Transaction_0_Coinbase) Size() uint32 {
	base := 1 +
		trs.Address.Size() +
		trs.Reward.Size() +
		trs.Message.Size() +
		trs.ExtendDataVersion.Size()
	// ----------- 按版本区分 ----------- //
	// Judge version number
	version := uint8(trs.ExtendDataVersion)
	if version == 0 {
		// No subsequent fields
	}
	if version >= 1 {
		base += 32                      // nonce
		base += 1                       // WitnessCount
		length := int(trs.WitnessCount) // WitnessSigs size
		base += uint32(length)
		for i := 0; i < length; i++ {
			base += trs.Witnesses[i].Size()
		}
	}
	return base
}

// Transaction unique hash value
func (trs *Transaction_0_Coinbase) HashWithFee() fields.Hash {
	stuff, _ := trs.Serialize()
	digest := fields.CalculateHash(stuff)
	return digest
}

func (trs *Transaction_0_Coinbase) Hash() fields.Hash {
	return trs.HashWithFee()
}

// Balance check required
func (trs *Transaction_0_Coinbase) RequestAddressBalance() ([][]byte, []big.Int, error) {
	panic("cannot RequestAddressBalance for Transaction_0_Coinbase")
}

// Get the address to be signed from actions
func (trs *Transaction_0_Coinbase) RequestSignAddresses([]fields.Address, bool) ([]fields.Address, error) {
	panic("never call Transaction_0_Coinbase.RequestSignAddresses")
	return []fields.Address{}, nil
}
func (trs *Transaction_0_Coinbase) VerifyTargetSigns([]fields.Address) (bool, error) {
	panic("never call Transaction_0_Coinbase.VerifyTargetSigns")
	return true, nil
}

// Clear all signatures
func (trs *Transaction_0_Coinbase) CleanSigns() {
	panic("cannot CleanSigns for Transaction_0_Coinbase")
}

// Return all signatures
func (trs *Transaction_0_Coinbase) GetSigns() []fields.Sign {
	panic("cannot GetSigns for Transaction_0_Coinbase")
}

// Set signature
func (trs *Transaction_0_Coinbase) SetSigns([]fields.Sign) {
	panic("cannot SetSigns for Transaction_0_Coinbase")
}

// Fill in all required signatures
func (trs *Transaction_0_Coinbase) FillTargetSign(*account.Account) error {
	panic("cannot FillTargetSign for Transaction_0_Coinbase")
}

// Fill in signature
func (trs *Transaction_0_Coinbase) FillNeedSigns(map[string][]byte, []fields.Address) error {
	panic("cannot FillNeedSigns for Transaction_0_Coinbase")
}

// Verify required signatures
func (trs *Transaction_0_Coinbase) VerifyAllNeedSigns() (bool, error) {
	return true, nil
}

func (trs *Transaction_0_Coinbase) WriteInChainState(state interfaces.ChainStateOperation) error {

	// Total supply statistics
	// reward
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBlockReward, trs.Reward.ToMei())
	// feeBurning
	if trs.TotalFeeMinerReceived.NotEqual(&trs.TotalFeeUserPayed) {
		// With destruction
		burnamt, e := trs.TotalFeeUserPayed.Sub(&trs.TotalFeeMinerReceived)
		if e != nil {
			return e
		}
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBurningTotal, burnamt.ToMei())
	}
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	// fmt.Printf("trs.TotalFee = %s\n", trs.TotalFee.ToFinString())
	rwd_and_txfee, _ := trs.Reward.Add(&trs.TotalFeeMinerReceived)
	// addr, _ := base58check.Encode(trs.Address)
	// fmt.Printf("coinbase.ChangeChainState,  %s  +=  %s\n", addr, rwd.ToFinString())
	return actions.DoAddBalanceFromChainState(state, trs.Address, *rwd_and_txfee)
}

// 修改 / 恢复 状态数据库
func (trs *Transaction_0_Coinbase) WriteinChainState(state interfacev2.ChainStateOperation) error {

	// Total supply statistics
	// reward
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBlockReward, trs.Reward.ToMei())
	// feeBurning
	if trs.TotalFeeMinerReceived.NotEqual(&trs.TotalFeeUserPayed) {
		// With destruction
		burnamt, e := trs.TotalFeeUserPayed.Sub(&trs.TotalFeeMinerReceived)
		if e != nil {
			return e
		}
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBurningTotal, burnamt.ToMei())
	}
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	// fmt.Printf("trs.TotalFee = %s\n", trs.TotalFee.ToFinString())
	rwd_and_txfee, _ := trs.Reward.Add(&trs.TotalFeeMinerReceived)
	// addr, _ := base58check.Encode(trs.Address)
	// fmt.Printf("coinbase.ChangeChainState,  %s  +=  %s\n", addr, rwd.ToFinString())
	return actions.DoAddBalanceFromChainStateV2(state, trs.Address, *rwd_and_txfee)
}

func (trs *Transaction_0_Coinbase) RecoverChainState(state interfacev2.ChainStateOperation) error {

	panic("RecoverChainState be deprecated")

	// Total supply statistics
	// reward
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	totalsupply.DoSub(stores.TotalSupplyStoreTypeOfBlockReward, trs.Reward.ToMei())
	// feeBurning
	if trs.TotalFeeMinerReceived.NotEqual(&trs.TotalFeeUserPayed) {
		// With destruction
		burnamt, e := trs.TotalFeeUserPayed.Sub(&trs.TotalFeeMinerReceived)
		if e != nil {
			return e
		}
		totalsupply.DoSub(stores.TotalSupplyStoreTypeOfBurningTotal, burnamt.ToMei())
	}
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	// Address balance
	rwd_and_txfee, _ := trs.Reward.Add(&trs.TotalFeeMinerReceived)
	return actions.DoSubBalanceFromChainStateV2(state, trs.Address, *rwd_and_txfee)
}

func (trs *Transaction_0_Coinbase) FeePurity() uint32 {
	panic("cannot GetFeePurity for Transaction_0_Coinbase")
}

// query
func (trs *Transaction_0_Coinbase) GetAddress() fields.Address {
	return trs.Address
}

func (trs *Transaction_0_Coinbase) SetAddress(addr fields.Address) {
	trs.Address = addr
}

func (trs *Transaction_0_Coinbase) GetFeeOfMinerRealReceived() *fields.Amount {
	return &trs.TotalFeeMinerReceived
}

func (trs *Transaction_0_Coinbase) GetFee() *fields.Amount {
	return &trs.TotalFeeUserPayed
}

func (trs *Transaction_0_Coinbase) SetFee(fee *fields.Amount) {
	panic("cannot SetFee for Transaction_0_Coinbase")
}

func (trs *Transaction_0_Coinbase) GetActions() []interfacev2.Action {
	return nil
}
func (trs *Transaction_0_Coinbase) GetActionList() []interfaces.Action {
	return nil
}

func (trs *Transaction_0_Coinbase) GetTimestamp() uint64 { // time stamp
	panic("cannot GetTimestamp for Transaction_0_Coinbase")
}

func (trs *Transaction_0_Coinbase) SetMessage(msg fields.TrimString16) {
	trs.Message = msg
}

func (trs *Transaction_0_Coinbase) GetMessage() fields.TrimString16 {
	return trs.Message
}
