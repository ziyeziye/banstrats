package adv

import (
	"github.com/banbox/banbot/entry"
)

func init() {
	entry.AddGroup("chart", "generate chart commands")

	entry.AddCmdJob(&entry.CmdJob{
		Name:    "kline",
		Parent:  "chart",
		Run:     genKlineChart,
		Options: []string{"pairs", "timeframes"},
		Help:    "generate klineChart for symbol",
	})

	entry.AddCmdJob(&entry.CmdJob{
		Name:   "demo",
		Parent: "chart",
		RunRaw: genAnyChart,
		Help:   "generate demoChart",
	})
}
