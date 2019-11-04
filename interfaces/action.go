package interfaces

import "github.com/hacash/core/fields"

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
}
