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
	SetLatestBlockHash(hx fields.Hash)
	GetImmutableBlockHeadMeta() interfaces.BlockHeadMetaRead
	SetImmutableBlockHeadMeta(interfaces.BlockHeadMetaRead)
}

type PendingStatus interface {
	interfaces.Field

	GetPendingBlockHead() interfaces.BlockHeadMetaRead
	GetPendingBlockHeight() uint64
	GetPendingBlockHash() fields.Hash
	GetWaitingSubmitDiamond() *stores.DiamondSmelt
	SetWaitingSubmitDiamond(*stores.DiamondSmelt)
	ClearWaitingSubmitDiamond()
}
