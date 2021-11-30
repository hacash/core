package interfacev3

import (
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
)

type BlockChain interface {
	Start() error

	InsertBlock(interfacev2.Block, string) error

	StateImmutable() ChainStateImmutable // 已确认区块状态
	State() interfacev2.ChainState       // 不成熟未确认的状态

	StateRead() interfaces.ChainStateOperationRead // 只读状态

	ValidateTransaction(interfacev2.Transaction, func(interfacev2.ChainState)) error
	ValidateDiamondCreateAction(interfacev2.Action) error
	CreateNextBlockByValidateTxs([]interfacev2.Transaction) (interfacev2.Block, []interfacev2.Transaction, uint32, error)

	SubscribeValidatedBlockOnInsert(chan interfacev2.Block)
	SubscribeDiamondOnCreate(chan *stores.DiamondSmelt)

	ReadPrev288BlockTimestamp(blockHeight uint64) (uint64, error)

	RollbackToBlockHeight(uint64) (uint64, error)
}
