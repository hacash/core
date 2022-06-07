package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

type BlockStoreRead interface {
	// read
	ReadBlockBytesByHash(fields.Hash) ([]byte, error)
	ReadBlockBytesByHeight(uint64) (fields.Hash, []byte, error)
	ReadBlockHashByHeight(uint64) (fields.Hash, error)

	ReadDiamond(fields.DiamondName) (*stores.DiamondSmelt, error)
	ReadDiamondByNumber(uint32) (*stores.DiamondSmelt, error)
	ReadDiamondNameByNumber(uint32) (fields.DiamondName, error)

	GetBTCMoveLogTotalPage() (int, error)                        // Number of data pages, 100 pieces of data per page
	GetBTCMoveLogPageData(int) ([]*stores.SatoshiGenesis, error) // Get data page

}
