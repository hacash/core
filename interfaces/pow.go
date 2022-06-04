package interfaces

type PowServer interface {
	Excavate(input Block, resCh chan Block) // find block nonce or change coinbase message
	StopMining()                            // stop all
}

// Mining feeding unit
type PowWorkerMiningStuffItem interface {

	// Set or obtain whether mining is successful
	SetMiningSuccessed(bool)
	GetMiningSuccessed() bool

	// Copy and set mining nonce
	CopyForMiningByRandomSetCoinbaseNonce() PowWorkerMiningStuffItem

	GetHeadMetaBlock() Block
	GetCoinbaseNonce() []byte
	GetHeadNonce() []byte
	SetHeadNonce(nonce []byte)
}

// Executive end
type PowWorker interface {
	InitStart() error              // initialization
	CloseUploadHashrate()          // Turn off force statistics
	SetPowDevice(device PowDevice) // Set mining equipment end
	Excavate(miningStuffCh chan PowWorkerMiningStuffItem, resultCh chan PowWorkerMiningStuffItem)
	DoNextMining(nextheight uint64) // to do next
	StopAllMining()                 // stop all
}

// Equipment end
type PowDevice interface {
	Init() error           // initialization
	CloseUploadHashrate()  // Turn off force statistics
	GetSuperveneWide() int // Concurrent number
	// Perform a mining operation
	DoMining(blockHeight uint64, reporthashrate bool, stopmark *byte, tarhashvalue []byte, blockheadmeta [][]byte) (bool, int, []byte, []byte)
}
