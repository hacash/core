package stores

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/fields"
)

const (
	LockblsSize = 21 + 5 + 3 + 8 + 8 + 8
)

type Lockbls struct {
	MasterAddress            fields.Address // 主地址（领取权）
	EffectBlockHeight        fields.VarInt5 // 生效（开始）区块
	LinearBlockNumber        fields.VarInt3 // 步进区块数 < 17000000 约 160年
	TotalLockAmountBytes     fields.Bytes8  // 总共存入额度
	LinearReleaseAmountBytes fields.Bytes8  // 每次释放额度
	BalanceAmountBytes       fields.Bytes8  // 有效余额（每次可以取出可取额度之内的任意数额）
}

func NewEmptyLockbls(addr fields.Address) *Lockbls {
	return &Lockbls{
		MasterAddress: addr[:],
	}
}

///////////////////////////////////////

func (this *Lockbls) PutAmount(key *fields.Bytes8, amt *fields.Amount) error {
	newBalanceAmountBytes, e5 := amt.Serialize()
	if e5 != nil {
		return e5
	}
	if len(newBalanceAmountBytes) > 8 {
		return fmt.Errorf("length cannot over 8 bytes.") // 储存大小精度
	}
	*key = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	copy(*key, newBalanceAmountBytes)
	return nil
}

func (this *Lockbls) GetAmount(key *fields.Bytes8) (*fields.Amount, error) {
	var amt = new(fields.Amount)
	_, e1 := amt.Parse(*key, 0)
	if e1 != nil {
		return nil, e1
	}
	return amt, nil
}

///////////////////////////////////////

func (this *Lockbls) Size() uint32 {
	return uint32(LockblsSize)
}

func (this *Lockbls) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := this.MasterAddress.Serialize()
	b2, _ := this.EffectBlockHeight.Serialize()
	b3, _ := this.LinearBlockNumber.Serialize()
	b4, _ := this.TotalLockAmountBytes.Serialize()
	b5, _ := this.LinearReleaseAmountBytes.Serialize()
	b6, _ := this.BalanceAmountBytes.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	return buffer.Bytes(), nil
}

func (this *Lockbls) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = this.MasterAddress.Parse(buf, seek)
	seek, _ = this.EffectBlockHeight.Parse(buf, seek)
	seek, _ = this.LinearBlockNumber.Parse(buf, seek)
	seek, _ = this.TotalLockAmountBytes.Parse(buf, seek)
	seek, _ = this.LinearReleaseAmountBytes.Parse(buf, seek)
	seek, _ = this.BalanceAmountBytes.Parse(buf, seek)
	return seek, nil
}
