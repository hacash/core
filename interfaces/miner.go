package interfaces

type Miner interface {
	Start() //
	Stop()  //

	SetBlockChain(BlockChain)

	SubmitTx(Transaction)

	//GetPrevDiamondHash() (uint32, fields.HashWithFee)
	//SetPrevDiamondHash(uint32, fields.HashWithFee)

}
