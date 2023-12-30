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
	typeSizeValid int = 22 // Currently available
	// Diamonds
	TotalSupplyStoreTypeOfDiamond uint8 = 0 // Number of diamonds excavated
	// BTC
	TotalSupplyStoreTypeOfTransferBitcoin uint8 = 1 // Number of BTCs transferred successfully
	// Circulation quantity
	TotalSupplyStoreTypeOfBlockReward                    uint8 = 2 // Block reward HAC accumulation
	TotalSupplyStoreTypeOfChannelInterest                uint8 = 3 // Channel interest HAC accumulation
	TotalSupplyStoreTypeOfBitcoinTransferUnlockSuccessed uint8 = 4 // HAC accumulation successfully unlocked by bitcoin transfer and issuance
	// data statistics
	TotalSupplyStoreTypeOfLocatedHACInChannel uint8 = 5 // Number of HACs currently effectively locked in the channel
	TotalSupplyStoreTypeOfLocatedSATInChannel uint8 = 6 // Number of SATs currently effectively locked in the channel
	TotalSupplyStoreTypeOfChannelOfOpening    uint8 = 7 // Number of channels currently on
	// Destruction fee
	TotalSupplyStoreTypeOfBurningTotal uint8 = 8 // Handling charge combustion destruction HAC accumulation
	// Diamond lending
	TotalSupplyStoreTypeOfSystemLendingDiamondCurrentMortgageCount      uint8 = 9  // Real time statistics of loan and mortgage amount of diamond system
	TotalSupplyStoreTypeOfSystemLendingDiamondCumulationLoanHacAmount   uint8 = 10 // HAC daily flow quantity of diamond system mortgage cumulative lending
	TotalSupplyStoreTypeOfSystemLendingDiamondCumulationRansomHacAmount uint8 = 11 // Cumulative redemption (destruction) HAC flow of diamond system mortgage
	// Bitcoin lending (0.01 BTC per share)
	TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCurrentMortgageCount      uint8 = 12 // Real time statistics of loan mortgage shares in bitcoin system
	TotalSupplyStoreTypeOfSystemLendingBitcoinPortionBurningInterestHacAmount  uint8 = 13 // Statistics of pre destruction interest of debit and credit in bitcoin system
	TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationLoanHacAmount   uint8 = 14 // Cumulative lending HAC daily flow quantity of debit and credit in bitcoin system
	TotalSupplyStoreTypeOfSystemLendingBitcoinPortionCumulationRansomHacAmount uint8 = 15 // HAC daily flow quantity of debit and credit cumulative redemption in bitcoin system
	// Inter user credit
	TotalSupplyStoreTypeOfUsersLendingCumulationDiamond                  uint8 = 16 // Daily accumulation of inter user loan diamond quantity
	TotalSupplyStoreTypeOfUsersLendingCumulationBitcoin                  uint8 = 17 // Daily accumulation of inter user debit and credit bitcoin quantity (unit: piece)
	TotalSupplyStoreTypeOfUsersLendingCumulationHacAmount                uint8 = 18 // 用户间借贷HAC借出额流水累计（借出累计而非归还累计）
	TotalSupplyStoreTypeOfUsersLendingBurningOnePercentInterestHacAmount uint8 = 19 // 1% interest statistics of inter user loan system destruction
	TotalSupplyStoreTypeOfDiamondBidBurningZhu                           uint8 = 20 // Diamond bidding fee burning part (unit:zhu)
	TotalSupplyStoreTypeOfDiamondEngravedBurning                         uint8 = 21 // Diamond Engraved burning
	TotalSupplyStoreTypeOfDiamondEngravedOperateCount                    uint8 = 22 // Diamond Engraved to do count

	// TotalSupplyStoreTypeOfUsersLendingLendersInterestHacAmountCumulation uint8 = ... // 用户间借贷贷出方赚取的利息流水累计

)

type TotalSupply struct {
	changeMark []bool
	dataBytes  []uint64 //fields.Bytes8
}

func NewTotalSupplyStoreData() *TotalSupply {
	return &TotalSupply{
		changeMark: make([]bool, typeSizeMax),
		dataBytes:  make([]uint64, typeSizeMax), // fields.Bytes8
	}
}

func (t *TotalSupply) GetUint(ty uint8) uint64 {
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

func (t *TotalSupply) Get(ty uint8) float64 {
	var intn = t.GetUint(ty)
	//var intbts = make([]byte, 8)
	//binary.BigEndian.PutUint64(intbts, intn)
	return math.Float64frombits(intn)
}

// set up
func (t *TotalSupply) SetUint(ty uint8, value uint64) {
	if ty > uint8(typeSizeValid) {
		panic("type error")
	}
	// check mark
	t.changeMark[ty] = true
	// preservation
	t.dataBytes[ty] = value
}

func (t *TotalSupply) Set(ty uint8, value float64) {
	intbits := math.Float64bits(value)
	t.SetUint(ty, intbits)
}

// increase

func (t *TotalSupply) DoAddUint(ty uint8, value uint64) uint64 {
	if ty > uint8(typeSizeValid) {
		panic("type error")
	}
	// check mark
	t.changeMark[ty] = true
	vbase := t.GetUint(ty)
	newv := vbase + value
	t.SetUint(ty, newv)
	return newv
}

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

// reduce

func (t *TotalSupply) DoSubUint(ty uint8, value uint64) uint64 {
	if ty > uint8(typeSizeValid) {
		panic("type error")
	}
	// check mark
	t.changeMark[ty] = true
	vbase := t.GetUint(ty)
	newv := vbase - value
	t.SetUint(ty, newv)
	return newv
}

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

// Overwrite save
func (t *TotalSupply) CoverCopySave(src *TotalSupply) {
	for i := 0; i < typeSizeMax; i++ {
		if src.changeMark[i] {
			t.changeMark[i] = true
			t.dataBytes[i] = src.dataBytes[i]
		}
	}
}

// Copy replication
func (t *TotalSupply) Clone() *TotalSupply {
	changeMark := []bool{}
	dataBytes := []uint64{}
	changeMark = append(changeMark, t.changeMark...)
	dataBytes = append(dataBytes, t.dataBytes...)
	return &TotalSupply{
		changeMark: changeMark,
		dataBytes:  dataBytes,
	}
}

// serialize
func (t *TotalSupply) Serialize() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{uint8(typeSizeMax)}) // 长度
	for i := 0; i < typeSizeMax; i++ {
		//intbits := math.Float64bits(t.dataBytes[i])
		var btstore = fields.Bytes8{0, 0, 0, 0, 0, 0, 0, 0}
		binary.BigEndian.PutUint64(btstore, t.dataBytes[i])
		buf.Write(btstore)
	}
	return buf.Bytes(), nil
}

// Deserialization
func (t *TotalSupply) Parse(buf []byte, seek uint32) (uint32, error) {
	if int(seek)+1 > len(buf) {
		return 0, fmt.Errorf("TotalSupply Parse: buf too short")
	}
	tysize := int(buf[seek])
	t.changeMark = make([]bool, tysize)
	seek += 1
	for i := 0; i < tysize; i++ {
		t.changeMark[i] = true
		intbts := binary.BigEndian.Uint64(buf[seek : seek+8])
		t.dataBytes[i] = intbts
		seek += 8
	}
	return seek, nil
}
