package channel

import (
	"bytes"
)

/**
 * 通道支付完整单据
 */

type ChannelPayCompleteDocuments struct {

	// 支付体
	Form *OffChainFormPaymentChannelTransfer

	// 通道列表对账数据
	Channels []*ChannelChainTransferProveBodyInfo
}

func (elm *ChannelPayCompleteDocuments) Size() uint32 {
	size := elm.Form.Size()
	for i := 0; i < int(elm.Form.ChannelCount); i++ {
		size += elm.Channels[i].Size()
	}
	// ok
	return size
}

func (elm *ChannelPayCompleteDocuments) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, e = elm.Form.Serialize()
	if e != nil {
		return nil, e
	}
	buffer.Write(bt1)
	for i := 0; i < int(elm.Form.ChannelCount); i++ {
		bt1, e := elm.Channels[i].Serialize()
		if e != nil {
			return nil, e
		}
		buffer.Write(bt1)
	}
	return buffer.Bytes(), nil
}

func (elm *ChannelPayCompleteDocuments) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.Form.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	// 地址
	scn := int(elm.Form.ChannelCount)
	elm.Channels = make([]*ChannelChainTransferProveBodyInfo, scn)
	for i := 0; i < scn; i++ {
		elm.Channels[i] = &ChannelChainTransferProveBodyInfo{}
		seek, e = elm.Channels[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}
