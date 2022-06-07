package fields

import (
	"bytes"
	"fmt"
	"strings"
)

type TrimString16 string
type TrimString34 string
type TrimString64 string

////////////////////////////////////////////////////////

func (elm TrimString16) Serialize() ([]byte, error) { return trimStringSerialize(string(elm), 16) }
func (elm TrimString34) Serialize() ([]byte, error) { return trimStringSerialize(string(elm), 34) }
func (elm TrimString64) Serialize() ([]byte, error) { return trimStringSerialize(string(elm), 64) }

func (elm *TrimString16) Parse(buf []byte, seek uint32) (uint32, error) {
	return trimStringParse(elm, buf, seek, 16)
}
func (elm *TrimString34) Parse(buf []byte, seek uint32) (uint32, error) {
	return trimStringParse(elm, buf, seek, 34)
}
func (elm *TrimString64) Parse(buf []byte, seek uint32) (uint32, error) {
	return trimStringParse(elm, buf, seek, 64)
}

func (elm TrimString16) Size() uint32 { return 16 }
func (elm TrimString34) Size() uint32 { return 34 }
func (elm TrimString64) Size() uint32 { return 64 }

func (elm TrimString16) ValueShow() string { return valueShow(string(elm)) }
func (elm TrimString34) ValueShow() string { return valueShow(string(elm)) }
func (elm TrimString64) ValueShow() string { return valueShow(string(elm)) }

func valueShow(str string) string {
	str = strings.TrimRight(str, " ")
	msg := ""
	for _, v := range []byte(str) {
		if v < 32 || v > 126 {
			break
		}
		if v == 34 { // Process conversion of double quotation marks to single quotation marks
			v = 39
		}
		msg += string([]byte{v})
	}
	return msg
}

////////////////////////////////////////////////////////

func trimStringParse(elm interface{}, buf []byte, seek uint32, maxlen uint32) (uint32, error) {
	if seek+maxlen > uint32(len(buf)) {
		return 0, fmt.Errorf("[trimStringParse] seek out of buf len.")
	}
	var nnnold = buf[seek : seek+maxlen]
	var addrbytes = make([]byte, len(nnnold))
	copy(addrbytes, nnnold)
	addrbytes = bytes.TrimRight(addrbytes, " ")
	var sd = string(addrbytes)
	switch a := elm.(type) {
	case *TrimString16:
		*a = (TrimString16)(sd)
	case *TrimString34:
		*a = (TrimString34)(sd)
	case *TrimString64:
		*a = (TrimString64)(sd)
	default:
		return 0, fmt.Errorf("not find type")
	}
	return seek + maxlen, nil
}

func trimStringSerialize(str string, maxlen int) ([]byte, error) {
	//var str = string(*elm)
	//fmt.Println("trimStringSerialize ---------", str, "===")
	// Space filling
	newbts := bytes.Repeat([]byte{byte(' ')}, maxlen)
	copy(newbts, str) // Copy by length
	return newbts, nil
}
