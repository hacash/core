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

	// 设置区块高度指向的区块哈希
	UpdateSetBlockHashReferToHeight(uint64, fields.Hash) error

	// tx
	//ReadTransactionBytesByHash(fields.Hash) (uint64, []byte, error)

	// diamond
	SaveDiamond(*stores.DiamondSmelt) error

	// 设置钻石数字指向的钻石名称
	UpdateSetDiamondNameReferToNumber(uint32, fields.DiamondName) error

	// btc move log
	RunDownLoadBTCMoveLog()
	SaveBTCMoveLogPageData(int, []*stores.SatoshiGenesis) error // 保存数据页

	LoadValidatedSatoshiGenesis(int64) (*stores.SatoshiGenesis, bool) // 获取已验证的BTC转移日志 & 是否需要验证

}
