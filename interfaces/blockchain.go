package interfaces

import "github.com/hacash/core/stores"

type BlockChain interface {
	InsertBlock(Block) error
	State() ChainStateOperation
	ValidateTransaction(Transaction) error
	CreateNextBlockByValidateTxs([]Transaction) (Block, error)

	SubscribeValidatedBlockOnInsert(chan Block)
	SubscribeDiamondOnCreate(chan *stores.DiamondSmelt)
}
