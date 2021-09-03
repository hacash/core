package channel

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/fields"
)

/**

跨节点支付对账票据

*/

type CrossNodeSimplePaymentReconciliationBill struct {

	// 地址
	LeftAddr  fields.Address
	RightAddr fields.Address

	// 本通道对账单
	ChannelChainTransferTargetProveBody ChannelChainTransferProveBodyInfo

	// 通道链支付数据
	ChannelChainTransferData OffChainFormPaymentChannelTransfer
}

func (c CrossNodeSimplePaymentReconciliationBill) TypeCode() uint8 {
	return BillTypeCodeSimplePay // 类型
}

func (c CrossNodeSimplePaymentReconciliationBill) Size() uint32 {
	return c.LeftAddr.Size() +
		c.RightAddr.Size() +
		c.ChannelChainTransferTargetProveBody.Size() +
		c.ChannelChainTransferData.Size()
}

func (c *CrossNodeSimplePaymentReconciliationBill) Parse(buf []byte, seek uint32) (uint32, error) {
	var err error
	seek, err = c.LeftAddr.Parse(buf, seek)
	if err != nil {
		return 0, err
	}
	seek, err = c.RightAddr.Parse(buf, seek)
	if err != nil {
		return 0, err
	}
	seek, err = c.ChannelChainTransferTargetProveBody.Parse(buf, seek)
	if err != nil {
		return 0, err
	}
	return c.ChannelChainTransferData.Parse(buf, seek)
}

func (c CrossNodeSimplePaymentReconciliationBill) Serialize() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	b1, err := c.LeftAddr.Serialize()
	if err != nil {
		return nil, err
	}
	buf.Write(b1)
	b2, err := c.RightAddr.Serialize()
	if err != nil {
		return nil, err
	}
	buf.Write(b2)
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

func (c CrossNodeSimplePaymentReconciliationBill) SerializeWithTypeCode() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{c.TypeCode()})
	b1, err := c.Serialize()
	if err != nil {
		return nil, err
	}
	buf.Write(b1)
	return buf.Bytes(), nil
}

func (c CrossNodeSimplePaymentReconciliationBill) ChannelId() fields.Bytes16 {
	return c.ChannelChainTransferTargetProveBody.ChannelId
}

func (c CrossNodeSimplePaymentReconciliationBill) LeftAmount() fields.Amount {
	return c.ChannelChainTransferTargetProveBody.LeftAmount
}

func (c CrossNodeSimplePaymentReconciliationBill) RightAmount() fields.Amount {
	return c.ChannelChainTransferTargetProveBody.RightAmount
}

func (c CrossNodeSimplePaymentReconciliationBill) LeftAddress() fields.Address {
	return c.LeftAddr
}

func (c CrossNodeSimplePaymentReconciliationBill) RightAddress() fields.Address {
	return c.RightAddr
}

func (c CrossNodeSimplePaymentReconciliationBill) Timestamp() uint64 {
	return uint64(c.ChannelChainTransferData.Timestamp)
}

func (c CrossNodeSimplePaymentReconciliationBill) ChannelReuseVersionAndAutoNumber() (uint32, uint64) {
	return uint32(c.ChannelChainTransferTargetProveBody.ChannelReuseVersion),
		uint64(c.ChannelChainTransferTargetProveBody.BillAutoNumber)
}

func (c CrossNodeSimplePaymentReconciliationBill) ChannelAutoNumber() uint64 {
	return uint64(c.ChannelChainTransferTargetProveBody.BillAutoNumber)
}

// 检查数据可用性
func (c CrossNodeSimplePaymentReconciliationBill) CheckValidity() error {
	var checkIsOk bool = false
	hxchecker := c.ChannelChainTransferTargetProveBody.SignStuffHashHalfChecker()
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
	// 对账模式
	if c.ChannelChainTransferTargetProveBody.Mode != ChannelTransferProveBodyPayModeNormal {
		return fmt.Errorf("ProveBody's pay Mode is not Normal.")
	}
	// 检查成功
	return nil
}

// 验证对票据的签名
func (c CrossNodeSimplePaymentReconciliationBill) VerifySignature() error {

	return c.ChannelChainTransferData.CheckMustAddressAndSigns()
}
