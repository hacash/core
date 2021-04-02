package interfaces

type PowServer interface {
	Excavate(input Block, resCh chan Block) // find block nonce or change coinbase message
	StopMining()                            // stop all
}

// 挖矿投喂单元
type PowWorkerMiningStuffItem interface {

	// 设置或获取挖矿是否成功
	SetMiningSuccessed(bool)
	GetMiningSuccessed() bool

	// 拷贝并设置挖矿nonce
	CopyForMiningByRandomSetCoinbaseNonce() PowWorkerMiningStuffItem

	GetHeadMetaBlock() Block
	GetCoinbaseNonce() []byte
	GetHeadNonce() []byte
	SetHeadNonce(nonce []byte)
}

// 执行端
type PowWorker interface {
	InitStart() error              // 初始化
	CloseUploadHashrate()          // 关闭算力统计
	SetPowDevice(device PowDevice) // 设置挖矿设备端
	Excavate(miningStuffCh chan PowWorkerMiningStuffItem, resultCh chan PowWorkerMiningStuffItem)
	DoNextMining(nextheight uint64) // to do next
	StopAllMining()                 // stop all
}

// 设备端
type PowDevice interface {
	Init() error                                                                                                                               // 初始化
	CloseUploadHashrate()                                                                                                                      // 关闭算力统计
	GetSuperveneWide() int                                                                                                                     // 并发数
	DoMining(blockHeight uint64, reporthashrate bool, stopmark *byte, tarhashvalue []byte, blockheadmeta [][]byte) (bool, int, []byte, []byte) // 执行一次挖矿
}
