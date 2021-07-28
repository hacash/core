package transactions

import (
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"testing"
	"time"
)

// 钻石借贷交易
func Test_diamond_lending(t *testing.T) {

	hash14, _ := hex.DecodeString("530dd68299cf6d2bd68299cf6d2b")

	// tx
	feeamt, _ := fields.NewAmountFromFinString("ㄜ1:246")
	mainaddr, _ := fields.CheckReadableAddress("1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9")
	tx, _ := NewEmptyTransaction_2_Simple(*mainaddr)
	tx.Fee = *feeamt
	tx.Timestamp = 1618839281

	// 创建钻石
	amt1 := fields.NewAmountSmall(16, 248)
	act1 := actions.Action_15_DiamondsSystemLendingCreate{
		LendingID: hash14,
		MortgageDiamondList: fields.DiamondListMaxLen200{
			Count:    2,
			Diamonds: []fields.DiamondName{[]byte("XXXYYY"), []byte("WWWMMM")},
		},
		LoanTotalAmount: *amt1,
		BorrowPeriod:    20,
	}
	tx.AppendAction(&act1)

	// 签名
	feeacc := account.CreateAccountByPassword("123456")
	addrPrivateKeys := map[string][]byte{}
	addrPrivateKeys[string(feeacc.Address)] = feeacc.PrivateKey
	tx.FillNeedSigns(addrPrivateKeys, nil)

	// 序列化
	txbody, _ := tx.Serialize()
	fmt.Println("tx body:", hex.EncodeToString(txbody))

}

// 创建钻石交易
func Test_create_diamond(t *testing.T) {

	hash8, _ := hex.DecodeString("530dd68299cf6d2b")
	hash32, _ := hex.DecodeString("000000000e8ca4376218601120e12b6724a8c174087b9614530dd68299cf6d2b")

	// tx
	feeamt, _ := fields.NewAmountFromFinString("ㄜ1:246")
	mainaddr, _ := fields.CheckReadableAddress("1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9")
	tx, _ := NewEmptyTransaction_2_Simple(*mainaddr)
	tx.Fee = *feeamt
	tx.Timestamp = 1618839281

	// 创建钻石
	act1 := actions.Action_4_DiamondCreate{
		Diamond:       fields.DiamondName("WWWMMM"),
		Number:        3,
		PrevHash:      hash32,
		Nonce:         hash8,
		Address:       *mainaddr,
		CustomMessage: hash32,
	}
	tx.AppendAction(&act1)

	// 签名
	feeacc := account.CreateAccountByPassword("123456")
	addrPrivateKeys := map[string][]byte{}
	addrPrivateKeys[string(feeacc.Address)] = feeacc.PrivateKey
	tx.FillNeedSigns(addrPrivateKeys, nil)

	// 序列化
	txbody, _ := tx.Serialize()
	fmt.Println("tx body:", hex.EncodeToString(txbody))

}

func Test_alltx(t *testing.T) {

	// 全类别交易 测试
	feeamt, _ := fields.NewAmountFromFinString("ㄜ1234:246")
	mainaddr, _ := fields.CheckReadableAddress("1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9")
	tx, _ := NewEmptyTransaction_2_Simple(*mainaddr)
	tx.Fee = *feeamt
	tx.Timestamp = 1618839281

	// 1 普通转账
	addr1, _ := fields.CheckReadableAddress("1AVRuFXNFi3rdMrPH4hdqSgFrEBnWisWaS")
	amt1, _ := fields.NewAmountFromFinString("ㄜ500:248")
	act1 := actions.NewAction_1_SimpleToTransfer(*addr1, amt1)
	tx.AppendAction(act1)

	// 2 开启通道
	addr2_1, _ := fields.CheckReadableAddress("1EDUeK8NAjrgYhgDFv9NJecn8dNyJJsu3y")
	addr2_2, _ := fields.CheckReadableAddress("1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9")
	channelid, _ := hex.DecodeString("277095b321f3ffe7e80f3dd328e2f338")
	amt2_1, _ := fields.NewAmountFromFinString("ㄜ500:248")
	amt2_2, _ := fields.NewAmountFromFinString("ㄜ239:248")
	act2 := actions.Action_2_OpenPaymentChannel{
		ChannelId:    channelid,
		LeftAddress:  *addr2_1,
		LeftAmount:   *amt2_1,
		RightAddress: *addr2_2,
		RightAmount:  *amt2_2,
	}
	tx.AppendAction(&act2)

	// 3 关闭通道

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
