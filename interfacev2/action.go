package interfacev2

import (
	"github.com/hacash/core/fields"
)

type Action interface {

	// the action type number
	Kind() uint16

	// assembling datas
	Serialize() ([]byte, error)
	Parse([]byte, uint32) (uint32, error)
	Size() uint32

	// Addresses that need to verify signatures
	RequestSignAddresses() []fields.Address

	// change chain state
	WriteinChainState(ChainStateOperation) error
	RecoverChainState(ChainStateOperation) error

	// help func
	SetBelongTransaction(Transaction)
	Describe() map[string]interface{} // json api

	// burning fees
	IsBurning90PersentTxFees() bool // Whether to destroy 90% of the transaction cost of this transaction
}
