package interfaces

import (
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
)

type Transaction interface {

	// base super
	Field

	Clone() Transaction

	// the transaction type number
	Type() uint8

	// hash data
	Hash() fields.Hash        // Hash without service charge
	HashWithFee() fields.Hash // inclued fee amount

	// Addresses that need to verify signatures
	RequestSignAddresses(appends []fields.Address, dropfeeaddr bool) ([]fields.Address, error)

	// sign
	CleanSigns()
	GetSigns() []fields.Sign // Return all signature data
	SetSigns([]fields.Sign)  // Set signature data
	// fill signatures
	FillTargetSign(signacc *account.Account) error           // Designated account signature
	FillNeedSigns(map[string][]byte, []fields.Address) error // All signatures
	// verify signatures
	VerifyAllNeedSigns() (bool, error)
	VerifyTargetSigns(reqaddr []fields.Address) (bool, error)

	// change chain state
	WriteInChainState(ChainStateOperation) error

	// fee
	FeePurity() uint64 // fee purity

	// get something
	GetAddress() fields.Address
	SetAddress(fields.Address)
	// Transaction costs actually received by miners
	// If some transactions receive less handling charges than they originally paid, the rest will be destroyed
	GetFeeOfMinerRealReceived() *fields.Amount
	GetFee() *fields.Amount
	SetFee(*fields.Amount)
	GetActionList() []Action
	GetTimestamp() uint64
	GetMessage() fields.TrimString16
	SetMessage(fields.TrimString16)
}
