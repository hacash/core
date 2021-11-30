package interfacev2

import (
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

type BlockChain interface {
	Start() error

	InsertBlock(Block, string) error
	State() ChainState
	StateRead() interfaces.ChainStateOperationRead

	ValidateTransactionForTxPool(Transaction) error
	ValidateDiamondCreateAction(Action) error
	CreateNextBlockByValidateTxs([]Transaction) (Block, []Transaction, uint32, error)

	SubscribeValidatedBlockOnInsert(chan Block)
	SubscribeDiamondOnCreate(chan *stores.DiamondSmelt)

	ReadPrev288BlockTimestamp(blockHeight uint64) (uint64, error)

	RollbackToBlockHeight(uint64) (uint64, error)
}
