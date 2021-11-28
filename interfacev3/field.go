package interfacev3

type Field interface {
	// assembling datas
	Size() uint32
	Serialize() ([]byte, error)
	Parse([]byte, uint32) (uint32, error)
}
