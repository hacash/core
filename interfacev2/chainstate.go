package interfacev2

type ChainState interface {
	ChainStateOperation

	Fork() (ChainState, error)
	Close()   // 关闭
	Destory() // 销毁，包括删除所有文件储存

}
