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
	// Signature hash calculation data part
	ChannelId fields.ChannelId // Channel ID

	ReuseVersion   fields.VarUint4 // Channel reuse sequence number
	BillAutoNumber fields.VarUint8 // Serial number of channel bill

	LeftBalance  fields.Amount // Real time amount on the left
	RightBalance fields.Amount // Right real time amount

	LeftSatoshi  fields.SatoshiVariation // Number of bitcoin sat on the left
	RightSatoshi fields.SatoshiVariation // Number of bitcoin sat on the right

	// Unsigned hash calculation data part
	LeftAddress  fields.Address // Left address
	RightAddress fields.Address // Right address

	Timestamp fields.BlockTxTimestamp // Reconciliation timestamp

	// Signature on both sides
	LeftSign  fields.Sign // Left address reconciliation signature
	RightSign fields.Sign // Right address reconciliation signature
}

// interface
// type
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
	bt, _ = elm.SerializeForSign() // Signature part data body
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

// serialize
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
	var conbt, _ = elm.SerializeForSign() // Data body
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

// Check signature
func (elm *OffChainFormPaymentChannelRealtimeReconciliation) CheckAddressAndSign() error {
	// Verify hash
	var conhx = elm.SignStuffHash() // Data body HX
	ok1, _ := account.CheckSignByHash32(conhx, elm.LeftSign.PublicKey, elm.LeftSign.Signature)
	if !ok1 {
		return fmt.Errorf("Left account %s verify signature fail.", elm.LeftAddress.ToReadable())
	}
	ok2, _ := account.CheckSignByHash32(conhx, elm.RightSign.PublicKey, elm.RightSign.Signature)
	if !ok2 {
		return fmt.Errorf("Right account %s verify signature fail.", elm.RightAddress.ToReadable())
	}
	// All checked successfully
	return nil
}

// Check data availability
func (elm *OffChainFormPaymentChannelRealtimeReconciliation) CheckValidity() error {
	return nil
}

// Verify signature on ticket
func (elm *OffChainFormPaymentChannelRealtimeReconciliation) VerifySignature() error {
	return elm.CheckAddressAndSign()
}

// 填充一方签名
func (elm *OffChainFormPaymentChannelRealtimeReconciliation) FillTargetSignature(acc *account.Account) (*fields.Sign, bool, error) {
	hx := elm.SignStuffHash()
	addrIsLeft := elm.LeftAddress.Equal(acc.Address)
	// Calculate signature
	signdata, e := acc.Private.Sign(hx)
	if e != nil {
		return nil, addrIsLeft, e // Signature error
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
	// Signature hash calculation data part
	ChannelId fields.ChannelId // Channel ID

	ReuseVersion   fields.VarUint4 // Channel reuse sequence number
	BillAutoNumber fields.VarUint8 // Serial number of channel bill

	LeftBalance  fields.Amount // Real time amount on the left
	RightBalance fields.Amount // Right real time amount

	LeftSatoshi  fields.SatoshiVariation // Number of bitcoin sat on the left
	RightSatoshi fields.SatoshiVariation // Number of bitcoin sat on the right

	// Signature on both sides
	LeftSign  fields.Sign // Left address reconciliation signature
	RightSign fields.Sign // Right address reconciliation signature
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
	bt, _ = elm.SerializeForSign() // Signature part data body
	buffer.Write(bt)
	bt, _ = elm.LeftSign.Serialize()
	buffer.Write(bt)
	bt, _ = elm.RightSign.Serialize()
	buffer.Write(bt)
	return buffer.Bytes(), nil
}

func (elm *OnChainArbitrationBasisReconciliation) SignStuffHash() fields.Hash {
	var conbt, _ = elm.SerializeForSign() // Data body
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

// Fill in signature
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

// Check signature
func (elm *OnChainArbitrationBasisReconciliation) CheckAddressAndSign(laddr, raddr fields.Address) error {
	// Verify hash
	var conhx = elm.SignStuffHash() // Data body HX
	ok1, _ := account.CheckSignByHash32(conhx, elm.LeftSign.PublicKey, elm.LeftSign.Signature)
	if !ok1 {
		return fmt.Errorf("Left account %s verify signature fail.", laddr.ToReadable())
	}
	ok2, _ := account.CheckSignByHash32(conhx, elm.RightSign.PublicKey, elm.RightSign.Signature)
	if !ok2 {
		return fmt.Errorf("Right account %s verify signature fail.", raddr.ToReadable())
	}
	// All checked successfully
	return nil
}
