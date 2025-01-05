package main

import (
	"github.com/reiver/logjam/cfg"
)

func main() {
	app := App{}

	app.Init(cfg.Config)
	app.Run()
}
