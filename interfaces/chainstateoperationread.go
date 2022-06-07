package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// Chain state operation

type ChainStateOperationRead interface {
	// Database upgrade mode
	IsDatabaseVersionRebuildMode() bool

	// status
	IsInTxPool() bool // No in the trading pool

	//GetPending() (PendingStatus)
	GetPendingBlockHeight() uint64
	GetPendingBlockHash() fields.Hash

	//LatestStatusRead() (LatestStatus, error)
	//ReadLastestBlockHeadMetaForRead() (BlockHeadMetaRead, error)
	ReadLastestDiamond() (*stores.DiamondSmelt, error)

	ReadTotalSupply() (*stores.TotalSupply, error)

	// store
	BlockStoreRead() BlockStoreRead

	// tx hash
	CheckTxHash(fields.Hash) (bool, error)                                      // Check whether the transaction has been linked
	ReadTxBelongHeightByHash(fields.Hash) (fields.BlockHeight, error)           // Check the block ID of the exchange
	ReadTransactionBytesByHash(fields.Hash) (fields.BlockHeight, []byte, error) // Read transaction content

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
