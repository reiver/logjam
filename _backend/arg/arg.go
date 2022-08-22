package arg

import (
	"flag"
	"os"
)

var (
	Values []string
)

var (
	Verbose bool
	VeryVerbose bool
	VeryVeryVerbose bool
	VeryVeryVeryVerbose bool
	VeryVeryVeryVeryVerbose bool
	VeryVeryVeryVeryVeryVerbose bool
)

var (
	Help bool
)

var (
	HttpPort uint64
)

var (
	TestLogs bool
)

func init() {
	flag.BoolVar(&Verbose,                     "v",      false,                          "verbose logs outputted")
	flag.BoolVar(&VeryVerbose,                 "vv",     false,                     "very verbose logs outputted")
	flag.BoolVar(&VeryVeryVerbose,             "vvv",    false,                "very very verbose logs outputted")
	flag.BoolVar(&VeryVeryVeryVerbose,         "vvvv",   false,           "very very very verbose logs outputted")
	flag.BoolVar(&VeryVeryVeryVeryVerbose,     "vvvvv",  false,      "very very very very verbose logs outputted")
	flag.BoolVar(&VeryVeryVeryVeryVeryVerbose, "vvvvvv", false, "very very very very very verbose logs outputted")

	flag.BoolVar(&Help, "help",   false, "outputs this help message")

	flag.Uint64Var(&HttpPort, "http-port", 8080, "HTTP port")
	flag.BoolVar(&TestLogs, "test-logs", false, "test logs")

	flag.Parse()

	Values = flag.Args()

	// --help
	if Help {
		flag.PrintDefaults()
		os.Exit(0)
	}
}
