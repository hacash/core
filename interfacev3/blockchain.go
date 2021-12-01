package interfacev3

/*
type BlockChain interface {
	Start() error

	InsertBlock(interfacev2.Block, string) error

	StateImmutable() ChainStateImmutable // 已确认区块状态

	StateRead() interfaces.ChainStateOperationRead // 只读状态

	ValidateTransactionForTxPool(interfacev2.Transaction) error
	ValidateDiamondCreateAction(interfacev2.Action) error
	CreateNextBlockByValidateTxs([]interfacev2.Transaction) (interfacev2.Block, []interfacev2.Transaction, uint32, error)

	SubscribeValidatedBlockOnInsert(chan interfacev2.Block)
	SubscribeDiamondOnCreate(chan *stores.DiamondSmelt)

	ReadPrev288BlockTimestamp(blockHeight uint64) (uint64, error)

	RollbackToBlockHeight(uint64) (uint64, error)
}


*/
