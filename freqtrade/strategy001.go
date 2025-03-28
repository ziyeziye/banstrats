package freqtrade

import (
	"github.com/banbox/banbot/config"
	"github.com/banbox/banbot/core"
	"github.com/banbox/banbot/strat"
	ta "github.com/banbox/banta"
)

/*
 Strategy 001
author@: Gerald Lonlas
github@: https://github.com/freqtrade/freqtrade-strategies/blob/main/user_data/strategies/Strategy001.py
*/

func Strategy001(p *config.RunPolicyConfig) *strat.TradeStrat {
	lenSml := int(p.Def("lenSml", 20, core.PNorm(10, 30)))
	midRate := p.Def("midRate", 2.5, core.PNorm(1.5, 5))
	bigRate := p.Def("bigRate", 2, core.PNorm(1.5, 3))
	lenMid := int(float64(lenSml) * midRate)
	lenBig := int(float64(lenMid) * bigRate)
	return &strat.TradeStrat{
		WarmupNum:   200,
		EachMaxLong: 1,
		OnBar: func(s *strat.StratJob) {
			e := s.Env
			maSml := ta.EMA(e.Close, lenSml)
			maMid := ta.EMA(e.Close, lenMid)
			maBig := ta.EMA(e.Close, lenBig)
			haOpenC, _, _, haCloseC := ta.HeikinAshi(e)
			haOpen := haOpenC.Get(0)
			haClose := haCloseC.Get(0)
			smlXmid := ta.Cross(maSml, maMid)
			midXBig := ta.Cross(maMid, maBig)

			if smlXmid == 1 && haClose > maSml.Get(0) && haOpen < haClose {
				s.OpenOrder(&strat.EnterReq{
					Tag:         "long",
					StopLossVal: haClose * 0.1,
				})
			} else if midXBig == 1 && haClose < maSml.Get(0) && haOpen > haClose {
				s.CloseOrders(&strat.ExitReq{Tag: "exit"})
			}
		},
	}
}
