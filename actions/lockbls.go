package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
	"math/big"
)

type Action_9_LockblsCreate struct {
	LockblsId           fields.LockblsId   // Linear lock ID
	PaymentAddress      fields.Address     // Payment address
	MasterAddress       fields.Address     // Main address (claim)
	EffectBlockHeight   fields.BlockHeight // Effective (start) block
	LinearBlockNumber   fields.VarUint3    // Number of stepping blocks < 17000000 about 160 years
	TotalStockAmount    fields.Amount      // Total deposit limit
	LinearReleaseAmount fields.Amount      // Limit released each time

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func NewAction_9_LockblsCreate() *Action_9_LockblsCreate {
	return &Action_9_LockblsCreate{}
}

func (elm *Action_9_LockblsCreate) Kind() uint16 {
	return 9
}

func (elm *Action_9_LockblsCreate) Size() uint32 {
	return 2 +
		elm.LockblsId.Size() +
		elm.PaymentAddress.Size() +
		elm.MasterAddress.Size() +
		elm.EffectBlockHeight.Size() +
		elm.LinearBlockNumber.Size() +
		elm.TotalStockAmount.Size() +
		elm.LinearReleaseAmount.Size()
}

// json api
func (elm *Action_9_LockblsCreate) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_9_LockblsCreate) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var b1, _ = elm.LockblsId.Serialize()
	var b2, _ = elm.PaymentAddress.Serialize()
	var b3, _ = elm.MasterAddress.Serialize()
	var b4, _ = elm.EffectBlockHeight.Serialize()
	var b5, _ = elm.LinearBlockNumber.Serialize()
	var b6, _ = elm.TotalStockAmount.Serialize()
	var b7, _ = elm.LinearReleaseAmount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	buffer.Write(b7)
	return buffer.Bytes(), nil
}

func (elm *Action_9_LockblsCreate) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	sk1, e := elm.LockblsId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	sk2, e := elm.PaymentAddress.Parse(buf, sk1)
	if e != nil {
		return 0, e
	}
	sk3, e := elm.MasterAddress.Parse(buf, sk2)
	if e != nil {
		return 0, e
	}
	sk4, e := elm.EffectBlockHeight.Parse(buf, sk3)
	if e != nil {
		return 0, e
	}
	sk5, e := elm.LinearBlockNumber.Parse(buf, sk4)
	if e != nil {
		return 0, e
	}
	sk6, e := elm.TotalStockAmount.Parse(buf, sk5)
	if e != nil {
		return 0, e
	}
	sk7, e := elm.LinearReleaseAmount.Parse(buf, sk6)
	if e != nil {
		return 0, e
	}
	return sk7, nil
}

func (act *Action_9_LockblsCreate) RequestSignAddresses() []fields.Address {
	return []fields.Address{
		act.PaymentAddress, // Signature is required for the payment account of warehouse lock
	}
}

func (act *Action_9_LockblsCreate) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Check the validity of ID value
	if len(act.LockblsId) != stores.LockblsIdLength || act.LockblsId[0] == 0 || act.LockblsId[stores.LockblsIdLength-1] == 0 {
		// The first and last digits of the lock ID created by the user cannot be zero
		// The first zero ID is the lockup ID of bitcoin one-way transfer
		return fmt.Errorf("LockblsId format error.")
	}
	// Check whether the key already exists
	haslock, e := state.Lockbls(act.LockblsId)
	if e != nil {
		return e
	}
	if haslock != nil {
		return fmt.Errorf("Lockbls id<%s> already.", hex.EncodeToString(act.LockblsId))
	}
	// Check step block number
	if act.LinearBlockNumber < 288 {
		return fmt.Errorf("LinearBlockNumber cannot less 288.")
	}
	if act.LinearBlockNumber > 1600*10000 {
		return fmt.Errorf("LinearBlockNumber cannot over 16000000.")
	}
	// Check amount
	if !act.TotalStockAmount.IsPositive() || !act.LinearReleaseAmount.IsPositive() {
		return fmt.Errorf("TotalStockAmount or LinearReleaseAmount error.")
	}
	// Check balance
	mainblsamt, e := state.Balance(act.PaymentAddress)
	if e != nil {
		return e
	}
	if mainblsamt == nil {
		return fmt.Errorf("Balance cannot empty.")
	}
	if mainblsamt.Hacash.LessThan(&act.TotalStockAmount) {
		return fmt.Errorf("Balance not enough.")
	}
	// Step cannot be greater than deposit amount
	if act.TotalStockAmount.LessThan(&act.LinearReleaseAmount) {
		return fmt.Errorf("LinearReleaseAmount cannot more than TotalStockAmount.")
	}

	// storage
	lockbls := stores.NewEmptyLockbls(act.MasterAddress)
	lockbls.EffectBlockHeight = act.EffectBlockHeight
	lockbls.LinearBlockNumber = act.LinearBlockNumber
	lockbls.TotalLockAmount = act.TotalStockAmount
	lockbls.BalanceAmount = act.TotalStockAmount
	lockbls.LinearReleaseAmount = act.LinearReleaseAmount
	// Deduct payment
	e1 := DoSubBalanceFromChainStateV3(state, act.PaymentAddress, act.TotalStockAmount)
	if e1 != nil {
		return e1
	}

	// Save lock
	e2 := state.LockblsCreate(act.LockblsId, lockbls)
	if e2 != nil {
		return e2
	}

	// ok
	return nil
}

func (act *Action_9_LockblsCreate) WriteinChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Check the validity of ID value
	if len(act.LockblsId) != stores.LockblsIdLength || act.LockblsId[0] == 0 || act.LockblsId[stores.LockblsIdLength-1] == 0 {
		// The first and last digits of the lock ID created by the user cannot be zero
		// The first zero ID is the lockup ID of bitcoin one-way transfer
		return fmt.Errorf("LockblsId format error.")
	}
	// Check whether the key already exists
	haslock, e := state.Lockbls(act.LockblsId)
	if e != nil {
		return e
	}
	if haslock != nil {
		return fmt.Errorf("Lockbls id<%s> already.", hex.EncodeToString(act.LockblsId))
	}
	// Check step block number
	if act.LinearBlockNumber < 288 {
		return fmt.Errorf("LinearBlockNumber cannot less 288.")
	}
	if act.LinearBlockNumber > 1600*10000 {
		return fmt.Errorf("LinearBlockNumber cannot over 16000000.")
	}
	// Check amount
	if !act.TotalStockAmount.IsPositive() || !act.LinearReleaseAmount.IsPositive() {
		return fmt.Errorf("TotalStockAmount or LinearReleaseAmount error.")
	}
	// Check balance
	mainblsamt, e := state.Balance(act.PaymentAddress)
	if e != nil {
		return e
	}
	if mainblsamt == nil {
		return fmt.Errorf("Balance cannot empty.")
	}
	if mainblsamt.Hacash.LessThan(&act.TotalStockAmount) {
		return fmt.Errorf("Balance not enough.")
	}
	// Step cannot be greater than deposit amount
	if act.TotalStockAmount.LessThan(&act.LinearReleaseAmount) {
		return fmt.Errorf("LinearReleaseAmount cannot more than TotalStockAmount.")
	}

	// storage
	lockbls := stores.NewEmptyLockbls(act.MasterAddress)
	lockbls.EffectBlockHeight = act.EffectBlockHeight
	lockbls.LinearBlockNumber = act.LinearBlockNumber
	lockbls.TotalLockAmount = act.TotalStockAmount
	lockbls.BalanceAmount = act.TotalStockAmount
	lockbls.LinearReleaseAmount = act.LinearReleaseAmount
	// Deduct payment
	e1 := DoSubBalanceFromChainState(state, act.PaymentAddress, act.TotalStockAmount)
	if e1 != nil {
		return e1
	}

	// Save lock
	e2 := state.LockblsCreate(act.LockblsId, lockbls)
	if e2 != nil {
		return e2
	}

	// ok
	return nil
}

func (act *Action_9_LockblsCreate) RecoverChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// Fallback HAC
	e1 := DoAddBalanceFromChainState(state, act.PaymentAddress, act.TotalStockAmount)
	if e1 != nil {
		return e1
	}
	// Delete lockbls
	e2 := state.LockblsDelete(act.LockblsId)
	if e2 != nil {
		return e2
	}

	// ok
	return nil
}

// Set belongs to long_ trs
func (act *Action_9_LockblsCreate) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}

func (act *Action_9_LockblsCreate) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_9_LockblsCreate) IsBurning90PersentTxFees() bool {
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////

type Action_10_LockblsRelease struct {
	LockblsId     fields.LockblsId // Linear lock ID
	ReleaseAmount fields.Amount    // Current withdrawal limit

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func NewAction_10_LockblsRelease() *Action_10_LockblsRelease {
	return &Action_10_LockblsRelease{}
}

func (elm *Action_10_LockblsRelease) Kind() uint16 {
	return 10
}

// json api
func (elm *Action_10_LockblsRelease) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_10_LockblsRelease) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var idBytes, _ = elm.LockblsId.Serialize()
	var amtBytes, _ = elm.ReleaseAmount.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(idBytes)
	buffer.Write(amtBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_10_LockblsRelease) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	moveseek, e := elm.LockblsId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	moveseek2, e := elm.ReleaseAmount.Parse(buf, moveseek)
	if e != nil {
		return 0, e
	}
	return moveseek2, nil
}

func (elm *Action_10_LockblsRelease) Size() uint32 {
	return 2 + elm.LockblsId.Size() + elm.ReleaseAmount.Size()
}

func (elm *Action_10_LockblsRelease) RequestSignAddresses() []fields.Address {
	return []fields.Address{}
}

func (act *Action_10_LockblsRelease) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Because only the specified address can be extracted, anyone can extract it without the signature of the lock up address
	// query
	lockbls, e := state.Lockbls(act.LockblsId)
	if e != nil {
		return e
	}
	if lockbls == nil {
		return fmt.Errorf("Lockbls id<%s> not find.", hex.EncodeToString(act.LockblsId))
	}
	// Extract
	currentBlockHeight := state.GetPendingBlockHeight()
	if currentBlockHeight < uint64(lockbls.EffectBlockHeight) {
		return fmt.Errorf("EffectBlockHeight be set %d", lockbls.EffectBlockHeight)
	}
	// Calculate withdrawal limit
	// Rlsnum = = extractable times
	rlsnum := (currentBlockHeight - uint64(lockbls.EffectBlockHeight)) / uint64(lockbls.LinearBlockNumber)
	if rlsnum == 0 {
		return fmt.Errorf("first release Block Height is %d, ", uint64(lockbls.EffectBlockHeight)+uint64(lockbls.LinearBlockNumber))
	}
	totalrlsamt := lockbls.TotalLockAmount
	steprlsamt := lockbls.LinearReleaseAmount
	// Effective withdrawable balance
	lockblsamt := lockbls.BalanceAmount
	// contrast
	if lockblsamt.LessThan(&act.ReleaseAmount) {
		return fmt.Errorf("BalanceAmount not enough.") // Sorry, your credit is running low
	}
	maxrlsamtbig := new(big.Int).Mul(steprlsamt.GetValue(), new(big.Int).SetUint64(rlsnum))
	currentMaxReleaseAmount, e3 := fields.NewAmountByBigInt(maxrlsamtbig)
	if e3 != nil {
		return e3
	}
	// The withdrawable balance shall be deducted from the withdrawn balance
	alreadyExtractedAmount, e9 := totalrlsamt.Sub(&lockblsamt) // 已经提走的余额
	if e9 != nil {
		return e9
	}
	// Effective withdrawable balance
	currentMaxReleaseAmount, e9 = currentMaxReleaseAmount.Sub(alreadyExtractedAmount)
	if e9 != nil {
		return e9
	}
	// Available balance judgment
	if currentMaxReleaseAmount.LessThan(&act.ReleaseAmount) {
		return fmt.Errorf("Current Max Release Amount not enough.") // The current available balance is insufficient
	}

	// Update lock in balance
	newBalanceAmount, e4 := lockblsamt.Sub(&act.ReleaseAmount)
	if e4 != nil {
		return e4
	}
	lockbls.BalanceAmount = *newBalanceAmount
	if newBalanceAmount.IsEmpty() {
		// All the locked warehouses have been extracted and deleted
		// Not deleted for rollback, but saved for block rollback
		e := state.LockblsUpdate(act.LockblsId, lockbls)
		if e != nil {
			return e
		}
	} else {
		// Deduct storage
		e := state.LockblsUpdate(act.LockblsId, lockbls)
		if e != nil {
			return e
		}
	}
	// Total supply statistics
	isbtcmoveunlock := act.LockblsId[0] == 0 // 第一位为 0 则是比特币转移的锁定
	if isbtcmoveunlock {
		totalsupply, e2 := state.ReadTotalSupply()
		if e2 != nil {
			return e2
		}
		// Cumulative unlocked HAC
		addamt := act.ReleaseAmount.ToMei()
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBitcoinTransferUnlockSuccessed, addamt)
		// update total supply
		e3 := state.UpdateSetTotalSupply(totalsupply)
		if e3 != nil {
			return e3
		}
	}
	// Plus balance
	return DoAddBalanceFromChainStateV3(state, lockbls.MasterAddress, act.ReleaseAmount)
}

func (act *Action_10_LockblsRelease) WriteinChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Because only the specified address can be extracted, anyone can extract it without the signature of the lock up address
	// query
	lockbls, e := state.Lockbls(act.LockblsId)
	if e != nil {
		return e
	}
	if lockbls == nil {
		return fmt.Errorf("Lockbls id<%s> not find.", hex.EncodeToString(act.LockblsId))
	}
	// Extract
	currentBlockHeight := state.GetPendingBlockHeight()
	if currentBlockHeight < uint64(lockbls.EffectBlockHeight) {
		return fmt.Errorf("EffectBlockHeight be set %d", lockbls.EffectBlockHeight)
	}
	// Calculate withdrawal limit
	// Rlsnum = = extractable times
	rlsnum := (currentBlockHeight - uint64(lockbls.EffectBlockHeight)) / uint64(lockbls.LinearBlockNumber)
	if rlsnum == 0 {
		return fmt.Errorf("first release Block Height is %d, ", uint64(lockbls.EffectBlockHeight)+uint64(lockbls.LinearBlockNumber))
	}
	totalrlsamt := lockbls.TotalLockAmount
	steprlsamt := lockbls.LinearReleaseAmount
	// Effective withdrawable balance
	lockblsamt := lockbls.BalanceAmount
	// contrast
	if lockblsamt.LessThan(&act.ReleaseAmount) {
		return fmt.Errorf("BalanceAmount not enough.") // Sorry, your credit is running low
	}
	maxrlsamtbig := new(big.Int).Mul(steprlsamt.GetValue(), new(big.Int).SetUint64(rlsnum))
	currentMaxReleaseAmount, e3 := fields.NewAmountByBigInt(maxrlsamtbig)
	if e3 != nil {
		return e3
	}
	// The withdrawable balance shall be deducted from the withdrawn balance
	alreadyExtractedAmount, e9 := totalrlsamt.Sub(&lockblsamt) // 已经提走的余额
	if e9 != nil {
		return e9
	}
	// Effective withdrawable balance
	currentMaxReleaseAmount, e9 = currentMaxReleaseAmount.Sub(alreadyExtractedAmount)
	if e9 != nil {
		return e9
	}
	// Available balance judgment
	if currentMaxReleaseAmount.LessThan(&act.ReleaseAmount) {
		return fmt.Errorf("Current Max Release Amount not enough.") // The current available balance is insufficient
	}

	// Update lock in balance
	newBalanceAmount, e4 := lockblsamt.Sub(&act.ReleaseAmount)
	if e4 != nil {
		return e4
	}
	lockbls.BalanceAmount = *newBalanceAmount
	if newBalanceAmount.IsEmpty() {
		// All the locked warehouses have been extracted and deleted
		// Not deleted for rollback, but saved for block rollback
		e := state.LockblsUpdate(act.LockblsId, lockbls)
		if e != nil {
			return e
		}
	} else {
		// Deduct storage
		e := state.LockblsUpdate(act.LockblsId, lockbls)
		if e != nil {
			return e
		}
	}
	// Total supply statistics
	isbtcmoveunlock := act.LockblsId[0] == 0 // 第一位为 0 则是比特币转移的锁定
	if isbtcmoveunlock {
		totalsupply, e2 := state.ReadTotalSupply()
		if e2 != nil {
			return e2
		}
		// Cumulative unlocked HAC
		addamt := act.ReleaseAmount.ToMei()
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBitcoinTransferUnlockSuccessed, addamt)
		// update total supply
		e3 := state.UpdateSetTotalSupply(totalsupply)
		if e3 != nil {
			return e3
		}
	}
	// Plus balance
	return DoAddBalanceFromChainState(state, lockbls.MasterAddress, act.ReleaseAmount)
}

func (act *Action_10_LockblsRelease) RecoverChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	lockbls, e := state.Lockbls(act.LockblsId)
	if e != nil {
		return e
	}
	if lockbls == nil {
		return fmt.Errorf("Lockbls id<%s> not find.", hex.EncodeToString(act.LockblsId))
	}
	// Lock back
	// Update lock in balance
	lockblsamt := lockbls.BalanceAmount
	oldBalanceAmount, e4 := lockblsamt.Add(&act.ReleaseAmount)
	if e4 != nil {
		return e4
	}
	lockbls.BalanceAmount = *oldBalanceAmount
	// Deduct storage
	state.LockblsUpdate(act.LockblsId, lockbls)
	// Total supply statistics
	isbtcmoveunlock := act.LockblsId[0] == 0 // 第一位为 0 则是比特币转移的锁定
	if isbtcmoveunlock {
		totalsupply, e2 := state.ReadTotalSupply()
		if e2 != nil {
			return e2
		}
		// Cumulative unlocked HAC
		addamt := act.ReleaseAmount.ToMei()
		totalsupply.DoSub(stores.TotalSupplyStoreTypeOfBitcoinTransferUnlockSuccessed, addamt)
		// update total supply
		e3 := state.UpdateSetTotalSupply(totalsupply)
		if e3 != nil {
			return e3
		}
	}
	// Return balance
	return DoSubBalanceFromChainState(state, lockbls.MasterAddress, act.ReleaseAmount)
}

// Set belongs to long_ trs
func (act *Action_10_LockblsRelease) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}
func (act *Action_10_LockblsRelease) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_10_LockblsRelease) IsBurning90PersentTxFees() bool {
	return false
}
