package actions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
	"math"
	"math/big"
	"time"
)

type Action_7_SatoshiGenesis struct {
	TransferNo               fields.VarUint4 // 转账流水编号
	BitcoinBlockHeight       fields.VarUint4 // 转账的比特币区块高度
	BitcoinBlockTimestamp    fields.VarUint5 // 转账的比特币区块时间戳
	BitcoinEffectiveGenesis  fields.VarUint4 // 在这笔之前已经成功转移的比特币数量
	BitcoinQuantity          fields.VarUint4 // 本笔转账的比特币数量（单位：枚）
	AdditionalTotalHacAmount fields.VarUint4 // 本次转账[总共]应该增发的 hac 数量 （单位：枚）
	OriginAddress            fields.Address  // 转出的比特币来源地址
	BitcoinTransferHash      fields.Hash     // 比特币转账交易哈希

	// data ptr
	belong_trs interfaces.Transaction
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

func (act *Action_7_SatoshiGenesis) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	//if act.belong_trs != nil {
	//	return fmt.Errorf("Not yet.")
	//}

	// 交易只能包含唯一一个action
	belongactionnum := len(act.belong_trs.GetActions())
	if 1 != belongactionnum {
		return fmt.Errorf("Satoshi Genesis tx need only one action but got %d actions.", belongactionnum)
	}
	// 检查已经记录的增发（避免已完成的增发重复执行）
	belongtxhx, berr2 := state.ReadMoveBTCTxHashByNumber(uint32(act.TransferNo))
	if berr2 != nil {
		return berr2
	}
	if belongtxhx != nil {
		// 增发已经完成
		return fmt.Errorf("Satoshi act TransferNo<%d> has been executed.", act.TransferNo)
	}

	// 请求验证数据
	checkact, mustcheck := state.LoadValidatedSatoshiGenesis(int64(act.TransferNo))
	if mustcheck {
		// 交易位于交易池时 和 设置了check url时， 必须验证
		if checkact == nil {
			// URL 未返回数据
			return fmt.Errorf("SatoshiGenesis btc move logs url return invalid.")
		}
		// 比较
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
		// 验证数据 （转移比特币数量 1 ～ 1万枚）
		if act.BitcoinQuantity < 1 && act.BitcoinQuantity > 10000 {
			return fmt.Errorf("SatoshiGenesis act BitcoinQuantity number is error (right is 1 ~ 10000).")
		}
		var ttHac int64 = 0
		for i := act.BitcoinEffectiveGenesis + 1; i <= act.BitcoinEffectiveGenesis+act.BitcoinQuantity; i++ {
			ttHac += moveBtcCoinRewardByIdx(int64(i)) // 累加总共HAC奖励
		}
		if ttHac != int64(act.AdditionalTotalHacAmount) {
			// 增发的 HAC 数量不对
			return fmt.Errorf("SatoshiGenesis act AdditionalTotalHacAmount need %d but got %d.", ttHac, act.AdditionalTotalHacAmount)
		}
		// 检查时间（延迟28天才能领取）
		targettime := time.Unix(int64(act.BitcoinBlockTimestamp), 0).AddDate(0, 0, 28)
		if time.Now().Before(targettime) {
			return fmt.Errorf("SatoshiGenesis submit tx time must over %s", targettime.Format("2006/01/02 15:04:05"))
		}
		// 检查成功！！！
	}
	// 直接加到解锁的统计
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}

	// 统计比特币转移数量
	totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfTransferBitcoin, float64(act.BitcoinQuantity))

	// 记录 标记 已完成的 转移增发
	stoerr := state.SaveMoveBTCBelongTxHash(uint32(act.TransferNo), act.belong_trs.Hash())
	if stoerr != nil {
		return stoerr
	}

	// 增发 hac
	hacmeibig := (new(big.Int)).SetUint64(uint64(act.AdditionalTotalHacAmount))
	totaladdhacamt, err := fields.NewAmountByBigIntWithUnit(hacmeibig, 248)
	if err != nil {
		return err
	}
	// 锁仓时间按最先一枚计算
	// 判断是否线性锁仓至 lockbls
	lockweek, weekhei := moveBtcLockWeekByIdx(int64(act.BitcoinEffectiveGenesis) + 1)

	// 开发者模式
	if sys.TestDebugLocalDevelopmentMark {
		weekhei = 10 // 开发者模式，十个区块就释放
	}

	if weekhei > 17000000 {
		return fmt.Errorf("SatoshiGenesis moveBtcLockWeekByIdx weekhei overflow.")
	}
	if lockweek > 0 {

		// 线性锁仓（周）
		lkblsid := GainLockblsIdByBtcMove(uint32(act.TransferNo))

		// 存储
		lockbls := stores.NewEmptyLockbls(act.OriginAddress)
		lockbls.EffectBlockHeight = fields.VarUint5(state.GetPendingBlockHeight())
		lockbls.LinearBlockNumber = fields.VarUint3(weekhei) // 2000
		// amts
		allamtstorebytes := make([]*fields.Bytes8, 3)
		allamtstorebytes[0] = &lockbls.TotalLockAmountBytes
		allamtstorebytes[1] = &lockbls.BalanceAmountBytes
		allamtstorebytes[2] = &lockbls.LinearReleaseAmountBytes
		allamts := make([]*fields.Amount, 3)
		// 锁的总额
		allamts[0] = totaladdhacamt
		// 余额
		allamts[1] = totaladdhacamt
		// 每周可以解锁的币
		hacweekbig := (new(big.Int)).SetUint64(uint64(act.AdditionalTotalHacAmount) / uint64(lockweek))
		wklkhacamt, _ := fields.NewAmountByBigIntWithUnit(hacweekbig, 248)
		allamts[2] = wklkhacamt // 每周解锁
		// 赋值
		for i := 0; i < 3; i++ {
			ea := lockbls.PutAmount(allamtstorebytes[i], allamts[i])
			if ea != nil {
				return ea
			}
		}
		// stores
		// 创建线性锁仓
		state.LockblsCreate(lkblsid, lockbls)

	} else {

		// 不锁仓，直接打到余额
		e1 := DoAddBalanceFromChainState(state, act.OriginAddress, *totaladdhacamt)
		if e1 != nil {
			return e1
		}

		// 累加解锁的HAC
		addamt := totaladdhacamt.ToMei()
		totalsupply.DoAdd(stores.TotalSupplyStoreTypeOfBitcoinTransferUnlockSuccessed, addamt)

	}

	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}

	// 发行 btc 到地址
	satBTC := uint64(act.BitcoinQuantity) * 10000 * 10000 // 单位： 聪 (SAT)
	return DoAddSatoshiFromChainState(state, act.OriginAddress, fields.VarUint8(satBTC))
}

func (act *Action_7_SatoshiGenesis) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 回退解锁的统计
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}

	// 回退比特币转移数量
	totalsupply.DoSub(stores.TotalSupplyStoreTypeOfTransferBitcoin, float64(act.BitcoinQuantity))

	// 回退 hac
	hacmeibig := (new(big.Int)).SetUint64(uint64(act.AdditionalTotalHacAmount))
	// 锁仓时间按最先一枚计算
	// 判断是否线性锁仓至 lockbls
	lockweek, weekhei := moveBtcLockWeekByIdx(int64(act.BitcoinEffectiveGenesis) + 1)
	if weekhei > 17000000 {
		return fmt.Errorf("moveBtcLockWeekByIdx weekhei overflow.")
	}
	if lockweek > 0 {

		// 回退锁仓
		lkblsid := GainLockblsIdByBtcMove(uint32(act.TransferNo))
		// 删除线性锁仓
		state.LockblsDelete(lkblsid)

	} else {

		// 回退 HAC 增发
		addhacamt, err := fields.NewAmountByBigIntWithUnit(hacmeibig, 248)
		if err != nil {
			return err
		}
		e1 := DoSubBalanceFromChainState(state, act.OriginAddress, *addhacamt)
		if e1 != nil {
			return e1
		}
		// 减去解锁的HAC
		addamt := addhacamt.ToMei()
		totalsupply.DoSub(stores.TotalSupplyStoreTypeOfBitcoinTransferUnlockSuccessed, addamt)

	}

	// update total supply
	e3 := state.UpdateSetTotalSupply(totalsupply)
	if e3 != nil {
		return e3
	}

	// 扣除 btc
	satBTC := uint64(act.BitcoinQuantity) * 10000 * 10000 // 单位 聪
	return DoSubSatoshiFromChainState(state, act.OriginAddress, fields.VarUint8(satBTC))
}

// 设置所属 belong_trs
func (act *Action_7_SatoshiGenesis) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
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

// 第几枚BTC增发HAC数量（单位：枚）
func moveBtcCoinRewardByIdx(btcidx int64) int64 {
	var lvn = 21
	if btcidx == 1 {
		return powf2(lvn - 1)
	}
	if btcidx > powf2(lvn)-1 {
		return 1 // 最后始终增发一枚
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

// 计算第几枚BTC锁仓信息
func moveBtcLockWeekByIdx(btcidx int64) (int64, int64) {
	var oneweekhei int64 = 2000   // 2000 / 288 = 6.9444天
	var mostlockweek int64 = 1024 // 1024周约等于 20 年
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

	// 自己创建的 lockbls key 不允许创建这样的的前面全为0的key!!!
	lockbleid := bytes.Repeat([]byte{0}, 4) // key size = 18
	binary.BigEndian.PutUint32(lockbleid, uint32(btcTransferNo))
	lockbleidbytes := bytes.NewBuffer(bytes.Repeat([]byte{0}, stores.LockblsIdLength-4))
	lockbleidbytes.Write(lockbleid)
	lkblsid := lockbleidbytes.Bytes()
	return lkblsid
}
