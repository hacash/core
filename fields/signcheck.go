package fields

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/account"
)

// Signature verification data
type SignCheckData struct {
	Signdata Sign
	Stuffstr StringMax65535
}

func CreateSignCheckData(str string) SignCheckData {
	l := len(str)
	if l > 65535 {
		l = 65535
	}
	return SignCheckData{
		Signdata: CreateEmptySign(),
		Stuffstr: CreateStringMax65535(str),
	}
}

// Fill in signature
func (elm *SignCheckData) FillSign(acc *account.Account) error {
	hash := CalculateHash([]byte(elm.Stuffstr.Value()))
	signature, e := acc.Private.Sign(hash)
	if e != nil {
		return e
	}
	elm.Signdata = Sign{
		PublicKey: acc.PublicKey,
		Signature: signature.Serialize64(),
	}
	// fill sign success
	return nil
}

func (elm SignCheckData) Size() uint32 {
	return elm.Signdata.Size() + elm.Stuffstr.Size()
}

func (elm SignCheckData) Serialize() ([]byte, error) {
	var buf = bytes.NewBuffer(nil)
	bt1, _ := elm.Signdata.Serialize()
	buf.Write(bt1)
	bt2, _ := elm.Stuffstr.Serialize()
	buf.Write(bt2)
	return buf.Bytes(), nil
}

func (elm *SignCheckData) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	if seek >= uint32(len(buf)) {
		return 0, fmt.Errorf("[SignCheckData.Parse] seek out of buf len.")
	}
	seek, e = elm.Signdata.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.Stuffstr.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

// check sign
func (elm SignCheckData) VerifySign() (bool, Address, error) {
	hash := CalculateHash([]byte(elm.Stuffstr.Value()))
	addr := elm.Signdata.GetAddress()
	ckok, e := account.CheckSignByHash32(hash, elm.Signdata.PublicKey, elm.Signdata.Signature)
	if e != nil {
		return false, nil, e
	}
	return ckok, addr, e
}
