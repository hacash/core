package channel

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/fields"
)

/**

跨节点支付对账票据

*/

type OffChainCrossNodeSimplePaymentReconciliationBill struct {

	// 本通道对账单
	ChannelChainTransferTargetProveBody ChannelChainTransferProveBodyInfo

	// 通道链支付数据
	ChannelChainTransferData OffChainFormPaymentChannelTransfer
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) TypeCode() uint8 {
	return BillTypeCodeSimplePay // 类型
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) Size() uint32 {
	return c.ChannelChainTransferTargetProveBody.Size() +
		c.ChannelChainTransferData.Size()
}

func (c *OffChainCrossNodeSimplePaymentReconciliationBill) Parse(buf []byte, seek uint32) (uint32, error) {
	var err error
	seek, err = c.ChannelChainTransferTargetProveBody.Parse(buf, seek)
	if err != nil {
		return 0, err
	}
	return c.ChannelChainTransferData.Parse(buf, seek)
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) Serialize() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	b3, err := c.ChannelChainTransferTargetProveBody.Serialize()
	if err != nil {
		return nil, err
	}
	buf.Write(b3)
	b4, err := c.ChannelChainTransferData.Serialize()
	if err != nil {
		return nil, err
	}
	buf.Write(b4)
	return buf.Bytes(), nil

}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) SerializeWithTypeCode() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{c.TypeCode()})
	b1, err := c.Serialize()
	if err != nil {
		return nil, err
	}
	buf.Write(b1)
	return buf.Bytes(), nil
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) GetChannelId() fields.ChannelId {
	return c.ChannelChainTransferTargetProveBody.ChannelId
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) GetLeftBalance() fields.Amount {
	return c.ChannelChainTransferTargetProveBody.LeftBalance
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) GetRightBalance() fields.Amount {
	return c.ChannelChainTransferTargetProveBody.RightBalance
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) GetLeftSatoshi() fields.Satoshi {
	return c.ChannelChainTransferTargetProveBody.LeftSatoshi.GetRealSatoshi()
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) GetRightSatoshi() fields.Satoshi {
	return c.ChannelChainTransferTargetProveBody.RightSatoshi.GetRealSatoshi()
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) GetLeftAddress() fields.Address {
	return c.ChannelChainTransferTargetProveBody.LeftAddress
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) GetRightAddress() fields.Address {
	return c.ChannelChainTransferTargetProveBody.RightAddress
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) GetTimestamp() uint64 {
	return uint64(c.ChannelChainTransferData.Timestamp)
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) GetReuseVersionAndAutoNumber() (uint32, uint64) {
	return uint32(c.ChannelChainTransferTargetProveBody.ReuseVersion),
		uint64(c.ChannelChainTransferTargetProveBody.BillAutoNumber)
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) GetReuseVersion() uint32 {
	return uint32(c.ChannelChainTransferTargetProveBody.ReuseVersion)
}

func (c OffChainCrossNodeSimplePaymentReconciliationBill) GetAutoNumber() uint64 {
	return uint64(c.ChannelChainTransferTargetProveBody.BillAutoNumber)
}

// 检查数据可用性
func (c OffChainCrossNodeSimplePaymentReconciliationBill) CheckValidity() error {
	var checkIsOk bool = false
	hxchecker := c.ChannelChainTransferTargetProveBody.GetSignStuffHashHalfChecker()
	cks := c.ChannelChainTransferData.ChannelTransferProveHashHalfCheckers
	for _, v := range cks {
		if v.Equal(hxchecker) {
			checkIsOk = true
			break
		}
	}
	// 是否包含
	if checkIsOk == false {
		return fmt.Errorf("ProveBody's HashHalfChecker <%s> not included in ChannelTransferProveHashHalfCheckers.",
			hxchecker.ToHex())
	}
	// 检查成功
	return nil
}

// 验证对票据的签名
func (c OffChainCrossNodeSimplePaymentReconciliationBill) VerifySignature() error {

	return c.ChannelChainTransferData.CheckMustAddressAndSigns()
}
