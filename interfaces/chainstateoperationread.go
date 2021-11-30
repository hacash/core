package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// chain state 操作

type ChainStateOperationRead interface {
	// 数据库升级模式
	IsDatabaseVersionRebuildMode() bool

	// status
	IsInMemTxPool() bool // 否在交易池

	//GetPending() (PendingStatus)
	GetPendingBlockHeight() uint64
	GetPendingBlockHash() fields.Hash

	//LatestStatusRead() (LatestStatus, error)
	ReadLastestBlockHeadMetaForRead() (BlockHeadMetaRead, error)
	ReadLastestDiamond() (*stores.DiamondSmelt, error)

	ReadTotalSupply() (*stores.TotalSupply, error)

	// store
	BlockStoreRead() BlockStoreRead

	// tx hash
	CheckTxHash(fields.Hash) (bool, error)                                      // 检查交易是否已经上链
	ReadTxBelongHeightByHash(fields.Hash) (fields.BlockHeight, error)           // 检查交易所属区块ID
	ReadTransactionBytesByHash(fields.Hash) (fields.BlockHeight, []byte, error) // 读取交易内容

	// query
	Balance(fields.Address) (*stores.Balance, error)
	Lockbls(fields.LockblsId) (*stores.Lockbls, error)
	Channel(fields.ChannelId) (*stores.Channel, error)
	Diamond(fields.DiamondName) (*stores.Diamond, error)
	DiamondSystemLending(fields.DiamondSyslendId) (*stores.DiamondSystemLending, error)
	BitcoinSystemLending(fields.BitcoinSyslendId) (*stores.BitcoinSystemLending, error)
	UserLending(fields.UserLendingId) (*stores.UserLending, error)
	Chaswap(fields.HashHalfChecker) (*stores.Chaswap, error)

	// movebtc
	ReadMoveBTCTxHashByTrsNo(trsno uint32) ([]byte, error)
	//LoadValidatedSatoshiGenesis(int64) (*stores.SatoshiGenesis, bool) // 获取已验证的BTC转移日志 & 是否需要验证

	// tx

}
