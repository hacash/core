package fields

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/account"
)

//////////////////////////////////////////////////////////////////

const (
	SignSize uint32 = 33 + 64
)

type Sign struct {
	PublicKey Bytes33
	Signature Bytes64
}

func (this *Sign) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.Write(this.PublicKey)
	buffer.Write(this.Signature)
	return buffer.Bytes(), nil
}

func (this *Sign) Parse(buf []byte, seek uint32) (uint32, error) {
	seek, _ = this.PublicKey.Parse(buf, seek)
	seek, _ = this.Signature.Parse(buf, seek)
	return seek, nil
}

func (this *Sign) Size() uint32 {
	return this.PublicKey.Size() + this.Signature.Size()
}

func (this *Sign) GetAddress() Address {
	return account.NewAddressFromPublicKeyV0(this.PublicKey)
}

func CreateEmptySign() Sign {
	b1 := bytes.Repeat([]byte{0}, 33)
	b2 := bytes.Repeat([]byte{0}, 64)
	return Sign{
		PublicKey: b1,
		Signature: b2,
	}
}

/********************************/

type SignListMax255 struct {
	Count VarUint1
	Signs []Sign
}

func CreateEmptySignListMax255() *SignListMax255 {
	return &SignListMax255{
		Count: 0,
		Signs: make([]Sign, 0),
	}
}

func (this *SignListMax255) Append(sign Sign) {
	this.Count += 1
	this.Signs = append(this.Signs, sign)
}

func (this *SignListMax255) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	bt, _ := this.Count.Serialize()
	buffer.Write(bt)
	for i := 0; i < int(this.Count); i++ {
		bt, _ := this.Signs[i].Serialize()
		buffer.Write(bt)
	}
	return buffer.Bytes(), nil
}

func (this *SignListMax255) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	if seek >= uint32(len(buf)) {
		return 0, fmt.Errorf("[Sign.Parse] seek out of buf len.")
	}
	this.Count = VarUint1(0)
	seek, e = this.Count.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if this.Count == 0 {
		return seek, nil // List is empty
	}
	this.Signs = make([]Sign, int(this.Count))
	for i := 0; i < int(this.Count); i++ {
		this.Signs[i] = Sign{}
		seek, e = this.Signs[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

func (this *SignListMax255) Size() uint32 {
	size := this.Count.Size()
	for i := 0; i < int(this.Count); i++ {
		size += this.Signs[i].Size()
	}
	return size
}

type SignListMax65535 struct {
	Count VarUint2
	Signs []Sign
}

func (this *SignListMax65535) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	bt, _ := this.Count.Serialize()
	buffer.Write(bt)
	for i := 0; i < int(this.Count); i++ {
		bt, _ := this.Signs[i].Serialize()
		buffer.Write(bt)
	}
	return buffer.Bytes(), nil
}

func (this *SignListMax65535) Parse(buf []byte, seek uint32) (uint32, error) {
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
	this.Signs = make([]Sign, int(this.Count))
	for i := 0; i < int(this.Count); i++ {
		this.Signs[i] = Sign{}
		seek, e = this.Signs[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

func (this *SignListMax65535) Size() uint32 {
	size := this.Count.Size()
	for i := 0; i < int(this.Count); i++ {
		size += this.Signs[i].Size()
	}
	return size
}

//////////////////////////////////////////////////////////////////

type Multisign2 struct {
	CondElem          uint8      // molecule
	CondBase          uint8      // denominator
	SignatureInds     []VarUint2 // Signature location
	BasePublicKeyInds []VarUint2 // Public key base location
}

type Multisign struct {
	CondElem      uint8 // molecule
	CondBase      uint8 // denominator
	PublicKeyList []Bytes33
	SignatureInds []uint8
	SignatureList []Bytes64
}

func (this *Multisign) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.Write([]byte{this.CondElem, this.CondElem})
	length1 := int(this.CondElem)
	length2 := int(this.CondBase)
	for i := 0; i < length2; i++ {
		buffer.Write(this.PublicKeyList[i])
	}
	for j := 0; j < length1; j++ {
		buffer.Write([]byte{this.SignatureInds[j]})
	}
	for k := 0; k < length1; k++ {
		buffer.Write(this.PublicKeyList[k])
	}
	return buffer.Bytes(), nil
}

func (this *Multisign) Parse(buf []byte, seek uint32) (uint32, error) {
	if int(seek)+2 > len(buf) {
		return 0, fmt.Errorf("buf len too short.")
	}
	this.CondElem = buf[seek]
	this.CondBase = buf[seek+1]
	seek = seek + 2
	length1 := int(this.CondElem)
	length2 := int(this.CondBase)
	this.PublicKeyList = make([]Bytes33, length2)
	this.SignatureInds = make([]uint8, length1)
	this.SignatureList = make([]Bytes64, length1)
	var e error
	for i := 0; i < length2; i++ {
		var b Bytes33
		seek, e = b.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		this.PublicKeyList[i] = b
		seek += b.Size()
	}
	for i := 0; i < length1; i++ {
		if int(seek) >= len(buf) {
			return 0, fmt.Errorf("buf len too short.")
		}
		this.SignatureInds[i] = buf[seek]
		seek += 1
	}
	for i := 0; i < length1; i++ {
		var b Bytes64
		seek, e = b.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
		this.SignatureList[i] = b
		seek += b.Size()
	}
	return seek, nil
}

func (this *Multisign) Size() uint32 {
	length := uint32(this.CondBase)
	return 1 + 1 + length*33 + length*1 + length*64
}
