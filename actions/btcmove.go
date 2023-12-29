package actions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
	"math"
	"time"
)

type Action_7_SatoshiGenesis struct {
	TransferNo               fields.VarUint4         // Transfer serial number
	BitcoinBlockHeight       fields.VarUint4         // Height of bitcoin block transferred
	BitcoinBlockTimestamp    fields.BlockTxTimestamp // Bitcoin block timestamp of transfer
	BitcoinEffectiveGenesis  fields.VarUint4         // The number of bitcoins successfully transferred before this
	BitcoinQuantity          fields.VarUint4         // Number of bitcoins transferred in this transaction (unit: piece)
	AdditionalTotalHacAmount fields.VarUint4         // 本次转账[总共]应该增发的 hac 数量 （单位：枚）
	OriginAddress            fields.Address          // Bitcoin source address transferred out
	BitcoinTransferHash      fields.Hash             // Bitcoin transfer transaction hash

	// data ptr
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func NewAction_7_SatoshiGenesis() *Action_7_SatoshiGenesis {
	return &Action_7_SatoshiGenesis{}
}

func (elm *Action_7_SatoshiGenesis) Kind() uint16 {
	return 7
}

func (elm *Action_7_SatoshiGenesis) Size() uint32 {
	return 2 +
		elm.TransferNo.Size() +
		elm.BitcoinBlockHeight.Size() +
		elm.BitcoinBlockTimestamp.Size() +
		elm.BitcoinEffectiveGenesis.Size() +
		elm.BitcoinQuantity.Size() +
		elm.AdditionalTotalHacAmount.Size() +
		elm.OriginAddress.Size() +
		elm.BitcoinTransferHash.Size()
}

// json api
func (elm *Action_7_SatoshiGenesis) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_7_SatoshiGenesis) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var b1, _ = elm.TransferNo.Serialize()
	var b2, _ = elm.BitcoinBlockHeight.Serialize()
	var b3, _ = elm.BitcoinBlockTimestamp.Serialize()
	var b4, _ = elm.BitcoinEffectiveGenesis.Serialize()
	var b5, _ = elm.BitcoinQuantity.Serialize()
	var b6, _ = elm.AdditionalTotalHacAmount.Serialize()
	var b7, _ = elm.OriginAddress.Serialize()
	var b8, _ = elm.BitcoinTransferHash.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	buffer.Write(b7)
	buffer.Write(b8)
	return buffer.Bytes(), nil
}

func (elm *Action_7_SatoshiGenesis) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.TransferNo.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.BitcoinBlockHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.BitcoinBlockTimestamp.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.BitcoinEffectiveGenesis.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.BitcoinQuantity.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.AdditionalTotalHacAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.OriginAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.BitcoinTransferHash.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (*Action_7_SatoshiGenesis) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_7_SatoshiGenesis) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	if sys.TestDebugLocalDevelopmentMark == false || sys.TransactionSystemCheckChainID == 0 {

		// Transaction can only contain one and only one action
		belongactionnum := len(act.belong_trs_v3.GetActionList())
		if 1 != belongactionnum {
			return fmt.Errorf("Satoshi Genesis tx need only one action but got %d actions.", belongactionnum)
		}
		// Check the additional issuance that has been recorded (to avoid repeated execution of the completed additional issuance)
		belongtxhx, berr2 := state.ReadMoveBTCTxHashByTrsNo(uint32(act.TransferNo))
		if berr2 != nil {
			return berr2
		}
		if belongtxhx != nil {
			// Additional issuance has been completed
			return fmt.Errorf("Satoshi act TransferNo<%d> has been executed.", act.TransferNo)
		}

		// Request validation data
		checkact, mustcheck := state.BlockStore().LoadValidatedSatoshiGenesis(int64(act.TransferNo))
		if false == mustcheck {
			// Must be verified when adding to the transaction pool
			// Prevent verified BTC transfer transactions from being packaged into blocks
			mustcheck = state.IsInTxPool() == true
		}
		if mustcheck {
			// When the transaction is in the transaction pool and the check URL is set, it must be verified
			if checkact == nil {
				// URL did not return data
				return fmt.Errorf("SatoshiGenesis btc move logs url return invalid.")
			}
			// compare
			if act.TransferNo != checkact.TransferNo ||
				act.BitcoinBlockHeight != checkact.BitcoinBlockHeight ||
				act.BitcoinBlockTimestamp != checkact.BitcoinBlockTimestamp ||
				act.BitcoinEffectiveGenesis != checkact.BitcoinEffectiveGenesis ||
				act.BitcoinQuantity != checkact.BitcoinQuantity ||
				act.AdditionalTotalHacAmount != checkact.AdditionalTotalHacAmount ||
				bytes.Compare(act.OriginAddress, checkact.OriginAddress) != 0 ||
				bytes.Compare(act.BitcoinTransferHash, checkact.BitcoinTransferHash) != 0 {
				return fmt.Errorf("Action_7_SatoshiGenesis act and check act is mismatch.")
			}
			// Verification data (transfer 10000 ~ 10000 bitcoins)
			if act.BitcoinQuantity < 1 && act.BitcoinQuantity > 10000 {
				return fmt.Errorf("SatoshiGenesis act BitcoinQuantity number is error (right is 1 ~ 10000).")
			}
			var ttHac int64 = 0
			for i := act.BitcoinEffectiveGenesis + 1; i <= act.BitcoinEffectiveGenesis+act.BitcoinQuantity; i++ {
				ttHac += moveBtcCoinRewardByIdx(int64(i)) // Cumulative total HAC Awards
			}
			if ttHac != int64(act.AdditionalTotalHacAmount) {
				// Incorrect number of additional HACs
				return fmt.Errorf("SatoshiGenesis act AdditionalTotalHacAmount need %d but got %d.", ttHac, act.AdditionalTotalHacAmount)
			}
			// Inspection time (28 days later)
			targettime := time.Unix(int64(act.BitcoinBlockTimestamp), 0).AddDate(0, 0, 28)
			if time.Now().Before(targettime) {
				return fmt.Errorf("SatoshiGenesis submit tx time must over %s", targettime.Format("2006/01/02 15:04:05"))
			}
			// Check succeeded!!!
		}

	}

	// Statistics directly added to unlock
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}

	// Count the number of bitcoin transfers
	totalsupply.DoAddUint(stores.TotalSupplyStoreTypeOfTransferBitcoin, uint64(act.BitcoinQuantity))

	// Record the transfer and additional issuance marked as completed
	stoerr := state.SaveMoveBTCBelongTxHash(uint32(act.TransferNo), act.belong_trs_v3.Hash())
	if stoerr != nil {
		return stoerr
	}

	// Additional issuance of HAC
	totaladdhacamt := fields.NewAmountByUnit248(int64(act.AdditionalTotalHacAmount))
	// The lock time shall be calculated according to the first one
	// Judge whether to lock the warehouse to lockbls linearly
	lockweek, weekhei := moveBtcLockWeekByIdx(int64(act.BitcoinEffectiveGenesis) + 1)

	// Developer mode
	if sys.TestDebugLocalDevelopmentMark {
		weekhei = 10 // Developer mode, release ten blocks
	}

	if weekhei > 17000000 {
		return fmt.Errorf("SatoshiGenesis moveBtcLockWeekByIdx weekhei overflow.")
	}
	if lockweek > 0 {

		// Linear lock (week)
		lkblsid := GainLockblsIdByBtcMove(uint32(act.TransferNo))

		// storage
		lockbls := stores.NewEmptyLockbls(act.OriginAddress)
		lockbls.EffectBlockHeight = fields.BlockHeight(state.GetPendingBlockHeight())
		lockbls.LinearBlockNumber = fields.VarUint3(weekhei) // 2000
		lockbls.TotalLockAmount = *totaladdhacamt            // General lock
		lockbls.BalanceAmount = *totaladdhacamt              // balance
		wklkhacamt := fields.NewAmountByUnit248(int64(act.AdditionalTotalHacAmount) / int64(lockweek))
		lockbls.LinearReleaseAmount = *wklkhacamt // Coins that can be unlocked every week
		// stores
		// Create linear lock
		e := state.LockblsCreate(lkblsid, lockbls)
		if e != nil {
			return e
		}

	} else {

		// Do not lock the position, and directly transfer to the balance
		e1 := DoAddBalanceFromChainStateV3(state, act.OriginAddress, *totaladdhacamt)
		if e1 != nil {
			return e1
		}

		// Cumulative unlocked HAC
		addamt := totaladdhacamt.ToMei()
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBitcoinTransferUnlockSuccessed, addamt)

	}

	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}

	// Issue BTC to address
	satBTC := uint64(act.BitcoinQuantity) * 10000 * 10000 // 单位： 聪 (SAT)
	return DoAddSatoshiFromChainStateV3(state, act.OriginAddress, fields.Satoshi(satBTC))
}

func (act *Action_7_SatoshiGenesis) WriteinChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Transaction can only contain one and only one action
	belongactionnum := len(act.belong_trs.GetActions())
	if 1 != belongactionnum {
		return fmt.Errorf("Satoshi Genesis tx need only one action but got %d actions.", belongactionnum)
	}
	// Check the additional issuance that has been recorded (to avoid repeated execution of the completed additional issuance)
	belongtxhx, berr2 := state.ReadMoveBTCTxHashByNumber(uint32(act.TransferNo))
	if berr2 != nil {
		return berr2
	}
	if belongtxhx != nil {
		// Additional issuance has been completed
		return fmt.Errorf("Satoshi act TransferNo<%d> has been executed.", act.TransferNo)
	}

	// Request validation data
	checkact, mustcheck := state.LoadValidatedSatoshiGenesis(int64(act.TransferNo))
	if mustcheck {
		// When the transaction is in the transaction pool and the check URL is set, it must be verified
		if checkact == nil {
			// URL did not return data
			return fmt.Errorf("SatoshiGenesis btc move logs url return invalid.")
		}
		// compare
		if act.TransferNo != checkact.TransferNo ||
			act.BitcoinBlockHeight != checkact.BitcoinBlockHeight ||
			act.BitcoinBlockTimestamp != checkact.BitcoinBlockTimestamp ||
			act.BitcoinEffectiveGenesis != checkact.BitcoinEffectiveGenesis ||
			act.BitcoinQuantity != checkact.BitcoinQuantity ||
			act.AdditionalTotalHacAmount != checkact.AdditionalTotalHacAmount ||
			bytes.Compare(act.OriginAddress, checkact.OriginAddress) != 0 ||
			bytes.Compare(act.BitcoinTransferHash, checkact.BitcoinTransferHash) != 0 {
			return fmt.Errorf("Action_7_SatoshiGenesis act and check act is mismatch.")
		}
		// Verification data (transfer 10000 ~ 10000 bitcoins)
		if act.BitcoinQuantity < 1 && act.BitcoinQuantity > 10000 {
			return fmt.Errorf("SatoshiGenesis act BitcoinQuantity number is error (right is 1 ~ 10000).")
		}
		var ttHac int64 = 0
		for i := act.BitcoinEffectiveGenesis + 1; i <= act.BitcoinEffectiveGenesis+act.BitcoinQuantity; i++ {
			ttHac += moveBtcCoinRewardByIdx(int64(i)) // Cumulative total HAC Awards
		}
		if ttHac != int64(act.AdditionalTotalHacAmount) {
			// Incorrect number of additional HACs
			return fmt.Errorf("SatoshiGenesis act AdditionalTotalHacAmount need %d but got %d.", ttHac, act.AdditionalTotalHacAmount)
		}
		// Inspection time (28 days later)
		targettime := time.Unix(int64(act.BitcoinBlockTimestamp), 0).AddDate(0, 0, 28)
		if time.Now().Before(targettime) {
			return fmt.Errorf("SatoshiGenesis submit tx time must over %s", targettime.Format("2006/01/02 15:04:05"))
		}
		// Check succeeded!!!
	}
	// Statistics directly added to unlock
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}

	// Count the number of bitcoin transfers
	totalsupply.DoAddUint(stores.TotalSupplyStoreTypeOfTransferBitcoin, uint64(act.BitcoinQuantity))

	// Record the transfer and additional issuance marked as completed
	stoerr := state.SaveMoveBTCBelongTxHash(uint32(act.TransferNo), act.belong_trs.Hash())
	if stoerr != nil {
		return stoerr
	}

	// Additional issuance of HAC
	/*
		hacmeibig := (new(big.Int)).SetUint64(uint64(act.AdditionalTotalHacAmount))
		totaladdhacamt, err := fields.NewAmountByBigIntWithUnit(hacmeibig, 248)
		if err != nil {
			return err
		}
	*/
	totaladdhacamt := fields.NewAmountByUnit248(int64(act.AdditionalTotalHacAmount))
	// The lock time shall be calculated according to the first one
	// Judge whether to lock the warehouse to lockbls linearly
	lockweek, weekhei := moveBtcLockWeekByIdx(int64(act.BitcoinEffectiveGenesis) + 1)

	// Developer mode
	if sys.TestDebugLocalDevelopmentMark {
		weekhei = 10 // Developer mode, release ten blocks
	}

	if weekhei > 17000000 {
		return fmt.Errorf("SatoshiGenesis moveBtcLockWeekByIdx weekhei overflow.")
	}
	if lockweek > 0 {

		// Linear lock (week)
		lkblsid := GainLockblsIdByBtcMove(uint32(act.TransferNo))

		// storage
		lockbls := stores.NewEmptyLockbls(act.OriginAddress)
		lockbls.EffectBlockHeight = fields.BlockHeight(state.GetPendingBlockHeight())
		lockbls.LinearBlockNumber = fields.VarUint3(weekhei) // 2000
		lockbls.TotalLockAmount = *totaladdhacamt            // General lock
		lockbls.BalanceAmount = *totaladdhacamt              // balance
		wklkhacamt := fields.NewAmountByUnit248(int64(act.AdditionalTotalHacAmount) / int64(lockweek))
		lockbls.LinearReleaseAmount = *wklkhacamt // Coins that can be unlocked every week
		// stores
		// Create linear lock
		state.LockblsCreate(lkblsid, lockbls)

	} else {

		// Do not lock the position, and directly transfer to the balance
		e1 := DoAddBalanceFromChainState(state, act.OriginAddress, *totaladdhacamt)
		if e1 != nil {
			return e1
		}

		// Cumulative unlocked HAC
		addamt := totaladdhacamt.ToMei()
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBitcoinTransferUnlockSuccessed, addamt)

	}

	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}

	// Issue BTC to address
	satBTC := uint64(act.BitcoinQuantity) * 10000 * 10000 // 单位： 聪 (SAT)
	return DoAddSatoshiFromChainState(state, act.OriginAddress, fields.Satoshi(satBTC))
}

func (act *Action_7_SatoshiGenesis) RecoverChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// Statistics of rollback unlocking
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}

	// Back bitcoin transfer quantity
	totalsupply.DoSubUint(stores.TotalSupplyStoreTypeOfTransferBitcoin, uint64(act.BitcoinQuantity))

	// Fallback HAC
	// The lock time shall be calculated according to the first one
	// Judge whether to lock the warehouse to lockbls linearly
	lockweek, weekhei := moveBtcLockWeekByIdx(int64(act.BitcoinEffectiveGenesis) + 1)
	if weekhei > 17000000 {
		return fmt.Errorf("moveBtcLockWeekByIdx weekhei overflow.")
	}
	if lockweek > 0 {

		// Back locking
		lkblsid := GainLockblsIdByBtcMove(uint32(act.TransferNo))
		// Delete linear lock
		state.LockblsDelete(lkblsid)

	} else {

		// Back HAC additional issue
		/*
			addhacamt, err := fields.NewAmountByBigIntWithUnit(hacmeibig, 248)
			if err != nil {
				return err
			}
		*/
		addhacamt := fields.NewAmountByUnit248(int64(act.AdditionalTotalHacAmount))
		e1 := DoSubBalanceFromChainState(state, act.OriginAddress, *addhacamt)
		if e1 != nil {
			return e1
		}
		// Subtract unlocked HAC
		addamt := addhacamt.ToMei()
		totalsupply.DoSub(stores.TotalSupplyStoreTypeOfBitcoinTransferUnlockSuccessed, addamt)

	}

	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}

	// Deduct BTC
	satBTC := uint64(act.BitcoinQuantity) * 10000 * 10000 // 单位 聪
	return DoSubSatoshiFromChainState(state, act.OriginAddress, fields.Satoshi(satBTC))
}

// Set belongs to long_ trs
func (act *Action_7_SatoshiGenesis) SetBelongTransaction(trs interfacev2.Transaction) {
	act.belong_trs = trs
}
func (act *Action_7_SatoshiGenesis) SetBelongTrs(trs interfaces.Transaction) {
	act.belong_trs_v3 = trs
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_7_SatoshiGenesis) IsBurning90PersentTxFees() bool {
	return false
}

//////////////////////////////

func powf2(n int) int64 {
	res := math.Pow(2.0, float64(n))
	return int64(res)
}

// The number of additional HACs issued by the BTC (unit: PCS)
func moveBtcCoinRewardByIdx(btcidx int64) int64 {
	var lvn = 21
	if btcidx == 1 {
		return powf2(lvn - 1)
	}
	if btcidx > powf2(lvn)-1 {
		return 1 // Finally, always issue an additional one
	}
	var tarlv int
	for i := 0; i < lvn; i++ {
		l := powf2(i) - 1
		r := powf2(i+1) - 1
		if btcidx > l && btcidx <= r {
			tarlv = i + 1
			break
		}
	}
	return powf2(lvn - tarlv)
}

// Calculate the lock information of the BTC
func moveBtcLockWeekByIdx(btcidx int64) (int64, int64) {
	var oneweekhei int64 = 2000   // 2000 / 288 = 6.9444天
	var mostlockweek int64 = 1024 // 1024 weeks is about 20 years
	if btcidx == 1 {
		return mostlockweek, oneweekhei
	}
	var lvn = 21
	var lockweek = mostlockweek
	for i := 0; i < lvn; i++ {
		l := powf2(i) - 1
		r := powf2(i+1) - 1
		if btcidx > l && btcidx <= r {
			break
		}
		lockweek /= 2
		if lockweek == 0 {
			return 0, oneweekhei
		}
	}
	return lockweek, oneweekhei
}

///////////////////////////

func GainLockblsIdByBtcMove(btcTransferNo uint32) []byte {

	// Self created lockbls keys are not allowed to create such keys with 0 in front!!!
	lockbleid := bytes.Repeat([]byte{0}, 4) // key size = 18
	binary.BigEndian.PutUint32(lockbleid, uint32(btcTransferNo))
	lockbleidbytes := bytes.NewBuffer(bytes.Repeat([]byte{0}, stores.LockblsIdLength-4))
	lockbleidbytes.Write(lockbleid)
	lkblsid := lockbleidbytes.Bytes()
	return lkblsid
}
