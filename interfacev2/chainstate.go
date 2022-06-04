package interfacev2

type ChainState interface {
	ChainStateOperation

	Fork() (ChainState, error)
	Close()   // close
	Destory() // Destruction, including deletion of all document stores

}
