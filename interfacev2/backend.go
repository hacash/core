package interfacev2

type Backend interface {
	BlockChain() BlockChain
	AllPeersDescribe() string
}
