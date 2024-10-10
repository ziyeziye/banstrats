package main

import (
	"github.com/banbox/banbot/entry"
	_ "github.com/banbox/banstrats/freqtrade"
	_ "github.com/banbox/banstrats/grid"
	_ "github.com/banbox/banstrats/ma"
)

func main() {
	entry.RunCmd()
}
