package actions

import (
	"bytes"
	"encoding/binary"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"math/big"
)

type Action_7_SatoshiGenesis struct {
	TransferNo               fields.VarInt4 // 转账流水编号
	BitcoinBlockHeight       fields.VarInt4 // 转账的比特币区块高度
	BitcoinBlockTimestamp    fields.VarInt4 // 转账的比特币区块时间戳
	BitcoinEffectiveGenesis  fields.VarInt4 // 在这笔之前已经成功转移的比特币数量
	BitcoinQuantity          fields.VarInt4 // 本笔转账的比特币数量（单位：枚）
	AdditionalTotalHacAmount fields.VarInt4 // 本次转账[总共]应该增发的 hac 数量 （单位：枚）
	OriginAddress            fields.Address // 转出的比特币来源地址
	BitcoinTransferHash      fields.Hash    // 比特币转账交易哈希

	// data ptr
	belong_trs interfaces.Transaction
}

func NewAction_7_SatoshiGenesis() *Action_7_SatoshiGenesis {
	return &Action_7_SatoshiGenesis{}
}

func (elm *Action_7_SatoshiGenesis) Kind() uint16 {
	return 7
}

func (elm *Action_7_SatoshiGenesis) Size() uint32 {
	return 2 +
		elm.TransferNo.Size() +
		elm.BitcoinBlockHeight.Size() +
		elm.BitcoinBlockTimestamp.Size() +
		elm.BitcoinEffectiveGenesis.Size() +
		elm.BitcoinQuantity.Size() +
		elm.AdditionalTotalHacAmount.Size() +
		elm.OriginAddress.Size() +
		elm.BitcoinTransferHash.Size()
}

func (elm *Action_7_SatoshiGenesis) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var b1, _ = elm.TransferNo.Serialize()
	var b2, _ = elm.BitcoinBlockHeight.Serialize()
	var b3, _ = elm.BitcoinBlockTimestamp.Serialize()
	var b4, _ = elm.BitcoinEffectiveGenesis.Serialize()
	var b5, _ = elm.BitcoinQuantity.Serialize()
	var b6, _ = elm.AdditionalTotalHacAmount.Serialize()
	var b7, _ = elm.OriginAddress.Serialize()
	var b8, _ = elm.BitcoinTransferHash.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	buffer.Write(b7)
	buffer.Write(b8)
	return buffer.Bytes(), nil
}

func (elm *Action_7_SatoshiGenesis) Parse(buf []byte, seek uint32) (uint32, error) {
	sk1, _ := elm.TransferNo.Parse(buf, seek)
	sk2, _ := elm.BitcoinBlockHeight.Parse(buf, sk1)
	sk3, _ := elm.BitcoinBlockTimestamp.Parse(buf, sk2)
	sk4, _ := elm.BitcoinEffectiveGenesis.Parse(buf, sk3)
	sk5, _ := elm.BitcoinQuantity.Parse(buf, sk4)
	sk6, _ := elm.AdditionalTotalHacAmount.Parse(buf, sk5)
	sk7, _ := elm.OriginAddress.Parse(buf, sk6)
	sk8, _ := elm.BitcoinTransferHash.Parse(buf, sk7)
	return sk8, nil
}

func (*Action_7_SatoshiGenesis) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_7_SatoshiGenesis) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 增发 hac
	hacmeibig := (new(big.Int)).SetUint64(uint64(act.AdditionalTotalHacAmount))
	addhacamt, err := fields.NewAmountByBigIntWithUnit(hacmeibig, 248)
	if err != nil {
		return err
	}
	e1 := DoAddBalanceFromChainState(state, act.OriginAddress, *addhacamt)
	if e1 != nil {
		return e1
	}
	// 发行 btc
	satBTC := uint64(act.BitcoinQuantity) * 10000 * 10000 // 单位 聪
	return DoAddSatoshiFromChainState(state, act.belong_trs.GetAddress(), fields.VarInt8(satBTC))
}

func (act *Action_7_SatoshiGenesis) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 回退 hac
	hacmeibig := (new(big.Int)).SetUint64(uint64(act.AdditionalTotalHacAmount))
	addhacamt, err := fields.NewAmountByBigIntWithUnit(hacmeibig, 248)
	if err != nil {
		return err
	}
	e1 := DoSubBalanceFromChainState(state, act.OriginAddress, *addhacamt)
	if e1 != nil {
		return e1
	}
	// 扣除 btc
	satBTC := uint64(act.BitcoinQuantity) * 10000 * 10000 // 单位 聪
	return DoSubSatoshiFromChainState(state, act.belong_trs.GetAddress(), fields.VarInt8(satBTC))
}

// 设置所属 belone_trs
func (act *Action_7_SatoshiGenesis) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}

///////////////////////////////////////////////////////////////////////////////////////////////

type Action_8_SimpleSatoshiTransfer struct {
	Address fields.Address
	Amount  fields.VarInt8

	// data ptr
	belong_trs interfaces.Transaction
}

func NewAction_8_SimpleSatoshiTransfer(addr fields.Address, amt fields.VarInt8) *Action_8_SimpleSatoshiTransfer {
	return &Action_8_SimpleSatoshiTransfer{
		Address: addr,
		Amount:  amt,
	}
}

func (elm *Action_8_SimpleSatoshiTransfer) Kind() uint16 {
	return 1
}

func (elm *Action_8_SimpleSatoshiTransfer) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var addrBytes, _ = elm.Address.Serialize()
	var amtBytes, _ = elm.Amount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(addrBytes)
	buffer.Write(amtBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_8_SimpleSatoshiTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var moveseek, _ = elm.Address.Parse(buf, seek)
	var moveseek2, _ = elm.Amount.Parse(buf, moveseek)
	return moveseek2, nil
}

func (elm *Action_8_SimpleSatoshiTransfer) Size() uint32 {
	return 2 + elm.Address.Size() + elm.Amount.Size()
}

func (*Action_8_SimpleSatoshiTransfer) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_8_SimpleSatoshiTransfer) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 转移
	return DoSimpleSatoshiTransferFromChainState(state, act.belong_trs.GetAddress(), act.Address, act.Amount)
}

func (act *Action_8_SimpleSatoshiTransfer) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 回退
	return DoSimpleSatoshiTransferFromChainState(state, act.Address, act.belong_trs.GetAddress(), act.Amount)
}

// 设置所属 belone_trs
func (act *Action_8_SimpleSatoshiTransfer) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}
