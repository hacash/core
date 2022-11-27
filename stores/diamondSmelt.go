package stores

import (
	"bytes"
	"encoding/hex"
	"github.com/hacash/core/fields"
	"strings"
)

const (
//DiamondSmeltSize = 6 + 3 + 5 + 32 + 32 + fields.AddressSize + 4 + 8 + 32 + 2
)

type DiamondSmelt struct {
	Diamond              fields.DiamondName // WTYUIAHXVMEKBSZN
	Number               fields.DiamondNumber
	ContainBlockHeight   fields.BlockHeight
	ContainBlockHash     fields.Hash // current pending block hash
	PrevContainBlockHash fields.Hash // prev block hash
	MinerAddress         fields.Address
	ApproxFeeOffer       fields.Amount  // Fee Amount
	Nonce                fields.Bytes8  // nonce
	CustomMessage        fields.Bytes32 // msg
	// data statistics
	AverageBidBurnPrice fields.VarUint2 // Average number of HAc destroyed by bidding, rounded down, the lowest one, the highest 65535
	// other data
	//VisualGene fields.Bytes10 // Visual appearance gene
	LifeGene fields.Hash
}

func (this *DiamondSmelt) Size() uint32 {
	return this.Diamond.Size() +
		this.Number.Size() +
		this.ContainBlockHeight.Size() +
		this.ContainBlockHash.Size() +
		this.PrevContainBlockHash.Size() +
		this.MinerAddress.Size() +
		this.ApproxFeeOffer.Size() +
		this.Nonce.Size() +
		this.CustomMessage.Size() +
		this.AverageBidBurnPrice.Size() +
		this.LifeGene.Size()
}

func (this *DiamondSmelt) GetApproxFeeOffer() *fields.Amount {
	return &this.ApproxFeeOffer
}

func (this *DiamondSmelt) ParseApproxFeeOffer(amt *fields.Amount) error {
	// Compress storage space
	approxfeeoffer, _, e11 := amt.CompressForMainNumLen(4, true)
	if e11 != nil {
		return e11
	}
	this.ApproxFeeOffer = *approxfeeoffer
	return nil
}

func (this *DiamondSmelt) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.Diamond.Serialize()
	b2, _ := this.Number.Serialize()
	b3, _ := this.ContainBlockHeight.Serialize()
	b4, _ := this.ContainBlockHash.Serialize()
	b5, _ := this.PrevContainBlockHash.Serialize()
	b6, _ := this.MinerAddress.Serialize()
	b7, _ := this.ApproxFeeOffer.Serialize()
	b8, _ := this.Nonce.Serialize()
	b9, _ := this.CustomMessage.Serialize()
	b10, _ := this.AverageBidBurnPrice.Serialize()
	b11, _ := this.LifeGene.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	buffer.Write(b7)
	buffer.Write(b8)
	buffer.Write(b9)
	buffer.Write(b10)
	buffer.Write(b11)
	return buffer.Bytes(), nil
}

func (this *DiamondSmelt) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = this.Diamond.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.Number.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.ContainBlockHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.ContainBlockHash.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.PrevContainBlockHash.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.MinerAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.ApproxFeeOffer.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.Nonce.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.CustomMessage.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.AverageBidBurnPrice.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.LifeGene.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

///////////////////

func (this *DiamondSmelt) GetVisualGene() fields.Bytes10 {

	var diamondstr = this.Diamond.Name()
	var vgenehash = this.LifeGene

	//fmt.Println(vgenehash.ToHex())

	genehexstr := make([]string, 18)
	// Top 6
	k := 0
	for i := 0; i < 6; i++ {
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
		return nil
	}
	// Last bit of hash as shape selection
	resbuf := bytes.NewBuffer([]byte{vgenehash[31]})
	resbuf.Write(resbts) // Color selector
	return resbuf.Bytes()
}
