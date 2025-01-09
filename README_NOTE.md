# TradeStrat 策略对象
```go
type TradeStrat struct {
	Name          string // 策略名称，无需设置，会自动设置为注册时的名称
	Version       int // 策略版本号
	WarmupNum     int // K线预热的长度
	MinTfScore    float64 // 最小时间周期质量，默认0.8
	WatchBook     bool // 是否监听订单簿实时深度信息
	DrawDownExit  bool // 是否启用回撤止损（即跟踪止损）
	BatchInOut    bool    // 是否批量执行入场/出场
	BatchInfo     bool    // 是否对OnInfoBar后执行批量处理
	StakeRate     float64 // 相对基础金额开单倍率
	StopEnterBars int // 限价单如果超过给定K线仍未入场则取消
	EachMaxLong   int      // 单个品种最大同时开多数量，默认0不限制
	EachMaxShort  int      // 单个品种最大同时开空数量，默认0不限制
	AllowTFs      []string // 允许运行的时间周期，不提供时使用全局配置
	Outputs       []string // 策略输出的文本文件内容，每个字符串是一行
	Policy        *config.RunPolicyConfig // 运行时配置

	OnPairInfos         func(s *StratJob) []*PairSub // 策略额外需要的品种或其他周期的数据
	OnStartUp           func(s *StratJob) // 启动时回调。初次执行前调用
	OnBar               func(s *StratJob) // 每个K线的回调函数
	OnInfoBar           func(s *StratJob, e *ta.BarEnv, pair, tf string)   // 其他依赖的bar数据
	OnTrades            func(s *StratJob, trades []*banexg.Trade)          // 逐笔交易数据
	OnBatchJobs         func(jobs []*StratJob)                             // 当前时间所有标的job，用于批量开单/平仓
	OnBatchInfos        func(jobs map[string]*StratJob)                    // 当前时间所有info标的job，用于批量处理
	OnCheckExit         func(s *StratJob, od *ormo.InOutOrder) *ExitReq     // 自定义订单退出逻辑
	OnOrderChange       func(s *StratJob, od *ormo.InOutOrder, chgType int) // 订单更新回调
	GetDrawDownExitRate CalcDDExitRate                                     // 计算跟踪止盈回撤退出的比率
	PickTimeFrame       PickTimeFrameFunc                                  // 为指定币选择适合的交易周期
	OnShutDown          func(s *StratJob)                                  // 机器人停止时回调
}

type EnterReq struct {
	Tag             string  // 入场信号
	StgyName        string  // 策略名称
	Short           bool    // 是否做空
	OrderType       int     // 订单类型, core.OrderType*
	Limit           float64 // 限价单入场价格，指定时订单将作为限价单提交
	CostRate        float64 // 开仓倍率、默认按配置1倍。用于计算LegalCost
	LegalCost       float64 // 花费法币金额。指定时忽略CostRate
	Leverage        float64 // 杠杆倍数
	Amount          float64 // 入场标的数量，由LegalCost和price计算
	StopLossVal     float64 // 入场价格到止损价格的距离，用于计算StopLoss
	StopLoss        float64 // 止损触发价格，不为空时在交易所提交一个止损单
	StopLossLimit   float64 // 止损限制价格，不提供使用StopLoss
	StopLossRate    float64 // 止损退出比例，0表示全部退出，需介于(0,1]之间
	StopLossTag     string  // 止损原因
	TakeProfitVal   float64 // 入场价格到止盈价格的距离，用于计算TakeProfit
	TakeProfit      float64 // 止盈触发价格，不为空时在交易所提交一个止盈单。
	TakeProfitLimit float64 // 止盈限制价格，不提供使用TakeProfit
	TakeProfitRate  float64 // 止盈退出比率，0表示全部退出，需介于(0,1]之间
	TakeProfitTag   string  // 止盈原因
	StopBars        int     // 入场限价单超过多少个bar未成交则取消
}

type ExitReq struct {
	Tag        string  // 退出信号
	StgyName   string  // 策略名称
	EnterTag   string  // 只退出入场信号为EnterTag的订单
	Dirt       int     // core.OdDirt* long/short/both
	OrderType  int     // 订单类型, core.OrderType*
	Limit      float64 // 限价单退出价格，指定时订单将作为限价单提交
	ExitRate   float64 // 退出比率，默认100%即所有订单全部退出
	Amount     float64 // 要退出的标的数量。指定时ExitRate无效
	OrderID    int64   // 只退出指定订单
	UnOpenOnly bool    // True时只退出尚未入场的订单
	FilledOnly bool    // True时只退出已入场的订单
	Force      bool    // 是否强制退出
}
```