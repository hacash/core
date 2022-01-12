package fields

const (
	DiamondNameSize = 6
)

type DiamondName = Bytes6
type DiamondNumber = VarUint3

func (d DiamondName) Name() string {
	return string(d)
}
