package actions

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/crypto/sha3"
	"github.com/hacash/core/fields"
)

////////////////////////////////

// 通道实时对账（链下签署）
type OffChainFormPaymentChannelRealtimeReconciliation struct {
	ChannelId   fields.Bytes16          // 通道id
	LeftAmount  fields.Amount           // 左侧实时金额
	RightAmount fields.Amount           // 右侧实时金额
	Timestamp   fields.BlockTxTimestamp // 对账时间戳
	AutoNumber  fields.VarUint8         // 对账序号
	// 两侧签名
	LeftSign  fields.Sign // 左侧地址对账签名
	RightSign fields.Sign // 右侧地址对账签名
}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) Size() uint32 {
	return elm.ChannelId.Size() +
		elm.LeftAmount.Size() +
		elm.RightAmount.Size() +
		elm.Timestamp.Size() +
		elm.AutoNumber.Size() +
		elm.LeftSign.Size() +
		elm.RightSign.Size()
}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) SerializeNoSign() ([]byte, error) {
	var bt1, _ = elm.ChannelId.Serialize()
	var bt2, _ = elm.LeftAmount.Serialize()
	var bt3, _ = elm.RightAmount.Serialize()
	var bt4, _ = elm.Timestamp.Serialize()
	var bt5, _ = elm.AutoNumber.Serialize()
	var buffer bytes.Buffer
	buffer.Write(bt1)
	buffer.Write(bt2)
	buffer.Write(bt3)
	buffer.Write(bt4)
	buffer.Write(bt5)
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

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) SignStuffHash() []byte {
	var conbt, _ = elm.SerializeNoSign() // 数据体
	conhx := sha3.Sum256(conbt)
	return conhx[:]
}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.Timestamp.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.AutoNumber.Parse(buf, seek)
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
func (elm *OffChainFormPaymentChannelRealtimeReconciliation) CheckAddressAndSign(leftAddress, rightAddress fields.Address) error {
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
type OffChainFormPaymentChannelTransfer struct {
	Timestamp     fields.BlockTxTimestamp // 时间戳
	OrderNoteHash fields.Hash             // 订单详情数据哈希

	ChannelCount              fields.VarUint1 // 途径通道数量，最大值 200
	ChannelTransferProveHashs []fields.Hash   // 通道转账证明哈希，顺序为从付款到最后收款

	SignCount fields.VarUint1 // 签名数量，最大值 200
	Signs     []fields.Sign   // 有序签名，顺序为从付款到最后收款
}

func (elm *OffChainFormPaymentChannelTransfer) Size() uint32 {
	size := elm.Timestamp.Size() +
		elm.OrderNoteHash.Size() +
		elm.ChannelCount.Size() +
		elm.SignCount.Size()

	size += uint32(len(elm.ChannelTransferProveHashs)) * (fields.Hash{}.Size())
	size += uint32(len(elm.Signs)) * fields.SignSize
	return size
}

func (elm *OffChainFormPaymentChannelTransfer) SerializeNoSign() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.Timestamp.Serialize()
	var bt2, _ = elm.OrderNoteHash.Serialize()
	buffer.Write(bt1)
	buffer.Write(bt2)
	var bt3, _ = elm.ChannelCount.Serialize()
	buffer.Write(bt3)
	for i := 0; i < len(elm.ChannelTransferProveHashs); i++ {
		var bt4, _ = elm.ChannelTransferProveHashs[i].Serialize()
		buffer.Write(bt4)
	}
	return buffer.Bytes(), nil

}

func (elm *OffChainFormPaymentChannelTransfer) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.SerializeNoSign() // 数据体
	buffer.Write(bt1)
	var bt5, _ = elm.SignCount.Serialize()
	buffer.Write(bt5)
	for i := 0; i < len(elm.Signs); i++ {
		var bt6, _ = elm.Signs[i].Serialize()
		buffer.Write(bt6)
	}
	return buffer.Bytes(), nil
}

func (elm *OffChainFormPaymentChannelTransfer) SignStuffHash() []byte {
	var conbt, _ = elm.SerializeNoSign() // 数据体
	conhx := sha3.Sum256(conbt)
	return conhx[:]
}

func (elm *OffChainFormPaymentChannelTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.Timestamp.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.OrderNoteHash.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	// 通道
	seek, e = elm.ChannelCount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	ccn := int(elm.ChannelCount)
	elm.ChannelTransferProveHashs = make([]fields.Hash, ccn)
	for i := 0; i < ccn; i++ {
		seek, e = elm.ChannelTransferProveHashs[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// 签名
	seek, e = elm.SignCount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	scn := int(elm.ChannelCount)
	elm.Signs = make([]fields.Sign, scn)
	for i := 0; i < scn; i++ {
		seek, e = elm.Signs[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// 完成
	return seek, nil
}

// 检查签名
func (elm *OffChainFormPaymentChannelTransfer) CheckAddressAndSign(leftAddress, rightAddress fields.Address) error {
	// 检查公钥和地址是否匹配
	// 全部检查成功
	return nil
}

////////////////////////////////

// 通道转账，数据体
type ChannelTransferProveBodyInfo struct {
	ChannelId fields.Bytes16 // 通道id

	Mode fields.VarUint1 // 模式，普通模式、快速转账模式等等

	// Mode = 1 // 普通模式（实时对账）可以上链仲裁，但有负值是则为信用透支额度，债权方仅可仲裁无信用部分，即债务清除
	LeftAmount  fields.Amount // 左侧实时金额
	RightAmount fields.Amount // 右侧实时金额

	// Mode = 2 // 快速支付模式（延迟对账）无法上链仲裁
	Direction fields.VarUint1 // 资金流动方向： 1.左侧支付给右侧； 2.右侧付给左侧
	PayAmount fields.Amount   // 左侧支付给右侧的金额，如果为负值则表示右侧支付给左侧

	AutoNumber fields.VarUint8 // 通道流水序号
}
