package indicator

import (
	"math"

	ta "github.com/banbox/banta"
)

// Nz 以系列中的指定数替换NaN值。
func Nz(vals float64, replacement float64) float64 {

	if math.IsNaN(vals) {
		return replacement
	}

	return vals
}

func IfCase[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// Crossover returns true if the last value of the series is greater than the last value of the reference series
func Crossover(s, ref *ta.Series) bool {
	return s.Get(0) > ref.Get(0) && s.Get(1) <= ref.Get(1)
}

// Crossunder returns true if the last value of the series is less than the last value of the reference series
func Crossunder(s, ref *ta.Series) bool {
	return s.Get(0) <= ref.Get(0) && s.Get(1) > ref.Get(1)
}

// Cross returns true if the last value of the series is greater than the last value of the
// reference series or less than the last value of the reference series
func Cross(s, ref *ta.Series) bool {
	return Crossover(s, ref) || Crossunder(s, ref)
}
