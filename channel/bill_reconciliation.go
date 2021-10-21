package channel

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
)

/**
 * 通道实时对账（链下签署）
 */
type OffChainFormPaymentChannelRealtimeReconciliation struct {
	// 签名哈希计算数据部分
	ChannelId fields.ChannelId // 通道id

	ReuseVersion   fields.VarUint4 // 通道重用序号
	BillAutoNumber fields.VarUint8 // 通道账单流水序号

	LeftBalance  fields.Amount // 左侧实时金额
	RightBalance fields.Amount // 右侧实时金额

	LeftSatoshi  fields.SatoshiVariation // 左侧比特币sat数量
	RightSatoshi fields.SatoshiVariation // 右侧比特币sat数量

	// 非签名哈希计算数据部分
	LeftAddress  fields.Address // 左侧地址
	RightAddress fields.Address // 右侧地址

	Timestamp fields.BlockTxTimestamp // 对账时间戳

	// 两侧签名
	LeftSign  fields.Sign // 左侧地址对账签名
	RightSign fields.Sign // 右侧地址对账签名
}

// interface
// 类型
func (e *OffChainFormPaymentChannelRealtimeReconciliation) TypeCode() uint8 {
	return BillTypeCodeReconciliation
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetChannelId() fields.ChannelId {
	return e.ChannelId
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetLeftBalance() fields.Amount {
	return e.LeftBalance
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetRightBalance() fields.Amount {
	return e.RightBalance
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetLeftSatoshi() fields.Satoshi {
	return e.LeftSatoshi.GetRealSatoshi()
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetRightSatoshi() fields.Satoshi {
	return e.RightSatoshi.GetRealSatoshi()
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetLeftAddress() fields.Address {
	return e.LeftAddress
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetRightAddress() fields.Address {
	return e.RightAddress
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetReuseVersion() uint32 {
	return uint32(e.ReuseVersion)
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetReuseVersionAndAutoNumber() (uint32, uint64) {
	return uint32(e.ReuseVersion), uint64(e.BillAutoNumber)
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetAutoNumber() uint64 {
	return uint64(e.BillAutoNumber)
}
func (e *OffChainFormPaymentChannelRealtimeReconciliation) GetTimestamp() uint64 {
	return uint64(e.Timestamp)
}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) Size() uint32 {
	return elm.ChannelId.Size() +
		elm.ReuseVersion.Size() +
		elm.BillAutoNumber.Size() +
		elm.LeftBalance.Size() +
		elm.RightBalance.Size() +
		elm.LeftSatoshi.Size() +
		elm.RightSatoshi.Size() +
		elm.LeftAddress.Size() +
		elm.RightAddress.Size() +
		elm.Timestamp.Size() +
		elm.LeftSign.Size() +
		elm.RightSign.Size()
}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) SerializeForSign() ([]byte, error) {
	var buffer bytes.Buffer
	var bt []byte
	bt, _ = elm.ChannelId.Serialize()
	buffer.Write(bt)
	bt, _ = elm.ReuseVersion.Serialize()
	buffer.Write(bt)
	bt, _ = elm.BillAutoNumber.Serialize()
	buffer.Write(bt)
	bt, _ = elm.LeftBalance.Serialize()
	buffer.Write(bt)
	bt, _ = elm.RightBalance.Serialize()
	buffer.Write(bt)
	bt, _ = elm.LeftSatoshi.Serialize()
	buffer.Write(bt)
	bt, _ = elm.RightSatoshi.Serialize()
	buffer.Write(bt)
	return buffer.Bytes(), nil

}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt []byte
	bt, _ = elm.SerializeForSign() // 签名部分数据体
	buffer.Write(bt)
	bt, _ = elm.LeftAddress.Serialize()
	buffer.Write(bt)
	bt, _ = elm.RightAddress.Serialize()
	buffer.Write(bt)
	bt, _ = elm.Timestamp.Serialize()
	buffer.Write(bt)
	bt, _ = elm.LeftSign.Serialize()
	buffer.Write(bt)
	bt, _ = elm.RightSign.Serialize()
	buffer.Write(bt)
	return buffer.Bytes(), nil
}

// 序列化
func (e *OffChainFormPaymentChannelRealtimeReconciliation) SerializeWithTypeCode() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{e.TypeCode()})
	b1, err := e.Serialize()
	if err != nil {
		return nil, err
	}
	buf.Write(b1)
	return buf.Bytes(), nil
}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) SignStuffHash() fields.Hash {
	var conbt, _ = elm.SerializeForSign() // 数据体
	return fields.CalculateHash(conbt)
}

func (elm *OffChainFormPaymentChannelRealtimeReconciliation) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ReuseVersion.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.BillAutoNumber.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftBalance.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightBalance.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftSatoshi.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightSatoshi.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAddress.Parse(buf, seek)
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
func (elm *OffChainFormPaymentChannelRealtimeReconciliation) CheckAddressAndSign() error {
	// 验证哈希
	var conhx = elm.SignStuffHash() // 数据体hx
	ok1, _ := account.CheckSignByHash32(conhx, elm.LeftSign.PublicKey, elm.LeftSign.Signature)
	if !ok1 {
		return fmt.Errorf("Left account %s verify signature fail.", elm.LeftAddress.ToReadable())
	}
	ok2, _ := account.CheckSignByHash32(conhx, elm.RightSign.PublicKey, elm.RightSign.Signature)
	if !ok2 {
		return fmt.Errorf("Right account %s verify signature fail.", elm.RightAddress.ToReadable())
	}
	// 全部检查成功
	return nil
}

// 检查数据可用性
func (elm *OffChainFormPaymentChannelRealtimeReconciliation) CheckValidity() error {
	return nil
}

// 验证对票据的签名
func (elm *OffChainFormPaymentChannelRealtimeReconciliation) VerifySignature() error {
	return elm.CheckAddressAndSign()
}

// 填充一方签名
func (elm *OffChainFormPaymentChannelRealtimeReconciliation) FillTargetSignature(acc *account.Account) (*fields.Sign, bool, error) {
	hx := elm.SignStuffHash()
	addrIsLeft := elm.LeftAddress.Equal(acc.Address)
	// 计算签名
	signdata, e := acc.Private.Sign(hx)
	if e != nil {
		return nil, addrIsLeft, e // 签名错误
	}
	signobj := fields.Sign{
		PublicKey: acc.PublicKey,
		Signature: signdata.Serialize64(),
	}
	if addrIsLeft {
		elm.LeftSign = signobj
	} else {
		elm.RightSign = signobj
	}
	return &signobj, addrIsLeft, nil
}

/********************************************************/

/**
 * 链上仲裁需要的对账单（链上仲裁）
 */

type OnChainArbitrationBasisReconciliation struct {
	// 签名哈希计算数据部分
	ChannelId fields.ChannelId // 通道id

	ReuseVersion   fields.VarUint4 // 通道重用序号
	BillAutoNumber fields.VarUint8 // 通道账单流水序号

	LeftBalance  fields.Amount // 左侧实时金额
	RightBalance fields.Amount // 右侧实时金额

	LeftSatoshi  fields.SatoshiVariation // 左侧比特币sat数量
	RightSatoshi fields.SatoshiVariation // 右侧比特币sat数量

	// 两侧签名
	LeftSign  fields.Sign // 左侧地址对账签名
	RightSign fields.Sign // 右侧地址对账签名
}

func (e *OnChainArbitrationBasisReconciliation) GetChannelId() fields.ChannelId {
	return e.ChannelId
}
func (e *OnChainArbitrationBasisReconciliation) GetLeftBalance() fields.Amount {
	return e.LeftBalance
}
func (e *OnChainArbitrationBasisReconciliation) GetRightBalance() fields.Amount {
	return e.RightBalance
}
func (e *OnChainArbitrationBasisReconciliation) GetLeftSatoshi() fields.Satoshi {
	return e.LeftSatoshi.GetRealSatoshi()
}
func (e *OnChainArbitrationBasisReconciliation) GetRightSatoshi() fields.Satoshi {
	return e.RightSatoshi.GetRealSatoshi()
}
func (e *OnChainArbitrationBasisReconciliation) GetReuseVersion() uint32 {
	return uint32(e.ReuseVersion)
}
func (e *OnChainArbitrationBasisReconciliation) GetAutoNumber() uint64 {
	return uint64(e.BillAutoNumber)
}
func (elm *OnChainArbitrationBasisReconciliation) Size() uint32 {
	return elm.ChannelId.Size() +
		elm.ReuseVersion.Size() +
		elm.BillAutoNumber.Size() +
		elm.LeftBalance.Size() +
		elm.RightBalance.Size() +
		elm.LeftSatoshi.Size() +
		elm.RightSatoshi.Size() +
		elm.LeftSign.Size() +
		elm.RightSign.Size()
}

func (elm *OnChainArbitrationBasisReconciliation) SerializeForSign() ([]byte, error) {
	var buffer bytes.Buffer
	var bt []byte
	bt, _ = elm.ChannelId.Serialize()
	buffer.Write(bt)
	bt, _ = elm.ReuseVersion.Serialize()
	buffer.Write(bt)
	bt, _ = elm.BillAutoNumber.Serialize()
	buffer.Write(bt)
	bt, _ = elm.LeftBalance.Serialize()
	buffer.Write(bt)
	bt, _ = elm.RightBalance.Serialize()
	buffer.Write(bt)
	bt, _ = elm.LeftSatoshi.Serialize()
	buffer.Write(bt)
	bt, _ = elm.RightSatoshi.Serialize()
	buffer.Write(bt)
	return buffer.Bytes(), nil

}

func (elm *OnChainArbitrationBasisReconciliation) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt []byte
	bt, _ = elm.SerializeForSign() // 签名部分数据体
	buffer.Write(bt)
	bt, _ = elm.LeftSign.Serialize()
	buffer.Write(bt)
	bt, _ = elm.RightSign.Serialize()
	buffer.Write(bt)
	return buffer.Bytes(), nil
}

func (elm *OnChainArbitrationBasisReconciliation) SignStuffHash() fields.Hash {
	var conbt, _ = elm.SerializeForSign() // 数据体
	return fields.CalculateHash(conbt)
}

func (elm *OnChainArbitrationBasisReconciliation) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ReuseVersion.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.BillAutoNumber.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftBalance.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightBalance.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftSatoshi.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightSatoshi.Parse(buf, seek)
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

// 填充签名
func (elm *OnChainArbitrationBasisReconciliation) FillSigns(lacc, racc *account.Account) error {

	txhx := elm.SignStuffHash()

	s1, e := lacc.Private.Sign(txhx)
	if e != nil {
		return e
	}
	s2, e := racc.Private.Sign(txhx)
	if e != nil {
		return e
	}
	elm.LeftSign = fields.Sign{
		PublicKey: lacc.PublicKey,
		Signature: s1.Serialize64(),
	}
	elm.RightSign = fields.Sign{
		PublicKey: racc.PublicKey,
		Signature: s2.Serialize64(),
	}

	return nil
}

// 检查签名
func (elm *OnChainArbitrationBasisReconciliation) CheckAddressAndSign(laddr, raddr fields.Address) error {
	// 验证哈希
	var conhx = elm.SignStuffHash() // 数据体hx
	ok1, _ := account.CheckSignByHash32(conhx, elm.LeftSign.PublicKey, elm.LeftSign.Signature)
	if !ok1 {
		return fmt.Errorf("Left account %s verify signature fail.", laddr.ToReadable())
	}
	ok2, _ := account.CheckSignByHash32(conhx, elm.RightSign.PublicKey, elm.RightSign.Signature)
	if !ok2 {
		return fmt.Errorf("Right account %s verify signature fail.", raddr.ToReadable())
	}
	// 全部检查成功
	return nil
}
