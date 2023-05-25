package fields

import (
	"bytes"
	"fmt"
)

type HashListMax65535 struct {
	Count VarUint2
	Hashs []Hash
}

func CreateHashListMax65535(hxs []Hash) HashListMax65535 {
	if hxs == nil {
		return HashListMax65535{
			Count: 0,
			Hashs: nil,
		}
	}
	count := len(hxs)
	return HashListMax65535{
		Count: VarUint2(count),
		Hashs: hxs,
	}

}

func (this *HashListMax65535) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	bt, _ := this.Count.Serialize()
	buffer.Write(bt)
	for i := 0; i < int(this.Count); i++ {
		bt, _ := this.Hashs[i].Serialize()
		buffer.Write(bt)
	}
	return buffer.Bytes(), nil
}

func (this *HashListMax65535) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	if seek >= uint32(len(buf)) {
		return 0, fmt.Errorf("[Sign.Parse] seek out of buf len.")
	}
	this.Count = VarUint2(0)
	seek, e = this.Count.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if this.Count == 0 {
		return seek, nil // List is empty
	}
	this.Hashs = make([]Hash, int(this.Count))
	for i := 0; i < int(this.Count); i++ {
		this.Hashs[i] = make([]byte, 32)
		seek, e = this.Hashs[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

func (this *HashListMax65535) Size() uint32 {
	size := this.Count.Size()
	for i := 0; i < int(this.Count); i++ {
		size += this.Hashs[i].Size()
	}
	return size
}
