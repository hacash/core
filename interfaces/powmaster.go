package interfaces

type PowMaster interface {
	Excavate(Block, chan Block) // find block nonce or change coinbase message

	StopMining() // stop all

}
