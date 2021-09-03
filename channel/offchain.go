package channel

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

// 实时对账单接口
type PaymentChannelRealtimeReconciliationInterface interface {
	GetChannelId() fields.Bytes16
	GetLeftAmount() fields.Amount                                                                // 左侧实时金额
	GetRightAmount() fields.Amount                                                               // 右侧实时金额
	GetChannelReuseVersion() fields.VarUint4                                                     // 通道重用序号
	GetBillAutoNumber() fields.VarUint8                                                          // 对账序号
	CheckAddressAndSign(state interfaces.ChainStateOperation, laddr, raddr fields.Address) error // 检查地址和签名
}

////////////////////////////////

const (
	ChannelTransferProveBodyPayModeNormal  fields.VarUint1 = 1
	ChannelTransferProveBodyPayModeFastPay fields.VarUint1 = 2
)

// 通道转账，数据体
type ChannelChainTransferProveBodyInfo struct {
	ChannelId fields.Bytes16 // 通道id

	Mode fields.VarUint1 // 模式，普通模式、快速转账模式等等

	// Mode = 1 // 普通模式（实时对账）可以上链仲裁，但有负值是则为信用透支额度，债权方仅可仲裁无信用部分，即债务清除
	LeftAmount  fields.Amount // 左侧实时金额
	RightAmount fields.Amount // 右侧实时金额

	// Mode = 2 // 快速支付模式（延迟对账）无法上链仲裁
	Direction fields.VarUint1 // 资金流动方向： 1.左侧支付给右侧； 2.右侧付给左侧
	PayAmount fields.Amount   // 左侧支付给右侧的金额，如果为负值则表示右侧支付给左侧

	ChannelReuseVersion fields.VarUint4 // 通道重用序号
	BillAutoNumber      fields.VarUint8 // 通道账单流水序号
}

// interface
func (e *ChannelChainTransferProveBodyInfo) GetChannelId() fields.Bytes16 {
	return e.ChannelId
}
func (e *ChannelChainTransferProveBodyInfo) GetLeftAmount() fields.Amount {
	return e.LeftAmount
}
func (e *ChannelChainTransferProveBodyInfo) GetRightAmount() fields.Amount {
	return e.RightAmount
}
func (e *ChannelChainTransferProveBodyInfo) GetChannelReuseVersion() fields.VarUint4 {
	return e.ChannelReuseVersion
}
func (e *ChannelChainTransferProveBodyInfo) GetBillAutoNumber() fields.VarUint8 {
	return e.BillAutoNumber
}

func (elm *ChannelChainTransferProveBodyInfo) Size() uint32 {
	size := elm.ChannelId.Size() +
		elm.Mode.Size()
	if elm.Mode == ChannelTransferProveBodyPayModeNormal {
		size += elm.LeftAmount.Size() +
			elm.RightAmount.Size()
	} else if elm.Mode == ChannelTransferProveBodyPayModeFastPay {
		size += elm.Direction.Size() +
			elm.PayAmount.Size()
	}
	size += elm.ChannelReuseVersion.Size() +
		elm.BillAutoNumber.Size()
	// ok
	return size
}

func (elm *ChannelChainTransferProveBodyInfo) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.ChannelId.Serialize()
	var bt2, _ = elm.Mode.Serialize()
	buffer.Write(bt1)
	buffer.Write(bt2)
	if elm.Mode == ChannelTransferProveBodyPayModeNormal {
		var bt1, _ = elm.LeftAmount.Serialize()
		var bt2, _ = elm.RightAmount.Serialize()
		buffer.Write(bt1)
		buffer.Write(bt2)
	} else if elm.Mode == ChannelTransferProveBodyPayModeFastPay {
		var bt1, _ = elm.Direction.Serialize()
		var bt2, _ = elm.PayAmount.Serialize()
		buffer.Write(bt1)
		buffer.Write(bt2)
	}
	var bt3, _ = elm.ChannelReuseVersion.Serialize()
	var bt4, _ = elm.BillAutoNumber.Serialize()
	buffer.Write(bt3)
	buffer.Write(bt4)
	return buffer.Bytes(), nil
}

func (elm *ChannelChainTransferProveBodyInfo) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.Mode.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if elm.Mode == ChannelTransferProveBodyPayModeNormal {
		seek, e = elm.LeftAmount.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		seek, e = elm.RightAmount.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	} else if elm.Mode == ChannelTransferProveBodyPayModeFastPay {
		seek, e = elm.Direction.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		seek, e = elm.PayAmount.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	seek, e = elm.ChannelReuseVersion.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.BillAutoNumber.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *ChannelChainTransferProveBodyInfo) SignStuff() []byte {
	var conbt, _ = elm.Serialize() // 数据体
	return conbt                   // 哈希
}
func (elm *ChannelChainTransferProveBodyInfo) SignStuffHashHalfChecker() fields.HashHalfChecker {
	var conbt = elm.SignStuff()                         // 数据体
	return fields.CalculateHash(conbt).GetHalfChecker() // 哈希检测
}

// 检查签名
func (elm *ChannelChainTransferProveBodyInfo) CheckAddressAndSign(state interfaces.ChainStateOperation, leftAddress, rightAddress fields.Address) error {
	// 全部检查成功
	return nil
}

////////////////////////////////

type ChannelPayProveBodyList struct {
	Count      fields.VarUint1
	ProveBodys []*ChannelChainTransferProveBodyInfo
}

func (c ChannelPayProveBodyList) Size() uint32 {
	size := c.Count.Size()
	for i := 0; i < int(c.Count); i++ {
		size += c.ProveBodys[i].Size()
	}
	// ok
	return size
}

func (c ChannelPayProveBodyList) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = c.Count.Serialize() // 数据体
	buffer.Write(bt1)
	for i := 0; i < len(c.ProveBodys); i++ {
		var bt6, _ = c.ProveBodys[i].Serialize()
		buffer.Write(bt6)
	}
	return buffer.Bytes(), nil
}

func (c *ChannelPayProveBodyList) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	// 通道
	seek, e = c.Count.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	ccn := int(c.Count)
	c.ProveBodys = make([]*ChannelChainTransferProveBodyInfo, ccn)
	for i := 0; i < ccn; i++ {
		c.ProveBodys[i] = &ChannelChainTransferProveBodyInfo{}
		seek, e = c.ProveBodys[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// 完成
	return seek, nil
}

////////////////////////////////

// 通道实时对账（链下签署）
type OffChainFormPaymentChannelRealtimeReconciliation struct {
	TranferProveBody ChannelChainTransferProveBodyInfo
	Timestamp        fields.BlockTxTimestamp // 对账时间戳

	// 两侧签名
	LeftSign  fields.Sign // 左侧地址对账签名
	RightSign fields.Sign // 右侧地址对账签名
}

// interface
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetChannelId() fields.Bytes16 {
	return e.TranferProveBody.ChannelId
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetLeftAmount() fields.Amount {
	return e.TranferProveBody.LeftAmount
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetRightAmount() fields.Amount {
	return e.TranferProveBody.RightAmount
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetChannelReuseVersion() fields.VarUint4 {
	return e.TranferProveBody.ChannelReuseVersion
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetBillAutoNumber() fields.VarUint8 {
	return e.TranferProveBody.BillAutoNumber
}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) Size() uint32 {
	return elm.TranferProveBody.Size() +
		elm.Timestamp.Size() +
		elm.LeftSign.Size() +
		elm.RightSign.Size()
}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) SerializeNoSign() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.TranferProveBody.Serialize()
	var bt2, _ = elm.Timestamp.Serialize()
	buffer.Write(bt1)
	buffer.Write(bt2)
	return buffer.Bytes(), nil

}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) Serialize() ([]byte, error) {
	var bt1, _ = elm.SerializeNoSign() // 数据体
	var bt2, _ = elm.LeftSign.Serialize()
	var bt3, _ = elm.RightSign.Serialize()
	var buffer bytes.Buffer
	buffer.Write(bt1)
	buffer.Write(bt2)
	buffer.Write(bt3)
	return buffer.Bytes(), nil
}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) SignStuffHash() fields.Hash {
	var conbt, _ = elm.SerializeNoSign() // 数据体
	return fields.CalculateHash(conbt)
}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.TranferProveBody.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.Timestamp.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftSign.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightSign.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

// 检查签名
func (elm *OffChainFormPaymentChannelRealtimeReconciliation) CheckAddressAndSign(state interfaces.ChainStateOperation, leftAddress, rightAddress fields.Address) error {
	// 检查公钥和地址是否匹配
	addr1 := account.NewAddressFromPublicKeyV0(elm.LeftSign.PublicKey)
	if leftAddress.NotEqual(fields.Address(addr1)) {
		return fmt.Errorf("Left sign address %s is not request address %s",
			fields.Address(addr1).ToReadable(),
			leftAddress.ToReadable())
	}
	addr2 := account.NewAddressFromPublicKeyV0(elm.RightSign.PublicKey)
	if rightAddress.NotEqual(fields.Address(addr2)) {
		return fmt.Errorf("Right sign address %s is not request address %s",
			fields.Address(addr2).ToReadable(),
			rightAddress.ToReadable())
	}
	// 验证哈希
	var conhx = elm.SignStuffHash() // 数据体hx
	ok1, _ := account.CheckSignByHash32(conhx, elm.LeftSign.PublicKey, elm.LeftSign.Signature)
	if !ok1 {
		return fmt.Errorf("Left account %s verify signature fail.", leftAddress.ToReadable())
	}
	ok2, _ := account.CheckSignByHash32(conhx, elm.LeftSign.PublicKey, elm.LeftSign.Signature)
	if !ok2 {
		return fmt.Errorf("Right account %s verify signature fail.", leftAddress.ToReadable())
	}
	// 全部检查成功
	return nil
}

////////////////////////////////

// 通道链支付

// 通道链转账交易（链下签署）可以上链仲裁
// 采用零知识证明模式
type OffChainFormPaymentChannelTransfer struct {
	Timestamp                fields.BlockTxTimestamp // 时间戳
	OrderNoteHashHalfChecker fields.HashHalfChecker  // 订单详情数据哈希  len = 16

	MustSignCount     fields.VarUint1  // 必须签名地址的数量，最大值 200
	MustSignAddresses []fields.Address // 顺序打乱/随机的通道必须签名的地址

	ChannelCount                         fields.VarUint1          // 途径通道数量，最大值 200
	ChannelTransferProveHashHalfCheckers []fields.HashHalfChecker // 通道转账证明哈希，顺序为从付款到最后收款，哈希  len = 16

	MustSigns []fields.Sign // 顺序打乱/随机的签名，顺序与地址相同
}

func (elm *OffChainFormPaymentChannelTransfer) Size() uint32 {
	size := elm.Timestamp.Size() +
		elm.OrderNoteHashHalfChecker.Size() +
		elm.ChannelCount.Size() +
		elm.MustSignCount.Size()

	size += uint32(len(elm.ChannelTransferProveHashHalfCheckers)) * (fields.HashHalfCheckerSize)
	size += uint32(len(elm.MustSignAddresses)) * (fields.AddressSize)
	size += uint32(len(elm.MustSigns)) * fields.SignSize
	return size
}

func (elm *OffChainFormPaymentChannelTransfer) SerializeForPrefixSignStuff() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.Timestamp.Serialize()
	var bt2, _ = elm.OrderNoteHashHalfChecker.Serialize()
	buffer.Write(bt1)
	buffer.Write(bt2)
	var bt3, _ = elm.MustSignCount.Serialize()
	buffer.Write(bt3)
	for i := 0; i < len(elm.MustSignAddresses); i++ {
		var bt4, _ = elm.MustSignAddresses[i].Serialize()
		buffer.Write(bt4)
	}
	var bt4, _ = elm.ChannelCount.Serialize()
	buffer.Write(bt4)
	return buffer.Bytes(), nil

}

func (elm *OffChainFormPaymentChannelTransfer) SerializeNoSign() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.SerializeForPrefixSignStuff() // 数据体
	buffer.Write(bt1)
	for i := 0; i < len(elm.ChannelTransferProveHashHalfCheckers); i++ {
		var bt4, _ = elm.ChannelTransferProveHashHalfCheckers[i].Serialize()
		buffer.Write(bt4)
	}
	return buffer.Bytes(), nil

}

func (elm *OffChainFormPaymentChannelTransfer) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.SerializeNoSign() // 数据体
	buffer.Write(bt1)
	for i := 0; i < len(elm.MustSigns); i++ {
		var bt6, _ = elm.MustSigns[i].Serialize()
		buffer.Write(bt6)
	}
	return buffer.Bytes(), nil
}

func (elm *OffChainFormPaymentChannelTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.Timestamp.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.OrderNoteHashHalfChecker.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	// 地址
	scn := int(elm.MustSignCount)
	seek, e = elm.MustSignCount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	elm.MustSignAddresses = make([]fields.Address, scn)
	for i := 0; i < scn; i++ {
		seek, e = elm.MustSignAddresses[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// 通道
	seek, e = elm.ChannelCount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	ccn := int(elm.ChannelCount)
	elm.ChannelTransferProveHashHalfCheckers = make([]fields.HashHalfChecker, ccn)
	for i := 0; i < ccn; i++ {
		seek, e = elm.ChannelTransferProveHashHalfCheckers[i].Parse(buf, seek)
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

// 按位置填充签名
func (elm *OffChainFormPaymentChannelTransfer) FillSignByPosition(sign fields.Sign) error {
	sgaddr := sign.GetAddress()
	sn := int(elm.MustSignCount)
	var istok = false
	for i := 0; i < sn; i++ {
		addr := elm.MustSignAddresses[i]
		if addr.Equal(sgaddr) {
			istok = true
			elm.MustSigns[i] = sign
		}
	}
	if istok == false {
		return fmt.Errorf(" sign address %s not find in must list.", sgaddr.ToReadable())
	}
	return nil
}

// 检查所有签名
func (elm *OffChainFormPaymentChannelTransfer) CheckMustAddressAndSigns() error {
	var e error

	stuff, e := elm.SerializeNoSign()
	if e != nil {
		return e
	}
	conhx := fields.CalculateHash(stuff)

	// 检查数量
	sn := int(elm.MustSignCount)
	if sn < 2 || sn > 200 {
		fmt.Errorf("MustSignCount error.")
	}
	if sn != len(elm.MustSignAddresses) || sn != len(elm.MustSigns) {
		fmt.Errorf("Addresses or Signs length error.")
	}

	// 签名按地址排列，检查所有地址和签名是否匹配
	for i := 0; i < sn; i++ {
		sign := elm.MustSigns[i]
		addr := elm.MustSignAddresses[i]
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

/*******************************************/

// 通道链支付票据集合
type ChannelPayBillAssemble struct {
	// 对账票据表
	ProveBodys *ChannelPayProveBodyList
	// 支付签名票据
	ChainPayment *OffChainFormPaymentChannelTransfer
}

func (c ChannelPayBillAssemble) Size() uint32 {
	return c.ProveBodys.Size() + c.ChainPayment.Size()
}

func (c ChannelPayBillAssemble) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = c.ProveBodys.Serialize()
	buffer.Write(bt1)
	var bt2, _ = c.ChainPayment.Serialize()
	buffer.Write(bt2)
	return buffer.Bytes(), nil
}

func (c *ChannelPayBillAssemble) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	// 通道
	c.ProveBodys = &ChannelPayProveBodyList{}
	seek, e = c.ProveBodys.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	c.ChainPayment = &OffChainFormPaymentChannelTransfer{}
	seek, e = c.ChainPayment.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}
