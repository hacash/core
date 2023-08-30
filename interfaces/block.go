package interfaces

import (
	"github.com/hacash/core/fields"
)

type Block interface {

	// base super
	Field

	// origin
	OriginMark() string // "", "sync", "discover", "mining"
	SetOriginMark(string)
	ArrivedTime() int64 // timestamp
	SetArrivedTime(int64)

	// copy
	CopyHeadMetaForMining() Block
	CopyForMining() Block

	// delete cache data
	Fresh()

	////////////////////////////////////

	// the block type number
	Version() uint8

	SerializeHead() ([]byte, error)
	ParseHead([]byte, uint32) (uint32, error)

	SerializeBody() ([]byte, error)
	ParseBody([]byte, uint32) (uint32, error)

	SerializeMeta() ([]byte, error)
	ParseMeta([]byte, uint32) (uint32, error)

	ParseTransactions([]byte, uint32) (uint32, error)

	SerializeExcludeTransactions() ([]byte, error)
	ParseExcludeTransactions([]byte, uint32) (uint32, error)

	// change chain state
	WriteInChainState(ChainStateOperation) error

	// hash
	Hash() fields.Hash
	HashFresh() fields.Hash

	// change struct data
	AddTrs(Transaction)
	SetTrsList([]Transaction)
	SetMrklRoot(fields.Hash)
	SetNonce(uint32)
	SetNonceByte(nonce []byte)

	// verify signatures
	VerifyNeedSigns() (bool, error)

	// get some datas

	GetTrsList() []Transaction
	GetHeight() uint64
	GetDifficulty() uint32
	GetWitnessStage() uint16
	GetNonce() uint32
	GetNonceByte() []byte
	GetPrevHash() fields.Hash
	GetTimestamp() uint64
	GetTransactionCount() uint32
	GetMrklRoot() fields.Hash

	// customer trs count
	GetCustomerTransactionCount() uint32
}
