package indicator

import (
	"fmt"

	"github.com/banbox/banbot/btime"
	"github.com/banbox/banstrats/trade/environment"
	ta "github.com/banbox/banta"
)

type BaseIndicator struct {
	isHeikin bool
	Value    *ta.Series
	Side     *ta.Series
}

func (e *BaseIndicator) UseHeikin() *BaseIndicator {
	e.isHeikin = true
	return e
}

func (e *BaseIndicator) Print(name string, noPrintHold ...bool) {
	if len(noPrintHold) > 0 && environment.Side(e.Side.Get(0)) != environment.BuySide && environment.Side(e.Side.Get(0)) != environment.SellSide {

	} else {
		fmt.Printf("\n%s %s SIDE:%s Value:%f 收盘: %f\n", btime.ToDateStr(e.Value.Env.TimeStart, ""), name, environment.Side(e.Side.Get(0)).PrintColorSide(), e.Value.Get(0), e.Value.Env.Close.Get(0))
	}
}
