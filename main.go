package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/allaman/werkzeugkasten/cli"
	"github.com/allaman/werkzeugkasten/tool"
	"github.com/allaman/werkzeugkasten/tui/model"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
	"github.com/lmittmann/tint"
)

func main() {
	cfg := cli.CLI()
	opts := &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.DateTime,
	}
	out := os.Stderr
	if cfg.Debug {
		opts.Level = slog.LevelDebug
		opts.AddSource = true
		out = os.Stdout
	}
	_, noColorEnv := os.LookupEnv("NO_COLOR")
	opts.NoColor = cfg.NoColor || noColorEnv || !term.IsTerminal(out.Fd())
	slog.SetDefault(slog.New(tint.NewTextHandler(out, opts)))
	tools, err := tool.CreateToolData()
	if err != nil {
		slog.Error("could not parse tools data", "error", err)
		os.Exit(1)
	}
	slog.Debug("download dir", "dir", cfg.DownloadDir)

	if cfg.Tools {
		tool.PrintTools(tools)
		slog.Debug("found tools", "count", len(tools.Tools))
		os.Exit(0)
	}
	if cfg.Categories {
		tool.PrintCategories(tool.GetCategories(tools))
		os.Exit(0)
	}
	if cfg.Category != "" {
		result := tool.GetToolsByCategory(cfg.Category, tools)
		if len(result.Tools) == 0 {
			slog.Warn("no results found", "category", cfg.Category)
			os.Exit(0)
		}
		tool.PrintTools(result)
		slog.Debug("found tools", "category", cfg.Category, "count", len(result.Tools))
		os.Exit(0)
	}
	// interactive mode
	if len(cfg.ToolList) == 0 {
		p := tea.NewProgram(model.InitialModel(tools, cfg))
		if _, err := p.Run(); err != nil {
			slog.Error("program exited with error", "error", err)
			os.Exit(1)
		}
	} else {
		// non-interactive mode
		if err := tool.InstallEget(); err != nil {
			slog.Error("could not install eget", "error", err)
			os.Exit(1)
		}
		failed := false
		for _, toolName := range cfg.ToolList {
			toolDef, ok := tools.Tools[toolName]
			if !ok {
				slog.Warn("unknown tool requested", "tool", toolName)
				failed = true
				continue
			}
			err = tool.DownloadToolWithEget(cfg.DownloadDir, toolDef, cfg.System)
			if err != nil {
				slog.Warn("could not download tool", "tool", toolName, "error", err)
				failed = true
				continue
			}
		}
		if failed {
			os.Exit(1)
		}
	}
}
