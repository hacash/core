package fields

import (
	"bytes"
	"fmt"
)

// 可选地址

type OptionalAddress struct {
	Exist Bool
	Addr  Address
}

func NewEmptyOptionalAddress() OptionalAddress {
	return OptionalAddress{
		Exist: CreateBool(false),
		Addr:  nil,
	}
}

func (elm OptionalAddress) Size() uint32 {
	if elm.Exist.Check() {
		return 1 + AddressSize
	} else {
		return 1
	}
}

func (elm OptionalAddress) Serialize() ([]byte, error) {
	var buf = bytes.NewBuffer(nil)
	bt1, _ := elm.Exist.Serialize()
	buf.Write(bt1)
	if elm.Exist.Check() {
		bt2, _ := elm.Addr.Serialize()
		buf.Write(bt2)
	}
	return buf.Bytes(), nil
}

func (elm *OptionalAddress) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	if seek >= uint32(len(buf)) {
		return 0, fmt.Errorf("[OptionalAddress.Parse] seek out of buf len.")
	}
	seek, e = elm.Exist.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if elm.Exist.Check() {
		seek, e = elm.Addr.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

// 显示地址
func (elm OptionalAddress) ShowReadableOrEmpty() string {
	if elm.Exist.Check() {
		return elm.Addr.ToReadable()
	} else {
		return ""
	}
}
