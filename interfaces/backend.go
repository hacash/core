package interfaces

type Backend interface {
	BlockChain() BlockChain
	AllPeersDescribe() string
}
