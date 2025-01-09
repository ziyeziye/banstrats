package main

import (
	"os"

	"github.com/banbox/banbot/entry"
	"github.com/banbox/banstrats/trade/strategies"
)

func main() {
	os.Setenv("BanDataDir", "./testdata")
	os.Setenv("BanStratDir", "./")

	strategies.Init()
	entry.RunCmd()
}
