package freqtrade

import (
	"github.com/banbox/banbot/config"
	"github.com/banbox/banbot/strat"
	ta "github.com/banbox/banta"
)

/*
Strategy 004
author@: Gerald Lonlas
github@: https://github.com/freqtrade/freqtrade-strategies/blob/main/user_data/strategies/Strategy004.py
*/
func Strategy004(p *config.RunPolicyConfig) *strat.TradeStrat {
	return &strat.TradeStrat{
		WarmupNum:   100,
		EachMaxLong: 1,
		OnBar: func(s *strat.StratJob) {
			e := s.Env
			adx := ta.ADX(e.High, e.Low, e.Close, 14).Cols[0].Get(0)
			adxSlow := ta.ADX(e.High, e.Low, e.Close, 35).Cols[0].Get(0)
			cci := ta.CCI(e.Close, 14).Get(0)
			fastD, _, fastK := ta.KDJBy(e.High, e.Low, e.Close, 5, 3, 3, "sma")
			fastDBig, _, fastKBig := ta.KDJBy(e.High, e.Low, e.Close, 50, 3, 3, "sma")
			ema := ta.EMA(e.Close, 5).Get(0)
			KDX := ta.Cross(fastK, fastD)

			adxOk := adx > 50 || adxSlow > 26
			fastKDOk := fastK.Get(1) < 20 && fastD.Get(1) < 20
			fastKDBigOk := fastKBig.Get(1) < 30 && fastDBig.Get(1) < 30

			fastKDgt70 := fastK.Get(0) > 70 || fastD.Get(0) > 70
			if adxOk && cci < -100 && fastKDOk && fastKDBigOk && KDX == 1 {
				s.OpenOrder(&strat.EnterReq{
					Tag:         "long",
					StopLossVal: ema * 0.1,
				})
			} else if adxSlow < 25 && fastKDgt70 && fastK.Get(1) < fastD.Get(1) && e.Close.Get(0) > ema {
				s.CloseOrders(&strat.ExitReq{Tag: "exit"})
			}
		},
	}
}
