package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// chain state 操作

type ChainStateOperation interface {

	// status
	IsInMemTxPool() bool // 否在交易池
	SetInMemTxPool(bool)

	LoadValidatedSatoshiGenesis(int64) (*stores.SatoshiGenesis, bool) // 获取已验证的BTC转移日志 & 是否需要验证

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

	// store
	BlockStore() BlockStore
	SetBlockStore(BlockStore) error

	// state

	// query

	Balance(fields.Address) *stores.Balance
	Satoshi(fields.Address) *stores.Satoshi
	Channel(fields.Bytes16) *stores.Channel
	Diamond(fields.Bytes6) *stores.Diamond

	// operate

	BalanceSet(fields.Address, *stores.Balance) error
	BalanceDel(fields.Address) error

	SatoshiSet(fields.Address, *stores.Satoshi) error
	SatoshiDel(fields.Address) error

	ChannelCreate(fields.Bytes16, *stores.Channel) error
	ChannelUpdate(fields.Bytes16, *stores.Channel) error
	ChannelDelete(fields.Bytes16) error

	DiamondSet(fields.Bytes6, *stores.Diamond) error
	DiamondDel(fields.Bytes6) error

	// movebtc
	SaveMoveBTCBelongTxHash(trsno uint32, txhash []byte) error
	ReadMoveBTCTxHashByNumber(trsno uint32) ([]byte, error)
}
