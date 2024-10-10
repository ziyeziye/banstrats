package freqtrade

/*
Strategy 005
author@: Gerald Lonlas
github@: https://github.com/freqtrade/freqtrade-strategies/blob/main/user_data/strategies/Strategy005.py
*/

//func Strategy005(p *config.RunPolicyConfig) *strat.TradeStrat {
//	buyVolAvg := int(p.Def("buyVolAvg", 150, core.PNorm(50, 300)))
//	buyRsi := p.Def("buyRsi", 26, core.PNorm(1, 100))
//	buyFastD := p.Def("buyFastD", 1, core.PNorm(1, 100))
//	buyFishNo := p.Def("buyFishNo", 5, core.PNorm(1, 100))
//	sellRsi := p.Def("sellRsi", 74, core.PNorm(1, 100))
//	sellMiDi := p.Def("sellMiDi", 4, core.PNorm(1, 100))
//	sellFiRsiN := p.Def("sellFiRsiN", 30, core.PNorm(1, 100))
//	return &strat.TradeStrat{
//		WarmupNum: 100,
//		OnBar: func(s *strat.StratJob) {
//			e := s.Env
//			macdCols := ta.MACD(e.Close, 12, 26, 9).Cols
//			macd, macdSig := macdCols[0], macdCols[1]
//
//		},
//	}
//}
