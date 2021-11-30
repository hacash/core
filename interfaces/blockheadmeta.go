package interfaces

import "github.com/hacash/core/fields"

type BlockHeadMetaRead interface {
	Hash() fields.Hash
	Version() uint8
	GetHeight() uint64
	GetDifficulty() uint32
	GetWitnessStage() uint16
	GetNonce() uint32
	GetNonceByte() []byte
	GetPrevHash() fields.Hash
	GetTimestamp() uint64
	GetTransactionCount() uint32
	GetMrklRoot() fields.Hash

	// customer trs count
	GetCustomerTransactionCount() uint32

	SerializeExcludeTransactions() ([]byte, error)
}
