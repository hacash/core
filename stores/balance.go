package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	BalanceSize = 3 + 8 + 11 // len = 22
)

type Balance struct {
	Diamond fields.VarUint3
	Satoshi fields.VarUint8
	Hacash  fields.Amount // len = 11
}

func NewEmptyBalance() *Balance {
	return &Balance{
		Diamond: 0,
		Satoshi: 0,
		Hacash: fields.Amount{
			Unit:    0,
			Dist:    0,
			Numeral: nil,
		},
	}
}

func NewBalanceWithAmount(amt *fields.Amount) *Balance {
	return &Balance{
		Diamond: 0,
		Satoshi: 0,
		Hacash:  *amt,
	}
}

///////////////////////////////////////

func (this *Balance) Size() uint32 {
	return uint32(BalanceSize)
}

func (this *Balance) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.Diamond.Serialize()
	b2, _ := this.Satoshi.Serialize()
	b3, _ := this.Hacash.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	return buffer.Bytes(), nil
}

func (this *Balance) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = this.Diamond.Parse(buf, seek)
	seek, _ = this.Satoshi.Parse(buf, seek)
	seek, _ = this.Hacash.Parse(buf, seek)
	return seek, nil
}
