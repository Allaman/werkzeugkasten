package main

import (
	"os"

	"github.com/charmbracelet/log"
)

var logger = log.New(os.Stderr)

// will be overwritten in release pipeline
var version = "dev"

func main() {
	cfg := cli()
	if cfg.debug {
		logger.SetReportCaller(true)
		logger.SetLevel(log.DebugLevel)
	}
	startUI(cfg)
}
