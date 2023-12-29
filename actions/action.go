package actions

import (
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/interfaces"
)

/* *********************************************************** */

func NewActionByKind(kind uint16) (interfaces.Action, error) {
	////////////////////   ACTIONS   ////////////////////
	switch kind {
	case 1:
		return new(Action_1_SimpleToTransfer), nil
	case 2:
		return new(Action_2_OpenPaymentChannel), nil
	case 3:
		return new(Action_3_ClosePaymentChannel), nil
	case 4:
		return new(Action_4_DiamondCreate), nil
	case 5:
		return new(Action_5_DiamondTransfer), nil
	case 6:
		return new(Action_6_OutfeeQuantityDiamondTransfer), nil
	case 7:
		return new(Action_7_SatoshiGenesis), nil
	case 8:
		return new(Action_8_SimpleSatoshiTransfer), nil
	case 9:
		return new(Action_9_LockblsCreate), nil
	case 10:
		return new(Action_10_LockblsRelease), nil
	case 11:
		return new(Action_11_FromToSatoshiTransfer), nil
	case 12:
		return new(Action_12_ClosePaymentChannelBySetupAmount), nil
	case 13:
		return new(Action_13_FromTransfer), nil
	case 14:
		return new(Action_14_FromToTransfer), nil
	case 15:
		return new(Action_15_DiamondsSystemLendingCreate), nil
	case 16:
		return new(Action_16_DiamondsSystemLendingRansom), nil
	case 17:
		return new(Action_17_BitcoinsSystemLendingCreate), nil
	case 18:
		return new(Action_18_BitcoinsSystemLendingRansom), nil
	case 19:
		return new(Action_19_UsersLendingCreate), nil
	case 20:
		return new(Action_20_UsersLendingRansom), nil
	case 21:
		return new(Action_21_ClosePaymentChannelBySetupOnlyLeftAmount), nil
	case 22:
		return new(Action_22_UnilateralClosePaymentChannelByNothing), nil
	case 23:
		return new(Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation), nil
	case 24:
		return new(Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody), nil
	case 25:
		return new(Action_25_PaymantChannelAndOnchainAtomicExchange), nil
	case 26:
		return new(Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange), nil
	case 27:
		return new(Action_27_ClosePaymentChannelByClaimDistribution), nil
	case 28:
		return new(Action_28_FromSatoshiTransfer), nil
	case 29:
		return new(Action_29_SubmitTimeLimit), nil
	case 30:
		return new(Action_30_SupportDistinguishForkChainID), nil
	case 31:
		return new(Action_31_OpenPaymentChannelWithSatoshi), nil
	case 32:
		return new(Action_32_DiamondsEngraved), nil
	case 33:
		return new(Action_33_DiamondsEngravedRecovery), nil

	}
	////////////////////    END      ////////////////////
	return nil, fmt.Errorf("Cannot find Action kind of %d.", +kind)
}

func ParseAction(buf []byte, seek uint32) (interfaces.Action, uint32, error) {
	if seek+2 >= uint32(len(buf)) {
		return nil, 0, fmt.Errorf("[ParseAction] seek out of buf len.")
	}
	var kind = binary.BigEndian.Uint16(buf[seek : seek+2])
	var act, e1 = NewActionByKind(kind)
	if e1 != nil {
		return nil, 0, e1
	}
	var mv, err = act.Parse(buf, seek+2)
	return act, mv, err
}
