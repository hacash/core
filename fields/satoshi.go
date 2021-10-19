package fields

import (
	"bytes"
	"fmt"
)

// 比特币一亿分之一
type Satoshi = VarUint8

func (elm Satoshi) GetSatoshiVariation() SatoshiVariation {
	if elm == 0 {
		return NewEmptySatoshiVariation()
	}
	return SatoshiVariation{
		NotEmpty: CreateBool(true),
		ValueSAT: elm,
	}
}

type SatoshiVariation struct {
	NotEmpty Bool
	ValueSAT Satoshi
}

func NewEmptySatoshiVariation() SatoshiVariation {
	return SatoshiVariation{
		NotEmpty: CreateBool(false),
		ValueSAT: 0,
	}
}

func (elm SatoshiVariation) Size() uint32 {
	if elm.NotEmpty.Check() {
		return 1 + 8
	} else {
		return 1
	}
}

func (elm SatoshiVariation) Serialize() ([]byte, error) {
	var buf = bytes.NewBuffer(nil)
	bt1, _ := elm.NotEmpty.Serialize()
	buf.Write(bt1)
	if elm.NotEmpty.Check() {
		bt2, _ := elm.ValueSAT.Serialize()
		buf.Write(bt2)
	}
	return buf.Bytes(), nil
}

func (elm *SatoshiVariation) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	if seek >= uint32(len(buf)) {
		return 0, fmt.Errorf("[SatoshiVariation.Parse] seek out of buf len.")
	}
	seek, e = elm.NotEmpty.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if elm.NotEmpty.Check() {
		seek, e = elm.ValueSAT.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

// 获得真实值
func (elm *SatoshiVariation) GetRealSatoshi() Satoshi {
	if elm.NotEmpty.Check() {
		return elm.ValueSAT
	}
	return Satoshi(0)
}
