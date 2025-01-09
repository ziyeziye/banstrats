package environment

import "github.com/banbox/banstrats/trade/utils"

type Side float64

const (
	HoldSide Side = iota
	BuySide
	SellSide
	BuyHoldSide
	SellHoldSide
)

func (s Side) String() string {
	return []string{"Hold", "Buy", "Sell", "BuyHold", "SellHold"}[int(s)]
}

func (s Side) Float64() float64 {
	return float64(s)
}

func (s Side) IsBuyTrend() bool {
	return s == BuySide || s == BuyHoldSide
}

func (s Side) IsSellTrend() bool {
	return s == SellSide || s == SellHoldSide
}

func (s Side) PrintColorSide() string {
	return getColorSide(s)
}

func getColorSide(s Side) string {
	switch s {
	case BuySide, BuyHoldSide:
		return utils.PrintColor(s.String(), utils.Green) // 绿色
	case SellSide, SellHoldSide:
		return utils.PrintColor(s.String(), utils.Red) // 红色
	default:
		return s.String()
	}
}
