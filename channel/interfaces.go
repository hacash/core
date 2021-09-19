package channel

import (
	"fmt"
	"github.com/hacash/core/fields"
)

/**

链上仲裁接口

*/

// 链上仲裁对账依据
type OnChainChannelPaymentArbitrationReconciliationBasis interface {
	GetChannelId() fields.ChannelId
	GetLeftBalance() fields.Amount  // 左侧实时金额
	GetRightBalance() fields.Amount // 右侧实时金额
	GetReuseVersion() uint32        // 通道重用序号
	GetAutoNumber() uint64
	// 检查地址和签名
	CheckAddressAndSign(laddr, raddr fields.Address) error
}

/*********************************************************/

/**

支付票据类型

*/
const (
	BillTypeCodeSimplePay      uint8 = 1 // 普通支付
	BillTypeCodeReconciliation uint8 = 2 // 对账
)

/**

通道链票据接口

*/

// 通用对账票据接口
type ReconciliationBalanceBill interface {
	Size() uint32
	Parse(buf []byte, seek uint32) (uint32, error) // 反序列化
	Serialize() ([]byte, error)                    // 序列化
	SerializeWithTypeCode() ([]byte, error)        // 序列化
	TypeCode() uint8                               // 类型

	GetChannelId() fields.ChannelId

	GetLeftAddress() fields.Address
	GetRightAddress() fields.Address

	GetLeftBalance() fields.Amount
	GetRightBalance() fields.Amount

	GetTimestamp() uint64 // 对账时间戳，秒

	// 通道重用序号 & 通道账单流水序号
	GetReuseVersionAndAutoNumber() (uint32, uint64)
	GetReuseVersion() uint32
	GetAutoNumber() uint64

	CheckValidity() error   // 检查数据可用性
	VerifySignature() error // 验证对票据的签名

}

/**

对账票据解析

*/

// 序列化
func SerializeReconciliationBalanceBillWithPrefixTypeCode(bill ReconciliationBalanceBill) ([]byte, error) {
	// 类型
	// 解析
	bts, e := bill.SerializeWithTypeCode()
	if e != nil {
		return nil, e
	}
	return bts, nil

}

// 反序列化
func ParseReconciliationBalanceBillByPrefixTypeCode(buf []byte, seek uint32) (ReconciliationBalanceBill, error) {
	ty := buf[seek]
	var bill ReconciliationBalanceBill = nil

	// 类型
	switch ty {
	case BillTypeCodeSimplePay: // 普通通道链支付
		bill = &OffChainCrossNodeSimplePaymentReconciliationBill{}
	case BillTypeCodeReconciliation: // 通道链对账
		bill = &OffChainFormPaymentChannelRealtimeReconciliation{}
	default:
		return nil, fmt.Errorf("Unsupported bill type <%d>", ty)
	}

	// 解析
	_, e := bill.Parse(buf, seek+1)
	if e != nil {
		return nil, e
	}
	return bill, nil

}
