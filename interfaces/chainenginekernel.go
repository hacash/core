package interfaces

import (
	"github.com/hacash/core/stores"
)

type ChainEngine interface {
	Start() error
	Close() error

	ChainStateIinitializeCall(func(ChainStateOperation))

	InsertBlock(Block, string) error
	GetRecentArrivedBlocks() []Block
	GetLatestAverageFeePurity() uint32

	StateRead() ChainStateOperationRead // Read-only status
	CurrentState() ChainState           // Latest status

	LatestBlock() (BlockHeadMetaRead, BlockHeadMetaRead, error) // Latest blocks (confirmed and immature)
	LatestDiamond() (*stores.DiamondSmelt, error)               // Latest block diamonds

	SubscribeValidatedBlockOnInsert(chan Block)
	SubscribeDiamondOnCreate(chan *stores.DiamondSmelt)

	//RollbackToBlockHeight(uint64) (uint64, error)
}
