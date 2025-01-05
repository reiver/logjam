package main

import (
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/flags"
)

func main() {
	app := App{}

	app.Init(flags.ProdMode, cfg.Config)
	app.Run()
}
