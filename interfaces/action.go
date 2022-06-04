package interfaces

import (
	"github.com/hacash/core/fields"
)

type Action interface {

	// base super
	Field

	// the action type number
	Kind() uint16

	// Addresses that need to verify signatures
	RequestSignAddresses() []fields.Address

	// change chain state
	WriteInChainState(ChainStateOperation) error

	// help func
	SetBelongTrs(Transaction)
	Describe() map[string]interface{} // json api

	// burning fees
	IsBurning90PersentTxFees() bool // Whether to destroy 90% of the transaction cost of this transaction
}
