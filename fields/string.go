package fields

import (
	"bytes"
	"fmt"
)

// max length = 255
type StringMax255 struct {
	Len VarUint1
	Str string
}

func CreateStringMax255(str string) StringMax255 {
	l := len(str)
	if l > 255 {
		l = 255
	}
	return StringMax255{
		Len: VarUint1(l),
		Str: str,
	}
}

func (elm StringMax255) Value() string {
	return elm.Str
}

func (elm StringMax255) Size() uint32 {
	return elm.Len.Size() + uint32(elm.Len)
}

func (elm StringMax255) Serialize() ([]byte, error) {
	var buf = bytes.NewBuffer(nil)
	bt1, _ := elm.Len.Serialize()
	buf.Write(bt1)
	if elm.Len > 0 {
		if len(elm.Str) != int(elm.Len) {
			return nil, fmt.Errorf("Str lenght error")
		}
		buf.Write([]byte(elm.Str))
	}
	return buf.Bytes(), nil
}

func (elm *StringMax255) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	if seek >= uint32(len(buf)) {
		return 0, fmt.Errorf("[StringMax255.Parse] seek out of buf len.")
	}
	seek, e = elm.Len.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if elm.Len > 0 {
		end := seek + uint32(elm.Len)
		if len(buf) < int(end) {
			return 0, fmt.Errorf("Str lenght error")
		}
		elm.Str = string(buf[seek:end])
		seek = end
	}
	return seek, nil
}

// max length = 65535
type StringMax65535 struct {
	Len VarUint2
	Str string
}

func CreateStringMax65535(str string) StringMax65535 {
	l := len(str)
	if l > 65535 {
		l = 65535
	}
	return StringMax65535{
		Len: VarUint2(l),
		Str: str,
	}
}

func (elm StringMax65535) Value() string {
	return elm.Str
}

func (elm StringMax65535) Size() uint32 {
	return elm.Len.Size() + uint32(elm.Len)
}

func (elm StringMax65535) Serialize() ([]byte, error) {
	var buf = bytes.NewBuffer(nil)
	bt1, _ := elm.Len.Serialize()
	buf.Write(bt1)
	if elm.Len > 0 {
		if len(elm.Str) != int(elm.Len) {
			return nil, fmt.Errorf("Str lenght error")
		}
		buf.Write([]byte(elm.Str))
	}
	return buf.Bytes(), nil
}

func (elm *StringMax65535) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	if seek >= uint32(len(buf)) {
		return 0, fmt.Errorf("[StringMax65535.Parse] seek out of buf len.")
	}
	seek, e = elm.Len.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if elm.Len > 0 {
		end := seek + uint32(elm.Len)
		if len(buf) < int(end) {
			return 0, fmt.Errorf("Str lenght error")
		}
		elm.Str = string(buf[seek:end])
		seek = end
	}
	return seek, nil
}

// max length = 16777215
type StringMax16777215 struct {
	Len VarUint3
	Str string
}

func CreateStringMax16777215(str string) StringMax16777215 {
	l := len(str)
	if l > 16777215 {
		l = 16777215
	}
	return StringMax16777215{
		Len: VarUint3(l),
		Str: str,
	}
}

func (elm StringMax16777215) Value() string {
	return elm.Str
}

func (elm StringMax16777215) Size() uint32 {
	return elm.Len.Size() + uint32(elm.Len)
}

func (elm StringMax16777215) Serialize() ([]byte, error) {
	var buf = bytes.NewBuffer(nil)
	bt1, _ := elm.Len.Serialize()
	buf.Write(bt1)
	if elm.Len > 0 {
		if len(elm.Str) != int(elm.Len) {
			return nil, fmt.Errorf("Str lenght error")
		}
		buf.Write([]byte(elm.Str))
	}
	return buf.Bytes(), nil
}

func (elm *StringMax16777215) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	if seek >= uint32(len(buf)) {
		return 0, fmt.Errorf("[StringMax16777215.Parse] seek out of buf len.")
	}
	seek, e = elm.Len.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if elm.Len > 0 {
		end := seek + uint32(elm.Len)
		if len(buf) < int(end) {
			return 0, fmt.Errorf("Str lenght error")
		}
		elm.Str = string(buf[seek:end])
		seek = end
	}
	return seek, nil
}
