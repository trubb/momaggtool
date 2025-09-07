package main

import (
	"fmt"
	"sort"
)

// TODO is this really needed, can't we call order and print from main?
func display(funds map[int]FundInfo) error {
	ordered := order(funds)

	print(ordered)

	return nil
}

func order(funds map[int]FundInfo) []FundInfo {
	toOrder := make([]FundInfo, 0, len(funds))

	for _, f := range funds {
		toOrder = append(toOrder, f)
	}

	sort.Slice(toOrder, func(i, j int) bool {
		return toOrder[i].ThreeMonthPerformance > toOrder[j].ThreeMonthPerformance
	})

	return toOrder
}

func print(fs []FundInfo) {
	for i, f := range fs {
		fmt.Printf(
			"\t%d: %s | 3m perf: %s%% | SMA%dd: %.2f (dist %s%%) | NAV: %f | ID %d\n",
			i,
			f.Name,
			highlightPercentage(f.ThreeMonthPerformance),
			f.SmaPeriod,
			f.Sma,
			highlightPercentage(f.SmaDistance),
			f.Ohlc[len(f.Ohlc)-1].Close,
			f.AzaID,
		)

		if i == 3 {
			fmt.Println("\t", repeat("-", 120))
		}
	}
}

// highlightPercentage returns a string with ANSI color codes
// to highlight positive (green) and negative (red) percentage values
func highlightPercentage(f float64) string {
	if f < 0 {
		return fmt.Sprintf("\033[31m%.2f%%\033[0m", f) // red
	}
	return fmt.Sprintf("\033[32m%.2f%%\033[0m", f) // green
}
