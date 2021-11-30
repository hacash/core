package interfacev2

import (
	"github.com/hacash/core/fields"
)

type TxPool interface {
	SetBlockChain(BlockChain)

	// 检查交易是否已经存在
	CheckTxExistByHash(fields.Hash) (Transaction, bool)
	CheckTxExist(Transaction) (Transaction, bool)
	// 添加交易
	AddTx(Transaction) error
	// 从交易池里查询一笔交易
	// FindTxByHash(fields.Hash) (Transaction, bool)
	// 获取全部交易，按手续费纯度高到低排序
	CopyTxsOrderByFeePurity(targetblockheight uint64, maxcount uint32, maxsize uint32) []Transaction
	// 过滤、清除交易
	RemoveTxs([]Transaction)
	RemoveTxsOnNextBlockArrive([]Transaction)
	SetAutomaticallyCleanInvalidTransactions(bool)
	// 添加交易成功事件订阅
	SubscribeOnAddTxSuccess(chan Transaction)
	PauseEventSubscribe()   // 暂停事件订阅
	RenewalEventSubscribe() // 重开事件订阅

	GetDiamondCreateTxs(int) []Transaction

	// 获取手续费最高的一笔交易
	// PopTxByHighestFee() Transaction
	// 订阅交易池加入新交易事件
	// SubscribeNewTx(chan<- []block.Transaction) event.Subscription
}
