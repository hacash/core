package coinbase

import (
	"fmt"
	"github.com/hacash/core/fields"
	"math/big"
)

// Calculate the amount of bitcoin redemption
func CalculationBitcoinSystemLendingRedeemAmount(
	feeAddr fields.Address, lendingMainAddr fields.Address,
	loanTotalAmount *fields.Amount, ransomBlockNumberBase uint64,
	lendingCreateBlockHeight uint64, penddingBlockHeight uint64,
) (uint8, *fields.Amount, error) {

	var e error = nil

	// Redemption period stage
	redeemStage := uint8(1) // 1. 私有期  2. 公共期  3. 拍卖期

	// Redemption period basis
	// ransomBlockNumberBase := uint64(100000) // 十万个区块约一年

	// Check private redemption period
	privateHeight := uint64(lendingCreateBlockHeight) + ransomBlockNumberBase
	if penddingBlockHeight <= privateHeight && feeAddr.NotEqual(lendingMainAddr) {
		// Before maturity, it can only be redeemed privately by the mortgagee
		return 0, nil, fmt.Errorf("It can only be redeemed privately by the mortgagor %s before the blockheight %d", lendingMainAddr.ToReadable(), privateHeight)
	}

	// Anyone can redeem publicly
	if penddingBlockHeight > privateHeight {
		redeemStage = 2
	}

	// The redemption amount is the original lending amount (because the interest is paid in advance)
	var realRansomAmt = loanTotalAmount

	// Check the public redemption period and calculate the real-time redemption amount
	publicHeight := privateHeight + ransomBlockNumberBase
	if penddingBlockHeight > publicHeight {
		// Greater than the public redemption period, start the Dutch auction mode of interest, and reduce the redemption amount to 0 in ten years (millions of blocks)
		maxDown := ransomBlockNumberBase * 10
		redeemStage = 3
		overhei := penddingBlockHeight - publicHeight
		if overhei >= maxDown {
			realRansomAmt = fields.NewEmptyAmount() // Redeemable reduced to 0
		} else {
			boli := (float64(maxDown) - float64(overhei)) / float64(maxDown) * loanTotalAmount.ToMei()
			boli *= 100000000
			realRansomAmt, e = fields.NewAmountByBigIntWithUnit(big.NewInt(int64(boli)), 240)
			if e != nil {
				return 0, nil, e
			}
		}
	}

	// Calculation successful return
	return redeemStage, realRansomAmt, nil
}

// Calculate the amount required for diamond redemption
func CalculationDiamondSystemLendingRedeemAmount(
	feeAddr fields.Address, lendingMainAddr fields.Address,
	borrowPeriod int64, lendingCreateBlockHeight int64,
	loanTotalAmountMei int64,
	dslbpbn int64, penddingBlockHeight int64,
) (uint8, *fields.Amount, error) {

	// Redemption period stage
	redeemStage := uint8(1) // 1. 私有期  2. 公共期  3. 拍卖期

	// Redemption period period
	// The redemption period is divided into private redemption period, public redemption period and Dutch auction period
	ransomBlockNumberBase := borrowPeriod * dslbpbn

	// Check private redemption period
	privateHeight := lendingCreateBlockHeight + ransomBlockNumberBase
	if penddingBlockHeight <= privateHeight && feeAddr.NotEqual(lendingMainAddr) {
		// Before maturity, it can only be redeemed privately by the mortgagee
		return 0, nil, fmt.Errorf("It can only be redeemed privately by the mortgagor %s before the blockheight %d", lendingMainAddr.ToReadable(), privateHeight)
	}

	// Anyone can redeem publicly
	if penddingBlockHeight > privateHeight {
		redeemStage = 2
	}

	// Calculate the redemption amount (each cycle represents 0.005, i.e. 0.5% interest)
	realRansom1qian := (1000 + (5 * int64(borrowPeriod))) * int64(loanTotalAmountMei)

	// Check public redemption period
	// The duration of public redemption period is the same as that of private redemption period. Assuming that borrowperiod = 20, then:
	// Private redemption period: 2 years, public redemption period: 2 years, Dutch auction period: 2+2=4 years
	publicHeight := privateHeight + ransomBlockNumberBase
	if penddingBlockHeight > publicHeight {
		redeemStage = 3
		// Greater than the public redemption period, start the Dutch auction mode of interest
		subcount := int64((penddingBlockHeight - publicHeight) / dslbpbn)
		maxsub := int64(borrowPeriod) * 2
		if subcount > maxsub {
			redeemStage = 4
			subcount = maxsub // The maximum deduction interest can only be equal to twice the original interest paid, that is, it provides a positive and negative symmetrical arbitrage space
		}
		// Deduct interest by time, auction: deduct 0.5% interest every 35 days
		realRansom1qian -= 5 * subcount * int64(loanTotalAmountMei)
	}

	// True and effective redemption amount
	validRansomAmt, e1 := fields.NewAmountByBigIntWithUnit(big.NewInt(realRansom1qian), 245)
	if e1 != nil {
		return 0, nil, e1
	}

	// Calculation succeeded, return
	return redeemStage, validRansomAmt, nil

}

// Bitcoin mortgage loan: calculate the quantity that can be borrowed
// Totallendingpercentage total percentage lent, unit:%
// Return lendable quantity and prepaid interest
func CalculationOfInterestBitcoinMortgageLoanAmount(totalLendingPercentage float64) (float64, float64) {

	ttp := totalLendingPercentage

	if ttp < 1 {
		ttp = 1 // Minimum 1
	}

	// Lendable quantity
	var loanhac float64 = 0
	loanhac = ((float64(100)/ttp)-1)*10 + 1

	// Contact quantity adjustment
	if loanhac > 200 {
		loanhac = 200 + (200 * ((5 - totalLendingPercentage) / 5))
	}

	// Advance interest payment quantity
	predeshac := float64(loanhac) * 0.02
	if predeshac < 1 {
		predeshac = 1
	}

	// Return lendable quantity and prepaid interest
	return loanhac, predeshac
}
