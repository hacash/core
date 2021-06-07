package coinbase

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
