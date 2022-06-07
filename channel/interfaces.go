package channel

import (
	"fmt"
	"github.com/hacash/core/fields"
)

/**

链上仲裁接口

*/

// Online arbitration reconciliation basis
type OnChainChannelPaymentArbitrationReconciliationBasis interface {
	GetChannelId() fields.ChannelId
	GetLeftBalance() fields.Amount   // Left HAC amount
	GetRightBalance() fields.Amount  // Right HAC amount
	GetLeftSatoshi() fields.Satoshi  // BTC sat amount allocated on the left
	GetRightSatoshi() fields.Satoshi // BTC sat amount allocated on the right
	GetReuseVersion() uint32         // Channel reuse sequence number
	GetAutoNumber() uint64
	// Check address and signature
	CheckAddressAndSign(laddr, raddr fields.Address) error
}

/*********************************************************/

/**

支付票据类型

*/
const (
	BillTypeCodeSimplePay      uint8 = 1 // Ordinary payment
	BillTypeCodeReconciliation uint8 = 2 // Reconciliation
)

/**

通道链票据接口

*/

// General reconciliation bill interface
type ReconciliationBalanceBill interface {
	Size() uint32
	Parse(buf []byte, seek uint32) (uint32, error) // Deserialization
	Serialize() ([]byte, error)                    // serialize
	SerializeWithTypeCode() ([]byte, error)        // serialize
	TypeCode() uint8                               // type

	GetChannelId() fields.ChannelId

	GetLeftSatoshi() fields.Satoshi
	GetRightSatoshi() fields.Satoshi

	GetLeftBalance() fields.Amount
	GetRightBalance() fields.Amount

	GetLeftAddress() fields.Address
	GetRightAddress() fields.Address

	GetTimestamp() uint64 // Reconciliation timestamp, seconds

	// Channel reuse serial number channel bill serial number
	GetReuseVersionAndAutoNumber() (uint32, uint64)
	GetReuseVersion() uint32
	GetAutoNumber() uint64

	CheckValidity() error   // Check data availability
	VerifySignature() error // Verify signature on ticket

}

/**

对账票据解析

*/

// serialize
func SerializeReconciliationBalanceBillWithPrefixTypeCode(bill ReconciliationBalanceBill) ([]byte, error) {
	// type
	// analysis
	bts, e := bill.SerializeWithTypeCode()
	if e != nil {
		return nil, e
	}
	return bts, nil

}

// Deserialization
func ParseReconciliationBalanceBillByPrefixTypeCode(buf []byte, seek uint32) (ReconciliationBalanceBill, uint32, error) {
	ty := buf[seek]
	var bill ReconciliationBalanceBill = nil

	// type
	switch ty {
	case BillTypeCodeSimplePay: // 普通通道链支付
		bill = &OffChainCrossNodeSimplePaymentReconciliationBill{}
	case BillTypeCodeReconciliation: // 通道链对账
		bill = &OffChainFormPaymentChannelRealtimeReconciliation{}
	default:
		return nil, 0, fmt.Errorf("Unsupported bill type <%d>", ty)
	}

	// analysis
	var e error
	seek, e = bill.Parse(buf, seek+1)
	if e != nil {
		return nil, 0, e
	}
	return bill, seek, nil

}
