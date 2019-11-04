package interfaces

import "github.com/hacash/core/fields"

type Miner interface {

	GetPrevDiamondHash() (uint32, fields.Hash)

	SetPrevDiamondHash(uint32, fields.Hash)

}
