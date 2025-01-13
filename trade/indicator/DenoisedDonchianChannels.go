package indicator

import (
	"fmt"
	"math"

	"github.com/banbox/banbot/btime"
	"github.com/banbox/banstrats/trade/environment"
	ta "github.com/banbox/banta"
)

type DenoisedDonchianChannels struct {
	*BaseIndicator
	Period        int
	Threshold     float64
	ChangePercent float64
	Upper         *ta.Series
	Lower         *ta.Series
	Basis         *ta.Series
}

func NewDenoisedDonchianChannels(period int, threshold float64, changePercent float64, useHeikin bool) *DenoisedDonchianChannels {
	e := &DenoisedDonchianChannels{
		Period:        period,
		Threshold:     threshold,
		ChangePercent: changePercent,
		BaseIndicator: &BaseIndicator{
			isHeikin: useHeikin,
		},
	}

	return e
}

func NewDefaultDenoisedDonchianChannels() *DenoisedDonchianChannels {
	return NewDenoisedDonchianChannels(20, 0.005, 0.2, false)
}

func (e *DenoisedDonchianChannels) Calculation(df *ta.BarEnv, useClose bool) *DenoisedDonchianChannels {
	keyVal := e.Period*100000 + int(e.Threshold*10000)
	upThreshold := 1 + e.Threshold
	downThreshold := 1 - e.Threshold

	e.Value = df.Close.To("_denoiseddonchianchannels", keyVal)
	if e.Value.Cached() {
		return e
	}

	e.Upper = e.Value.To("_upper", 0)
	e.Lower = e.Value.To("_lower", 0)
	e.Basis = e.Value.To("_basis", 0)
	e.Side = e.Value.To("_side", 0)
	trendData := e.Value.To("_trend", 0)

	trend := 0.0
	side := environment.HoldSide
	open, high, low, close := GetSources(df, e.isHeikin)

	if df.High.Len() > e.Period {
		var upper, lower *ta.Series
		if useClose {
			upper = ta.Highest(close, e.Period)
			lower = ta.Lowest(close, e.Period)
		} else {
			upper = ta.Highest(high, e.Period)
			lower = ta.Lowest(low, e.Period)
		}

		upperVal, lowerVal, upperVal1, lowerVal1 := upper.Get(0), lower.Get(0), e.Upper.Get(0), e.Lower.Get(0)
		close0, open0 := close.Get(0), open.Get(0)
		basisVal := (upperVal + lowerVal) / 2

		if !Na(upperVal1) {
			denoiseLower := lowerVal1 * downThreshold
			denoiseUpper := upperVal1 * upThreshold

			if lowerVal < denoiseLower || lowerVal > lowerVal1 {
				// lowerVal = lowerVal
			} else {
				lowerVal = lowerVal1
			}

			if upperVal < upperVal1 || upperVal > denoiseUpper {
				// upperVal = upperVal
			} else {
				upperVal = upperVal1
			}

			candleHeightPer := math.Abs(close0-open0) / open0 * 100
			buy := IfCase(upperVal > upperVal1 && close0 > open0 && candleHeightPer > e.ChangePercent, true, false)
			sell := IfCase(lowerVal < lowerVal1 && close0 < open0 && candleHeightPer > e.ChangePercent, true, false)

			trend = IfCase(buy, 1.0, IfCase(sell, -1.0, trendData.Get(0)))
			trend = IfCase((trend == 1 && upperVal < upperVal1) || (trend == -1 && lowerVal > lowerVal1), 0, trend)

			if trend == 1 {
				trend = IfCase(close0 > basisVal, 1.0, 0.0)
			} else if trend == -1 {
				trend = IfCase(close0 < basisVal, -1.0, 0.0)
			}

			if buy {
				side = environment.BuySide
			} else if sell {
				side = environment.SellSide
			} else {
				if trend == 1 {
					side = environment.BuyHoldSide
				} else if trend == -1 {
					side = environment.SellHoldSide
				} else {
					side = environment.HoldSide
				}
			}
		}
		e.Basis.Append(basisVal)
		e.Lower.Append(lowerVal)
		e.Upper.Append(upperVal)

	} else {
		e.Basis.Append(0.0)
		e.Lower.Append(0.0)
		e.Upper.Append(0.0)
	}

	trendData.Append(trend)
	e.Side.Append(side.Float64())

	return e
}

func (e *DenoisedDonchianChannels) Print(noPrintHold ...bool) {
	// if len(noPrintHold) > 0 && environment.Side(e.Side.Get(0)) != environment.BuySide && environment.Side(e.Side.Get(0)) != environment.SellSide {

	// } else {
	upper, lower := e.Upper.Get(0), e.Lower.Get(0)
	fmt.Printf("\n%s DenoisedDonchianChannels SIDE:%s Upper:%f Lower:%f Basis:%f 收盘: %f\n", btime.ToDateStr(e.Side.Env.TimeStart, ""), environment.Side(e.Side.Get(0)).PrintColorSide(), upper, lower, e.Basis.Get(0), e.Side.Env.Close.Get(0))
	// }
}
