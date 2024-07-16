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
	tools, err := createToolData()
	if err != nil {
		logger.Fatal("could not parse tools data", "error", err)
	}
	logger.Debug("download dir", "dir", cfg.downloadDir)

	installDir, err := normalizePath(cfg.downloadDir)
	if err != nil {
		logger.Fatal("could not normalize path")
	}
	if cfg.tools {
		printTools(tools)
		os.Exit(0)
	}
	if cfg.categories {
		printCategories(getCategories(tools))
		os.Exit(0)
	}
	// interactive mode
	if len(cfg.toolList) == 0 {
		startUI(cfg, tools)
	} else {
		// non-interactive mode
		installEget(cfg.downloadDir)
		for _, toolName := range cfg.toolList {
			err = downloadToolWithEget(installDir, tools.Tools[toolName])
			if err != nil {
				logger.Warn("could not download tool", "tool", toolName, "error", err)
				continue
			}
		}
	}
}
