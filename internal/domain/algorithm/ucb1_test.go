package algorithm

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUcb1(t *testing.T) {
	t.Run("choose first element with zero attempts", func(t *testing.T) {
		stats := []Ucb1ElStat{{
			AvgIncome: 0,
			Attempts:  1,
		}, {
			AvgIncome: 1,
			Attempts:  100,
		}, {
			AvgIncome: 0,
			Attempts:  0,
		}, {
			AvgIncome: 0,
			Attempts:  0,
		}}

		number := Ucb1(stats)

		require.Equal(t, 2, number)
	})

	t.Run("choose element with min attempts", func(t *testing.T) {
		stats := []Ucb1ElStat{{
			AvgIncome: 0,
			Attempts:  5,
		}, {
			AvgIncome: 0,
			Attempts:  4,
		}, {
			AvgIncome: 0,
			Attempts:  3,
		}, {
			AvgIncome: 0,
			Attempts:  2,
		}, {
			AvgIncome: 0,
			Attempts:  1,
		}}

		number := Ucb1(stats)

		require.Equal(t, 4, number)
	})

	t.Run("choose element with max average income", func(t *testing.T) {
		stats := []Ucb1ElStat{{
			AvgIncome: 0.1,
			Attempts:  100,
		}, {
			AvgIncome: 0.2,
			Attempts:  100,
		}, {
			AvgIncome: 0.3,
			Attempts:  100,
		}, {
			AvgIncome: 0.4,
			Attempts:  100,
		}, {
			AvgIncome: 0.5,
			Attempts:  100,
		}}

		number := Ucb1(stats)

		require.Equal(t, 4, number)
	})
}
