package main

import (
	"os"

	"github.com/charmbracelet/log"
)

var logger = log.New(os.Stderr)

func main() {
	debugMode := os.Getenv("WK_DEBUG") != ""
	if debugMode {
		logger.SetReportCaller(true)
		logger.SetLevel(log.DebugLevel)
	}
	startUI()
}
