package channel

import (
	"bytes"
)

/**
 * 通道支付完整单据
 */

// 通道链支付票据
type ChannelPayCompleteDocuments struct {
	// 对账票据表
	ProveBodys *ChannelPayProveBodyList
	// 支付签名票据
	ChainPayment *OffChainFormPaymentChannelTransfer
}

func (c ChannelPayCompleteDocuments) Size() uint32 {
	return c.ProveBodys.Size() + c.ChainPayment.Size()
}

func (c ChannelPayCompleteDocuments) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = c.ProveBodys.Serialize()
	buffer.Write(bt1)
	var bt2, _ = c.ChainPayment.Serialize()
	buffer.Write(bt2)
	return buffer.Bytes(), nil
}

func (c *ChannelPayCompleteDocuments) Parse(buf []byte, seek uint32) (uint32, error) {
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
