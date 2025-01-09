package ma

import (
	"fmt"
	"math"

	"github.com/banbox/banbot/btime"

	"github.com/banbox/banbot/biz"
	"github.com/banbox/banbot/config"
	"github.com/banbox/banbot/orm/ormo"
	"github.com/banbox/banbot/strat"
	"github.com/banbox/banstrats/trade/environment"
	"github.com/banbox/banstrats/trade/indicator"
	"github.com/banbox/banstrats/trade/utils"
	ta "github.com/banbox/banta"
)

func init() {
	// 注册策略到banbot中，后续在配置文件中使用ma:demo即可引用此策略
	// `init`函数是go中的特殊函数，会在当前包被导入时立刻执行
	strat.AddStratGroup("atr", map[string]strat.FuncMakeStrat{
		"demo": Demo,
	})
}

type atrDemo struct {
	atr     *indicator.AtrTrailingStop
	wallets *biz.BanWallets
}

func Demo(pol *config.RunPolicyConfig) *strat.TradeStrat {
	// smlLen := int(pol.Def("smlLen", 5, core.PNorm(3, 10)))
	// bigLen := int(pol.Def("bigLen", 20, core.PNorm(10, 40)))
	return &strat.TradeStrat{
		WarmupNum: 100,
		OnStartUp: func(s *strat.StratJob) {
			s.More = &atrDemo{}
			s.More.(*atrDemo).wallets = biz.GetWallets(config.DefAcc)
			s.More.(*atrDemo).atr = indicator.NewDefaultAtrTrailingStop().UseHeikin()
		},
		OnPairInfos: func(s *strat.StratJob) []*strat.PairSub {
			return []*strat.PairSub{
				{"_cur_", "1m", 10},
			}
		},
		OnBar: func(s *strat.StratJob) {
			e := s.Env
			// if e.BarNum < smlLen {
			// 	return
			// }

			atr := s.More.(*atrDemo).atr.Calculation(e)
			side, side1, stopLoss := environment.Side(atr.Side.Get(0)), environment.Side(atr.Side.Get(1)), atr.Value.Get(0)
			// atr.Print(true)

			if side == environment.BuySide {
				if side1.IsSellTrend() && s.GetOrderNum(-1) > 0 {
					s.CloseOrders(&strat.ExitReq{Tag: "ATR指标转向-多", EnterTag: "Short"})
				}

				OpenOrder(s, environment.BuySide)

			} else if side == environment.SellSide {
				if side1.IsBuyTrend() && s.GetOrderNum(1) > 0 {
					s.CloseOrders(&strat.ExitReq{Tag: "ATR指标转向-空", EnterTag: "Long"})
				}

				OpenOrder(s, environment.SellSide)
			}

			orders := s.GetOrders(0)
			for _, od := range orders {
				od.SetStopLoss(&ormo.ExitTrigger{
					Price: stopLoss,
					Tag:   "ATR指标跟踪止损/盈",
				})
			}

		},
		OnOrderChange: func(s *strat.StratJob, od *ormo.InOutOrder, chgType int) {
			ava := s.More.(*atrDemo).wallets.Get("USDT").Available
			totalAssets := s.More.(*atrDemo).wallets.TotalLegal(nil, true)

			if chgType == strat.OdChgEnterFill {
				fmt.Printf("【%s】 ID:%d 建仓:【%s】, 倍数: %f, 可用: %f, 总资产: %f, 每单最大本金: %f, 价格: %f, Fee: %f, 价值: %f, 数量: %f\n", btime.ToDateStr(od.EnterAt, ""), od.ID, od.EnterTag, od.Leverage, ava, totalAssets, od.QuoteCost/od.Leverage, od.Enter.Average, od.Enter.Fee, od.Enter.Average*od.Enter.Amount, od.Enter.Amount)
			} else if chgType == strat.OdChgExitFill {
				holdNum := int(math.Round(s.Env.BarCount(od.EnterAt)))
				fmt.Printf("【%s】 ID:%d 平仓:%s, 持仓K线:%d, 可用: %f, 总资产: %f, 价格: %f, Fee: %f, 价值: %f, 数量: %f, 利润率:%s%%, 盈亏: %s USDT\n\n", btime.ToDateStr(od.ExitAt, ""), od.ID, utils.PrintMsgByAmountColor("【"+od.ExitTag+"】", od.Profit), holdNum, ava, totalAssets, od.Exit.Average, od.Exit.Fee, od.Exit.Average*od.Exit.Amount, od.Exit.Amount, utils.PrintAmountColor(od.ProfitRate*100), utils.PrintAmountColor(od.Profit))
			}
		},
		OnInfoBar: func(s *strat.StratJob, e *ta.BarEnv, pair, tf string) {
			c := e.Close.Get(0)
			if tf == "1m" {
				orders := s.GetOrders(0)
				for _, od := range orders {
					holdNum := int(math.Round(s.Env.BarCount(od.EnterAt)))
					if (od.ProfitRate < -0.05 && holdNum >= 5) || od.ProfitRate <= -0.2 {
						// fmt.Printf("ID:%d 【%s】Order 持仓K线: %d, 当前价格: %f, 入场价格: %f, profitRate: %f\n", od.ID, carbon.CreateFromTimestampMilli(e.Close.Time).ToDateTimeString(), holdNum, c, od.Enter.Average, od.ProfitRate)

						tag := "多单主动止损"
						if od.Short {
							tag = "空单主动止损"
						}
						s.CloseOrders(&strat.ExitReq{Tag: tag, OrderID: od.ID})
						od.SetStopLoss(&ormo.ExitTrigger{
							Price: c,
							Tag:   tag,
						})
					}

				}
			}
		},
		// OnCheckExit: func(s *strat.StratJob, od *ormo.InOutOrder) *strat.ExitReq {
		// 	// if od.ProfitRate > 0.5 {
		// 	// 	// 盈利超过10%退出
		// 	// 	return &strat.ExitReq{Tag: "profit"}
		// 	// }
		// 	if od.ProfitRate < -0.1 {
		// 		return &strat.ExitReq{Tag: "StopLoss"}
		// 	}
		// 	return nil
		// },
		// GetDrawDownExitRate: func(s *strat.StratJob, od *ormo.InOutOrder, maxChg float64) float64 {
		// 	if maxChg > 0.1 {
		// 		// 订单最佳盈利超过10%后，回撤50%退出
		// 		return 0.5
		// 	}
		// 	return 0
		// },
	}
}

func OpenOrder(s *strat.StratJob, side environment.Side) {
	// 总资产
	// legalValue := wallets.TotalLegal(nil, true)
	ava := s.More.(*atrDemo).wallets.Get("USDT").Available
	// 对于合约市场，百分比开单应基于带杠杆的名义资产价值
	legalValue := ava * config.Leverage
	// Round to the nearest tenth place
	// 四舍五入到十位
	pctAmt := math.Round(legalValue*config.StakePct/1000) * 10

	account, ok := config.Accounts[s.More.(*atrDemo).wallets.Account]
	if ok {
		account.StakePctAmt = pctAmt
	}

	if side == environment.BuySide {
		s.OpenOrder(&strat.EnterReq{Tag: "Long"})
	} else if side == environment.SellSide {
		s.OpenOrder(&strat.EnterReq{Tag: "Short", Short: true})
	}

}
