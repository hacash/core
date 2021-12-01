package interfacev3

import "github.com/hacash/core/interfaces"

type Backend interface {
	BlockChain() interfaces.BlockChain
	AllPeersDescribe() string
}
