package rpc_ai

import (
	"github.com/banbox/banbot/biz"
	"github.com/banbox/banbot/btime"
	"github.com/banbox/banbot/core"
	"github.com/banbox/banbot/orm"
	"github.com/banbox/banexg"
	"github.com/banbox/banexg/errs"
	"github.com/banbox/banexg/log"
	utils2 "github.com/banbox/banexg/utils"
	"github.com/banbox/banta"
	"go.uber.org/zap"
	"gonum.org/v1/gonum/mat"
)

type AIFeaState struct {
	feasBig *mat.Dense
}

func pubAiFea(exsList []*orm.ExSymbol, req *biz.SubReq, rsp biz.FeaFeeder_SubFeaturesServer) error {
	var sample = int(req.Sample)
	if sample <= 0 {
		sample = 10
	}
	startDt := btime.ToDateStr(req.Start, core.DefaultDateFmt)
	log.Info("start ai fea2", zap.Int("pair", len(exsList)), zap.String("start", startDt),
		zap.Int("sample", sample))
	var seqNum = 50 // 特征序列长度
	tfSml, tfBig := "5m", "1h"
	tfMSecs := int64(utils2.TFToSecs(tfSml) * 1000)
	tfMSecsBig := int64(utils2.TFToSecs(tfBig) * 1000)

	states := make(map[string]*AIFeaState)
	for _, exs := range exsList {
		states[exs.Symbol] = &AIFeaState{}
	}

	task := &GenFeaTask{
		TfSml:      tfSml,
		TfBig:      tfBig,
		WarmSml:    300,
		WarmBig:    100,
		TfMSecs:    tfMSecs,
		TfMSecsBig: tfMSecsBig,
		NextObsNum: 30,
	}

	batch := NewBatchMap(len(exsList))

	task.OnBar = func(e *banta.BarEnv, b *banexg.Kline, futs []*banexg.Kline) *errs.Error {
		//if !core.IsWarmUp {
		//	log.Info("OnBar", zap.Int64("time", e.TimeStop))
		//}
		state := states[e.Symbol]
		feas := onAiFeatures(e, seqNum)
		if e.BarNum%sample > 0 || len(futs) == 0 {
			// 预热时futs为nil
			return nil
		}
		// 收集特征，发送到监听端
		minChgVal := banta.ATR(e.High, e.Low, e.Close, 30).Get(0) * 5
		// 生成未来行情的分类概率：不确定、上升、下降
		// generate classify probs for future trend: unsure, up, down
		futProbs := genFutProbs(futs, minChgVal, 0.4)
		info := make([]float64, 0, 9)
		info = append(info, float64(b.Time/1000), b.Open, b.High, b.Low, b.Close, b.Volume)
		info = append(info, futProbs[0], futProbs[1], futProbs[2])
		return batch.Add(rsp, e, b.Time, map[string]*biz.NumArr{
			"feas":  matToNumArr(feas),
			"feas2": matToNumArr(state.feasBig),
			"info":  {Data: info, Shape: []int32{int32(len(info))}},
		})
	}
	task.OnInfoBar = func(e *banta.BarEnv, b *banexg.Kline) *errs.Error {
		//if !core.IsWarmUp {
		//	log.Info("OnInfoBar", zap.Int64("time", e.TimeStop))
		//}
		state := states[e.Symbol]
		state.feasBig = onAiFeatures(e, seqNum)
		return nil
	}
	task.OnEnvEnd = func(bar *banexg.PairTFKline, adj *orm.AdjInfo) {
		if bar == nil {
			return
		}
		//state := states[bar.Symbol]
	}
	return pubFeaBase(exsList, req, task)
}
