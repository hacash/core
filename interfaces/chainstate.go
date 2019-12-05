package interfaces

import (
	"github.com/hacash/core/fields"
)

type ChainState interface {

	// save
	SaveBlockUniteTransactions(Block) error

	// block
	ReadBlockHeadBytesByHash(fields.Hash) ([]byte, error)
	ReadBlockBytesByHash(fields.Hash) ([]byte, error)
	ReadBlockHeadBytesByHeight(fields.VarInt5) ([]byte, error)
	ReadBlockBytesByHeight(fields.VarInt5) ([]byte, error)

	// tx
	ReadTransactionBytes(fields.Hash) ([]byte, error)
	TransactionIsExist(fields.Hash) (bool, error)

	// status
	SetLastestBlockHead(Block) error
	ReadLastestBlockHead() ([]byte, error)
	SetLastestDiamondStatus(uint32, fields.Hash) error
	ReadLastestDiamondStatus() (uint32, fields.Hash, error)
}
