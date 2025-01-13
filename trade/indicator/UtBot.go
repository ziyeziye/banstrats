package indicator

import (
	"fmt"
	"math"

	"github.com/banbox/banbot/btime"
	"github.com/banbox/banstrats/trade/environment"
	ta "github.com/banbox/banta"
)

type UtBot struct {
	*BaseIndicator
	Period int
	Mult   float64
	Upper  *ta.Series
	Lower  *ta.Series
}

func NewUtBot(period int, mult float64, useHeikin bool) *UtBot {
	e := &UtBot{
		Period: period,
		Mult:   mult,
		BaseIndicator: &BaseIndicator{
			isHeikin: useHeikin,
		},
	}

	return e
}

func NewDefaultUtBot() *UtBot {
	return NewUtBot(10, 1, false)
}

func (e *UtBot) Calculation(df *ta.BarEnv, srcType environment.SourceType) *UtBot {
	keyVal := e.Period*100000 + int(e.Mult*10000)
	// xAtrStop
	e.Value = df.Close.To("_utbot", keyVal)
	if e.Value.Cached() {
		return e
	}

	e.Side = e.Value.To("_side", 0)

	xStopLoss := 0.0
	side := environment.HoldSide
	if df.High.Len() > e.Period {
		xAtr := ta.ATR(df.High, df.Low, df.Close, e.Period)
		nLoss := xAtr.Mul(e.Mult)
		src := GetSource(df, srcType, e.isHeikin)

		src0, src1, nLoss0, xStopLoss1 := src.Get(0), src.Get(1), nLoss.Get(0), e.Value.Get(0)

		val := IfCase(src0 > xStopLoss1, src0-nLoss0, src0+nLoss0)
		val = IfCase(src0 < xStopLoss1 && src1 < xStopLoss1, math.Min(xStopLoss1, src0+nLoss0), val)
		xStopLoss = IfCase(src0 > xStopLoss1 && src1 > xStopLoss1, math.Max(xStopLoss1, src0-nLoss0), val)
		e.Value.Append(xStopLoss)

		ema := ta.EMA(src, 1)
		cross := ta.Cross(ema, e.Value)

		if cross == 1 && src0 > xStopLoss {
			side = environment.BuySide
		} else if cross == -1 && src0 < xStopLoss {
			side = environment.SellSide
		} else {
			if cross == 0 {
				side = environment.HoldSide
			} else {
				side = IfCase(cross > 0, environment.BuyHoldSide, environment.SellHoldSide)
			}
		}

	} else {
		e.Value.Append(0)
	}

	e.Side.Append(side.Float64())

	return e
}

func (e *UtBot) Print(noPrintHold ...bool) {
	if len(noPrintHold) > 0 && environment.Side(e.Side.Get(0)) != environment.BuySide && environment.Side(e.Side.Get(0)) != environment.SellSide {

	} else {
		fmt.Printf("%s UT BOT SIDE:%s Value:%f 收盘: %f\n", btime.ToDateStr(e.Side.Env.TimeStart, ""), environment.Side(e.Side.Get(0)).PrintColorSide(), e.Value.Get(0), e.Side.Env.Close.Get(0))
	}
}
