package transactions

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/crypto/sha3"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"math/big"
)

type Transaction_0_Coinbase struct {
	// type
	Address fields.Address
	Reward  fields.Amount
	Message fields.TrimString16
	// 版本号
	BodyVersion fields.VarUint1 // 220000高度之前都等于 0

	// 当 BodyVersion >= 1 时 具有以下字段：
	MinerNonce   fields.Bytes32
	WitnessCount fields.VarUint1 // 投票见证人数量
	WitnessSigs  []uint8         // 见证人指定哈希尾数
	Witnesses    []fields.Sign   // 对prev区块hash的签名，投票分叉

	/* -------- -------- */

	// cache data
	TotalFeeUserPayed     fields.Amount // 区块总交易手续费
	TotalFeeMinerReceived fields.Amount // 区块总交易手续费
}

func NewTransaction_0_CoinbaseV0() *Transaction_0_Coinbase {
	return &Transaction_0_Coinbase{
		BodyVersion: 0,
	}
}

func NewTransaction_0_CoinbaseV1() *Transaction_0_Coinbase {
	return &Transaction_0_Coinbase{
		BodyVersion:  1,
		MinerNonce:   make([]byte, 32),
		WitnessCount: 0,
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
	bodyVersion := uint8(trs.BodyVersion)
	cbinfo["body_version"] = bodyVersion
	if bodyVersion >= 1 {
		if !isForMining {
			cbinfo["miner_nonce"] = trs.MinerNonce.ToHex() // 不需要
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

func (trs *Transaction_0_Coinbase) Copy() interfaces.Transaction {
	// copy
	bodys, _ := trs.Serialize()
	newtrsbts := make([]byte, len(bodys))
	copy(newtrsbts, bodys)
	// create
	var newtrs = new(Transaction_0_Coinbase)
	newtrs.Parse(newtrsbts, 1) // over type
	return newtrs
	/*
		return &Transaction_0_Coinbase{
			Address:      append([]byte{}, trs.Address...),
			Reward:       *trs.Reward.Copy(),
			Message:      fields.TrimString16(string(append([]byte{}, trs.Message...))),
			WitnessCount: trs.WitnessCount,
			WitnessSigs:  append([]uint8{}, trs.WitnessSigs...),
			Witnesses:    append([]fields.Sign{}, trs.Witnesses...),
		}
	*/
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
	b4, _ := trs.BodyVersion.Serialize()
	// fmt.Println("trs.Message=", trs.Message)
	buffer.Write([]byte{trs.Type()}) // type
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	// 版本号
	version := uint8(trs.BodyVersion)
	// ----------- 按版本区分 ----------- //
	if version == 0 {
		// 没有后续字段
	}
	if version >= 1 {
		// 附加的 nonce 值
		buffer.Write(trs.MinerNonce)
		// 见证人
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
	seek, e = trs.BodyVersion.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	// ----------- 按版本区分 ----------- //
	// 判断版本号
	version := uint8(trs.BodyVersion)
	if version == 0 {
		// 没有后续字段
	}
	if version >= 1 {
		// nonce值
		seek, e = trs.MinerNonce.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		// 见证人
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
		trs.BodyVersion.Size()
	// ----------- 按版本区分 ----------- //
	// 判断版本号
	version := uint8(trs.BodyVersion)
	if version == 0 {
		// 没有后续字段
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

// 交易唯一哈希值
func (trs *Transaction_0_Coinbase) HashWithFee() fields.Hash {
	stuff, _ := trs.Serialize()
	digest := sha3.Sum256(stuff)
	return digest[:]
}

func (trs *Transaction_0_Coinbase) Hash() fields.Hash {
	return trs.HashWithFee()
}

// 需要的余额检查
func (trs *Transaction_0_Coinbase) RequestAddressBalance() ([][]byte, []big.Int, error) {
	panic("cannot RequestAddressBalance for Transaction_0_Coinbase")
}

// 从 actions 拿出需要签名的地址
func (trs *Transaction_0_Coinbase) RequestSignAddresses([]fields.Address, bool) ([]fields.Address, error) {
	panic("never call Transaction_0_Coinbase.RequestSignAddresses")
	return []fields.Address{}, nil
}
func (trs *Transaction_0_Coinbase) VerifyTargetSigns([]fields.Address) (bool, error) {
	panic("never call Transaction_0_Coinbase.VerifyTargetSigns")
	return true, nil
}

// 清除所有签名
func (trs *Transaction_0_Coinbase) CleanSigns() {
	panic("cannot CleanSigns for Transaction_0_Coinbase")
}

// 返回所有签名
func (trs *Transaction_0_Coinbase) GetSigns() []fields.Sign {
	panic("cannot GetSigns for Transaction_0_Coinbase")
}

// 设置签名
func (trs *Transaction_0_Coinbase) SetSigns([]fields.Sign) {
	panic("cannot SetSigns for Transaction_0_Coinbase")
}

// 填充全部需要的签名
func (trs *Transaction_0_Coinbase) FillTargetSign(*account.Account) error {
	panic("cannot FillTargetSign for Transaction_0_Coinbase")
}

// 填充签名
func (trs *Transaction_0_Coinbase) FillNeedSigns(map[string][]byte, []fields.Address) error {
	panic("cannot FillNeedSigns for Transaction_0_Coinbase")
}

// 验证需要的签名
func (trs *Transaction_0_Coinbase) VerifyAllNeedSigns() (bool, error) {
	return true, nil
}

// 修改 / 恢复 状态数据库
func (trs *Transaction_0_Coinbase) WriteinChainState(state interfaces.ChainStateOperation) error {

	// total supply 统计
	// reward
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBlockMinerReward, trs.Reward.ToMei())
	// feeBurning
	if trs.TotalFeeMinerReceived.Equal(&trs.TotalFeeUserPayed) == false {
		// 有销毁
		burnamt, e := trs.TotalFeeUserPayed.Sub(&trs.TotalFeeMinerReceived)
		if e != nil {
			return e
		}
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBurningFee, burnamt.ToMei())
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

func (trs *Transaction_0_Coinbase) RecoverChainState(state interfaces.ChainStateOperation) error {

	// total supply 统计
	// reward
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	totalsupply.DoSub(stores.TotalSupplyStoreTypeOfBlockMinerReward, trs.Reward.ToMei())
	// feeBurning
	if trs.TotalFeeMinerReceived.Equal(&trs.TotalFeeUserPayed) == false {
		// 有销毁
		burnamt, e := trs.TotalFeeUserPayed.Sub(&trs.TotalFeeMinerReceived)
		if e != nil {
			return e
		}
		totalsupply.DoSub(stores.TotalSupplyStoreTypeOfBurningFee, burnamt.ToMei())
	}
	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}
	// 地址余额
	rwd_and_txfee, _ := trs.Reward.Add(&trs.TotalFeeMinerReceived)
	return actions.DoSubBalanceFromChainState(state, trs.Address, *rwd_and_txfee)
}

func (trs *Transaction_0_Coinbase) FeePurity() uint64 {
	panic("cannot GetFeePurity for Transaction_0_Coinbase")
}

// 查询
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

func (trs *Transaction_0_Coinbase) GetActions() []interfaces.Action {
	return nil
}

func (trs *Transaction_0_Coinbase) GetTimestamp() uint64 { // 时间戳
	panic("cannot GetTimestamp for Transaction_0_Coinbase")
}

func (trs *Transaction_0_Coinbase) SetMessage(msg fields.TrimString16) {
	trs.Message = msg
}

func (trs *Transaction_0_Coinbase) GetMessage() fields.TrimString16 {
	return trs.Message
}
