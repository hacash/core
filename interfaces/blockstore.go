package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

type BlockStore interface {
	BlockStoreRead

	// close
	Close()

	// save
	SaveBlock(Block) error

	// Set the block hash that the block height points to
	UpdateSetBlockHashReferToHeight(uint64, fields.Hash) error

	// tx
	//ReadTransactionBytesByHash(fields.Hash) (uint64, []byte, error)

	// diamond
	SaveDiamond(*stores.DiamondSmelt) error

	// Set the diamond name pointed by the diamond number
	UpdateSetDiamondNameReferToNumber(uint32, fields.DiamondName) error

	// btc move log
	RunDownLoadBTCMoveLog()
	SaveBTCMoveLogPageData(int, []*stores.SatoshiGenesis) error // Save data page

	LoadValidatedSatoshiGenesis(int64) (*stores.SatoshiGenesis, bool) // Get verified BTC transfer logs  获取已验证的BTC转移日志 & 是否需要验证 whether verification is required

}
