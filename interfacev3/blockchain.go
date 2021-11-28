package interfacev3

import (
	"github.com/hacash/core/stores"
)

type BlockChain interface {
	Start() error

	InsertBlock(Block, string) error

	StateImmutable() ChainStateImmutable // 已确认区块状态
	State() ChainState                   // 不成熟未确认的状态

	ValidateTransaction(Transaction, func(ChainState)) error
	ValidateDiamondCreateAction(Action) error
	CreateNextBlockByValidateTxs([]Transaction) (Block, []Transaction, uint32, error)

	SubscribeValidatedBlockOnInsert(chan Block)
	SubscribeDiamondOnCreate(chan *stores.DiamondSmelt)

	ReadPrev288BlockTimestamp(blockHeight uint64) (uint64, error)

	RollbackToBlockHeight(uint64) (uint64, error)
}
