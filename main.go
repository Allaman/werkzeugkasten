package main

import (
	"fmt"
	"os"

	"github.com/allaman/werkzeugkasten/tui/model"

	"github.com/charmbracelet/log"

	tea "github.com/charmbracelet/bubbletea"
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

	if cfg.tools {
		printTools(tools)
		logger.Debug("found tools", "count", len(tools.Tools))
		os.Exit(0)
	}
	if cfg.categories {
		printCategories(getCategories(tools))
		os.Exit(0)
	}
	if cfg.category != "" {
		result := getToolsByCategory(cfg.category, tools)
		if len(result.Tools) == 0 {
			logger.Warn("no results found", "category", cfg.category)
			os.Exit(0)
		}
		printTools(result)
		logger.Debug("found tools", "category", cfg.category, "count", len(result.Tools))
		os.Exit(0)
	}
	// interactive mode
	if len(cfg.toolList) == 0 {
		p := tea.NewProgram(model.InitialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	} else {
		// non-interactive mode
		installEget(cfg.downloadDir)
		for _, toolName := range cfg.toolList {
			err = downloadToolWithEget(cfg.downloadDir, tools.Tools[toolName])
			if err != nil {
				logger.Warn("could not download tool", "tool", toolName, "error", err)
				continue
			}
		}
	}
}
