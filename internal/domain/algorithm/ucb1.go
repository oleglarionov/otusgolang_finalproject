package algorithm

import "math"

type Ucb1ElStat struct {
	AvgIncome float64
	Attempts  uint64
}

func Ucb1(stats []Ucb1ElStat) int {
	var totalAttempts uint64 = 0
	for j, stat := range stats {
		if stat.Attempts == 0 {
			return j
		}
		totalAttempts += stat.Attempts
	}

	var maxJ int
	var maxValue float64
	for j, stat := range stats {
		value := stat.AvgIncome + math.Sqrt(2*math.Log(float64(totalAttempts))/float64(stat.Attempts))
		if value > maxValue {
			maxValue = value
			maxJ = j
		}
	}

	return maxJ
}
