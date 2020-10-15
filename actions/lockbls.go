package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"math/big"
)

type Action_9_LockblsCreate struct {
	LockblsId           fields.Bytes18  // 线性锁仓id
	PaymentAddress      fields.Address  // 付款地址
	MasterAddress       fields.Address  // 主地址（领取权）
	EffectBlockHeight   fields.VarUint5 // 生效（开始）区块
	LinearBlockNumber   fields.VarUint3 // 步进区块数 < 17000000 约 160年
	TotalStockAmount    fields.Amount   // 总共存入额度
	LinearReleaseAmount fields.Amount   // 每次释放额度

	// data ptr
	belong_trs interfaces.Transaction
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
	sk1, _ := elm.LockblsId.Parse(buf, seek)
	sk2, _ := elm.PaymentAddress.Parse(buf, sk1)
	sk3, _ := elm.MasterAddress.Parse(buf, sk2)
	sk4, _ := elm.EffectBlockHeight.Parse(buf, sk3)
	sk5, _ := elm.LinearBlockNumber.Parse(buf, sk4)
	sk6, _ := elm.TotalStockAmount.Parse(buf, sk5)
	sk7, _ := elm.LinearReleaseAmount.Parse(buf, sk6)
	return sk7, nil
}

func (act *Action_9_LockblsCreate) RequestSignAddresses() []fields.Address {
	return []fields.Address{
		act.PaymentAddress, // 主账户需要签名
	}
}

func (act *Action_9_LockblsCreate) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 检查key值合法性
	if act.LockblsId[0] == 0 || act.LockblsId[stores.LockblsIdLength-1] == 0 {
		return fmt.Errorf("LockblsId format error.")
	}
	// 检查是否key已经存在
	haslock := state.Lockbls(act.LockblsId)
	if haslock != nil {
		return fmt.Errorf("Lockbls id<%s> already.", hex.EncodeToString(act.LockblsId))
	}
	// 检查 步进 block number
	if act.LinearBlockNumber < 288 {
		return fmt.Errorf("LinearBlockNumber less 288.")
	}
	if act.LinearBlockNumber > 1700*10000 {
		return fmt.Errorf("LinearBlockNumber over 17000000.")
	}
	// 检查数额
	if !act.TotalStockAmount.IsPositive() || !act.LinearReleaseAmount.IsPositive() {
		return fmt.Errorf("TotalStockAmount or LinearReleaseAmount error.")
	}
	// 检查余额
	mainblsamt := state.Balance(act.PaymentAddress)
	if mainblsamt == nil {
		return fmt.Errorf("Balance cannot empty.")
	}
	if mainblsamt.Amount.LessThan(&act.TotalStockAmount) {
		return fmt.Errorf("Balance not enough.")
	}
	// 步进不能大于存入额
	if act.TotalStockAmount.LessThan(&act.LinearReleaseAmount) {
		return fmt.Errorf("LinearReleaseAmount cannot more than TotalStockAmount.")
	}

	// 存储
	lockbls := stores.NewEmptyLockbls(act.MasterAddress)
	lockbls.EffectBlockHeight = act.EffectBlockHeight
	lockbls.LinearBlockNumber = act.LinearBlockNumber
	// 保存
	allamtstorebytes := make([]*fields.Bytes8, 3)
	allamtstorebytes[0] = &lockbls.TotalLockAmountBytes     // 锁的总额
	allamtstorebytes[1] = &lockbls.LinearReleaseAmountBytes // 每次解锁
	allamtstorebytes[2] = &lockbls.BalanceAmountBytes       // 有效（未提取）余额
	allamts := make([]*fields.Amount, 3)
	allamts[0] = &act.TotalStockAmount    // 锁的总额
	allamts[1] = &act.LinearReleaseAmount // 每次解锁
	allamts[2] = &act.TotalStockAmount    // 有效（未提取）余额
	for i := 0; i < 3; i++ {
		e5 := lockbls.PutAmount(allamtstorebytes[i], allamts[i])
		if e5 != nil {
			return e5
		}
	}
	// 扣除 payment
	e1 := DoSubBalanceFromChainState(state, act.PaymentAddress, act.TotalStockAmount)
	if e1 != nil {
		return e1
	}

	// 保存锁仓
	e2 := state.LockblsCreate(act.LockblsId, lockbls)
	if e2 != nil {
		return e2
	}

	// ok
	return nil
}

func (act *Action_9_LockblsCreate) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 回退 hac
	e1 := DoAddBalanceFromChainState(state, act.PaymentAddress, act.TotalStockAmount)
	if e1 != nil {
		return e1
	}
	// 删除 lockbls
	e2 := state.LockblsDelete(act.LockblsId)
	if e2 != nil {
		return e2
	}

	// ok
	return nil
}

// 设置所属 belone_trs
func (act *Action_9_LockblsCreate) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}

///////////////////////////////////////////////////////////////////////////////////////////////

type Action_10_LockblsRelease struct {
	LockblsId     fields.Bytes18 // 线性锁仓id
	ReleaseAmount fields.Amount  // 本次提取额度

	// data ptr
	belong_trs interfaces.Transaction
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
	var moveseek, _ = elm.LockblsId.Parse(buf, seek)
	var moveseek2, _ = elm.ReleaseAmount.Parse(buf, moveseek)
	return moveseek2, nil
}

func (elm *Action_10_LockblsRelease) Size() uint32 {
	return 2 + elm.LockblsId.Size() + elm.ReleaseAmount.Size()
}

func (elm *Action_10_LockblsRelease) RequestSignAddresses() []fields.Address {
	return []fields.Address{}
}

func (act *Action_10_LockblsRelease) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 查询
	lockbls := state.Lockbls(act.LockblsId)
	if lockbls == nil {
		return fmt.Errorf("Lockbls id<%s> not find.", hex.EncodeToString(act.LockblsId))
	}
	// 提取出来
	currentBlockHeight := state.GetPendingBlockHeight()
	if currentBlockHeight < uint64(lockbls.EffectBlockHeight) {
		return fmt.Errorf("EffectBlockHeight be set %d", lockbls.EffectBlockHeight)
	}
	// 计算提取额度
	rlsnum := (currentBlockHeight - uint64(lockbls.EffectBlockHeight)) / uint64(lockbls.LinearBlockNumber)
	if rlsnum == 0 {
		return fmt.Errorf("first release Block Height is %d, ", uint64(lockbls.EffectBlockHeight)+uint64(lockbls.LinearBlockNumber))
	}
	totalrlsamt, e0 := lockbls.GetAmount(&lockbls.TotalLockAmountBytes)
	if e0 != nil {
		return e0
	}
	steprlsamt, e1 := lockbls.GetAmount(&lockbls.LinearReleaseAmountBytes)
	if e1 != nil {
		return e1
	}
	// 有效余额
	lockblsamt, e2 := lockbls.GetAmount(&lockbls.BalanceAmountBytes)
	if e2 != nil {
		return e2
	}
	if lockblsamt.LessThan(&act.ReleaseAmount) {
		return fmt.Errorf("BalanceAmount not enough.") // 余额不足
	}
	maxrlsamtbig := new(big.Int).Mul(steprlsamt.GetValue(), new(big.Int).SetUint64(rlsnum))
	currentMaxReleaseAmount, e3 := fields.NewAmountByBigInt(maxrlsamtbig)
	if e3 != nil {
		return e3
	}
	// 可提余额要减除掉已经提走的
	alreadyExtractedAmount, e9 := totalrlsamt.Sub(lockblsamt) // 已经提走的余额
	if e9 != nil {
		return e9
	}
	// 有效可提余额
	currentMaxReleaseAmount, e9 = currentMaxReleaseAmount.Sub(alreadyExtractedAmount)
	if e9 != nil {
		return e9
	}
	// 可提余额判断
	if currentMaxReleaseAmount.LessThan(&act.ReleaseAmount) {
		return fmt.Errorf("Current Max Release Amount not enough.") // 目前可提余额不足
	}

	// 更新锁仓余额
	newBalanceAmount, e4 := lockblsamt.Sub(&act.ReleaseAmount)
	if e4 != nil {
		return e4
	}
	e5 := lockbls.PutAmount(&lockbls.BalanceAmountBytes, newBalanceAmount)
	if e5 != nil {
		return e5
	}
	if newBalanceAmount.IsEmpty() {
		// 锁仓已经全部提取，删除
		// 为了回退暂不删除，而是为区块回退而暂时保存
		state.LockblsUpdate(act.LockblsId, lockbls)
	} else {
		// 扣除 储存
		state.LockblsUpdate(act.LockblsId, lockbls)
	}
	// 加上余额
	return DoAddBalanceFromChainState(state, lockbls.MasterAddress, act.ReleaseAmount)
}

func (act *Action_10_LockblsRelease) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	lockbls := state.Lockbls(act.LockblsId)
	if lockbls == nil {
		return fmt.Errorf("Lockbls id<%s> not find.", hex.EncodeToString(act.LockblsId))
	}
	// 锁仓回退
	// 更新锁仓余额
	lockblsamt, e2 := lockbls.GetAmount(&lockbls.BalanceAmountBytes)
	if e2 != nil {
		return e2
	}
	oldBalanceAmount, e4 := lockblsamt.Add(&act.ReleaseAmount)
	if e4 != nil {
		return e4
	}
	e5 := lockbls.PutAmount(&lockbls.BalanceAmountBytes, oldBalanceAmount)
	if e5 != nil {
		return e5
	}
	// 扣除 储存
	state.LockblsUpdate(act.LockblsId, lockbls)
	// 回退余额
	return DoSubBalanceFromChainState(state, lockbls.MasterAddress, act.ReleaseAmount)
}

// 设置所属 belone_trs
func (act *Action_10_LockblsRelease) SetBelongTransaction(trs interfaces.Transaction) {
	act.belong_trs = trs
}
