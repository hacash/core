package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

type BlockStore interface {

	// save
	SaveBlockUniteTransactions(Block) error

	// block
	ReadBlockHeadBytesByHeight(uint64) ([]byte, error)
	ReadBlockHeadBytesByHash(fields.Hash) ([]byte, error)
	ReadBlockBytesByHeight(uint64, uint32) ([]byte, []byte, error)
	ReadBlockBytesByHash(fields.Hash, uint32) ([]byte, error)
	ReadBlockHashByHeight(uint64) (fields.Hash, error)

	// tx
	ReadTransactionBytesByHash(fields.Hash) (uint64, []byte, error)
	TransactionIsExist(fields.Hash) (bool, error)

	// diamond
	SaveDiamond(*stores.DiamondSmelt) error
	ReadDiamond(fields.Bytes6) (*stores.DiamondSmelt, error)
	ReadDiamondByNumber(uint32) (*stores.DiamondSmelt, error)
}
