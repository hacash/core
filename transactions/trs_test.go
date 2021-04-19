package transactions

import (
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"testing"
	"time"
)

func Test_alltx(t *testing.T) {

	// 全类别交易 测试
	feeamt, _ := fields.NewAmountFromFinString("ㄜ1234:246")
	mainaddr, _ := fields.CheckReadableAddress("1AVRuFXNFi3rdMrPH4hdqSgFrEBnWisWaS")
	tx, _ := NewEmptyTransaction_2_Simple(*mainaddr)
	tx.Fee = *feeamt
	tx.Timestamp = 1618839281

	// 1 普通转账
	addr1, _ := fields.CheckReadableAddress("1EDUeK8NAjrgYhgDFv9NJecn8dNyJJsu3y")
	amt1, _ := fields.NewAmountFromFinString("ㄜ500:248")
	act1 := actions.NewAction_1_SimpleTransfer(*addr1, amt1)
	tx.AppendAction(act1)

	// 打印
	fmt.Println("tx hash:", tx.Hash().ToHex())
	fmt.Println("tx hash with fee:", tx.HashWithFee().ToHex())
	txtime := time.Unix(int64(tx.Timestamp), 0)
	fmt.Println("tx timestamp:", tx.GetTimestamp(), txtime.Format("2006-01-02 15:04:05"))
	fmt.Println("--------")
	txbody, _ := tx.Serialize()
	fmt.Println("tx body:", hex.EncodeToString(txbody))
	fmt.Println("--------")

}

func Test_coinbaseCopy(t *testing.T) {

	cbtrs := NewTransaction_0_CoinbaseV0()
	addr, _ := fields.CheckReadableAddress("1AVRuFXNFi3rdMrPH4hdqSgFrEBnWisWaS")
	cbtrs.Address = *addr
	reward := fields.NewAmountNumSmallCoin(1)
	cbtrs.Reward = *reward
	cbtrs.Message = "ABC123"

	fmt.Println(cbtrs.Serialize())

	clonetrs := cbtrs.Copy()

	fmt.Println(clonetrs.Serialize())

}
