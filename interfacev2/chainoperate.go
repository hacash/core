package interfacev2

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// Chain state operation

type ChainStateOperation interface {
	// Database upgrade mode
	IsDatabaseVersionRebuildMode() bool
	SetDatabaseVersionRebuildMode(bool)

	// status
	IsInTxPool() bool // No in the trading pool
	SetInTxPool(bool)

	// status
	GetPendingBlockHeight() uint64
	SetPendingBlockHeight(uint64)
	GetPendingBlockHash() fields.Hash
	SetPendingBlockHash(fields.Hash)

	GetPendingSubmitStoreDiamond() (*stores.DiamondSmelt, error)
	SetPendingSubmitStoreDiamond(*stores.DiamondSmelt) error

	// status
	SetLastestBlockHeadAndMeta(Block) error
	ReadLastestBlockHeadAndMeta() (Block, error)
	SetLastestDiamond(*stores.DiamondSmelt) error
	ReadLastestDiamond() (*stores.DiamondSmelt, error)

	UpdateSetTotalSupply(totalobj *stores.TotalSupply) error
	ReadTotalSupply() (*stores.TotalSupply, error)

	// store
	BlockStore() BlockStore
	SetBlockStore(BlockStore) error

	// tx hash
	ContainTxHash(fields.Hash, fields.BlockHeight) error // Write include transaction hash
	RemoveTxHash(fields.Hash) error                      // Remove transaction
	CheckTxHash(fields.Hash) (bool, error)               // Check whether the transaction has been linked

	// query

	Balance(fields.Address) (*stores.Balance, error)
	//Satoshi(fields.Address) *stores.Satoshi
	Lockbls(fields.LockblsId) (*stores.Lockbls, error)
	Channel(fields.ChannelId) (*stores.Channel, error)
	Diamond(fields.DiamondName) (*stores.Diamond, error)
	DiamondSystemLending(fields.DiamondSyslendId) (*stores.DiamondSystemLending, error)
	BitcoinSystemLending(fields.BitcoinSyslendId) (*stores.BitcoinSystemLending, error)
	UserLending(fields.UserLendingId) (*stores.UserLending, error)
	Chaswap(fields.HashHalfChecker) (*stores.Chaswap, error)

	// operate

	BalanceSet(fields.Address, *stores.Balance) error
	BalanceDel(fields.Address) error

	//SatoshiSet(fields.Address, *stores.Satoshi) error
	//SatoshiDel(fields.Address) error

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
	ReadMoveBTCTxHashByNumber(trsno uint32) ([]byte, error)
	LoadValidatedSatoshiGenesis(int64) (*stores.SatoshiGenesis, bool) // Get verified BTC transfer logs  获取已验证的BTC转移日志 & 是否需要验证 whether verification is required

}
