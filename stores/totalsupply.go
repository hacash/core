package stores

import (
	"bytes"
	"encoding/binary"
	"github.com/hacash/core/fields"
	"math"
)

const (
	typeSizeMax   int = 20
	typeSizeValid int = 6 // 当前可用的

	// 钻石
	TotalSupplyStoreTypeOfDiamond uint8 = 0 // 已挖掘出的钻石数量

	// 流通数量
	TotalSupplyStoreTypeOfBlockMinerReward               uint8 = 1 // 区块奖励HAC累计
	TotalSupplyStoreTypeOfChannelInterest                uint8 = 2 // 通道利息HAC累计
	TotalSupplyStoreTypeOfBitcoinTransferUnlockSuccessed uint8 = 3 // 比特币转移增发成功解锁的HAC累计

	// 数据统计
	TotalSupplyStoreTypeOfLocatedInChannel uint8 = 4 // 当前有效锁定在通道内的HAC数量

	// 销毁
	TotalSupplyStoreTypeOfBurningFee uint8 = 5 // 手续费燃烧销毁HAC累计
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
	if ty >= uint8(typeSizeValid) {
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
	if ty >= uint8(typeSizeValid) {
		panic("type error")
	}
	// check mark
	t.changeMark[ty] = true
	// 保存
	t.dataBytes[ty] = value
}

// 增加
func (t *TotalSupply) DoAdd(ty uint8, value float64) float64 {
	if ty >= uint8(typeSizeValid) {
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
	if ty >= uint8(typeSizeValid) {
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
