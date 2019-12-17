package interfaces

const (
	PowMasterResultsReturnStatusContinue uint8 = 0
	PowMasterResultsReturnStatusSuccess  uint8 = 1
	PowMasterResultsReturnStatusStop     uint8 = 2
	PowMasterResultsReturnStatusError    uint8 = 3
)

type PowMasterResultsReturn struct {
	Status         uint8 //
	CoinbaseMsgNum uint32
	NonceBytes     []byte
	BlockHeadMeta  Block
}

type PowMaster interface {
	SetCoinbaseMsgNum(uint32)

	Excavate(headmeta Block, resCh chan PowMasterResultsReturn) // find block nonce or change coinbase message

	StopMining() // stop all

}

type PowServer interface {
	Excavate(input Block, resCh chan Block) // find block nonce or change coinbase message

	StopMining() // stop all

}
