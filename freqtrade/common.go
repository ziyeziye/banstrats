package freqtrade

import "github.com/banbox/banbot/strat"

func init() {
	strat.AddStratGroup("freqtrade", map[string]strat.FuncMakeStrat{
		"Strategy001": Strategy001,
		"Strategy004": Strategy004,
	})
}
