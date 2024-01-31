package fields

import (
	"bytes"
	"fmt"
)

///////////////////////////////////////

// max length = 255
type StringMax255List255 struct {
	Count VarUint1
	Lists []StringMax255
}

func CreateEmptyStringMax255List255() StringMax255List255 {
	return StringMax255List255{
		Count: VarUint1(0),
		Lists: make([]StringMax255, 0),
	}
}

func (this *StringMax255List255) Append(str *StringMax255) error {
	var num = int(this.Count)
	if num >= 255 {
		return fmt.Errorf("StringMax255List255 size overflow")
	}
	this.Count += 1
	this.Lists = append(this.Lists, *str)
	return nil
}

func (this *StringMax255List255) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	bt, _ := this.Count.Serialize()
	buffer.Write(bt)
	for i := 0; i < int(this.Count); i++ {
		bt, _ := this.Lists[i].Serialize()
		buffer.Write(bt)
	}
	return buffer.Bytes(), nil
}

func (this *StringMax255List255) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	if seek >= uint32(len(buf)) {
		return 0, fmt.Errorf("[Sign.Parse] seek out of buf len.")
	}
	this.Count = VarUint1(0)
	seek, e = this.Count.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	this.Lists = make([]StringMax255, int(this.Count))
	if this.Count == 0 {
		return seek, nil // List is empty
	}
	for i := 0; i < int(this.Count); i++ {
		this.Lists[i] = CreateStringMax255("")
		seek, e = this.Lists[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

func (this *StringMax255List255) Size() uint32 {
	size := this.Count.Size()
	for i := 0; i < int(this.Count); i++ {
		size += this.Lists[i].Size()
	}
	return size
}
