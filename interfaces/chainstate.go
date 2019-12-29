package interfaces

type ChainState interface {
	ChainStateOperation

	Fork() (ChainState, error)
	Destory()
}
