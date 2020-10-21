package stores

/*
import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	SatoshiSize = 4 + 8
)

type Satoshi struct {
	Config fields.Bytes4
	Amount fields.VarUint8
}

func NewEmptySatoshi() *Satoshi {
	return &Satoshi{
		Config: []byte{0, 0, 0, 0},
		Amount: fields.VarUint8(0),
	}
}

///////////////////////////////////////

func (this *Satoshi) Size() uint32 {
	return uint32(SatoshiSize)
}

func (this *Satoshi) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.Config.Serialize()
	b2, _ := this.Amount.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	return buffer.Bytes(), nil
}

func (this *Satoshi) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = this.Config.Parse(buf, seek)
	seek, _ = this.Amount.Parse(buf, seek)
	return seek, nil
}


*/
