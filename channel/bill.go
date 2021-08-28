package channel

import (
	"fmt"
	"github.com/hacash/core/fields"
)

/**

支付票据类型

*/
const (
	BillTypeCodeSimplePay uint8 = 1 // 普通支付
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

	ChannelId() fields.Bytes16

	LeftAddress() fields.Address
	RightAddress() fields.Address

	LeftAmount() fields.Amount
	RightAmount() fields.Amount

	Timestamp() uint64 // 对账时间戳，秒

	// 通道重用序号 & 通道账单流水序号
	ChannelReuseVersionAndAutoNumber() (uint32, uint64)

	CheckValidity() error   // 检查数据可用性
	VerifySignature() error // 验证对票据的签名

}

/**

对账票据解析

*/

// 反序列化
func ParseReconciliationBalanceBillByPrefixTypeCode(buf []byte, seek uint32) (ReconciliationBalanceBill, error) {
	ty := buf[seek]
	var bill ReconciliationBalanceBill = nil

	// 类型
	switch ty {
	case BillTypeCodeSimplePay:
		bill = &CrossNodeSimplePaymentReconciliationBill{}
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
