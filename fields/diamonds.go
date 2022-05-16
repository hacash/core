package fields

import (
	"bytes"
	"fmt"
	"strings"
)

type DiamondListMaxLen200 struct {
	Count    VarUint1
	Diamonds []DiamondName
}

func NewEmptyDiamondListMaxLen200() *DiamondListMaxLen200 {
	return &DiamondListMaxLen200{
		Count:    VarUint1(0),
		Diamonds: []DiamondName{},
	}
}

func (elm DiamondListMaxLen200) Size() uint32 {
	return 1 + uint32(len(elm.Diamonds))*6
}

func (elm DiamondListMaxLen200) Serialize() ([]byte, error) {
	if int(elm.Count) != len(elm.Diamonds) {
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
	if elm.Count == 0 {
		return seek, nil // 列表为空
	}
	elm.Diamonds = make([]DiamondName, int(elm.Count))
	for i := 0; i < int(elm.Count); i++ {
		elm.Diamonds[i] = DiamondName{}
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
	// 去除空格和换行符
	hacdlistsplitcomma = strings.Replace(hacdlistsplitcomma, " ", "", -1)
	hacdlistsplitcomma = strings.Replace(hacdlistsplitcomma, "\n", "", -1)
	hacdlistsplitcomma = strings.Trim(hacdlistsplitcomma, ",")
	// 分割
	diamonds := strings.Split(hacdlistsplitcomma, ",")
	dianum := len(diamonds)
	if dianum > 200 {
		return fmt.Errorf("diamonds quantity cannot over 200")
	}
	dianamesmap := make(map[string]bool) // 去重
	diamondsbytes := make([]DiamondName, dianum)
	for i, v := range diamonds {
		dok := IsDiamondValueString(v)
		if !dok {
			return fmt.Errorf("<%s> not a valid diamond name", v)
		}
		diamondsbytes[i] = []byte(v)
		if dianamesmap[v] {
			return fmt.Errorf("<%s> appear in the list repeatedly", v)
		}
		dianamesmap[v] = true
	}
	elm.Count = VarUint1(dianum)
	elm.Diamonds = diamondsbytes
	// 成功返回
	return nil
}
