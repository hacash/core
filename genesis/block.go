package genesis

import (
	"bytes"
	"encoding/hex"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/transactions"
	"sync"
	"time"
)

var (
	genesisBlock *blocks.Block_v1 = nil
)

var genesisblockLock sync.RWMutex

/**
 * 创世区块
 */
func GetGenesisBlock() *blocks.Block_v1 {
	genesisblockLock.Lock()
	defer genesisblockLock.Unlock()

	if genesisBlock != nil {
		return genesisBlock
	}
	genesis := blocks.NewEmptyBlock_v1(nil)
	//loc, _ := time.LoadLocation("Asia/Chongqing")
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	loc_chongqing := time.FixedZone("Asia/Chongqing", secondsEastOfUTC)
	//fmt.Println(time.Now().In(loc))
	ttt := time.Date(2019, time.February, 4, 11, 25, 0, 0, loc_chongqing).Unix()
	//fmt.Println( ttt )
	genesis.Timestamp = fields.BlockTxTimestamp(ttt)
	genesis.Nonce = fields.VarUint4(160117829)
	genesis.Difficulty = fields.VarUint4(0)
	// coinbase
	addrreadble := "1271438866CSDpJUqrnchoJAiGGBFSQhjd"
	addr, e0 := fields.CheckReadableAddress(addrreadble)
	if e0 != nil {
		panic(e0)
	}
	reward := fields.NewAmountNumSmallCoin(1)
	coinbase := transactions.NewTransaction_0_CoinbaseV0()
	coinbase.Address = *addr
	coinbase.Reward = *reward
	coinbase.Message = "hardertodobetter"
	genesis.TransactionCount = 1
	genesis.Transactions = make([]interfacev2.Transaction, 1)
	genesis.Transactions[0] = coinbase
	root := blocks.CalculateMrklRoot(genesis.GetTransactions())
	//fmt.Println( hex.EncodeToString(root) )
	genesis.SetMrklRoot(root)
	hash := genesis.HashFresh()
	check_hash := "000000077790ba2fcdeaef4a4299d9b667135bac577ce204dee8388f1b97f7e6"
	check, _ := hex.DecodeString(check_hash)
	if 0 != bytes.Compare(hash, check) {
		panic("Genesis Block HashWithFee Error: need " + check_hash + ", but give " + hex.EncodeToString(hash))
	}
	genesisBlock = genesis
	//bbb, _ := genesisBlock.Serialize()
	//fmt.Println( hex.EncodeToString(bbb) )
	return genesisBlock
}
