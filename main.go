package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/allaman/werkzeugkasten/tool"
	"github.com/allaman/werkzeugkasten/tui/model"

	tea "github.com/charmbracelet/bubbletea"
)

// will be overwritten in release pipeline
var version = "dev"

func main() {
	cfg := cli()
	if cfg.debug {
		opts := &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					if t, ok := a.Value.Any().(time.Time); ok {
						return slog.String(a.Key, t.Format("2006-01-02 15:04:05"))
					}
				}
				return a
			},
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
		slog.SetDefault(logger)
	}
	tools, err := tool.CreateToolData()
	if err != nil {
		slog.Error("could not parse tools data", "error", err)
		os.Exit(1)
	}
	slog.Debug("download dir", "dir", cfg.downloadDir)

	if cfg.tools {
		tool.PrintTools(tools)
		slog.Debug("found tools", "count", len(tools.Tools))
		os.Exit(0)
	}
	if cfg.categories {
		tool.PrintCategories(tool.GetCategories(tools))
		os.Exit(0)
	}
	if cfg.category != "" {
		result := tool.GetToolsByCategory(cfg.category, tools)
		if len(result.Tools) == 0 {
			slog.Warn("no results found", "category", cfg.category)
			os.Exit(0)
		}
		tool.PrintTools(result)
		slog.Debug("found tools", "category", cfg.category, "count", len(result.Tools))
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
		tool.InstallEget(cfg.downloadDir)
		for _, toolName := range cfg.toolList {
			err = tool.DownloadToolWithEget(cfg.downloadDir, tools.Tools[toolName])
			if err != nil {
				slog.Warn("could not download tool", "tool", toolName, "error", err)
				continue
			}
		}
	}
}
