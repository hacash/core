package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
	"github.com/hacash/x16rs"
	"strings"
)

/**
 * 钻石交易类型
 */

// 第 20001 个钻石开始，启用 32 位的 msg byte
const DiamondCreateCustomMessageAboveNumber uint32 = 20000

// 第 30001 个钻石开始，销毁 90% 的竞价费用
const DiamondCreateBurning90PercentTxFeesAboveNumber uint32 = 30000

// 采用 30001 ~ 40000 枚钻石平均竞价费用，之前的设定为 10 枚
const DiamondStatisticsAverageBiddingBurningPriceAboveNumber uint32 = 40000

// 第 40001 个钻石，开始用 sha3_hash(diamondreshash + blockhash) 决定钻石形状和配色
const DiamondResourceHashAndContainBlockHashDecideVisualGeneAboveNumber uint32 = 40000

// 第 41001 个钻石，开始用 sha3_hash(diamondreshash + blockhash + bidfee) 包括竞价费参与决定钻石形状配色
const DiamondResourceAppendBiddingFeeDecideVisualGeneAboveNumber uint32 = 41000

// 挖出钻石
type Action_4_DiamondCreate struct {
	Diamond  fields.DiamondName   // 钻石字面量 WTYUIAHXVMEKBSZN
	Number   fields.DiamondNumber // 钻石序号，用于难度检查
	PrevHash fields.Hash          // 上一个包含钻石的区块hash
	Nonce    fields.Bytes8        // 随机数
	Address  fields.Address       // 所属账户
	// 客户消息
	CustomMessage fields.Bytes32

	// 所属交易
	belong_trs interfaces.Transaction
}

func (elm *Action_4_DiamondCreate) Kind() uint16 {
	return 4
}

// json api
func (elm *Action_4_DiamondCreate) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_4_DiamondCreate) Size() uint32 {
	size := 2 +
		elm.Diamond.Size() +
		elm.Number.Size() +
		elm.PrevHash.Size() +
		elm.Nonce.Size() +
		elm.Address.Size()
	// 加上 msg byte
	if uint32(elm.Number) > DiamondCreateCustomMessageAboveNumber {
		size += elm.CustomMessage.Size()
	}
	return size
}

func (elm *Action_4_DiamondCreate) GetRealCustomMessage() []byte {
	if uint32(elm.Number) > DiamondCreateCustomMessageAboveNumber {
		var msgBytes, _ = elm.CustomMessage.Serialize()
		return msgBytes
	}
	return []byte{}
}

func (elm *Action_4_DiamondCreate) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var diamondBytes, _ = elm.Diamond.Serialize()
	var numberBytes, _ = elm.Number.Serialize()
	var prevBytes, _ = elm.PrevHash.Serialize()
	var nonceBytes, _ = elm.Nonce.Serialize()
	var addrBytes, _ = elm.Address.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(diamondBytes)
	buffer.Write(numberBytes)
	buffer.Write(prevBytes)
	buffer.Write(nonceBytes)
	buffer.Write(addrBytes)
	// 加上 msg byte
	if uint32(elm.Number) > DiamondCreateCustomMessageAboveNumber {
		var msgBytes, _ = elm.CustomMessage.Serialize()
		buffer.Write(msgBytes)
	}
	return buffer.Bytes(), nil
}

func (elm *Action_4_DiamondCreate) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	moveseek1, e := elm.Diamond.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	moveseek2, e := elm.Number.Parse(buf, moveseek1)
	if e != nil {
		return 0, e
	}
	moveseek3, e := elm.PrevHash.Parse(buf, moveseek2)
	if e != nil {
		return 0, e
	}
	moveseek4, e := elm.Nonce.Parse(buf, moveseek3)
	if e != nil {
		return 0, e
	}
	moveseek5, e := elm.Address.Parse(buf, moveseek4)
	if e != nil {
		return 0, e
	}
	// 加上 msg byte
	if uint32(elm.Number) > DiamondCreateCustomMessageAboveNumber {
		moveseek5, e = elm.CustomMessage.Parse(buf, moveseek5)
		if e != nil {
			return 0, e
		}
	}
	return moveseek5, nil
}

func (elm *Action_4_DiamondCreate) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // no sign
}

func (act *Action_4_DiamondCreate) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	//区块高度
	blkhei := state.GetPendingBlockHeight()
	blkhash := state.GetPendingBlockHash()
	diamondVisualUseContainBlockHash := blkhash
	if diamondVisualUseContainBlockHash == nil || len(diamondVisualUseContainBlockHash) != 32 {
		diamondVisualUseContainBlockHash = bytes.Repeat([]byte{0}, 32)
	}

	// 是否必须全面检查
	var mustDoAllCheck = true

	if sys.TestDebugLocalDevelopmentMark {
		mustDoAllCheck = false // 开发者模式 不检查
	}
	//fmt.Println(state.IsDatabaseVersionRebuildMode(), "-------------------------")
	if state.IsDatabaseVersionRebuildMode() {
		mustDoAllCheck = false // 数据库升级模式 不检查
	}

	// 计算钻石哈希
	diamondResHash, diamondStr := x16rs.Diamond(uint32(act.Number), act.PrevHash, act.Nonce, act.Address, act.GetRealCustomMessage())
	// 是否做全面的检查
	if mustDoAllCheck {
		// 交易只能包含唯一一个action
		belongactionnum := len(act.belong_trs.GetActions())
		if 1 != belongactionnum {
			return fmt.Errorf("Diamond create tx need only one action but got %d actions.", belongactionnum)
		}
		// 检查区块高度
		// 检查区块高度值是否为5的倍数
		// {BACKTOPOOL} 表示扔回交易池等待下个区块再次处理
		if blkhei%5 != 0 {
			return fmt.Errorf("{BACKTOPOOL} Diamond must be in block height multiple of 5.")
		}

		// 矿工状态检查
		lastdiamond, err := state.ReadLastestDiamond()
		if err != nil {
			return err
		}
		if lastdiamond != nil {
			//fmt.Println(lastdiamond.Diamond)
			//fmt.Println(lastdiamond.Number)
			//fmt.Println(lastdiamond.ContainBlockHash.ToHex())
			//fmt.Println(lastdiamond.PrevContainBlockHash.ToHex())
			prevdiamondnum, prevdiamondhash := uint32(lastdiamond.Number), lastdiamond.ContainBlockHash
			// 检查钻石是否是从上一个区块得来
			if act.PrevHash.Equal(prevdiamondhash) != true {
				return fmt.Errorf("Diamond prev hash must be <%s> but got <%s>.", hex.EncodeToString(prevdiamondhash), hex.EncodeToString(act.PrevHash))
			}
			if prevdiamondnum+1 != uint32(act.Number) {
				return fmt.Errorf("Diamond number must be <%d> but got <%d>.", prevdiamondnum+1, act.Number)
			}
		}
		// 检查钻石挖矿计算
		diamondstrval, isdia := x16rs.IsDiamondHashResultString(diamondStr)
		if !isdia {
			return fmt.Errorf("String <%s> is not diamond.", diamondStr)
		}
		if strings.Compare(diamondstrval, string(act.Diamond)) != 0 {
			return fmt.Errorf("Diamond need <%s> but got <%s>", act.Diamond, diamondstrval)
		}
		// 检查钻石难度值
		difok := x16rs.CheckDiamondDifficulty(uint32(act.Number), diamondResHash)
		if !difok {
			return fmt.Errorf("Diamond difficulty not meet the requirements.")
		}
		// 查询钻石是否已经存在
		hasaddr := state.Diamond(act.Diamond)
		if hasaddr != nil {
			return fmt.Errorf("Diamond <%s> already exist.", string(act.Diamond))
		}
		// 检查一个区块只能包含一枚钻石
		pendingdiamond, e2 := state.GetPendingSubmitStoreDiamond()
		if e2 != nil {
			return e2
		}
		if pendingdiamond != nil {
			return fmt.Errorf("This block height:%d has already exist diamond:<%s> .", blkhei, pendingdiamond.Diamond)
		}
		// 全部条件检查成功
	}

	// 存入钻石
	//fmt.Println(act.Address.ToReadable())
	var diastore = stores.NewDiamond(act.Address)
	diastore.Address = act.Address                // 钻石所属地址
	e3 := state.DiamondSet(act.Diamond, diastore) // 保存
	if e3 != nil {
		return e3
	}
	// 增加钻石余额 +1
	e9 := DoAddDiamondFromChainState(state, act.Address, 1)
	if e9 != nil {
		return e9
	}

	// 设置矿工状态
	// 标记本区块已经包含钻石
	// 存储对象，计算视觉基因
	visualGene, e15 := calculateVisualGeneByDiamondStuffHash(act.belong_trs, uint32(act.Number), diamondResHash, diamondStr, diamondVisualUseContainBlockHash)
	if e15 != nil {
		return e15
	}

	var diamondstore = &stores.DiamondSmelt{
		Diamond:              act.Diamond,
		Number:               act.Number,
		ContainBlockHeight:   fields.BlockHeight(blkhei),
		ContainBlockHash:     nil, // current block not exist !!!
		PrevContainBlockHash: act.PrevHash,
		MinerAddress:         act.Address,
		Nonce:                act.Nonce,
		CustomMessage:        act.GetRealCustomMessage(),
		VisualGene:           visualGene,
	}

	// 写入手续费报价
	feeoffer := act.belong_trs.GetFee()
	e11 := diamondstore.ParseApproxFeeOffer(feeoffer)
	if e11 != nil {
		return e11
	}
	// total supply 统计
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}

	// 计算平均竞价HAC枚数
	if uint32(act.Number) <= DiamondStatisticsAverageBiddingBurningPriceAboveNumber {
		diamondstore.AverageBidBurnPrice = 10 // 固定设为 10 枚
	} else {
		bsnum := uint32(act.Number) - DiamondCreateBurning90PercentTxFeesAboveNumber
		burnhac := totalsupply.Get(stores.TotalSupplyStoreTypeOfBurningFee)
		bidprice := uint64(burnhac/float64(bsnum) + 0.99999999) // 向上取整
		setprice := fields.VarUint2(bidprice)
		if setprice < 1 {
			setprice = 1 // 最小为1
		}
		diamondstore.AverageBidBurnPrice = setprice
	}

	// 设置最新状态
	e4 := state.SetLastestDiamond(diamondstore)
	if e4 != nil {
		return e4
	}
	e5 := state.SetPendingSubmitStoreDiamond(diamondstore)
	if e5 != nil {
		return e5
	}

	totalsupply.Set(stores.TotalSupplyStoreTypeOfDiamond, float64(act.Number))
	// update total supply
	e7 := state.UpdateSetTotalSupply(totalsupply)
	if e7 != nil {
		return e7
	}

	//fmt.Println("Action_4_DiamondCreate:", diamondstore.Number, string(diamondstore.Diamond), diamondstore.MinerAddress.ToReadable())
	//fmt.Print(string(diamondstore.Diamond)+",")

	//fmt.Println("Action_4_DiamondCreate:", act.Nonce)

	return nil
}

func (act *Action_4_DiamondCreate) RecoverChainState(state interfaces.ChainStateOperation) error {

	panic("RecoverChainState be deprecated")

	//chainstate := state.BlockStore()
	//if chainstate == nil {
	//	panic("Action get state.Miner() cannot be nil !")
	//}
	// 删除钻石
	e1 := state.DiamondDel(act.Diamond)
	if e1 != nil {
		return e1
	}
	// 回退矿工状态
	chainstore := state.BlockStore()
	if chainstore == nil {
		return fmt.Errorf("not find BlockStore object.")

	}
	prevDiamond, e2 := chainstore.ReadDiamondByNumber(uint32(act.Number) - 1)
	if e2 != nil {
		return e1
	}
	// setPrev
	e3 := state.SetLastestDiamond(prevDiamond)
	if e3 != nil {
		return e3
	}
	// 扣除钻石余额 -1
	e9 := DoSubDiamondFromChainState(state, act.Address, 1)
	if e9 != nil {
		return e9
	}
	// total supply 统计
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}
	totalsupply.Set(stores.TotalSupplyStoreTypeOfDiamond, float64(uint32(act.Number)-1))
	// update total supply
	e7 := state.UpdateSetTotalSupply(totalsupply)
	if e7 != nil {
		return e7
	}

	return nil
}

func (elm *Action_4_DiamondCreate) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_4_DiamondCreate) IsBurning90PersentTxFees() bool {
	if uint32(act.Number) > DiamondCreateBurning90PercentTxFeesAboveNumber {
		// 从第 30001 钻石开始，销毁本笔交易的 90% 的费用
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////

// 计算钻石的可视化基因
func calculateVisualGeneByDiamondStuffHash(belong_trs interfaces.Transaction, number uint32, stuffhx []byte, diamondstr string, peddingblkhash []byte) (fields.Bytes10, error) {
	if len(stuffhx) != 32 || len(peddingblkhash) != 32 {
		return nil, fmt.Errorf("stuffhx and peddingblkhash length must 32")
	}
	if len(diamondstr) != 16 {
		return nil, fmt.Errorf("diamondstr length must 16")
	}
	vgenehash := make([]byte, 32)
	copy(vgenehash, stuffhx)
	if number > DiamondResourceHashAndContainBlockHashDecideVisualGeneAboveNumber {
		// 第 40001 个钻石，开始用 sha3_hash(diamondreshash, blockhash) 决定钻石形状和配色
		vgenestuff := bytes.NewBuffer(stuffhx)
		vgenestuff.Write(peddingblkhash)
		if number > DiamondResourceAppendBiddingFeeDecideVisualGeneAboveNumber {
			bidfeebts, e := belong_trs.GetFee().Serialize() // 竞价手续费
			if e != nil {
				return nil, e // 返回错误
			}
			vgenestuff.Write(bidfeebts) // 竞价费参与决定钻石形状和配色
		}
		vgenehash = fields.CalculateHash(vgenestuff.Bytes()) // 开盲盒
		// 跟区块哈希一样是随机的，需要等待钻石确认的那一刻才能知晓形状和配色
		// fmt.Println(hex.EncodeToString(vgenestuff.Bytes()))
	}
	// fmt.Printf("Calculate Visual Gene #%d, vgenehash: %s, stuffhx: %s, peddingblkhash: %s\n", number, hex.EncodeToString(vgenehash), hex.EncodeToString(stuffhx), hex.EncodeToString(peddingblkhash))

	genehexstr := make([]string, 18)
	// 前6位
	k := 0
	for i := 10; i < 16; i++ {
		s := diamondstr[i]
		e := "0"
		switch s {
		case 'W': // WTYUIAHXVMEKBSZN
			e = "0"
		case 'T':
			e = "1"
		case 'Y':
			e = "2"
		case 'U':
			e = "3"
		case 'I':
			e = "4"
		case 'A':
			e = "5"
		case 'H':
			e = "6"
		case 'X':
			e = "7"
		case 'V':
			e = "8"
		case 'M':
			e = "9"
		case 'E':
			e = "A"
		case 'K':
			e = "B"
		case 'B':
			e = "C"
		case 'S':
			e = "D"
		case 'Z':
			e = "E"
		case 'N':
			e = "F"
		}
		genehexstr[k] = e
		k++
	}
	// 后11位
	for i := 20; i < 31; i++ {
		x := vgenehash[i]
		x = x % 16
		e := "0"
		switch x {
		case 0:
			e = "0"
		case 1:
			e = "1"
		case 2:
			e = "2"
		case 3:
			e = "3"
		case 4:
			e = "4"
		case 5:
			e = "5"
		case 6:
			e = "6"
		case 7:
			e = "7"
		case 8:
			e = "8"
		case 9:
			e = "9"
		case 10:
			e = "A"
		case 11:
			e = "B"
		case 12:
			e = "C"
		case 13:
			e = "D"
		case 14:
			e = "E"
		case 15:
			e = "F"
		}
		genehexstr[k] = e
		k++
	}
	// 补齐最后一位
	genehexstr[17] = "0"
	resbts, e1 := hex.DecodeString(strings.Join(genehexstr, ""))
	if e1 != nil {
		return nil, e1
	}
	// 哈希的最后一位作为形状选择
	resbuf := bytes.NewBuffer([]byte{vgenehash[31]})
	resbuf.Write(resbts) // 颜色选择器
	return resbuf.Bytes(), nil
}

///////////////////////////////////////////////////////////////

// 转移钻石
type Action_5_DiamondTransfer struct {
	Diamond   fields.DiamondName // 钻石字面量 WTYUIAHXVMEKBSZN
	ToAddress fields.Address     // 收钻方账户

	// 数据指针
	// 所属交易
	trs interfaces.Transaction
}

func (elm *Action_5_DiamondTransfer) Kind() uint16 {
	return 5
}

// json api
func (elm *Action_5_DiamondTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_5_DiamondTransfer) Size() uint32 {
	return 2 + elm.Diamond.Size() + elm.ToAddress.Size()
}

func (elm *Action_5_DiamondTransfer) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var diamondBytes, _ = elm.Diamond.Serialize()
	var addrBytes, _ = elm.ToAddress.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(diamondBytes)
	buffer.Write(addrBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_5_DiamondTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var moveseek1, _ = elm.Diamond.Parse(buf, seek)
	var moveseek2, _ = elm.ToAddress.Parse(buf, moveseek1)
	return moveseek2, nil
}

func (elm *Action_5_DiamondTransfer) RequestSignAddresses() []fields.Address {
	return []fields.Address{} // not sign
}

func (act *Action_5_DiamondTransfer) WriteinChainState(state interfaces.ChainStateOperation) error {

	if act.trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	trsMainAddress := act.trs.GetAddress()

	//fmt.Println("Action_5_DiamondTransfer:", trsMainAddress.ToReadable(), act.Address.ToReadable(), string(act.Diamond))

	// 自己不能转给自己
	if bytes.Compare(act.ToAddress, trsMainAddress) == 0 {
		return fmt.Errorf("Cannot transfer to self.")
	}
	// 查询钻石是否已经存在
	diaitem := state.Diamond(act.Diamond)
	if diaitem == nil {
		return fmt.Errorf("Diamond <%s> not exist.", string(act.Diamond))
	}
	item := diaitem
	// 检查是否抵押，是否可以转账
	if diaitem.Status != stores.DiamondStatusNormal {
		return fmt.Errorf("Diamond <%s> has been mortgaged and cannot be transferred.", string(act.Diamond))
	}
	// 检查所属
	if bytes.Compare(item.Address, trsMainAddress) != 0 {
		return fmt.Errorf("Diamond <%s> not belong to belong_trs address.", string(act.Diamond))
	}
	// 转移钻石
	item.Address = act.ToAddress
	err := state.DiamondSet(act.Diamond, item)
	if err != nil {
		return err
	}
	// 转移钻石余额
	e9 := DoSimpleDiamondTransferFromChainState(state, trsMainAddress, act.ToAddress, 1)
	if e9 != nil {
		return e9
	}
	return nil
}

func (act *Action_5_DiamondTransfer) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	trsMainAddress := act.trs.GetAddress()
	// get diamond
	diaitem := state.Diamond(act.Diamond)
	if diaitem == nil {
		return fmt.Errorf("Diamond <%s> not exist.", string(act.Diamond))
	}
	item := diaitem
	// 回退钻石
	item.Address = act.trs.GetAddress()
	err := state.DiamondSet(act.Diamond, item)
	if err != nil {
		return err
	}
	// 回退钻石余额
	e9 := DoSimpleDiamondTransferFromChainState(state, act.ToAddress, trsMainAddress, 1)
	if e9 != nil {
		return e9
	}
	return nil
}

func (elm *Action_5_DiamondTransfer) SetBelongTransaction(t interfaces.Transaction) {
	elm.trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_5_DiamondTransfer) IsBurning90PersentTxFees() bool {
	return false
}

///////////////////////////////////////////////////////////////

// 批量转移钻石
type Action_6_OutfeeQuantityDiamondTransfer struct {
	FromAddress fields.Address              // 拥有钻石的账户
	ToAddress   fields.Address              // 收钻方账户
	DiamondList fields.DiamondListMaxLen200 // 钻石列表

	// 数据指针
	// 所属交易
	trs interfaces.Transaction
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) Kind() uint16 {
	return 6
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) Size() uint32 {
	return 2 +
		elm.FromAddress.Size() +
		elm.ToAddress.Size() +
		elm.DiamondList.Size() // 每个钻石长6位
}

// json api
func (elm *Action_6_OutfeeQuantityDiamondTransfer) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var addr1Bytes, _ = elm.FromAddress.Serialize()
	var addr2Bytes, _ = elm.ToAddress.Serialize()
	var diaBytes, e = elm.DiamondList.Serialize()
	if e != nil {
		return nil, e
	}
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(addr1Bytes)
	buffer.Write(addr2Bytes)
	buffer.Write(diaBytes)
	return buffer.Bytes(), nil
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = elm.FromAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ToAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.DiamondList.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) RequestSignAddresses() []fields.Address {
	reqs := make([]fields.Address, 1) // 需from签名
	reqs[0] = elm.FromAddress
	return reqs
}

func (act *Action_6_OutfeeQuantityDiamondTransfer) WriteinChainState(state interfaces.ChainStateOperation) error {
	if act.trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 数量检查
	dianum := int(act.DiamondList.Count)
	if dianum == 0 || dianum != len(act.DiamondList.Diamonds) {
		return fmt.Errorf("Diamonds quantity error")
	}
	if dianum > 200 {
		return fmt.Errorf("Diamonds quantity cannot over 200")
	}
	// 自己不能转给自己
	if bytes.Compare(act.FromAddress, act.ToAddress) == 0 {
		return fmt.Errorf("Cannot transfer to self.")
	}
	// 批量转移钻石
	for i := 0; i < len(act.DiamondList.Diamonds); i++ {
		diamond := act.DiamondList.Diamonds[i]

		//fmt.Println("Action_6_OutfeeQuantityDiamondTransfer:", act.FromAddress.ToReadable(), act.ToAddress.ToReadable(), string(diamond))

		// fmt.Println("--- " + string(diamond))
		// 查询钻石是否已经存在
		diaitem := state.Diamond(diamond)
		if diaitem == nil {
			//panic("Quantity Diamond <%s> not exist. " + string(diamond))
			return fmt.Errorf("Quantity Diamond <%s> not exist.", string(diamond))
		}
		item := diaitem
		// 检查是否抵押，是否可以转账
		if diaitem.Status != stores.DiamondStatusNormal {
			return fmt.Errorf("Diamond <%s> has been mortgaged and cannot be transferred.", string(diamond))
		}
		// 检查所属
		if bytes.Compare(item.Address, act.FromAddress) != 0 {
			return fmt.Errorf("Diamond <%s> not belong to address '%s'", string(diamond), act.FromAddress.ToReadable())
		}
		// 转移钻石
		item.Address = act.ToAddress
		e5 := state.DiamondSet(diamond, item)
		if e5 != nil {
			return e5
		}
	}
	// 转移钻石余额
	e9 := DoSimpleDiamondTransferFromChainState(state, act.FromAddress, act.ToAddress, fields.DiamondNumber(dianum))
	if e9 != nil {
		return e9
	}
	return nil
}

func (act *Action_6_OutfeeQuantityDiamondTransfer) RecoverChainState(state interfaces.ChainStateOperation) error {
	if act.trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 批量回退钻石
	for i := 0; i < len(act.DiamondList.Diamonds); i++ {
		diamond := act.DiamondList.Diamonds[i]
		// get diamond
		diaitem := state.Diamond(diamond)
		if diaitem == nil {
			return fmt.Errorf("Diamond <%s> not exist.", string(diamond))
		}
		item := diaitem
		// 回退钻石
		item.Address = act.FromAddress
		err := state.DiamondSet(diamond, item)
		if err != nil {
			return err
		}
	}
	// 回退钻石余额
	e9 := DoSimpleDiamondTransferFromChainState(state, act.ToAddress, act.FromAddress, fields.DiamondNumber(act.DiamondList.Count))
	if e9 != nil {
		return e9
	}
	return nil
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) SetBelongTransaction(t interfaces.Transaction) {
	elm.trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_6_OutfeeQuantityDiamondTransfer) IsBurning90PersentTxFees() bool {
	return false
}

// 获取区块钻石的名称列表
func (elm *Action_6_OutfeeQuantityDiamondTransfer) GetDiamondNamesSplitByComma() string {
	return elm.DiamondList.SerializeHACDlistToCommaSplitString()
}

///////////////////////////////////////////////////////////////////////
