package interfaces

type PowServer interface {
	Excavate(input Block, resCh chan Block) // find block nonce or change coinbase message

	StopMining() // stop all

}
