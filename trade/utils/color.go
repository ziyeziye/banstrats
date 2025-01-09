package utils

import "fmt"

var (
	Green = string([]byte{27, 91, 51, 50, 109})
	Red   = string([]byte{27, 91, 51, 49, 109})
	Reset = string([]byte{27, 91, 48, 109})
)

func PrintColor(s string, color string) string {
	return fmt.Sprintf("%s%s%s", color, s, Reset)
}

func PrintAmountColor(amount float64) string {
	amountStr := fmt.Sprintf("%f", amount)
	color := Green
	if amount < 0 {
		color = Red
	}
	return PrintColor(amountStr, color)
}

func PrintMsgByAmountColor(msg string, amount float64) string {
	color := Green
	if amount < 0 {
		color = Red
	}
	return PrintColor(msg, color)
}
