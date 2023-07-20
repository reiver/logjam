package main

import (
	"flag"
	"os"
)

type S struct {
	Name string
}

func main() {
	src := flag.String("src", "0.0.0.0:8090", "source listen address")
	prodMode := flag.Bool("prod", false, "enable production mode ( its in dev mode by default )")
	help := flag.Bool("h", false, "print help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	app := App{}

	app.Init(*src, *prodMode)
	app.Run()
}
