package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	workingDir    string
	selectedTools []string
)

func formatToolString(name string, tool Tool) string {
	// Define the styles
	toolNameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF4757"))

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1E90FF"))

	descriptionStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#2ED573"))

	categoriesStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFAAA"))

	styledToolName := toolNameStyle.Render(name)
	styledVersion := versionStyle.Render(tool.Tag)
	styledDescription := descriptionStyle.Render(tool.Description)
	styledCategories := categoriesStyle.Render(strings.Join(tool.Categories, ","))

	return fmt.Sprintf("%s:%s - %s [%s]", styledToolName, styledVersion, styledDescription, styledCategories)

}

// TODO: sort alphabetically?
func createToolOptions(tools Tools) []huh.Option[string] {
	options := make([]huh.Option[string], 0, len(tools.Tools))
	for name, tool := range tools.Tools {
		option := huh.NewOption(formatToolString(name, tool), name)
		options = append(options, option)
	}
	return options
}

func startUI() {
	accessibleMode := os.Getenv("WK_ACCESSIBLE") != ""
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
				Prompt("#").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("Please provide a path!")
					}
					return nil
				}).
				Value(&workingDir),
		),

		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Which tools do you want to install?").
				Options(createToolOptions(tools)...).
				Validate(func(t []string) error {
					if len(t) == 0 {
						return errors.New("You must select at least one tool")
					}
					return nil
				}).
				Value(&selectedTools),
		),
	).WithAccessible(accessibleMode)

	form.WithTheme(huh.ThemeBase16())

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

	_ = spinner.New().Title("Downloading tools ...").Accessible(accessibleMode).Action(start).Run()
	logger.Print(fmt.Sprintf("Run 'export PATH=$PATH:%s'", installDir))
}