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
	Hash() fields.Hash
	HashNoFee() fields.Hash // 无手续费的哈希

	// Addresses that need to verify signatures
	RequestSignAddresses([][]byte) ([][]byte, error)

	// fill signatures
	FillNeedSigns(map[string][]byte, [][]byte) error

	// verify signatures
	VerifyNeedSigns([][]byte) (bool, error)

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
}
