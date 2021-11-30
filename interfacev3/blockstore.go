package interfacev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

type BlockStore interface {

	// close
	Close()

	// save
	SaveBlock(Block) error

	// read
	ReadBlockBytesByHash(fields.Hash) ([]byte, error)
	ReadBlockBytesByHeight(uint64) (fields.Hash, []byte, error)
	ReadBlockHashByHeight(uint64) (fields.Hash, error)

	// 设置区块高度指向的区块哈希
	UpdateSetBlockHashReferToHeight(uint64, fields.Hash) error

	// tx
	//ReadTransactionBytesByHash(fields.Hash) (uint64, []byte, error)

	// diamond
	SaveDiamond(*stores.DiamondSmelt) error
	ReadDiamond(fields.DiamondName) (*stores.DiamondSmelt, error)
	ReadDiamondByNumber(uint32) (*stores.DiamondSmelt, error)

	// 设置钻石数字指向的钻石名称
	UpdateSetDiamondNameReferToNumber(uint32, fields.DiamondName) error

	// btc move log
	RunDownLoadBTCMoveLog()
	GetBTCMoveLogTotalPage() (int, error)                        // 数据页数，每页100条数据
	GetBTCMoveLogPageData(int) ([]*stores.SatoshiGenesis, error) // 获取数据页
	SaveBTCMoveLogPageData(int, []*stores.SatoshiGenesis) error  // 保存数据页

	LoadValidatedSatoshiGenesis(int64) (*stores.SatoshiGenesis, bool) // 获取已验证的BTC转移日志 & 是否需要验证

}
