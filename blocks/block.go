package blocks

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/crypto/sha3"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"

	"github.com/hacash/x16rs"
)

const (
	BlockHeadSize   = 1 + 5 + 5 + 32 + 32 + 4 // = 79
	BlockMetaSizeV1 = 4 + 4 + 2               // = 10
)

// protocol
const (
	BlockVersion    = fields.VarUint1(1)  // uint8
	TransactionType = fields.VarUint1(2)  // uint8
	ActionKind      = fields.VarUint2(11) // uint16
	RepairVersion   = fields.VarUint2(2)  // uint16
)

////////////////////////////////////////////////////////////

func NewBlockByVersion(ty uint8) (interfaces.Block, error) {
	switch ty {
	////////////////////  BLOCK  ////////////////////
	case 1:
		return new(Block_v1), nil
		////////////////////   END   ////////////////////
	}
	return nil, fmt.Errorf("Cannot find Block type of " + string(ty))
}

func ParseBlock(buf []byte, seek uint32) (interfaces.Block, uint32, error) {
	if len(buf) < 1 {
		return nil, 0, fmt.Errorf("buf too short")
	}
	version := uint8(buf[seek])
	var blk, e = NewBlockByVersion(version)
	if e != nil {
		return nil, 0, e
	}
	var mv, err = blk.Parse(buf, seek+1)
	return blk, mv, err
}

func ParseBlockHead(buf []byte, seek uint32) (interfaces.Block, uint32, error) {
	version := uint8(buf[seek])
	var blk, ee = NewBlockByVersion(version)
	if ee != nil {
		fmt.Println("Block not Find. Version:", version)
	}
	var mv, err = blk.ParseHead(buf, seek+1)
	return blk, mv, err
}
func ParseExcludeTransactions(buf []byte, seek uint32) (interfaces.Block, uint32, error) {
	version := uint8(buf[seek])
	var blk, ee = NewBlockByVersion(version)
	if ee != nil {
		fmt.Println("Block not Find. Version:", version)
	}
	var mv, err = blk.ParseExcludeTransactions(buf, seek+1)
	return blk, mv, err
}

//////////////////////////////////

func CalculateBlockHash(block interfaces.Block) fields.Hash {
	stuff := CalculateBlockHashBaseStuff(block)
	return x16rs.CalculateBlockHash(block.GetHeight(), stuff)
}

/*
func CalculateBlockHashByStuff(loopnum int, stuff []byte) fields.Hash {
	hashbase := sha3.Sum256(stuff)
	return x16rs.HashX16RS_Optimize(loopnum, hashbase[:])
}
*/

func CalculateBlockHashBaseStuff(block interfaces.Block) []byte {
	var buffer bytes.Buffer
	head, _ := block.SerializeHead()
	meta, _ := block.SerializeMeta()
	buffer.Write(head)
	buffer.Write(meta)
	return buffer.Bytes()
}

func CalculateMrklRoot(transactions []interfaces.Transaction) fields.Hash {
	trslen := len(transactions)
	if trslen == 0 {
		return fields.EmptyZeroBytes32
	}
	hashs := make([]fields.Hash, trslen)
	for i := 0; i < trslen; i++ {
		hashs[i] = transactions[i].HashWithFee()
	}
	for true {
		if len(hashs) == 1 {
			return hashs[0]
		}
		hashs = hashMerge(hashs) // 两两归并
	}
	return nil
}

func hashMerge(hashs []fields.Hash) []fields.Hash {
	length := len(hashs)
	mgsize := length / 2
	if length%2 == 1 {
		mgsize = (length + 1) / 2
	}
	var mergehashs = make([]fields.Hash, mgsize)
	for m := 0; m < length; m += 2 {
		var buf bytes.Buffer
		h1 := hashs[m]
		buf.Write(hashs[m])
		if m+1 < length {
			h2 := hashs[m+1]
			buf.Write(h2)
		} else {
			buf.Write(h1) // repeat h1
		}
		digest := sha3.Sum256(buf.Bytes())
		mergehashs[m/2] = digest[:]
	}
	return mergehashs
}

// 通过修改 coinbase tx 哈希 来重新计算默克尔根
func CalculateMrklRootByCoinbaseTxModify(coinbasetxhx fields.Hash, mdftree []fields.Hash) fields.Hash {
	mrklroot := append([]byte{}, coinbasetxhx...)
	for i := 0; i < len(mdftree); i++ {
		mghx := hashMerge([]fields.Hash{mrklroot, mdftree[i]})
		mrklroot = mghx[0]
	}
	return mrklroot
}

// 计算并获得与 coinbase tx 修改有关的默克尔相关哈希列表
func PickMrklListForCoinbaseTxModify(transactions []interfaces.Transaction) []fields.Hash {
	hxlist := make([]fields.Hash, 0)
	trslen := len(transactions)
	if trslen == 0 {
		panic("len(transactions) must not be empty") // 不能为空
	}
	if trslen == 1 {
		return hxlist
	}
	if trslen == 2 {
		hxlist = append(hxlist, transactions[1].HashWithFee())
		return hxlist
	}
	// 计算哈希关系树
	hashs := make([]fields.Hash, trslen)
	for i := 0; i < trslen; i++ {
		hashs[i] = transactions[i].HashWithFee()
	}
	for true {
		lhx := len(hashs)
		if lhx == 1 {
			break
		}
		if lhx >= 2 {
			hxlist = append(hxlist, hashs[1]) // 获取关系哈希
		}
		hashs = hashMerge(hashs) // 两两归并
	}
	// ok
	return hxlist

}
