package interfaces

import "github.com/hacash/core/fields"

type Action interface {

	// the action type number
	Kind() uint16

	// assembling datas
	Serialize() ([]byte, error)
	Parse([]byte, uint32) (uint32, error)
	Size() uint32

	// Addresses that need to verify signatures
	RequestSignAddresses() []fields.Address

	// change chain state
	WriteinChainState(ChainStateOperation) error
	RecoverChainState(ChainStateOperation) error

	// help func
	SetBelongTransaction(Transaction)
	Describe() map[string]interface{} // json api

	// burning fees
	IsBurning90PersentTxFees() bool // 是否销毁本笔交易的 90% 的交易费用
}
