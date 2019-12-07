package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

type ChainStore interface {

	// save
	SaveBlockUniteTransactions(Block) error

	// block
	ReadBlockHeadBytesByHash(fields.Hash) ([]byte, error)
	ReadBlockBytesByHash(fields.Hash) ([]byte, error)
	ReadBlockHeadBytesByHeight(fields.VarInt5) ([]byte, error)
	ReadBlockBytesByHeight(fields.VarInt5) ([]byte, error)

	// tx
	ReadTransactionDataByHash(fields.Hash) ([]byte, error)
	TransactionIsExist(fields.Hash) (bool, error)

	// diamond
	SaveDiamond(*stores.DiamondSmelt) error
	ReadDiamond(fields.Bytes6) (*stores.DiamondSmelt, error)
	ReadDiamondByNumber(uint32) (*stores.DiamondSmelt, error)
}
