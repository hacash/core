package interfaces

type Miner interface {
	Start() //

	StartMining() //
	StopMining()  //

	SetBlockChain(BlockChain)
	SetTxPool(TxPool)
	SetPowMaster(PowMaster)

	SubmitTx(Transaction)

	//GetPrevDiamondHash() (uint32, fields.HashWithFee)
	//SetPrevDiamondHash(uint32, fields.HashWithFee)

}
