package interfacev2

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

type BlockStore interface {

	// close
	Close()

	// save
	SaveBlockUniteTransactions(Block) error
	// cancel
	CancelUniteTransactions(Block) error

	// block
	ReadBlockHeadBytesByHeight(uint64) ([]byte, error)
	ReadBlockHeadBytesByHash(fields.Hash) ([]byte, error)
	ReadBlockBytesLengthByHeight(uint64, uint32) (fields.Hash, []byte, error)
	ReadBlockBytesByHeight(uint64) (fields.Hash, []byte, error)
	ReadBlockBytesLengthByHash(fields.Hash, uint32) ([]byte, error)
	ReadBlockBytesByHash(fields.Hash) ([]byte, error)
	ReadBlockHashByHeight(uint64) (fields.Hash, error)

	// tx
	ReadTransactionBytesByHash(fields.Hash) (uint64, []byte, error)
	TransactionIsExist(fields.Hash) (bool, error)

	// diamond
	SaveDiamond(*stores.DiamondSmelt) error
	ReadDiamond(fields.DiamondName) (*stores.DiamondSmelt, error)
	ReadDiamondByNumber(uint32) (*stores.DiamondSmelt, error)

	// delete
	DeleteTransactionByHash(fields.Hash) error

	// btc move log
	GetBTCMoveLogTotalPage() (int, error)                        // 数据页数，每页100条数据
	GetBTCMoveLogPageData(int) ([]*stores.SatoshiGenesis, error) // 获取数据页
	SaveBTCMoveLogPageData(int, []*stores.SatoshiGenesis) error  // 保存数据页

}
