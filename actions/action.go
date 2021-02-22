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
		return new(Action_1_SimpleTransfer), nil
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
	}
	////////////////////    END      ////////////////////
	return nil, fmt.Errorf("Cannot find Action kind of " + string(kind))
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
