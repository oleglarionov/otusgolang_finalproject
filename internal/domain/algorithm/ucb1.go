package algorithm

import "math"

func Ucb1(avgIncome []float64, nj []uint64, n uint64) int {
	l := len(avgIncome)

	var maxJ int
	var maxValue float64
	for j := 0; j < l; j++ {
		value := avgIncome[j] + math.Sqrt(2*math.Log(float64(n))/float64(nj[j]))
		if value > maxValue {
			maxValue = value
			maxJ = j
		}
	}

	return maxJ
}
