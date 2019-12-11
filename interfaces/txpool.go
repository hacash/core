package interfaces

import "github.com/hacash/core/fields"

type TxPool interface {
	// 检查交易是否已经存在
	CheckTxExist(Transaction) bool
	// 添加交易
	AddTx(Transaction) error
	// 从交易池里查询一笔交易
	FindTxByHash(fields.Hash) (Transaction, bool)
	// 获取全部交易，按手续费纯度高到低排序
	CopyTxsOrderByFeePurity() []Transaction
	// 过滤、清除交易
	RemoveTxs([]Transaction)
	// 获取手续费最高的一笔交易
	// PopTxByHighestFee() Transaction
	// 订阅交易池加入新交易事件
	// SubscribeNewTx(chan<- []block.Transaction) event.Subscription
}
