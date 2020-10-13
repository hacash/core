package actions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"math"
	"math/big"
	"time"
)

type Action_7_SatoshiGenesis struct {
	TransferNo               fields.VarUint4 // 转账流水编号
	BitcoinBlockHeight       fields.VarUint4 // 转账的比特币区块高度
	BitcoinBlockTimestamp    fields.VarUint4 // 转账的比特币区块时间戳
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
	sk1, _ := elm.TransferNo.Parse(buf, seek)
	sk2, _ := elm.BitcoinBlockHeight.Parse(buf, sk1)
	sk3, _ := elm.BitcoinBlockTimestamp.Parse(buf, sk2)
	sk4, _ := elm.BitcoinEffectiveGenesis.Parse(buf, sk3)
	sk5, _ := elm.BitcoinQuantity.Parse(buf, sk4)
	sk6, _ := elm.AdditionalTotalHacAmount.Parse(buf, sk5)
	sk7, _ := elm.OriginAddress.Parse(buf, sk6)
	sk8, _ := elm.BitcoinTransferHash.Parse(buf, sk7)
	return sk8, nil
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

	// 检查已经记录的增发
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
			return fmt.Errorf("Satoshi btc move logs url return invalid.")
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
			return fmt.Errorf("Satoshi act BitcoinQuantity number is error (right is 1 ~ 10000).")
		}
		var ttHac int64 = 0
		for i := act.BitcoinEffectiveGenesis + 1; i <= act.BitcoinEffectiveGenesis+act.BitcoinQuantity; i++ {
			ttHac += act.moveBtcCoinRewardByIdx(int64(i))
		}
		if ttHac != int64(act.AdditionalTotalHacAmount) {
			// 增发的 HAC 数量不对
			return fmt.Errorf("Satoshi act AdditionalTotalHacAmount need %d but goot %d.", ttHac, act.AdditionalTotalHacAmount)
		}
		// 检查时间（延迟28天才能领取）
		targettime := time.Unix(int64(act.BitcoinBlockTimestamp), 0).AddDate(0, 0, 28)
		if targettime.After(time.Now()) {
			return fmt.Errorf("SatoshiGenesis time must over %s", targettime.Format("2006/01/02 15:04:05"))
		}
		// 检查成功！！！
	}

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
	lockweek, weekhei := act.moveBtcLockWeekByIdx(int64(act.BitcoinEffectiveGenesis) + 1)
	if weekhei > 17000000 {
		return fmt.Errorf("moveBtcLockWeekByIdx weekhei overflow.")
	}
	if lockweek > 0 {

		// 线性锁仓（周）
		// 自己创建的 lockbls key 不允许创建这样的的前面全为0的key!!!
		lockbleid := bytes.Repeat([]byte{0}, 4) // key size = 18
		binary.BigEndian.PutUint32(lockbleid, uint32(act.TransferNo))
		lockbleidbytes := bytes.NewBuffer(bytes.Repeat([]byte{0}, stores.LockblsIdLength-4))
		lockbleidbytes.Write(lockbleid)
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
		// 每周解锁的币
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
		state.LockblsCreate(lockbleidbytes.Bytes(), lockbls)

	} else {

		// 不锁仓，直接打到余额
		e1 := DoAddBalanceFromChainState(state, act.OriginAddress, *totaladdhacamt)
		if e1 != nil {
			return e1
		}
	}

	// 发行 btc 到地址
	satBTC := uint64(act.BitcoinQuantity) * 10000 * 10000 // 单位： 聪 (SAT)
	return DoAddSatoshiFromChainState(state, act.OriginAddress, fields.VarUint8(satBTC))
}

func (act *Action_7_SatoshiGenesis) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 回退 hac
	hacmeibig := (new(big.Int)).SetUint64(uint64(act.AdditionalTotalHacAmount))
	addhacamt, err := fields.NewAmountByBigIntWithUnit(hacmeibig, 248)
	if err != nil {
		return err
	}
	e1 := DoSubBalanceFromChainState(state, act.OriginAddress, *addhacamt)
	if e1 != nil {
		return e1
	}
	// 扣除 btc
	satBTC := uint64(act.BitcoinQuantity) * 10000 * 10000 // 单位 聪
	return DoSubSatoshiFromChainState(state, act.OriginAddress, fields.VarUint8(satBTC))
}

// 设置所属 belone_trs
func (act *Action_7_SatoshiGenesis) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}

func (act Action_7_SatoshiGenesis) powf2(n int) int64 {
	res := math.Pow(2.0, float64(n))
	return int64(res)
}

// 第几枚BTC增发HAC数量（单位：枚）
func (act Action_7_SatoshiGenesis) moveBtcCoinRewardByIdx(btcidx int64) int64 {
	var lvn = 21
	if btcidx == 1 {
		return act.powf2(lvn - 1)
	}
	if btcidx > act.powf2(lvn)-1 {
		return 1 // 最后始终增发一枚
	}
	var tarlv int
	for i := 0; i < lvn; i++ {
		l := act.powf2(i) - 1
		r := act.powf2(i+1) - 1
		if btcidx > l && btcidx <= r {
			tarlv = i + 1
			break
		}
	}
	return act.powf2(lvn - tarlv)
}

// 第几枚BTC锁仓信息
func (act Action_7_SatoshiGenesis) moveBtcLockWeekByIdx(btcidx int64) (int64, int64) {
	var oneweekhei int64 = 2000   // 2000 / 288 = 6.9444天
	var mostlockweek int64 = 1024 // 1024周约等于 20 年
	if btcidx == 1 {
		return mostlockweek, oneweekhei
	}
	var lvn = 21
	var lockweek = mostlockweek
	for i := 0; i < lvn; i++ {
		l := act.powf2(i) - 1
		r := act.powf2(i+1) - 1
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
