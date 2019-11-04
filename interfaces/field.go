package interfaces

type Field interface {

	// assembling datas
	Serialize() ([]byte, error)
	Parse([]byte, uint32) (uint32, error)
	Size() uint32

}
