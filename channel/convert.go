package channel

import "github.com/hacash/core/fields"

/**
 * 转换支付票据为实时对账单
 */
func (bill OffChainCrossNodeSimplePaymentReconciliationBill) ConvertToRealtimeReconciliation() *OffChainFormPaymentChannelRealtimeReconciliation {

	// Create reconciliation bill
	recbill := &OffChainFormPaymentChannelRealtimeReconciliation{
		ChannelId:      bill.GetChannelId(),
		ReuseVersion:   fields.VarUint4(bill.GetReuseVersion()),
		BillAutoNumber: fields.VarUint8(bill.GetAutoNumber()),
		LeftBalance:    bill.GetLeftBalance(),
		RightBalance:   bill.GetRightBalance(),
		LeftSatoshi:    bill.GetLeftSatoshi().GetSatoshiVariation(),
		RightSatoshi:   bill.GetRightSatoshi().GetSatoshiVariation(),
		LeftAddress:    bill.GetLeftAddress(),
		RightAddress:   bill.GetRightAddress(),
		Timestamp:      fields.BlockTxTimestamp(bill.GetTimestamp()),
		LeftSign:       fields.Sign{},
		RightSign:      fields.Sign{},
	}

	return recbill

}

/**
 * 转换链下实时对账单为链上仲裁单据
 */
func (bill OffChainFormPaymentChannelRealtimeReconciliation) ConvertToOnChain() *OnChainArbitrationBasisReconciliation {
	return &OnChainArbitrationBasisReconciliation{
		ChannelId:      bill.ChannelId,
		ReuseVersion:   bill.ReuseVersion,
		BillAutoNumber: bill.BillAutoNumber,
		LeftBalance:    bill.LeftBalance,
		RightBalance:   bill.RightBalance,
		LeftSatoshi:    bill.LeftSatoshi,
		RightSatoshi:   bill.RightSatoshi,
		LeftSign:       bill.LeftSign,
		RightSign:      bill.RightSign,
	}
}
