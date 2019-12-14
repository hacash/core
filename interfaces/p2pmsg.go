package interfaces

type MsgPeer interface {
	AddKnowledge(KnowledgeKey string, KnowledgeValue string) bool
	SendDataMsg(msgty uint16, msgbody []byte)
	Describe() string
	Disconnect()
}

type MsgCommunicator interface {
	PeerLen() int
	FindAnyOnePeerBetterBePublic() MsgPeer
	BroadcastMessageToUnawarePeers(ty uint16, msgbody []byte, KnowledgeKey string, KnowledgeValue string)
}

type MsgDataHandler interface {
	OnConnected(MsgCommunicator, MsgPeer)
	OnMsgData(mc MsgCommunicator, p MsgPeer, msgty uint16, msgbody []byte)
	OnDisconnected(MsgPeer)
}
