package interfaces

import (
	"github.com/hacash/core/fields"
)

type TxPool interface {
	SetBlockChain(BlockChain)

	// Check whether the transaction already exists
	CheckTxExistByHash(fields.Hash) (Transaction, bool)
	CheckTxExist(Transaction) (Transaction, bool)
	// Add transaction
	AddTx(Transaction) error
	// Query a transaction from the trading pool
	// FindTxByHash(fields.Hash) (Transaction, bool)
	// Obtain all transactions, and sort them according to the service fee purity from high to low
	CopyTxsOrderByFeePurity(targetblockheight uint64, maxcount uint32, maxsize uint32) []Transaction
	// Filter and clear transactions
	RemoveTxs([]Transaction)
	RemoveTxsOnNextBlockArrive([]Transaction)
	SetAutomaticallyCleanInvalidTransactions(bool)
	// Add transaction success event subscription
	SubscribeOnAddTxSuccess(chan Transaction)
	PauseEventSubscribe()   // Pause event subscription
	RenewalEventSubscribe() // Reopen event subscription

	GetDiamondCreateTxs(int) []Transaction

	// Get the transaction with the highest handling fee
	// PopTxByHighestFee() Transaction
	// Subscribe to the transaction pool and add a new transaction event
	// SubscribeNewTx(chan<- []block.Transaction) event.Subscription
}
