package interfacev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

type LatestStatus interface {
	interfaces.Field

	SetLastestDiamond(*stores.DiamondSmelt)
	ReadLastestDiamond() *stores.DiamondSmelt
	GetImmatureBlockHashList() []fields.Hash
	SetImmatureBlockHashList([]fields.Hash)
	GetLatestBlockHash() fields.Hash
	GetImmutableBlockHeadMeta() Block
	SetImmutableBlockHeadMeta(Block)
}

type PendingStatus interface {
	interfaces.Field

	GetPendingBlockHead() Block
	GetPendingBlockHeight() uint64
	GetPendingBlockHash() fields.Hash
	GetWaitingSubmitDiamond() *stores.DiamondSmelt
	SetWaitingSubmitDiamond(*stores.DiamondSmelt)
	ClearWaitingSubmitDiamond()
}
