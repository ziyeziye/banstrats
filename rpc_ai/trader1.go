package rpc_ai

import (
	"context"
	"fmt"
	"github.com/banbox/banbot/biz"
	"github.com/banbox/banbot/config"
	"github.com/banbox/banbot/core"
	"github.com/banbox/banbot/strat"
	"github.com/banbox/banexg/log"
	ta "github.com/banbox/banta"
	"go.uber.org/zap"
	"gonum.org/v1/gonum/mat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math"
)

type AIMore struct {
	feasBig *mat.Dense // 大周期特征
	feas    *mat.Dense // 小周期特征
	info    []float64
	atr     float64
	age     int // 已持仓bar数
}

func AITrade(pol *config.RunPolicyConfig) *strat.TradeStrat {
	var seqNum = 50    // 特征序列长度
	var maxOdNum = 1.0 // 单币种单方向最多开1单
	const maxMsgSize = 100 * 1024 * 1024
	const maxHoldAge = 300       // 持仓age超过此配置，强制平仓
	addr := "192.168.1.253:8080" // python model grpc addr
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err_ := grpc.NewClient(addr, creds, grpc.WithDefaultCallOptions(
		grpc.MaxCallSendMsgSize(maxMsgSize),
		grpc.MaxCallRecvMsgSize(maxMsgSize),
	))
	if err_ != nil {
		panic(err_)
	}
	client := biz.NewAInferClient(conn)
	return &strat.TradeStrat{
		BatchInOut: true,
		WarmupNum:  300,
		OnPairInfos: func(s *strat.StratJob) []*strat.PairSub {
			return []*strat.PairSub{
				{"_cur_", "1h", 300},
			}
		},
		OnBar: func(s *strat.StratJob) {
			m, _ := s.More.(*AIMore)
			e := s.Env
			//if !core.IsWarmUp {
			//	log.Info("OnBar", zap.Int64("time", e.TimeStop))
			//	defer log.Info("OnBarEnd", zap.Int64("time", e.TimeStop))
			//}
			m.feas = onAiFeatures(e, seqNum)
			o, h, l, c, v := e.Open.Get(0), e.High.Get(0), e.Low.Get(0), e.Close.Get(0), e.Volume.Get(0)
			info := []float64{float64(e.TimeStart / 1000), o, h, l, c, v}
			info = append(info, 0, 0, 0)
			m.info = info
			m.atr = ta.ATR(e.High, e.Low, e.Close, 30).Get(0)
		},
		OnInfoBar: func(s *strat.StratJob, e *ta.BarEnv, pair, tf string) {
			//if !core.IsWarmUp {
			//	log.Info("OnInfoBar", zap.Int64("time", e.TimeStop))
			//	defer log.Info("OnInfoBarEnd", zap.Int64("time", e.TimeStop))
			//}
			m, _ := s.More.(*AIMore)
			m.feasBig = onAiFeatures(e, seqNum)
		},
		OnBatchJobs: func(jobs []*strat.StratJob) {
			var valids []*strat.StratJob
			var feas1, feas2, info []float64
			var feasLen, feasDepth, infoDepth int
			for _, j := range jobs {
				m, _ := j.More.(*AIMore)
				if len(m.info) != 23 {
					continue
				}
				if len(valids) == 0 {
					feasLen, feasDepth = m.feas.Dims()
					infoDepth = len(m.info)
				}
				valids = append(valids, j)
				feas1 = append(feas1, m.feas.RawMatrix().Data...)
				feas2 = append(feas2, m.feasBig.RawMatrix().Data...)
				info = append(info, m.info...)
			}
			//if !core.IsWarmUp {
			//	log.Info("OnBatchJobs", zap.Int64("time", timeStop))
			//	defer log.Info("OnBatchJobs", zap.Int64("time", timeStop))
			//}
			if len(valids) == 0 {
				return
			}
			bSize := len(valids)
			feaShape := []int32{int32(bSize), int32(feasLen), int32(feasDepth)}
			ctx := context.Background()

			trend, err_ := client.Trend(ctx, &biz.ArrMap{
				Mats: map[string]*biz.NumArr{
					"feas":  {Data: feas1, Shape: feaShape},
					"feas2": {Data: feas2, Shape: feaShape},
					"info":  {Data: info, Shape: []int32{int32(bSize), int32(infoDepth)}},
				},
			})
			if err_ != nil {
				log.Panic("call ai trend fail", zap.Int("bSize", bSize), zap.Error(err_))
			}
			preds := trend.Mats["pred"].Data
			if len(preds) != bSize {
				panic(fmt.Sprintf("preds length error, expect: %v, got: %v", len(preds), bSize))
			}
			for i, j := range valids {
				m, _ := j.More.(*AIMore)
				pred := int(math.Round(preds[i])) // 1: long  2: short
				truncated := false
				if len(j.LongOrders) > 0 || len(j.ShortOrders) > 0 {
					m.age += 1
					if m.age >= maxHoldAge {
						truncated = true
					}
				}
				if (pred == 2 || truncated) && len(j.LongOrders) > 0 {
					_ = j.CloseOrders(&strat.ExitReq{
						Tag:  "exit_long",
						Dirt: core.OdDirtLong,
					})
				}
				if truncated || pred == 1 && len(j.ShortOrders) > 0 {
					_ = j.CloseOrders(&strat.ExitReq{
						Tag:  "exit_short",
						Dirt: core.OdDirtShort,
					})
				} else if pred == 1 && len(j.LongOrders) < int(maxOdNum) {
					_ = j.OpenOrder(&strat.EnterReq{
						Tag: "long",
					})
				} else if pred == 2 && len(j.ShortOrders) < int(maxOdNum) {
					_ = j.OpenOrder(&strat.EnterReq{
						Tag:   "short",
						Short: true,
					})
				}
			}
		},
	}
}
