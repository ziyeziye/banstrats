package demo

import (
	_ "embed"
	"fmt"
	"github.com/banbox/banbot/btime"
	"github.com/banbox/banbot/orm"
	"github.com/banbox/banbot/utils"
	"github.com/banbox/banexg"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/banbox/banbot/biz"
	"github.com/banbox/banbot/config"
	"github.com/banbox/banbot/core"
	"github.com/banbox/banbot/exg"
	"github.com/banbox/banbot/opt"
	"github.com/banbox/banexg/errs"
	"github.com/banbox/banexg/log"
	"go.uber.org/zap"
)

//go:embed chartLine.html
var chartTpl []byte

// GenAnyChart 生成任意K线图表示例
func TestGenAnyChart(t *testing.T) {
	exgName, market, symbol, tf, start, end := "binance", "spot", "BTC/USDT", "1d", "20241201", "20250101"
	klines := getKline(exgName, market, symbol, tf, start, end)

	// 4. 准备图表数据
	labels := make([]string, len(klines))
	prices := make([]float64, len(klines))

	for i, k := range klines {
		labels[i] = time.Unix(k.Time/1000, 0).Format("2006-01-02")
		prices[i] = k.Close
	}

	// 5. 创建数据集
	datasets := []*opt.ChartDs{
		{
			Label:       fmt.Sprintf("%s Close Price", symbol),
			Data:        prices,
			BorderColor: "rgb(75, 192, 192)",
		},
	}

	// 6. 生成图表
	title := fmt.Sprintf("%s %s %s Chart", exgName, symbol, tf)
	g := &opt.LineGraph{
		TplData:   chartTpl,
		Title:     title,
		Labels:    labels,
		Datasets:  datasets,
		Precision: 2,
	}

	// 7. 保存并打开文件
	workDir, err_ := os.Getwd()
	if err_ != nil {
		panic(err_)
	}
	outDir := filepath.Join(workDir, "charts")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		panic(errs.New(errs.CodeIOWriteFail, err))
	}

	symbolClean := strings.ReplaceAll(strings.ReplaceAll(symbol, "/", ""), ":", "_")
	outPath := filepath.Join(outDir, fmt.Sprintf("%s_%s_%s.html", exgName, symbolClean, tf))

	if err := g.DumpFile(outPath); err != nil {
		panic(err)
	}

	log.Info("chart generated", zap.String("path", outPath))
	_ = utils.OpenBrowser("file://" + outPath)
}

func getKline(exgName, market, symbol, tf, start, end string) []*banexg.Kline {
	// 1. 初始化
	core.SetRunMode(core.RunModeOther)
	if err := biz.SetupComs(&config.CmdArgs{}); err != nil {
		panic(err)
	}

	// 3. 获取K线数据
	exchange, err := exg.GetWith(exgName, market, "")
	if err != nil {
		panic(err)
	}
	err = orm.InitExg(exchange)
	if err != nil {
		panic(err)
	}
	exs, err := orm.GetExSymbol(exchange, symbol)
	if err != nil {
		panic(err)
	}

	startMS := btime.ParseTimeMS(start)
	endMS := btime.ParseTimeMS(end)
	// auto download and get klines
	_, klines, err := orm.AutoFetchOHLCV(exchange, exs, tf, startMS, endMS, 0, false, nil)
	if err != nil {
		panic(err)
	}
	if len(klines) == 0 {
		panic(errs.NewMsg(errs.CodeRunTime, "no kline data got"))
	}
	return klines
}
