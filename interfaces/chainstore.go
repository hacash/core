package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

type ChainStore interface {

	// save
	SaveBlockUniteTransactions(Block) error

	// block
	ReadBlockHeadBytesByHeight(uint64) ([]byte, error)
	ReadBlockHeadBytesByHash(fields.Hash) ([]byte, error)
	ReadBlockBytesByHeight(uint64) ([]byte, error)
	ReadBlockBytesByHash(fields.Hash) ([]byte, error)

	// tx
	ReadTransactionDataByHash(fields.Hash) ([]byte, error)
	TransactionIsExist(fields.Hash) (bool, error)

	// diamond
	SaveDiamond(*stores.DiamondSmelt) error
	ReadDiamond(fields.Bytes6) (*stores.DiamondSmelt, error)
	ReadDiamondByNumber(uint32) (*stores.DiamondSmelt, error)
}
