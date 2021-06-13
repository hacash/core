package coinbase

import (
	"fmt"
	"github.com/hacash/core/fields"
	"math/big"
)

// 计算比特币赎回所系金额
func CalculationBitcoinSystemLendingRedeemAmount(
	feeAddr fields.Address, lendingMainAddr fields.Address,
	loanTotalAmount *fields.Amount, ransomBlockNumberBase uint64,
	lendingCreateBlockHeight uint64, penddingBlockHeight uint64,
) (uint8, *fields.Amount, error) {

	var e error = nil

	// 赎回期限阶段
	redeemStage := uint8(1) // 1. 私有期  2. 公共期  3. 拍卖期

	// 赎回期基础
	// ransomBlockNumberBase := uint64(100000) // 十万个区块约一年

	// 检查私有赎回期
	privateHeight := uint64(lendingCreateBlockHeight) + ransomBlockNumberBase
	if penddingBlockHeight <= privateHeight && feeAddr.NotEqual(lendingMainAddr) {
		// 未到期之前只能被抵押者私下赎回
		return 0, nil, fmt.Errorf("It can only be redeemed privately by the mortgagor %s before the blockheight %d", lendingMainAddr.ToReadable(), privateHeight)
	}

	// 任何人可以公开赎回
	if penddingBlockHeight > privateHeight {
		redeemStage = 2
	}

	// 赎回金额就是原始借出金额（因为利息被预先支付了）
	var realRansomAmt = loanTotalAmount

	// 检查公共赎回期，计算实时赎回金额
	publicHeight := privateHeight + ransomBlockNumberBase
	if penddingBlockHeight > publicHeight {
		// 大于公共赎回期，开始利息荷兰拍卖模式，用十年（百万区块）将赎回金额降低到0
		maxDown := ransomBlockNumberBase * 10
		redeemStage = 3
		overhei := penddingBlockHeight - publicHeight
		if overhei >= maxDown {
			realRansomAmt = fields.NewEmptyAmount() // 可赎回已降低至0
		} else {
			boli := (float64(maxDown) - float64(overhei)) / float64(maxDown) * loanTotalAmount.ToMei()
			boli *= 100000000
			realRansomAmt, e = fields.NewAmountByBigIntWithUnit(big.NewInt(int64(boli)), 240)
			if e != nil {
				return 0, nil, e
			}
		}
	}

	// 计算成功返回
	return redeemStage, realRansomAmt, nil
}

// 计算钻石赎回所需金额
func CalculationDiamondSystemLendingRedeemAmount(
	feeAddr fields.Address, lendingMainAddr fields.Address,
	borrowPeriod int64, lendingCreateBlockHeight int64,
	loanTotalAmountMei int64,
	dslbpbn int64, penddingBlockHeight int64,
) (uint8, *fields.Amount, error) {

	// 赎回期限阶段
	redeemStage := uint8(1) // 1. 私有期  2. 公共期  3. 拍卖期

	// 赎回期阶段期限
	// 赎回期分为：私有赎回期，公开赎回期，荷兰拍卖期
	ransomBlockNumberBase := borrowPeriod * dslbpbn

	// 检查私有赎回期
	privateHeight := lendingCreateBlockHeight + ransomBlockNumberBase
	if penddingBlockHeight <= privateHeight && feeAddr.NotEqual(lendingMainAddr) {
		// 未到期之前只能被抵押者私下赎回
		return 0, nil, fmt.Errorf("It can only be redeemed privately by the mortgagor %s before the blockheight %d", lendingMainAddr.ToReadable(), privateHeight)
	}

	// 任何人可以公开赎回
	if penddingBlockHeight > privateHeight {
		redeemStage = 2
	}

	// 计算赎回金额（每一个周期表示 0.005 即 0.5% 的利息）
	realRansom1qian := (1000 + (5 * int64(borrowPeriod))) * int64(loanTotalAmountMei)

	// 检查公共赎回期
	// 公共赎回期于私有赎回期持续时间相同，假设 BorrowPeriod = 20 ，则：
	// 私有赎回期 2 年，公共赎回期 2 年， 荷兰拍卖期 2+2=4 年
	publicHeight := privateHeight + ransomBlockNumberBase
	if penddingBlockHeight > publicHeight {
		redeemStage = 3
		// 大于公共赎回期，开始利息荷兰拍卖模式
		subcount := int64((penddingBlockHeight - publicHeight) / dslbpbn)
		maxsub := int64(borrowPeriod) * 2
		if subcount > maxsub {
			redeemStage = 4
			subcount = maxsub // 减扣利息最大只能等于原始支付利息的两倍，即提供正负对称的套利空间
		}
		// 按时间减去利息，拍卖：每35天减扣 0.5%利息
		realRansom1qian -= 5 * subcount * int64(loanTotalAmountMei)
	}

	// 真实有效赎回金额
	validRansomAmt, e1 := fields.NewAmountByBigIntWithUnit(big.NewInt(realRansom1qian), 245)
	if e1 != nil {
		return 0, nil, e1
	}

	// 计算成功，返回
	return redeemStage, validRansomAmt, nil

}

// 比特币抵押借贷： 计算可借数量
// totalLendingPercentage 已经借出的总百分比
// 返回可借数量和预付利息
func CalculationOfInterestBitcoinMortgageLoanAmount(totalLendingPercentage float64) (float64, float64) {

	ttp := totalLendingPercentage

	if ttp < 1 {
		ttp = 1 // 最低为1
	}

	// 可借出数量
	var loanhac float64 = 0
	loanhac = ((float64(100)/ttp)-1)*10 + 1

	// 接触数量调整
	if loanhac > 200 {
		loanhac = 200 + (200 * ((5 - totalLendingPercentage) / 5))
	}

	// 预先付息数量
	predeshac := float64(loanhac) * 0.02
	if predeshac < 1 {
		predeshac = 1
	}

	// 返回可借数量和预付利息
	return loanhac, predeshac
}
