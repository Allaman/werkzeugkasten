package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedTools []string
)

func formatToolString(theme *huh.Theme, name string, tool Tool) string {

	toolNameStyle := lipgloss.NewStyle().
		Foreground(theme.Focused.Title.GetForeground())

	descriptionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Help.ShortDesc.GetForeground())

	categoriesStyle := lipgloss.NewStyle().
		Foreground(theme.Blurred.MultiSelectSelector.GetForeground())

	styledToolName := toolNameStyle.Render(name)
	// TODO: How to handle long descriptions?
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

func createToolOptions(theme *huh.Theme, tools Tools) []huh.Option[string] {
	sortedTools := sortTools(tools)
	options := make([]huh.Option[string], 0, len(tools.Tools))
	for _, name := range sortedTools {
		tool := tools.Tools[name]
		option := huh.NewOption(formatToolString(theme, name, tool), name)
		options = append(options, option)
	}
	return options
}

func createForm(theme *huh.Theme, tools Tools) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Which tools do you want to install?").
				Description("Chose one or more tools to be downloaded.").
				Options(createToolOptions(theme, tools)...).
				Validate(func(t []string) error {
					if len(t) == 0 {
						return errors.New("you must select at least one tool")
					}
					return nil
				}).
				Value(&selectedTools),
		),
	)

	return form
}

func processSelectedTools(cfg cliConfig, tools Tools) func() {
	return func() {
		installEget(cfg.downloadDir)
		for _, t := range selectedTools {
			err := downloadToolWithEget(cfg.downloadDir, tools.Tools[t])
			if err != nil {
				logger.Warn("could not download tool", "tool", t, "error", err)
				continue
			}
		}
		logger.Info(fmt.Sprintf("Run 'export PATH=$PATH:%s' to add your tools to the PATH", cfg.downloadDir))
	}
}

func startUI(cfg cliConfig, tools Tools) {
	var theme *huh.Theme
	switch strings.ToLower(cfg.theme) {
	case "base16":
		theme = huh.ThemeBase16()
	case "catppuccin":
		theme = huh.ThemeCatppuccin()
	case "charm":
		theme = huh.ThemeCharm()
	case "dracula":
		theme = huh.ThemeDracula()
	default:
		logger.Warn("unknown theme. valid themes are 'base16', 'catppuccin' (default), 'charm', and 'dracula'")
	}

	form := createForm(theme, tools)
	form.WithAccessible(cfg.accessible)
	form.WithTheme(theme)

	err := form.Run()

	if err != nil {
		logger.Fatal(err)
	}

	start := processSelectedTools(cfg, tools)

	_ = spinner.New().Title("Downloading tools ...").Accessible(cfg.accessible).Action(start).Run()
}
