package interfaces

type P2PManager interface {
	Start() error
	SetMsgHandler(P2PMsgDataHandler)
	AddKnowledge(KnowledgeKind string, KnowledgeKey string) bool   // Return false as known
	CheckKnowledge(KnowledgeKind string, KnowledgeKey string) bool // Return true to know
	GetConfigOfBootNodeFastSync() bool
}

type P2PMsgPeer interface {
	AddKnowledge(KnowledgeKind string, KnowledgeKey string) bool   // Return false as known
	CheckKnowledge(KnowledgeKind string, KnowledgeKey string) bool // Return true to know
	SendDataMsg(msgty uint16, msgbody []byte) error
	Describe() string
	Disconnect()
}

type P2PMsgCommunicator interface {
	PeerLen() int
	GetAllPeers() []P2PMsgPeer
	FindAnyOnePeerBetterBePublic() P2PMsgPeer
	BroadcastDataMessageToUnawarePeers(ty uint16, msgbody []byte, KnowledgeKind string, KnowledgeKey string)
}

type P2PMsgDataHandler interface {
	OnConnected(P2PMsgCommunicator, P2PMsgPeer)
	OnMsgData(mc P2PMsgCommunicator, p P2PMsgPeer, msgty uint16, msgbody []byte)
	OnDisconnected(P2PMsgPeer)
}
