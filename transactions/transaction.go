package transactions

import (
	"fmt"
	"github.com/hacash/core/interfaces"
)

////////////////////////////////////////////////////////////////////////

func NewTransactionByType(ty uint8) (interfaces.Transaction, error) {
	switch ty {
	////////////////////  TRANSATION  ////////////////////
	case 0:
		return new(Transaction_0_Coinbase), nil
	case 1:
		return new(Transaction_1_DO_NOT_USE_WITH_BUG), nil // [signed bug, obsolete!!!]
	case 2:
		return new(Transaction_2_Simple), nil
		////////////////////     END      ////////////////////
	}
	return nil, fmt.Errorf("Cannot find Transaction type of " + string(ty))
}

func ParseTransaction(buf []byte, seek uint32) (interfaces.Transaction, uint32, error) {
	if seek >= uint32(len(buf)) {
		return nil, 0, fmt.Errorf("buf length over range")
	}
	ty := uint8(buf[seek])
	var trx, e1 = NewTransactionByType(ty)
	if e1 != nil {
		return nil, 0, e1
	}
	var mv, err = trx.Parse(buf, seek+1)
	return trx, mv, err
}
