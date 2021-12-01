package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

type LatestStatus interface {
	Field

	SetLastestDiamond(*stores.DiamondSmelt)
	ReadLastestDiamond() *stores.DiamondSmelt
}

type ImmutableStatus interface {
	Field

	GetImmatureBlockHashList() []fields.Hash
	SetImmatureBlockHashList([]fields.Hash)
	GetLatestBlockHash() fields.Hash
	SetLatestBlockHash(hx fields.Hash)
	GetImmutableBlockHeadMeta() BlockHeadMetaRead
	SetImmutableBlockHeadMeta(BlockHeadMetaRead)
}

type PendingStatus interface {
	Field

	GetPendingBlockHead() BlockHeadMetaRead
	GetPendingBlockHeight() uint64
	GetPendingBlockHash() fields.Hash
	GetWaitingSubmitDiamond() *stores.DiamondSmelt
	SetWaitingSubmitDiamond(*stores.DiamondSmelt)
	ClearWaitingSubmitDiamond()
}
