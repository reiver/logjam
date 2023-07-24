package main

import (
	"flag"
	"os"
)

func main() {
	src := flag.String("src", "0.0.0.0:8090", "source listen address")
	prodMode := flag.Bool("prod", false, "enable production mode ( its in dev mode by default )")
	anSVCAddr := flag.String("auxiliarynode-svc-addr", "http://example.com/", "auxiliary node service rest api address")
	help := flag.Bool("h", false, "print help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	app := App{}

	app.Init(*src, *prodMode, *anSVCAddr)
	app.Run()
}
