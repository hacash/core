package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	BalanceSize = 4 + 8 + 20 // len = 32
)

type Balance struct {
	Config      fields.Bytes4
	ExtendMarks fields.Bytes8
	Amount      fields.Amount // max len = 20
}

func NewEmptyBalance() *Balance {
	return &Balance{
		Config:      []byte{0, 0, 0, 0},
		ExtendMarks: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		Amount:      *fields.NewEmptyAmount(),
	}
}

func NewBalanceWithAmount(amt *fields.Amount) *Balance {
	return &Balance{
		Config:      []byte{0, 0, 0, 0},
		ExtendMarks: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		Amount:      *amt,
	}
}

///////////////////////////////////////

func (this *Balance) Size() uint32 {
	return uint32(DiamondSize)
}

func (this *Balance) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.Config.Serialize()
	b2, _ := this.ExtendMarks.Serialize()
	b3, _ := this.Amount.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	return buffer.Bytes(), nil
}

func (this *Balance) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = this.Config.Parse(buf, seek)
	seek, _ = this.ExtendMarks.Parse(buf, seek)
	seek, _ = this.Amount.Parse(buf, seek)
	return seek, nil
}
