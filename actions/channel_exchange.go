package actions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
)

////////////////////////////////

// 通道资金与链上资金相互转移，原子互换

type ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange struct {
	ChannelTranferProveBodyHashChecker fields.HashHalfChecker

	OnChainTranferToAddress fields.Address // 链上转账收款地址
	OnChainTranferAmount    fields.Amount  // 链上转账数额

	AddressCount                            fields.VarUint1 // 签名数量，只能取值 2 或 3
	OnchainTransferFromAndMustSignAddresses []fields.Address
	// 两个或三个地址，其中第一个地址必须为链上转账 From 地址
	// 地址列表里必须包含通道双方的地址，否则提交挑战和仲裁时将验证失败

	// 地址对应的签名
	MustSigns []fields.Sign // 顺序与 []address 顺序必须一致
}

func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) Size() uint32 {
	size := elm.ChannelTranferProveBodyHashChecker.Size() +
		elm.OnChainTranferToAddress.Size() +
		elm.OnChainTranferAmount.Size() +
		elm.AddressCount.Size()
	size += uint32(len(elm.OnchainTransferFromAndMustSignAddresses)) * fields.AddressSize
	size += uint32(len(elm.MustSigns)) * fields.SignSize
	return size
}

func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) SerializeNoSign() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.ChannelTranferProveBodyHashChecker.Serialize()
	var bt2, _ = elm.OnChainTranferToAddress.Serialize()
	var bt3, _ = elm.OnChainTranferAmount.Serialize()
	var bt4, _ = elm.AddressCount.Serialize()
	buffer.Write(bt1)
	buffer.Write(bt2)
	buffer.Write(bt3)
	buffer.Write(bt4)
	for _, addr := range elm.OnchainTransferFromAndMustSignAddresses {
		var bt1, _ = addr.Serialize()
		buffer.Write(bt1)
	}
	return buffer.Bytes(), nil
}

func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) SignStuffHash() (fields.Hash, error) {
	var conbt, e = elm.SerializeNoSign() // 数据体
	if e != nil {
		return nil, e
	}
	return fields.CalculateHash(conbt), nil
}

func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.SerializeNoSign() // 数据体
	buffer.Write(bt1)
	for _, sign := range elm.MustSigns {
		var bt1, _ = sign.Serialize()
		buffer.Write(bt1)
	}
	return buffer.Bytes(), nil
}

func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelTranferProveBodyHashChecker.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.OnChainTranferToAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.OnChainTranferAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	// 地址
	seek, e = elm.AddressCount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	scn := int(elm.AddressCount)
	elm.OnchainTransferFromAndMustSignAddresses = make([]fields.Address, scn)
	for i := 0; i < scn; i++ {
		seek, e = elm.OnchainTransferFromAndMustSignAddresses[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// 签名
	elm.MustSigns = make([]fields.Sign, scn)
	for i := 0; i < scn; i++ {
		seek, e = elm.MustSigns[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// 完成
	return seek, nil
}

// 检查所有签名
func (elm *ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange) CheckMustAddressAndSigns() error {
	var e error

	// 计算哈希
	conhx, e := elm.SignStuffHash()
	if e != nil {
		return e
	}

	// 检查地址最低数量
	sgmn := len(elm.OnchainTransferFromAndMustSignAddresses)
	if sgmn < 2 || sgmn > 3 || sgmn != int(elm.AddressCount) || sgmn != len(elm.MustSigns) {
		return fmt.Errorf("Address or Sign length error, need 2~3 but got %d, %d, %d.",
			sgmn, int(elm.AddressCount), len(elm.MustSigns))
	}

	// 签名按地址排列，检查所有地址和签名是否匹配
	for i := 0; i < sgmn; i++ {
		sign := elm.MustSigns[i]
		addr := elm.OnchainTransferFromAndMustSignAddresses[i]
		sgaddr := account.NewAddressFromPublicKeyV0(sign.PublicKey)
		// 判断地址顺序
		if addr.NotEqual(sgaddr) {
			return fmt.Errorf("Address not match, need %s nut got %s.",
				addr.ToReadable(), fields.Address(sgaddr).ToReadable())
		}
		// 检查签名
		ok, _ := account.CheckSignByHash32(conhx, sign.PublicKey, sign.Signature)
		if !ok {
			return fmt.Errorf("Left account %s verify signature fail.", addr.ToReadable())
		}
	}

	// 全部签名验证成功
	return nil
}

////////////////////////////////////////////////////////

// 通道与链上原子互换
// 提交通道与链上互换交易
type Action_25_PaymantChannelAndOnchainAtomicExchange struct {

	// 原子互换交易凭证
	ExchangeEvidence ChannelAmountAndOnChainAmountTransferEachOtherByAtomicExchange

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) Kind() uint16 {
	return 25
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) Size() uint32 {
	return 2 + elm.ExchangeEvidence.Size()
}

// json api
func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.ExchangeEvidence.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	return buffer.Bytes(), nil
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ExchangeEvidence.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) RequestSignAddresses() []fields.Address {
	// action 内部判断签名
	return []fields.Address{}
}

func (act *Action_25_PaymantChannelAndOnchainAtomicExchange) WriteinChainState(state interfaces.ChainStateOperation) error {

	var e error

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 查询是否为重复提交
	swaphx := act.ExchangeEvidence.ChannelTranferProveBodyHashChecker
	chaswap, e := state.Chaswap(swaphx)
	if e != nil {
		return e
	}
	if chaswap != nil {
		// 已经存在，不可重复提交
		// 否则将会导致多次重复转账
		return fmt.Errorf("ChannelTranferProveBodyHashChecker <%s> is existence.",
			swaphx.ToHex())
	}

	if len(act.ExchangeEvidence.OnchainTransferFromAndMustSignAddresses) < 2 {
		return fmt.Errorf("Address lenght error.")
	}

	// 提交时不验证通道相关内容，仅仅操作链上转账
	e = act.ExchangeEvidence.CheckMustAddressAndSigns()
	if e != nil {
		return e
	}

	// 创建，保存凭证
	objsto := stores.Chaswap{
		IsBeUsed:                                fields.CreateBool(false), // 未使用过
		AddressCount:                            act.ExchangeEvidence.AddressCount,
		OnchainTransferFromAndMustSignAddresses: act.ExchangeEvidence.OnchainTransferFromAndMustSignAddresses,
	}
	e = state.ChaswapCreate(swaphx, &objsto)
	if e != nil {
		return e
	}

	// 转账
	fromAddr := act.ExchangeEvidence.OnchainTransferFromAndMustSignAddresses[0]
	toAddr := act.ExchangeEvidence.OnChainTranferToAddress
	trsAmt := act.ExchangeEvidence.OnChainTranferAmount
	return DoSimpleTransferFromChainState(state, fromAddr, toAddr, trsAmt)
}

func (act *Action_25_PaymantChannelAndOnchainAtomicExchange) RecoverChainState(state interfaces.ChainStateOperation) error {

	// 取消保存
	swaphx := act.ExchangeEvidence.ChannelTranferProveBodyHashChecker
	state.ChaswapDelete(swaphx)

	// 转账回退
	fromAddr := act.ExchangeEvidence.OnchainTransferFromAndMustSignAddresses[0]
	toAddr := act.ExchangeEvidence.OnChainTranferToAddress
	trsAmt := act.ExchangeEvidence.OnChainTranferAmount
	return DoSimpleTransferFromChainState(state, toAddr, fromAddr, trsAmt)
}

func (elm *Action_25_PaymantChannelAndOnchainAtomicExchange) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_25_PaymantChannelAndOnchainAtomicExchange) IsBurning90PersentTxFees() bool {
	return false
}
