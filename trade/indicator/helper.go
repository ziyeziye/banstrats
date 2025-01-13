package indicator

import (
	"math"

	"github.com/banbox/banstrats/trade/environment"
	ta "github.com/banbox/banta"
)

func Na(vals float64) bool {
	return math.IsNaN(vals)
}

// Nz 以系列中的指定数替换NaN值。
func Nz(vals float64, replacement float64) float64 {
	if Na(vals) {
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

// GetSource 获取指定来源的序列
func GetSource(df *ta.BarEnv, srcType environment.SourceType, isHeikin ...bool) *ta.Series {
	var open, high, low, close *ta.Series
	if len(isHeikin) > 0 && isHeikin[0] {
		haCols := ta.HeikinAshi(df).Cols
		open, high, low, close = haCols[0], haCols[1], haCols[2], haCols[3]
	} else {
		open, high, low, close = df.Open, df.High, df.Low, df.Close
	}

	var src *ta.Series
	switch srcType {
	case environment.OpenSource:
		src = open
	case environment.HighSource:
		src = high
	case environment.LowSource:
		src = low
	default:
		src = close
	}

	return src
}

func GetSources(df *ta.BarEnv, isHeikin ...bool) (*ta.Series, *ta.Series, *ta.Series, *ta.Series) {
	var open, high, low, close *ta.Series
	if len(isHeikin) > 0 && isHeikin[0] {
		haCols := ta.HeikinAshi(df).Cols
		open, high, low, close = haCols[0], haCols[1], haCols[2], haCols[3]
	} else {
		open, high, low, close = df.Open, df.High, df.Low, df.Close
	}

	return open, high, low, close
}
