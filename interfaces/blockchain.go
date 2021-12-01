package interfaces

type BlockChain interface {
	Start() error

	GetChainEngineKernel() ChainEngine
	SetChainEngineKernel(ChainEngine)

	ValidateTransactionForTxPool(Transaction) error
	ValidateDiamondCreateAction(Action) error
	CreateNextBlockByValidateTxs([]Transaction) (Block, []Transaction, uint32, error)

	//ReadPrev288BlockTimestamp(blockHeight uint64) (uint64, error)
}
