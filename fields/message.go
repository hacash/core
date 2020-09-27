package fields

import "fmt"

type ExtendMessageMaxLen255 struct {
	Count   VarUint1
	Message []byte
}

func (e *ExtendMessageMaxLen255) Size() uint32 {
	return 1 + uint32(len(e.Message))
}

func (e *ExtendMessageMaxLen255) Serialize() ([]byte, error) {
	bts := make([]byte, 0, e.Size())
	bts = append(bts, byte(e.Count))
	bts = append(bts, e.Message...)
	return bts, nil
}

func (e *ExtendMessageMaxLen255) Parse(buf []byte, seek uint32) (uint32, error) {
	e.Count = VarUint1(buf[int(seek)])
	seek++
	start := seek
	end := start + uint32(e.Count)
	if len(buf) < int(e.Count)-1 {
		return 0, fmt.Errorf("buf is too short.")
	}
	e.Message = buf[start:end]
	return end, nil
}
