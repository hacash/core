package interfaces

import (
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
)

type Transaction interface {
	Copy() Transaction

	// the transaction type number
	Type() uint8

	// assembling datas
	Serialize() ([]byte, error)
	Parse([]byte, uint32) (uint32, error)
	Size() uint32

	// hash data
	Hash() fields.Hash        // 无手续费的哈希
	HashWithFee() fields.Hash // inclued fee amount

	// Addresses that need to verify signatures
	RequestSignAddresses(appends []fields.Address, dropfeeaddr bool) ([]fields.Address, error)

	// sign
	CleanSigns()
	GetSigns() []fields.Sign // 返回所有签名数据
	SetSigns([]fields.Sign)  // 设置签名数据
	// fill signatures
	FillTargetSign(signacc *account.Account) error           // 指定账户签名
	FillNeedSigns(map[string][]byte, []fields.Address) error // 全部签名
	// verify signatures
	VerifyNeedSigns([]fields.Address) (bool, error)
	VerifyTargetSign(reqaddr fields.Address) (bool, error)

	// change chain state
	WriteinChainState(ChainStateOperation) error
	RecoverChainState(ChainStateOperation) error

	// fee
	FeePurity() uint64 // fee purity

	// get something
	GetAddress() fields.Address
	SetAddress(fields.Address)
	// 矿工实际拿到的交易费用
	// 如果某些交易拿到手续费的比原本付出的少，则剩下的被销毁了
	GetFeeOfMinerRealReceived() *fields.Amount
	GetFee() *fields.Amount
	SetFee(*fields.Amount)
	GetActions() []Action
	GetTimestamp() uint64
	GetMessage() fields.TrimString16
	SetMessage(fields.TrimString16)
}
