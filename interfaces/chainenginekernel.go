package interfaces

import (
	"github.com/hacash/core/stores"
)

type ChainEngineKernel interface {
	Start() error

	ChainStateIinitializeCall(func(ChainStateOperation))

	InsertBlock(Block, string) error

	StateRead() ChainStateOperationRead // 只读状态
	CurrentState() ChainState           // 最新的状态

	LatestBlock() (BlockHeadMetaRead, BlockHeadMetaRead, error) // 最新的区块（已确认的，和未成熟的）
	LatestDiamond() (*stores.DiamondSmelt, error)               // 最新的区块钻石

	SubscribeValidatedBlockOnInsert(chan Block)
	SubscribeDiamondOnCreate(chan *stores.DiamondSmelt)

	//RollbackToBlockHeight(uint64) (uint64, error)
}
