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

	GetBTCMoveLogTotalPage() (int, error)                        // 数据页数，每页100条数据
	GetBTCMoveLogPageData(int) ([]*stores.SatoshiGenesis, error) // 获取数据页

}
