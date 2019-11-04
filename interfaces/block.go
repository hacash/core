package interfaces

import (
	"github.com/hacash/core/fields"
)

type Block interface {

	// the block type number
	Version() uint8

	// assembling datas
	Serialize() ([]byte, error)
	Parse([]byte, uint32) (uint32, error)
	Size() uint32

	SerializeHead() ([]byte, error)
	ParseHead([]byte, uint32) (uint32, error)

	SerializeBody() ([]byte, error)
	ParseBody([]byte, uint32) (uint32, error)

	SerializeMeta() ([]byte, error)
	ParseMeta([]byte, uint32) (uint32, error)

	SerializeTransactions(SerializeTransactionsIterator) ([]byte, error)
	ParseTransactions([]byte, uint32) (uint32, error)

	ParseExcludeTransactions([]byte, uint32) (uint32, error)

	// change chain state
	WriteinChainState(ChainStateOperation) error
	RecoverChainState(ChainStateOperation) error

	// hash
	Hash() fields.Hash
	HashFresh() fields.Hash

	// delete cache data
	Fresh()

	// change struct data
	AddTransaction(Transaction)
	SetMrklRoot(fields.Hash)
	SetNonce(uint32)

	// verify signatures
	VerifyNeedSigns() (bool, error)

	// diamond state
	CheckHasHaveDiamond(string) bool
	DoMarkHaveDiamond(string)

	// get some datas
	GetTransactions() []Transaction
	GetHeight() uint64
	GetDifficulty() uint32
	GetNonce() uint32
	GetPrevHash() fields.Hash
	GetTimestamp() uint64
	GetTransactionCount() uint32
	GetMrklRoot() fields.Hash
}



type SerializeTransactionsIterator interface {
	Init(uint32)
	FinishOneTrs(uint32, Transaction, []byte)
}
