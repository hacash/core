package interfacev3

import "github.com/hacash/core/fields"

type ChainState interface {
	ChainStateOperation

	// 获得父级状态
	GetParent() *ChainState
	// 获得所有子状态
	GetChilds() []*ChainState

	// 启动一个子状态
	ForkNextBlock(uint64, fields.Hash) (*ChainState, error)
	ForkSubChild() (*ChainState, error)

	// 获取指向的区块
	GetReferBlock() (uint64, fields.Hash)

	// 销毁
	Destory() // 销毁，包括删除所有子状态、缓存、状态数据等

	// 判断类型
	IsImmutable() bool

	// 保存在磁盘
	ImmutableWriteToDisk() (ChainStateImmutable, error)
}

// 不可变、不可回退的锁定状态数据
type ChainStateImmutable interface {
	ChainState

	Close() // 关闭文件句柄等
}
