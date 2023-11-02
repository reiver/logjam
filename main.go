package main

import (
	"flag"
	"os"
)

func main() {
	src := flag.String("src", "0.0.0.0:8080", "source listen address")
	goldGorillaSVCAddr := flag.String("goldgorilla-svc-addr", "http://localhost:8080", "goldgorilla service address baseurl")
	prodMode := flag.Bool("prod", false, "enable production mode ( its in dev mode by default )")
	help := flag.Bool("h", false, "print help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	app := App{}

	app.Init(*src, *prodMode, *goldGorillaSVCAddr)
	app.Run()
}
