package stores

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/fields"
	"math"
)

const (
	typeSizeMax   int = 32
	typeSizeValid int = 19 // 当前可用的
	// 钻石
	TotalSupplyStoreTypeOfDiamond uint8 = 0 // 已挖掘出的钻石数量
	// BTC
	TotalSupplyStoreTypeOfTransferBitcoin uint8 = 1 // 已成功转移过来的 BTC 枚数
	// 流通数量
	TotalSupplyStoreTypeOfBlockReward                    uint8 = 2 // 区块奖励HAC累计
	TotalSupplyStoreTypeOfChannelInterest                uint8 = 3 // 通道利息HAC累计
	TotalSupplyStoreTypeOfBitcoinTransferUnlockSuccessed uint8 = 4 // 比特币转移增发成功解锁的HAC累计
	// 数据统计
	TotalSupplyStoreTypeOfLocatedHACInChannel uint8 = 5 // 当前有效锁定在通道内的HAC数量
	TotalSupplyStoreTypeOfLocatedSATInChannel uint8 = 6 // 当前有效锁定在通道内的HAC数量
	TotalSupplyStoreTypeOfChannelOfOpening    uint8 = 7 // 当前处于开启状态的通道数量
	// 销毁手续费
	TotalSupplyStoreTypeOfBurningFee uint8 = 8 // 手续费燃烧销毁HAC累计
	// 钻石借贷
	TotalSupplyStoreTypeOfSystemLendingDiamondCurrentMortgageCount      uint8 = 9  // 钻石系统借贷抵押数量实时统计
	TotalSupplyStoreTypeOfSystemLendingDiamondCumulationLoanHacAmount   uint8 = 10 // 钻石系统抵押累计借出HAC流水数量
	TotalSupplyStoreTypeOfSystemLendingDiamondCumulationRansomHacAmount uint8 = 11 // 钻石系统抵押累计赎回（销毁）HAC流水数量
	// 比特币借贷（每一份为 0.01 BTC）
	TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount      uint8 = 12 // 比特币系统借贷抵押份数实时统计
	TotalSupplyStoreTypeOfSystemLendingBitcoinPortionBurningInterestHacAmount  uint8 = 13 // 比特币系统借贷预销毁利息统计
	TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationLoanHacAmount   uint8 = 14 // 比特币系统借贷累计借出HAC流水数量
	TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationRansomHacAmount uint8 = 15 // 比特币系统借贷累计赎回HAC流水数量
	// 用户间借贷
	TotalSupplyStoreTypeOfUsersLendingCumulationDiamond                  uint8 = 16 // 用户间借贷钻石数量流水累计
	TotalSupplyStoreTypeOfUsersLendingCumulationBitcoin                  uint8 = 17 // 用户间借贷比特币数量流水累计（单位：枚）
	TotalSupplyStoreTypeOfUsersLendingCumulationHacAmount                uint8 = 18 // 用户间借贷HAC借出额流水累计（借出累计而非归还累计）
	TotalSupplyStoreTypeOfUsersLendingBurningOnePercentInterestHacAmount uint8 = 19 // 用户间借贷系统销毁的1%利息统计
	// TotalSupplyStoreTypeOfUsersLendingLendersInterestHacAmountCumulation uint8 = ... // 用户间借贷贷出方赚取的利息流水累计

)

type TotalSupply struct {
	changeMark []bool
	dataBytes  []float64 //fields.Bytes8
}

func NewTotalSupplyStoreData() *TotalSupply {
	return &TotalSupply{
		changeMark: make([]bool, typeSizeMax),
		dataBytes:  make([]float64, typeSizeMax), // fields.Bytes8
	}
}

func (t *TotalSupply) Get(ty uint8) float64 {
	if ty > uint8(typeSizeValid) {
		panic("type error")
	}
	// check mark
	if t.changeMark[ty] {
		return t.dataBytes[ty]
	}
	// none
	return 0
}

// 设置
func (t *TotalSupply) Set(ty uint8, value float64) {
	if ty > uint8(typeSizeValid) {
		panic("type error")
	}
	// check mark
	t.changeMark[ty] = true
	// 保存
	t.dataBytes[ty] = value
}

// 增加
func (t *TotalSupply) DoAdd(ty uint8, value float64) float64 {
	if ty > uint8(typeSizeValid) {
		panic("type error")
	}
	// check mark
	t.changeMark[ty] = true
	vbase := t.Get(ty)
	newv := vbase + value
	t.Set(ty, newv)
	return newv
}

// 减少
func (t *TotalSupply) DoSub(ty uint8, value float64) float64 {
	if ty > uint8(typeSizeValid) {
		panic("type error")
	}
	// check mark
	t.changeMark[ty] = true
	vbase := t.Get(ty)
	newv := vbase - value
	t.Set(ty, newv)
	return newv
}

// 覆盖保存
func (t *TotalSupply) CoverCopySave(src *TotalSupply) {
	for i := 0; i < typeSizeMax; i++ {
		if src.changeMark[i] {
			t.changeMark[i] = true
			t.dataBytes[i] = src.dataBytes[i]
		}
	}
}

// 拷贝复制
func (t *TotalSupply) Clone() *TotalSupply {
	changeMark := []bool{}
	dataBytes := []float64{}
	changeMark = append(changeMark, t.changeMark...)
	dataBytes = append(dataBytes, t.dataBytes...)
	return &TotalSupply{
		changeMark: changeMark,
		dataBytes:  dataBytes,
	}
}

// 序列化
func (t *TotalSupply) Serialize() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{uint8(typeSizeMax)}) // 长度
	for i := 0; i < typeSizeMax; i++ {
		intbits := math.Float64bits(t.dataBytes[i])
		var btstore = fields.Bytes8{0, 0, 0, 0, 0, 0, 0, 0}
		binary.BigEndian.PutUint64(btstore, intbits)
		buf.Write(btstore)
	}
	return buf.Bytes(), nil
}

// 反序列化
func (t *TotalSupply) Parse(buf []byte, seek uint32) (uint32, error) {
	if int(seek)+1 > len(buf) {
		return 0, fmt.Errorf("buf too short")
	}
	tysize := int(buf[seek])
	t.changeMark = make([]bool, tysize)
	seek += 1
	for i := 0; i < tysize; i++ {
		t.changeMark[i] = true
		intbts := binary.BigEndian.Uint64(buf[seek : seek+8])
		t.dataBytes[i] = math.Float64frombits(intbts)
		seek += 8
	}
	return seek, nil
}
