package interfacev2

type StoreItem interface {

	// assembling datas
	Serialize() ([]byte, error)
	Parse([]byte, uint32) (uint32, error)
	Size() uint32
}
