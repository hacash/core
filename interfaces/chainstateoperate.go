package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// Chain state operation

type ChainStateOperation interface {
	ChainStateOperationRead

	// Database upgrade mode
	SetDatabaseVersionRebuildMode(bool)

	// status
	SetInTxPool(bool)

	GetPending() PendingStatus
	SetPending(PendingStatus) error

	LatestStatusRead() (LatestStatus, error)
	LatestStatusSet(LatestStatus) error

	UpdateSetTotalSupply(totalobj *stores.TotalSupply) error

	// store
	BlockStore() BlockStore

	// tx hash
	ContainTxHash(fields.Hash, fields.BlockHeight) error // Write include transaction hash
	RemoveTxHash(fields.Hash) error                      // Remove transaction

	// operate

	BalanceSet(fields.Address, *stores.Balance) error
	BalanceDel(fields.Address) error

	LockblsCreate(fields.LockblsId, *stores.Lockbls) error // Create linear lock
	LockblsUpdate(fields.LockblsId, *stores.Lockbls) error // Update: release (take out any available quota)
	LockblsDelete(fields.LockblsId) error                  // Automatically delete after release

	ChannelCreate(fields.ChannelId, *stores.Channel) error
	ChannelUpdate(fields.ChannelId, *stores.Channel) error
	ChannelDelete(fields.ChannelId) error

	DiamondSet(fields.DiamondName, *stores.Diamond) error
	DiamondDel(fields.DiamondName) error

	DiamondLendingCreate(fields.DiamondSyslendId, *stores.DiamondSystemLending) error
	DiamondLendingUpdate(fields.DiamondSyslendId, *stores.DiamondSystemLending) error
	DiamondLendingDelete(fields.DiamondSyslendId) error

	BitcoinLendingCreate(fields.BitcoinSyslendId, *stores.BitcoinSystemLending) error
	BitcoinLendingUpdate(fields.BitcoinSyslendId, *stores.BitcoinSystemLending) error
	BitcoinLendingDelete(fields.BitcoinSyslendId) error

	UserLendingCreate(fields.UserLendingId, *stores.UserLending) error
	UserLendingUpdate(fields.UserLendingId, *stores.UserLending) error
	UserLendingDelete(fields.UserLendingId) error

	ChaswapCreate(fields.HashHalfChecker, *stores.Chaswap) error
	ChaswapUpdate(fields.HashHalfChecker, *stores.Chaswap) error
	ChaswapDelete(fields.HashHalfChecker) error

	// movebtc
	SaveMoveBTCBelongTxHash(trsno uint32, txhash []byte) error
	//ReadMoveBTCTxHashByNumber(trsno uint32) ([]byte, error)
	//ReadMoveBTCTxHashByTrsNo(trsno uint32) ([]byte, error)
	//LoadValidatedSatoshiGenesis(int64) (*stores.SatoshiGenesis, bool) // 获取已验证的BTC转移日志 & 是否需要验证

}
