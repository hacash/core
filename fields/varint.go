package fields

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

type Bool uint8
type VarUint1 uint8
type VarUint2 uint16
type VarUint3 uint32
type VarUint4 uint32
type VarUint5 uint64
type VarUint6 uint64
type VarUint7 uint64
type VarUint8 uint64

////////////////////////////////////////////////////////

func (elm *Bool) Serialize() ([]byte, error)     { return varIntSerialize(uint64(*elm), 1) }
func (elm *VarUint1) Serialize() ([]byte, error) { return varIntSerialize(uint64(*elm), 1) }
func (elm *VarUint2) Serialize() ([]byte, error) { return varIntSerialize(uint64(*elm), 2) }
func (elm *VarUint3) Serialize() ([]byte, error) { return varIntSerialize(uint64(*elm), 3) }
func (elm *VarUint4) Serialize() ([]byte, error) { return varIntSerialize(uint64(*elm), 4) }
func (elm *VarUint5) Serialize() ([]byte, error) { return varIntSerialize(uint64(*elm), 5) }
func (elm *VarUint6) Serialize() ([]byte, error) { return varIntSerialize(uint64(*elm), 6) }
func (elm *VarUint7) Serialize() ([]byte, error) { return varIntSerialize(uint64(*elm), 7) }
func (elm *VarUint8) Serialize() ([]byte, error) { return varIntSerialize(uint64(*elm), 8) }

func (elm *Bool) Parse(buf []byte, seek uint32) (uint32, error) {
	return varIntParse(elm, buf, seek, 1)
}
func (elm *VarUint1) Parse(buf []byte, seek uint32) (uint32, error) {
	return varIntParse(elm, buf, seek, 1)
}
func (elm *VarUint2) Parse(buf []byte, seek uint32) (uint32, error) {
	return varIntParse(elm, buf, seek, 2)
}
func (elm *VarUint3) Parse(buf []byte, seek uint32) (uint32, error) {
	return varIntParse(elm, buf, seek, 3)
}
func (elm *VarUint4) Parse(buf []byte, seek uint32) (uint32, error) {
	return varIntParse(elm, buf, seek, 4)
}
func (elm *VarUint5) Parse(buf []byte, seek uint32) (uint32, error) {
	return varIntParse(elm, buf, seek, 5)
}
func (elm *VarUint6) Parse(buf []byte, seek uint32) (uint32, error) {
	return varIntParse(elm, buf, seek, 6)
}
func (elm *VarUint7) Parse(buf []byte, seek uint32) (uint32, error) {
	return varIntParse(elm, buf, seek, 7)
}
func (elm *VarUint8) Parse(buf []byte, seek uint32) (uint32, error) {
	return varIntParse(elm, buf, seek, 8)
}

func (elm *Bool) Size() uint32     { return 1 }
func (elm *VarUint1) Size() uint32 { return 1 }
func (elm *VarUint2) Size() uint32 { return 2 }
func (elm *VarUint3) Size() uint32 { return 3 }
func (elm *VarUint4) Size() uint32 { return 4 }
func (elm *VarUint5) Size() uint32 { return 5 }
func (elm *VarUint6) Size() uint32 { return 6 }
func (elm *VarUint7) Size() uint32 { return 7 }
func (elm *VarUint8) Size() uint32 { return 8 }

// 判断
func (elm *Bool) Is(v bool) bool { return elm.Check() == v }
func (elm *Bool) Check() bool    { return int(*elm) != 0 }
func CreateBoolPtr(v bool) *Bool {
	b := CreateBool(v)
	return &b
}
func CreateBool(v bool) Bool {
	if v {
		return 1
	} else {
		return 0
	}
}

////////////////////////////////////////////////////////

func varIntSerialize(val uint64, maxlen uint32) ([]byte, error) {
	var intbytes = bytes.Repeat([]byte{0}, 8)
	binary.BigEndian.PutUint64(intbytes, val)
	byyy := intbytes[8-maxlen : 8]
	//fmt.Println(intbytes)
	//fmt.Println("---- %d", maxlen)
	//fmt.Println(byyy)
	return byyy, nil
}

func varIntParse(elm interface{}, buf []byte, seek uint32, maxlen uint32) (uint32, error) {
	// fmt.Println("xxx",*buf)
	if seek+maxlen > uint32(len(buf)) {
		return 0, fmt.Errorf("[varIntParse] seek out of buf len.")
	}
	nnnold := buf[seek : seek+maxlen]
	var intbytes = make([]byte, len(nnnold))
	copy(intbytes, nnnold)
	// fmt.Println(intbytes)
	padbytes := bytes.Repeat([]byte{0}, int(8-maxlen))
	intbytes = append(padbytes, intbytes...)
	//addrbytes = bytes.TrimRight(addrbytes, " ")
	val := binary.BigEndian.Uint64(intbytes)
	// fmt.Println(intbytes)
	// fmt.Println("====== %d", val)
	switch a := elm.(type) {
	case *Bool:
		// v:= (val)>>56
		// fmt.Println("**** %d", v)
		*a = *(*Bool)(unsafe.Pointer(&val))
		// fmt.Println("------- %d", *a)
	case *VarUint1:
		// v:= (val)>>56
		// fmt.Println("**** %d", v)
		*a = *(*VarUint1)(unsafe.Pointer(&val))
		// fmt.Println("------- %d", *a)
	case *VarUint2:
		// v:= val>>48
		// fmt.Println("**** %d", v)
		*a = *(*VarUint2)(unsafe.Pointer(&val))
		// fmt.Println("------- %d", *a)
	case *VarUint3:
		// v:= val>>48
		// fmt.Println("**** %d", v)
		*a = *(*VarUint3)(unsafe.Pointer(&val))
		// fmt.Println("------- %d", *a)
	case *VarUint4:
		// v:= val>>32
		// fmt.Println("**** %d", v)
		*a = *(*VarUint4)(unsafe.Pointer(&val))
		// fmt.Println("------- %d", *a)
	case *VarUint5:
		*a = *(*VarUint5)(unsafe.Pointer(&val))
		// fmt.Println("------- %d", *a)
	case *VarUint6:
		*a = *(*VarUint6)(unsafe.Pointer(&val))
		// fmt.Println("------- %d", *a)
	case *VarUint7:
		*a = *(*VarUint7)(unsafe.Pointer(&val))
		// fmt.Println("------- %d", *a)
	case *VarUint8:
		*a = *(*VarUint8)(unsafe.Pointer(&val))
		// fmt.Println("------- %d", *a)
	default:
		//fmt.Println("")
		return 0, fmt.Errorf("not find type")
	}

	return seek + maxlen, nil
}
