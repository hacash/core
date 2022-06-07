package stores

import (
	"bytes"
	"github.com/hacash/core/fields"
)

const (
	LockblsIdLength = 18
)

type Lockbls struct {
	MasterAddress       fields.Address     // Main address (claim)
	EffectBlockHeight   fields.BlockHeight // Effective (start) block
	LinearBlockNumber   fields.VarUint3    // Number of stepping blocks < 17000000 about 160 years
	TotalLockAmount     fields.Amount      // Total deposit limit
	LinearReleaseAmount fields.Amount      // Limit released each time
	BalanceAmount       fields.Amount      // Effective balance (any amount within the available limit can be withdrawn each time)
}

func NewEmptyLockbls(addr fields.Address) *Lockbls {
	return &Lockbls{
		MasterAddress: addr[:],
	}
}

func (this *Lockbls) Size() uint32 {
	return this.MasterAddress.Size() +
		this.EffectBlockHeight.Size() +
		this.LinearBlockNumber.Size() +
		this.TotalLockAmount.Size() +
		this.LinearReleaseAmount.Size() +
		this.BalanceAmount.Size()
}

func (this *Lockbls) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.MasterAddress.Serialize()
	b2, _ := this.EffectBlockHeight.Serialize()
	b3, _ := this.LinearBlockNumber.Serialize()
	b4, _ := this.TotalLockAmount.Serialize()
	b5, _ := this.LinearReleaseAmount.Serialize()
	b6, _ := this.BalanceAmount.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	return buffer.Bytes(), nil
}

func (this *Lockbls) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = this.MasterAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.EffectBlockHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.LinearBlockNumber.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.TotalLockAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.LinearReleaseAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = this.BalanceAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}
