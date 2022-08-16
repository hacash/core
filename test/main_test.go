package test

import (
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/genesis"
	"github.com/hacash/core/sys/inicnf"
	"os"
	"testing"
)

func Test_t187362873(t *testing.T) {
	genesis.GetGenesisBlock()

	//hx := x16rs.CalculateBlockHash(12, []byte{0,0,0,0})
	//fmt.Println(hex.EncodeToString(hx))

}

func Test_t1(t *testing.T) {

	testcnffn := "/home/shiqiujie/Desktop/Hacash/go/src/github.com/hacash/core/test/config.ini"

	cnf, e := inicnf.Load(testcnffn)
	if e != nil {
		fmt.Println(e)
		return
	}

	data_dir := cnf.Section("").Key("data_dir").MustString("~/.hacash_mainnet")

	fmt.Println(data_dir)

	fmt.Println(os.Getenv("HOME"))

}

// Test 1bitcoinmovetohacashneverback
func Test_1BitcoinMoveToHacashNeverBack(t *testing.T) {
	// const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	//                        1Aq1iP9ddFSnsj3PyTuGA2N2fZtdvxVSbc
	//bts, _ := b58decode("1Aq1iP9ddFSnsj3PyTuGA2N2fZtdvxVSbc")
	//fmt.Println(len(bts), bts)
	//addr := fields.Address(bts)
	//fmt.Println(addr.ToReadable())
	//                           1RuinBtcToHacashNeverBack8879XQar
	//
	//              base58check.DecodeTestPrint("1RuinBtcToHacashNeverBack8879XQar")
	//                                           1BGxnkVmNQECbvBYbvdquFQZBLgBpVLkYr
	//                                           1Lycj6rjsJQL6xAseaeCxAoPBLpjNtuiKB
	account.Base58CheckDecodeTestPrint("1RuinBtcToHacashNeverBack8879XQar")

}
