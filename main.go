package main

import (
	"flag"
	"os"
)

func main() {
	listenHost := flag.String("listen-host", "0.0.0.0", "interface to listen on")
	listenPort := flag.Int("listen-port", 8090, "http listen port")
	prodMode := flag.Bool("prod", false, "enable production mode ( its in dev mode by default )")
	help := flag.Bool("h", false, "print help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	app := App{}

	app.Init(*listenHost, *listenPort, *prodMode)
	app.Run()
}
