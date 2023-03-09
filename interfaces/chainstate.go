package interfaces

import "github.com/hacash/core/fields"

type ChainState interface {
	ChainStateOperation

	// Get parent status
	GetParent() ChainState
	// Get all child States
	GetChilds() map[uint64]ChainState

	// Start a sub state
	ForkNextBlock(uint64, fields.Hash, Block) (ChainState, error)
	ForkSubChild() (ChainState, error)

	TraversalCopy(ChainState) error

	//GetReferBlock() (uint64, fields.Hash)
	SearchBaseStateByBlockHash(fields.Hash) (ChainState, error)

	// Destruction
	Destory() // Destroy, including deleting all sub States, caches, status data, etc

	// Judgment type
	IsImmutable() bool

	// Save on disk
	ImmutableWriteToDisk() (ChainStateImmutable, error)

	GetTotalNonEmptyAccountStatistics() []int64
}

// Immutable, non fallback lock status data
type ChainStateImmutable interface {
	ChainState

	// Traversing immature block hash
	SeekImmatureBlockHashs() ([]fields.Hash, error)

	Close() // Close file handle, etc
}
