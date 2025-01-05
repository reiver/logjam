package main

import (
	"github.com/reiver/logjam/flags"
)

func main() {
	app := App{}

	app.Init(flags.Src, flags.ProdMode, flags.GoldGorillaSVCAddr)
	app.Run()
}
