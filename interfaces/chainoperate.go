package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// chain state 操作

type ChainStateOperation interface {
	// 数据库升级模式
	IsDatabaseVersionRebuildMode() bool

	// status
	IsInMemTxPool() bool // 否在交易池
	SetInMemTxPool(bool)

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

	// state

	// query

	Balance(fields.Address) *stores.Balance
	//Satoshi(fields.Address) *stores.Satoshi
	Lockbls(fields.Bytes18) *stores.Lockbls
	Channel(fields.Bytes16) *stores.Channel
	Diamond(fields.Bytes6) *stores.Diamond
	DiamondSystemLending(fields.Bytes14) *stores.DiamondSystemLending
	BitcoinSystemLending(fields.Bytes15) *stores.BitcoinSystemLending
	UserLending(fields.Bytes17) *stores.UserLending

	// operate

	BalanceSet(fields.Address, *stores.Balance) error
	BalanceDel(fields.Address) error

	//SatoshiSet(fields.Address, *stores.Satoshi) error
	//SatoshiDel(fields.Address) error

	LockblsCreate(fields.Bytes18, *stores.Lockbls) error // 创建线性锁仓
	LockblsUpdate(fields.Bytes18, *stores.Lockbls) error // 更新：释放（取出部分任意可取额度）
	LockblsDelete(fields.Bytes18) error                  // 释放完毕后自动删除

	ChannelCreate(fields.Bytes16, *stores.Channel) error
	ChannelUpdate(fields.Bytes16, *stores.Channel) error
	ChannelDelete(fields.Bytes16) error

	DiamondSet(fields.Bytes6, *stores.Diamond) error
	DiamondDel(fields.Bytes6) error

	DiamondLendingCreate(fields.Bytes14, *stores.DiamondSystemLending) error
	DiamondLendingUpdate(fields.Bytes14, *stores.DiamondSystemLending) error
	DiamondLendingDelete(fields.Bytes14) error

	BitcoinLendingCreate(fields.Bytes15, *stores.BitcoinSystemLending) error
	BitcoinLendingUpdate(fields.Bytes15, *stores.BitcoinSystemLending) error
	BitcoinLendingDelete(fields.Bytes15) error

	UserLendingCreate(fields.Bytes17, *stores.UserLending) error
	UserLendingUpdate(fields.Bytes17, *stores.UserLending) error
	UserLendingDelete(fields.Bytes17) error

	// movebtc
	SaveMoveBTCBelongTxHash(trsno uint32, txhash []byte) error
	ReadMoveBTCTxHashByNumber(trsno uint32) ([]byte, error)
	LoadValidatedSatoshiGenesis(int64) (*stores.SatoshiGenesis, bool) // 获取已验证的BTC转移日志 & 是否需要验证

}
