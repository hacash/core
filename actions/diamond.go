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
	"github.com/hacash/core/sys"
	"github.com/hacash/x16rs"
	"strings"
)

/**
 * diamond tx type
 */

// Start with the 20001st diamond and enable the 32-bit MSG byte
const DiamondCreateCustomMessageAboveNumber uint32 = 20000

// Starting from the 30001st diamond, destroy 90% of the bidding cost
const DiamondCreateBurning90PercentTxFeesAboveNumber uint32 = 30000

// The average bidding cost of 30001 ~ 40000 diamonds is adopted, and the previous setting is 10 diamonds
const DiamondStatisticsAverageBiddingBurningPriceAboveNumber uint32 = 40000

// 40001 diamond, start with Sha3_ Hash (diamondreshash + blockhash) determines diamond shape and color matching
const DiamondResourceHashAndContainBlockHashDecideVisualGeneAboveNumber uint32 = 40000

// 41001 diamond, start with Sha3_ Hash (diamondreshash + blockhash + bidfee) includes the bidding fee to participate in the decision of diamond shape color matching
const DiamondResourceAppendBiddingFeeDecideVisualGeneAboveNumber uint32 = 41000

// Dig out diamonds
type Action_4_DiamondCreate struct {
	Diamond  fields.DiamondName   // Diamond literal wtyuiahxvmekbszn
	Number   fields.DiamondNumber // Diamond serial number for difficulty check
	PrevHash fields.Hash          // Previous block hash containing diamond
	Nonce    fields.Bytes8        // random number
	Address  fields.Address       // Account
	// Customer message
	CustomMessage fields.Bytes32

	// Transaction
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
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
	// Add MSG byte
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
	// Add MSG byte
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
	// Add MSG byte
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

func (act *Action_4_DiamondCreate) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	blockstore := state.BlockStore()

	//区块高度
	pending := state.GetPending()
	blkhei := state.GetPendingBlockHeight()
	blkhash := state.GetPendingBlockHash()
	diamondVisualUseContainBlockHash := blkhash
	if diamondVisualUseContainBlockHash == nil || len(diamondVisualUseContainBlockHash) != 32 {
		diamondVisualUseContainBlockHash = bytes.Repeat([]byte{0}, 32)
	}

	// Whether comprehensive inspection is necessary
	var mustDoAllCheck = true

	lateststatus, e := state.LatestStatusRead()
	if e != nil {
		return e
	}

	if sys.TestDebugLocalDevelopmentMark {
		mustDoAllCheck = false // Developer mode does not check
	}
	//fmt.Println(state.IsDatabaseVersionRebuildMode(), "-------------------------")
	if state.IsDatabaseVersionRebuildMode() {
		mustDoAllCheck = false // Database upgrade mode does not check
	}

	// Calculate diamond hash
	sha3hash, diamondResHash, diamondStr := x16rs.Diamond(uint32(act.Number), act.PrevHash, act.Nonce, act.Address, act.GetRealCustomMessage())
	// Whether to conduct comprehensive inspection
	if mustDoAllCheck {
		// Transaction can only contain one and only one action
		belongactionnum := len(act.belong_trs_v3.GetActionList())
		if 1 != belongactionnum {
			return fmt.Errorf("Diamond create tx need only one action but got %d actions.", belongactionnum)
		}
		// Check block height
		// Check whether the block height value is a multiple of 5
		// {backtopool} means to throw it back to the trading pool and wait for the next block to be processed again
		if blkhei%5 != 0 {
			return fmt.Errorf("{BACKTOPOOL} Diamond must be in block height multiple of 5.")
		}

		// Miner status inspection
		lastdiamond := lateststatus.ReadLastestDiamond()
		if lastdiamond == nil && act.Number > 1 {
			// Not the first diamond but can not read the information of the last diamond
			return fmt.Errorf("Cannot ReadLastestDiamond() to get last diamond of act.Number > 1.")
		}
		if lastdiamond != nil {
			//fmt.Println(lastdiamond.Diamond)
			//fmt.Println(lastdiamond.Number)
			//fmt.Println(lastdiamond.ContainBlockHash.ToHex())
			//fmt.Println(lastdiamond.PrevContainBlockHash.ToHex())
			prevdiamondnum, prevdiamondhash := uint32(lastdiamond.Number), lastdiamond.ContainBlockHash
			if prevdiamondhash == nil {
				return fmt.Errorf("lastdiamond.ContainBlockHash is nil.")
			}
			// Check if the diamond is from the previous block
			if act.PrevHash.Equal(prevdiamondhash) != true {
				return fmt.Errorf("Diamond prev hash must be <%s> but got <%s>.", hex.EncodeToString(prevdiamondhash), hex.EncodeToString(act.PrevHash))
			}
			if prevdiamondnum+1 != uint32(act.Number) {
				return fmt.Errorf("Diamond number must be <%d> but got <%d>.", prevdiamondnum+1, act.Number)
			}
		}
		// Check diamond mining calculation
		diamondstrval, isdia := x16rs.IsDiamondHashResultString(diamondStr)
		if !isdia {
			return fmt.Errorf("String <%s> is not diamond.", diamondStr)
		}
		if strings.Compare(diamondstrval, string(act.Diamond)) != 0 {
			return fmt.Errorf("Diamond need <%s> but got <%s>", act.Diamond, diamondstrval)
		}
		// Check diamond difficulty value
		difok := x16rs.CheckDiamondDifficulty(uint32(act.Number), sha3hash, diamondResHash)
		if !difok {
			return fmt.Errorf("block %d Diamond difficulty not meet the requirements.", blkhei)
		}
		// Query whether the diamond already exists
		hasaddr, e := state.Diamond(act.Diamond)
		if e != nil {
			return e
		}
		if hasaddr != nil {
			return fmt.Errorf("Diamond <%s> already exist.", string(act.Diamond))
		}
		// Check that a block can contain only one diamond
		pendingdiamond := pending.GetWaitingSubmitDiamond()
		if pendingdiamond != nil {
			return fmt.Errorf("This block height:%d has already exist diamond:<%s> .", blkhei, pendingdiamond.Diamond)
		}
		// All conditions checked successfully
	}

	// Deposit diamonds
	//fmt.Println(act.Address.ToReadable())
	var diastore = stores.NewDiamond(act.Address)
	diastore.Address = act.Address // Diamond address
	e3 := state.DiamondSet(act.Diamond, diastore)
	if e3 != nil {
		return e3
	}
	// Increase diamond balance +1
	e9 := DoAddDiamondFromChainStateV3(state, act.Address, 1)
	if e9 != nil {
		return e9
	}

	// Set miner status
	// Mark that this block already contains diamonds
	// Storing objects, computing visual genes
	lifeGene, e15 := calculateLifeGeneByDiamondStuffHashV3(act.belong_trs_v3, uint32(act.Number), diamondResHash, diamondStr, diamondVisualUseContainBlockHash)
	if e15 != nil {
		return e15
	}

	var diamondstore = &stores.DiamondSmelt{
		Diamond:              act.Diamond,
		Number:               act.Number,
		ContainBlockHeight:   fields.BlockHeight(blkhei),
		ContainBlockHash:     state.GetPendingBlockHash(),
		PrevContainBlockHash: act.PrevHash,
		MinerAddress:         act.Address,
		Nonce:                act.Nonce,
		CustomMessage:        act.GetRealCustomMessage(),
		LifeGene:             lifeGene,
	}

	// Write service charge quotation
	feeoffer := act.belong_trs_v3.GetFee()
	e11 := diamondstore.ParseApproxFeeOffer(feeoffer)
	if e11 != nil {
		return e11
	}
	// Total supply statistics
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}

	// Calculate the average number of HACs for bidding
	if uint32(act.Number) <= DiamondStatisticsAverageBiddingBurningPriceAboveNumber {
		diamondstore.AverageBidBurnPrice = 10 // Fixed to 10
	} else {
		bsnum := uint32(act.Number) - DiamondCreateBurning90PercentTxFeesAboveNumber
		burnhac := totalsupply.Get(stores.TotalSupplyStoreTypeOfBurningFee)
		bidprice := uint64(burnhac/float64(bsnum) + 0.99999999) // up 1
		setprice := fields.VarUint2(bidprice)
		if setprice < 1 {
			setprice = 1 // Minimum 1
		}
		diamondstore.AverageBidBurnPrice = setprice
	}

	// Update block status
	lateststatus.SetLastestDiamond(diamondstore)
	pending.SetWaitingSubmitDiamond(diamondstore)
	// Save status
	e = state.LatestStatusSet(lateststatus)
	if e != nil {
		return e
	}

	// Save diamonds
	if state.IsInTxPool() == false {
		e = blockstore.SaveDiamond(diamondstore)
		if e != nil {
			return e
		}
		e = blockstore.UpdateSetDiamondNameReferToNumber(uint32(diamondstore.Number), diamondstore.Diamond)
		if e != nil {
			return e
		}
	}

	totalsupply.Set(stores.TotalSupplyStoreTypeOfDiamond, float64(act.Number))
	// update total supply
	e7 := state.UpdateSetTotalSupply(totalsupply)
	if e7 != nil {
		return e7
	}

	// update
	//e = state.LatestStatusSet(lateststatus)
	//if e != nil {
	//	return e
	//}
	e = state.SetPending(pending)
	if e != nil {
		return e
	}

	//fmt.Println("Action_4_DiamondCreate:", diamondstore.Number, string(diamondstore.Diamond), diamondstore.MinerAddress.ToReadable())
	//fmt.Print(string(diamondstore.Diamond)+",")

	//fmt.Println("Action_4_DiamondCreate:", act.Nonce)

	return nil
}

func (act *Action_4_DiamondCreate) WriteinChainState(state interfacev2.ChainStateOperation) error {
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

	// Whether comprehensive inspection is necessary
	var mustDoAllCheck = true

	if sys.TestDebugLocalDevelopmentMark {
		mustDoAllCheck = false // Developer mode does not check
	}
	//fmt.Println(state.IsDatabaseVersionRebuildMode(), "-------------------------")
	if state.IsDatabaseVersionRebuildMode() {
		mustDoAllCheck = false // Database upgrade mode does not check
	}

	// Calculate diamond hash
	sha3hash, diamondResHash, diamondStr := x16rs.Diamond(uint32(act.Number), act.PrevHash, act.Nonce, act.Address, act.GetRealCustomMessage())
	// Whether to conduct comprehensive inspection
	if mustDoAllCheck {
		// Transaction can only contain one and only one action
		belongactionnum := len(act.belong_trs.GetActions())
		if 1 != belongactionnum {
			return fmt.Errorf("Diamond create tx need only one action but got %d actions.", belongactionnum)
		}
		// Check block height
		// Check whether the block height value is a multiple of 5
		// {backtopool} means to throw it back to the trading pool and wait for the next block to be processed again
		if blkhei%5 != 0 {
			return fmt.Errorf("{BACKTOPOOL} Diamond must be in block height multiple of 5.")
		}

		// Miner status inspection
		lastdiamond, err := state.ReadLastestDiamond()
		if err != nil {
			return err
		}
		if lastdiamond == nil && act.Number > 1 {
			// Not the first diamond but can not read the information of the last diamond
			return fmt.Errorf("Cannot ReadLastestDiamond() to get last diamond of act.Number > 1.")
		}
		if lastdiamond != nil {
			//fmt.Println(lastdiamond.Diamond)
			//fmt.Println(lastdiamond.Number)
			//fmt.Println(lastdiamond.ContainBlockHash.ToHex())
			//fmt.Println(lastdiamond.PrevContainBlockHash.ToHex())
			prevdiamondnum, prevdiamondhash := uint32(lastdiamond.Number), lastdiamond.ContainBlockHash
			if prevdiamondhash == nil {
				return fmt.Errorf("lastdiamond.ContainBlockHash is nil.")
			}
			// Check if the diamond is from the previous block
			if act.PrevHash.Equal(prevdiamondhash) != true {
				return fmt.Errorf("Diamond prev hash must be <%s> but got <%s>.", hex.EncodeToString(prevdiamondhash), hex.EncodeToString(act.PrevHash))
			}
			if prevdiamondnum+1 != uint32(act.Number) {
				return fmt.Errorf("Diamond number must be <%d> but got <%d>.", prevdiamondnum+1, act.Number)
			}
		}
		// Check diamond mining calculation
		diamondstrval, isdia := x16rs.IsDiamondHashResultString(diamondStr)
		if !isdia {
			return fmt.Errorf("String <%s> is not diamond.", diamondStr)
		}
		if strings.Compare(diamondstrval, string(act.Diamond)) != 0 {
			return fmt.Errorf("Diamond need <%s> but got <%s>", act.Diamond, diamondstrval)
		}
		// Check diamond difficulty value
		difok := x16rs.CheckDiamondDifficulty(uint32(act.Number), sha3hash, diamondResHash)
		if !difok {
			return fmt.Errorf("Diamond difficulty not meet the requirements.")
		}
		// Query whether the diamond already exists
		hasaddr, e := state.Diamond(act.Diamond)
		if e != nil {
			return e
		}
		if hasaddr != nil {
			return fmt.Errorf("Diamond <%s> already exist.", string(act.Diamond))
		}
		// Check that a block can contain only one diamond
		pendingdiamond, e2 := state.GetPendingSubmitStoreDiamond()
		if e2 != nil {
			return e2
		}
		if pendingdiamond != nil {
			return fmt.Errorf("This block height:%d has already exist diamond:<%s> .", blkhei, pendingdiamond.Diamond)
		}
		// All conditions checked successfully
	}

	// Deposit diamonds
	//fmt.Println(act.Address.ToReadable())
	var diastore = stores.NewDiamond(act.Address)
	diastore.Address = act.Address // Diamond address
	e3 := state.DiamondSet(act.Diamond, diastore)
	if e3 != nil {
		return e3
	}
	// Increase diamond balance +1
	e9 := DoAddDiamondFromChainState(state, act.Address, 1)
	if e9 != nil {
		return e9
	}

	// Set miner status
	// Mark that this block already contains diamonds
	// Storing objects, computing visual genes
	lifeGene, e15 := calculateLifeGeneByDiamondStuffHash(act.belong_trs, uint32(act.Number), diamondResHash, diamondStr, diamondVisualUseContainBlockHash)
	if e15 != nil {
		return e15
	}

	var diamondstore = &stores.DiamondSmelt{
		Diamond:              act.Diamond,
		Number:               act.Number,
		ContainBlockHeight:   fields.BlockHeight(blkhei),
		ContainBlockHash:     blkhash,
		PrevContainBlockHash: act.PrevHash,
		MinerAddress:         act.Address,
		Nonce:                act.Nonce,
		CustomMessage:        act.GetRealCustomMessage(),
		LifeGene:             lifeGene,
	}

	// Write service charge quotation
	feeoffer := act.belong_trs.GetFee()
	e11 := diamondstore.ParseApproxFeeOffer(feeoffer)
	if e11 != nil {
		return e11
	}
	// Total supply statistics
	totalsupply, e2 := state.ReadTotalSupply()
	if e2 != nil {
		return e2
	}

	// Calculate the average number of HACs for bidding
	if uint32(act.Number) <= DiamondStatisticsAverageBiddingBurningPriceAboveNumber {
		diamondstore.AverageBidBurnPrice = 10 // Fixed to 10
	} else {
		bsnum := uint32(act.Number) - DiamondCreateBurning90PercentTxFeesAboveNumber
		burnhac := totalsupply.Get(stores.TotalSupplyStoreTypeOfBurningFee)
		bidprice := uint64(burnhac/float64(bsnum) + 0.99999999) // up 1
		setprice := fields.VarUint2(bidprice)
		if setprice < 1 {
			setprice = 1 // Minimum 1
		}
		diamondstore.AverageBidBurnPrice = setprice
	}

	// Set latest status
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

func (act *Action_4_DiamondCreate) RecoverChainState(state interfacev2.ChainStateOperation) error {

	panic("RecoverChainState be deprecated")

	//chainstate := state.BlockStore()
	//if chainstate == nil {
	//	panic("Action get state.Miner() cannot be nil !")
	//}
	// Delete diamond
	e1 := state.DiamondDel(act.Diamond)
	if e1 != nil {
		return e1
	}
	// Fallback miner status
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
	// Deduct diamond balance -1
	e9 := DoSubDiamondFromChainState(state, act.Address, 1)
	if e9 != nil {
		return e9
	}
	// Total supply statistics
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

func (elm *Action_4_DiamondCreate) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_4_DiamondCreate) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // IsBurning 90 PersentTxFees
func (act *Action_4_DiamondCreate) IsBurning90PersentTxFees() bool {
	if uint32(act.Number) > DiamondCreateBurning90PercentTxFeesAboveNumber {
		// 90% of the cost of destroying this transaction from the 30001 diamond
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////

// Calculate the visual gene of diamond
func calculateLifeGeneByDiamondStuffHash(belong_trs interfacev2.Transaction, number uint32, stuffhx []byte, diamondstr string, peddingblkhash []byte) (fields.Hash, error) {
	if len(stuffhx) != 32 || len(peddingblkhash) != 32 {
		return nil, fmt.Errorf("stuffhx and peddingblkhash length must 32")
	}
	if len(diamondstr) != 16 {
		return nil, fmt.Errorf("diamondstr length must 16")
	}
	vgenehash := make([]byte, 32)
	copy(vgenehash, stuffhx)
	if number > DiamondResourceHashAndContainBlockHashDecideVisualGeneAboveNumber {
		// 40001 diamond, start with Sha3_ Hash (diamondreshash, blockhash) determines diamond shape and color matching
		vgenestuff := bytes.NewBuffer(stuffhx)
		vgenestuff.Write(peddingblkhash)
		if number > DiamondResourceAppendBiddingFeeDecideVisualGeneAboveNumber {
			bidfeebts, e := belong_trs.GetFee().Serialize() // bid fee
			if e != nil {
				return nil, e // Return error
			}
			vgenestuff.Write(bidfeebts) // Bidding fee participates in determining diamond shape and color matching
		}
		vgenehash = fields.CalculateHash(vgenestuff.Bytes()) // Open blind box
		// Like block hash, it is random, and the shape and color matching can only be known when the diamond is confirmed
		// fmt.Println(hex.EncodeToString(vgenestuff.Bytes()))
	}
	// fmt.Printf("Calculate Visual Gene #%d, vgenehash: %s, stuffhx: %s, peddingblkhash: %s\n", number, hex.EncodeToString(vgenehash), hex.EncodeToString(stuffhx), hex.EncodeToString(peddingblkhash))

	return vgenehash, nil

	/*
		genehexstr := make([]string, 18)
		// Top 6
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
		// Last 11 digits
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
		// Make up the last digit
		genehexstr[17] = "0"
		resbts, e1 := hex.DecodeString(strings.Join(genehexstr, ""))
		if e1 != nil {
			return nil, e1
		}
		// Last bit of hash as shape selection
		resbuf := bytes.NewBuffer([]byte{vgenehash[31]})
		resbuf.Write(resbts) // Color selector
		return resbuf.Bytes(), nil
	*/
}

// Calculate the visual gene of diamond
func calculateLifeGeneByDiamondStuffHashV3(belong_trs interfaces.Transaction, number uint32, stuffhx []byte, diamondstr string, peddingblkhash []byte) (fields.Hash, error) {
	if len(stuffhx) != 32 || len(peddingblkhash) != 32 {
		return nil, fmt.Errorf("stuffhx and peddingblkhash length must 32")
	}
	if len(diamondstr) != 16 {
		return nil, fmt.Errorf("diamondstr length must 16")
	}
	vgenehash := make([]byte, 32)
	copy(vgenehash, stuffhx)
	if number > DiamondResourceHashAndContainBlockHashDecideVisualGeneAboveNumber {
		// 40001 diamond, start with Sha3_ Hash (diamondreshash, blockhash) determines diamond shape and color matching
		vgenestuff := bytes.NewBuffer(stuffhx)
		vgenestuff.Write(peddingblkhash)
		if number > DiamondResourceAppendBiddingFeeDecideVisualGeneAboveNumber {
			bidfeebts, e := belong_trs.GetFee().Serialize() // bid fee
			if e != nil {
				return nil, e // Return error
			}
			vgenestuff.Write(bidfeebts) // Bidding fee participates in determining diamond shape and color matching
		}
		vgenehash = fields.CalculateHash(vgenestuff.Bytes()) // Open blind box
		// Like block hash, it is random, and the shape and color matching can only be known when the diamond is confirmed
		// fmt.Println(hex.EncodeToString(vgenestuff.Bytes()))
	}

	return vgenehash, nil

	/*
		// fmt.Printf("Calculate Visual Gene #%d, vgenehash: %s, stuffhx: %s, peddingblkhash: %s\n", number, hex.EncodeToString(vgenehash), hex.EncodeToString(stuffhx), hex.EncodeToString(peddingblkhash))

		genehexstr := make([]string, 18)
		// Top 6
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
		// Last 11 digits
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
		// Make up the last digit
		genehexstr[17] = "0"
		resbts, e1 := hex.DecodeString(strings.Join(genehexstr, ""))
		if e1 != nil {
			return nil, e1
		}
		// Last bit of hash as shape selection
		resbuf := bytes.NewBuffer([]byte{vgenehash[31]})
		resbuf.Write(resbts) // Color selector
		return resbuf.Bytes(), nil
	*/
}

///////////////////////////////////////////////////////////////

// Transfer diamond
type Action_5_DiamondTransfer struct {
	Diamond   fields.DiamondName // Diamond literal wtyuiahxvmekbszn
	ToAddress fields.Address     // receive address

	// Data pointer
	// Transaction
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
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

func (act *Action_5_DiamondTransfer) WriteInChainState(state interfaces.ChainStateOperation) error {

	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}

	trsMainAddress := act.belong_trs_v3.GetAddress()

	//fmt.Println("Action_5_DiamondTransfer:", trsMainAddress.ToReadable(), act.Address.ToReadable(), string(act.Diamond))

	// cannot trs to self
	if bytes.Compare(act.ToAddress, trsMainAddress) == 0 {
		return fmt.Errorf("Cannot transfer to self.")
	}
	// Query whether the diamond already exists
	diaitem, e := state.Diamond(act.Diamond)
	if e != nil {
		return e
	}
	if diaitem == nil {
		return fmt.Errorf("Diamond <%s> not exist.", string(act.Diamond))
	}
	item := diaitem
	// Check whether it is mortgaged and whether it can be transferred
	if diaitem.Status != stores.DiamondStatusNormal {
		return fmt.Errorf("Diamond <%s> has been mortgaged and cannot be transferred.", string(act.Diamond))
	}
	// Check which
	if bytes.Compare(item.Address, trsMainAddress) != 0 {
		return fmt.Errorf("Diamond <%s> not belong to belong_trs address.", string(act.Diamond))
	}
	// Transfer diamond
	item.Address = act.ToAddress
	err := state.DiamondSet(act.Diamond, item)
	if err != nil {
		return err
	}
	// Transfer diamond balance
	e9 := DoSimpleDiamondTransferFromChainStateV3(state, trsMainAddress, act.ToAddress, 1)
	if e9 != nil {
		return e9
	}
	return nil
}

func (act *Action_5_DiamondTransfer) WriteinChainState(state interfacev2.ChainStateOperation) error {

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	trsMainAddress := act.belong_trs.GetAddress()

	//fmt.Println("Action_5_DiamondTransfer:", trsMainAddress.ToReadable(), act.Address.ToReadable(), string(act.Diamond))

	// cannot trs to self
	if bytes.Compare(act.ToAddress, trsMainAddress) == 0 {
		return fmt.Errorf("Cannot transfer to self.")
	}
	// Query whether the diamond already exists
	diaitem, e := state.Diamond(act.Diamond)
	if e != nil {
		return e
	}
	if diaitem == nil {
		return fmt.Errorf("Diamond <%s> not exist.", string(act.Diamond))
	}
	item := diaitem
	// Check whether it is mortgaged and whether it can be transferred
	if diaitem.Status != stores.DiamondStatusNormal {
		return fmt.Errorf("Diamond <%s> has been mortgaged and cannot be transferred.", string(act.Diamond))
	}
	// Check which
	if bytes.Compare(item.Address, trsMainAddress) != 0 {
		return fmt.Errorf("Diamond <%s> not belong to belong_trs address.", string(act.Diamond))
	}
	// Transfer diamond
	item.Address = act.ToAddress
	err := state.DiamondSet(act.Diamond, item)
	if err != nil {
		return err
	}
	// Transfer diamond balance
	e9 := DoSimpleDiamondTransferFromChainState(state, trsMainAddress, act.ToAddress, 1)
	if e9 != nil {
		return e9
	}
	return nil
}

func (act *Action_5_DiamondTransfer) RecoverChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	trsMainAddress := act.belong_trs.GetAddress()
	// get diamond
	diaitem, e := state.Diamond(act.Diamond)
	if e != nil {
		return e
	}
	if diaitem == nil {
		return fmt.Errorf("Diamond <%s> not exist.", string(act.Diamond))
	}
	item := diaitem
	// Back off diamond
	item.Address = act.belong_trs.GetAddress()
	err := state.DiamondSet(act.Diamond, item)
	if err != nil {
		return err
	}
	// Return diamond balance
	e9 := DoSimpleDiamondTransferFromChainState(state, act.ToAddress, trsMainAddress, 1)
	if e9 != nil {
		return e9
	}
	return nil
}

func (elm *Action_5_DiamondTransfer) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_5_DiamondTransfer) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // IsBurning 90 PersentTxFees
func (act *Action_5_DiamondTransfer) IsBurning90PersentTxFees() bool {
	return false
}

///////////////////////////////////////////////////////////////

// Bulk transfer of diamonds
type Action_6_OutfeeQuantityDiamondTransfer struct {
	FromAddress fields.Address              // Accounts with diamonds
	ToAddress   fields.Address              // receive address
	DiamondList fields.DiamondListMaxLen200 // Diamond list

	// Data pointer
	// Transaction
	belong_trs    interfacev2.Transaction
	belong_trs_v3 interfaces.Transaction
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) Kind() uint16 {
	return 6
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) Size() uint32 {
	return 2 +
		elm.FromAddress.Size() +
		elm.ToAddress.Size() +
		elm.DiamondList.Size() // Each diamond is 6 digits long
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
	reqs := make([]fields.Address, 1) // ed from address sign
	reqs[0] = elm.FromAddress
	return reqs
}

func (act *Action_6_OutfeeQuantityDiamondTransfer) WriteInChainState(state interfaces.ChainStateOperation) error {
	if act.belong_trs_v3 == nil {
		panic("Action belong to transaction not be nil !")
	}
	// Quantity check
	dianum := int(act.DiamondList.Count)
	if dianum == 0 || dianum != len(act.DiamondList.Diamonds) {
		return fmt.Errorf("Diamonds quantity error")
	}
	if dianum > 200 {
		return fmt.Errorf("Diamonds quantity cannot over 200")
	}
	// annot trs to self
	if bytes.Compare(act.FromAddress, act.ToAddress) == 0 {
		return fmt.Errorf("Cannot transfer to self.")
	}
	// Bulk transfer of diamonds
	for i := 0; i < len(act.DiamondList.Diamonds); i++ {
		diamond := act.DiamondList.Diamonds[i]

		//fmt.Println("Action_6_OutfeeQuantityDiamondTransfer:", act.FromAddress.ToReadable(), act.ToAddress.ToReadable(), string(diamond))

		// fmt.Println("--- " + string(diamond))
		// Query whether the diamond already exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		if diaitem == nil {
			//panic("Quantity Diamond <%s> not exist. " + string(diamond))
			return fmt.Errorf("Quantity Diamond <%s> not exist.", string(diamond))
		}
		item := diaitem
		// Check whether it is mortgaged and whether it can be transferred
		if diaitem.Status != stores.DiamondStatusNormal {
			return fmt.Errorf("Diamond <%s> has been mortgaged and cannot be transferred.", string(diamond))
		}
		// Check which
		if bytes.Compare(item.Address, act.FromAddress) != 0 {
			return fmt.Errorf("Diamond <%s> not belong to address '%s'", string(diamond), act.FromAddress.ToReadable())
		}
		// Transfer diamond
		item.Address = act.ToAddress
		e5 := state.DiamondSet(diamond, item)
		if e5 != nil {
			return e5
		}
	}
	// Transfer diamond balance
	e9 := DoSimpleDiamondTransferFromChainStateV3(state, act.FromAddress, act.ToAddress, fields.DiamondNumber(dianum))
	if e9 != nil {
		return e9
	}
	return nil
}

func (act *Action_6_OutfeeQuantityDiamondTransfer) WriteinChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// Quantity check
	dianum := int(act.DiamondList.Count)
	if dianum == 0 || dianum != len(act.DiamondList.Diamonds) {
		return fmt.Errorf("Diamonds quantity error")
	}
	if dianum > 200 {
		return fmt.Errorf("Diamonds quantity cannot over 200")
	}
	// cannot trs to self
	if bytes.Compare(act.FromAddress, act.ToAddress) == 0 {
		return fmt.Errorf("Cannot transfer to self.")
	}
	// Bulk transfer of diamonds
	for i := 0; i < len(act.DiamondList.Diamonds); i++ {
		diamond := act.DiamondList.Diamonds[i]

		//fmt.Println("Action_6_OutfeeQuantityDiamondTransfer:", act.FromAddress.ToReadable(), act.ToAddress.ToReadable(), string(diamond))

		// fmt.Println("--- " + string(diamond))
		// Query whether the diamond already exists
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		if diaitem == nil {
			//panic("Quantity Diamond <%s> not exist. " + string(diamond))
			return fmt.Errorf("Quantity Diamond <%s> not exist.", string(diamond))
		}
		item := diaitem
		// Check whether it is mortgaged and whether it can be transferred
		if diaitem.Status != stores.DiamondStatusNormal {
			return fmt.Errorf("Diamond <%s> has been mortgaged and cannot be transferred.", string(diamond))
		}
		// Check which
		if bytes.Compare(item.Address, act.FromAddress) != 0 {
			return fmt.Errorf("Diamond <%s> not belong to address '%s'", string(diamond), act.FromAddress.ToReadable())
		}
		// Transfer diamond
		item.Address = act.ToAddress
		e5 := state.DiamondSet(diamond, item)
		if e5 != nil {
			return e5
		}
	}
	// Transfer diamond balance
	e9 := DoSimpleDiamondTransferFromChainState(state, act.FromAddress, act.ToAddress, fields.DiamondNumber(dianum))
	if e9 != nil {
		return e9
	}
	return nil
}

func (act *Action_6_OutfeeQuantityDiamondTransfer) RecoverChainState(state interfacev2.ChainStateOperation) error {
	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// Batch return of diamonds
	for i := 0; i < len(act.DiamondList.Diamonds); i++ {
		diamond := act.DiamondList.Diamonds[i]
		// get diamond
		diaitem, e := state.Diamond(diamond)
		if e != nil {
			return e
		}
		if diaitem == nil {
			return fmt.Errorf("Diamond <%s> not exist.", string(diamond))
		}
		item := diaitem
		// Back off diamond
		item.Address = act.FromAddress
		err := state.DiamondSet(diamond, item)
		if err != nil {
			return err
		}
	}
	// Return diamond balance
	e9 := DoSimpleDiamondTransferFromChainState(state, act.ToAddress, act.FromAddress, fields.DiamondNumber(act.DiamondList.Count))
	if e9 != nil {
		return e9
	}
	return nil
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) SetBelongTransaction(t interfacev2.Transaction) {
	elm.belong_trs = t
}

func (elm *Action_6_OutfeeQuantityDiamondTransfer) SetBelongTrs(t interfaces.Transaction) {
	elm.belong_trs_v3 = t
}

// burning fees  // IsBurning 90 PersentTxFees
func (act *Action_6_OutfeeQuantityDiamondTransfer) IsBurning90PersentTxFees() bool {
	return false
}

// Get the name list of block diamonds
func (elm *Action_6_OutfeeQuantityDiamondTransfer) GetDiamondNamesSplitByComma() string {
	return elm.DiamondList.SerializeHACDlistToCommaSplitString()
}

///////////////////////////////////////////////////////////////////////
