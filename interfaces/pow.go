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
	InitStart() error  // 初始化
	CloseUploadPower() // 关闭算力统计
	Excavate(miningStuffCh chan PowWorkerMiningStuffItem, resultCh chan PowWorkerMiningStuffItem)
	NextMining(nextheight uint64) // to do next
	StopAllMining()               // stop all
}
