package interfaces

type ConfirmTxIndexer interface {
	Init()
	// if ret= 0: next Tx
	// if ret= 1: next block
	ScanTx(block BlockHeadMetaRead, tx Transaction) int8
}
