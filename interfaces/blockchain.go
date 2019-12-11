package interfaces

type BlockChain interface {
	InsertBlock(Block) error

	State() ChainStateOperation
}
