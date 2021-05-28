package fields

import (
	"bytes"
	"fmt"
	"github.com/hacash/x16rs"
	"strings"
)

type DiamondListMaxLen200 struct {
	Count    VarUint1
	Diamonds []Bytes6
}

func NewEmptyDiamondListMaxLen200() *DiamondListMaxLen200 {
	return &DiamondListMaxLen200{
		Count:    VarUint1(0),
		Diamonds: []Bytes6{},
	}
}

func (elm DiamondListMaxLen200) Size() uint32 {
	return 1 + uint32(len(elm.Diamonds))*6
}

func (elm DiamondListMaxLen200) Serialize() ([]byte, error) {
	if int(elm.Count) != len(elm.Diamonds) || len(elm.Diamonds) == 0 {
		return nil, fmt.Errorf("diamond list number quantity count error")
	}
	if len(elm.Diamonds) > 200 {
		return nil, fmt.Errorf("diamonds quantity cannot over 200")
	}
	var buf = bytes.NewBuffer([]byte{byte(elm.Count)})
	for _, v := range elm.Diamonds {
		buf.Write(v)
	}
	return buf.Bytes(), nil
}

func (elm *DiamondListMaxLen200) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	if seek >= uint32(len(buf)) {
		return 0, fmt.Errorf("[DiamondListMaxLen200.Parse] seek out of buf len.")
	}
	elm.Count = VarUint1(0)
	seek, e = elm.Count.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if len(buf) < int(elm.Count)-1 {
		return 0, fmt.Errorf("buf is too short.")
	}
	elm.Diamonds = make([]Bytes6, int(elm.Count))
	for i := 0; i < int(elm.Count); i++ {
		elm.Diamonds[i] = Bytes6{}
		seek, e = elm.Diamonds[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

// 获取区块钻石的名称列表
func (elm *DiamondListMaxLen200) SerializeHACDlistToCommaSplitString() string {
	var names = make([]string, len(elm.Diamonds))
	for i, v := range elm.Diamonds {
		names[i] = string(v)
	}
	return strings.Join(names, ",")
}

// 创建钻石
func (elm *DiamondListMaxLen200) ParseHACDlistBySplitCommaFromString(hacdlistsplitcomma string) error {
	diamonds := strings.Split(hacdlistsplitcomma, ",")
	dianum := len(diamonds)
	if dianum > 200 {
		return fmt.Errorf("diamonds quantity cannot over 200")
	}
	diamondsbytes := make([]Bytes6, dianum)
	for i, v := range diamonds {
		dok := x16rs.IsDiamondValueString(v)
		if !dok {
			return fmt.Errorf("<%s> not a valid diamond name", v)
		}
		diamondsbytes[i] = []byte(v)
	}
	elm.Count = VarUint1(dianum)
	elm.Diamonds = diamondsbytes
	// 成功返回
	return nil
}
