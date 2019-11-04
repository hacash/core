package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// chain state 操作

type ChainStateOperation interface {

	// status

	Block() Block
	SetBlock(Block)
	Miner() Miner
	SetMiner(Miner)

	// state

	GetPrevDiamondHash() (uint32, fields.Hash)
	SetPrevDiamondHash(uint32, fields.Hash)

	// query

	Balance(fields.Address) fields.Amount
	Channel(fields.Bytes16) *stores.Channel
	Diamond(fields.Bytes6) *stores.Diamond

	// operate

	BalanceSet(fields.Address, fields.Amount)
	BalanceDel(fields.Address)

	ChannelCreate(fields.Bytes16, *stores.Channel)
	ChannelDelete(fields.Bytes16)

	DiamondSet(fields.Bytes6, *stores.Diamond)
	DiamondMove(fields.Bytes6, fields.Address)
	DiamondDel(fields.Bytes6)

}
