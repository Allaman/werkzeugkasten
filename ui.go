package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	workingDir    string
	selectedTools []string
)

// TODO: make configurable
var theme = huh.ThemeCatppuccin()

func formatToolString(name string, tool Tool) string {

	toolNameStyle := lipgloss.NewStyle().
		Foreground(theme.Focused.Title.GetForeground())

	descriptionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Help.ShortDesc.GetForeground())

	categoriesStyle := lipgloss.NewStyle().
		Foreground(theme.Blurred.MultiSelectSelector.GetForeground())

	styledToolName := toolNameStyle.Render(name)
	styledDescription := descriptionStyle.Render(tool.Description)
	styledCategories := categoriesStyle.Render(strings.Join(tool.Categories, ","))

	// when a tool version is explicitly set
	if tool.Tag != "" {
		versionStyle := lipgloss.NewStyle().
			Foreground(theme.Form.GetForeground())
		styledVersion := versionStyle.Render(tool.Tag)
		return fmt.Sprintf("%s:%s - %s [%s]", styledToolName, styledVersion, styledDescription, styledCategories)
	}

	return fmt.Sprintf("%s - %s [%s]", styledToolName, styledDescription, styledCategories)

}

func createToolOptions(tools Tools) []huh.Option[string] {
	sortedTools := make([]string, 0, len(tools.Tools))
	for k := range tools.Tools {
		sortedTools = append(sortedTools, k)
	}
	sort.Strings(sortedTools)
	options := make([]huh.Option[string], 0, len(tools.Tools))
	for _, name := range sortedTools {
		tool := tools.Tools[name]
		option := huh.NewOption(formatToolString(name, tool), name)
		options = append(options, option)
	}
	return options
}

func startUI(cfg cliConfig) {
	tools, err := createDefaultTools()
	if err != nil {
		logger.Error("could not parse tools data", "error", err)
		os.Exit(1)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Your Werkzeugkasten").
				Description("Installing awesome tools with ease!"),

			huh.NewInput().
				Title("Where should your tools be installed?").
				Description("Provide an absolute or relative path where your tools should be downloaded to. All directories are created if needed.").
				Prompt("#").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("you must provide a path")
					}
					return nil
				}).
				Value(&workingDir),
		),

		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Which tools do you want to install?").
				Description("Chose one or more tools to be downloaded to the specified path.").
				Options(createToolOptions(tools)...).
				Validate(func(t []string) error {
					if len(t) == 0 {
						return errors.New("you must select at least one tool")
					}
					return nil
				}).
				Value(&selectedTools),
		),
	).WithAccessible(cfg.accessible)

	form.WithTheme(theme)

	err = form.Run()

	if err != nil {
		logger.Fatal(err)
	}

	installDir, err := normalizePath(workingDir)
	if err != nil {
		logger.Error("could not normalize path")
		os.Exit(1)
	}

	start := func() {
		config := newDefaultConfig()
		if os.Getenv("WK_EGET_VERSION") != "" {
			version := os.Getenv("WK_EGET_VERSION")
			logger.Debug("setting eget version", "version", version)
			config.version = version
		}
		err := downloadEgetBinary(installDir, config)
		if err != nil {
			logger.Error("could not download eget binary", "error", err)
			os.Exit(1)
		}
		for _, t := range selectedTools {
			err = downloadToolWithEget(installDir, tools.Tools[t])
			if err != nil {
				logger.Warn("could not download tool", "error", err)
				continue
			}
		}
	}

	_ = spinner.New().Title("Downloading tools ...").Accessible(cfg.accessible).Action(start).Run()
	logger.Print(fmt.Sprintf("Run 'export PATH=$PATH:%s' to add your tools to the PATH", installDir))
}
