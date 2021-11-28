package interfacev3

type Backend interface {
	BlockChain() BlockChain
	AllPeersDescribe() string
}
