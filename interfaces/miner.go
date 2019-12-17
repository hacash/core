package interfaces

type Miner interface {
	Start() //

	StartMining() //
	StopMining()  //

	SetBlockChain(BlockChain)
	SetTxPool(TxPool)
	SetPowServer(PowServer)

	SubmitTx(Transaction)

	//GetPrevDiamondHash() (uint32, fields.HashWithFee)
	//SetPrevDiamondHash(uint32, fields.HashWithFee)

}
