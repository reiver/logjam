package main

import (
	"flag"
	"os"
)

func main() {
	listenHost := flag.String("listen-host", "0.0.0.0", "interface to listen on")
	listenPort := flag.Int("listen-port", 8090, "http listen port")
	help := flag.Bool("h", false, "print help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	app := App{}

	app.Init(*listenHost, *listenPort, "dev|prod")
	app.Run()
}
