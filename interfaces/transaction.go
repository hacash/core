package interfaces

import "github.com/hacash/core/fields"

type Transaction interface {

	// the transaction type number
	Type() uint8

	// assembling datas
	Serialize() ([]byte, error)
	Parse([]byte, uint32) (uint32, error)
	Size() uint32

	// hash data
	Hash() fields.Hash        // 无手续费的哈希
	HashWithFee() fields.Hash // inclued fee amount

	// Addresses that need to verify signatures
	RequestSignAddresses([]fields.Address) ([]fields.Address, error)

	// fill signatures
	FillNeedSigns(map[string][]byte, []fields.Address) error

	// verify signatures
	VerifyNeedSigns([]fields.Address) (bool, error)

	// change chain state
	WriteinChainState(ChainStateOperation) error
	RecoverChainState(ChainStateOperation) error

	// fee
	FeePurity() uint64 // fee purity

	// get something
	GetAddress() fields.Address
	GetFee() fields.Amount
	GetActions() []Action
	GetTimestamp() uint64
	GetMessage() fields.TrimString16
	SetMessage(fields.TrimString16)
}
