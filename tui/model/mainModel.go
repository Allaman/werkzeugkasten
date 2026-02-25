package model

import (
	"fmt"
	"strings"

	"github.com/allaman/werkzeugkasten/cli"
	"github.com/allaman/werkzeugkasten/tool"
	"github.com/allaman/werkzeugkasten/tui/item"
	"github.com/allaman/werkzeugkasten/tui/keys"
	"github.com/allaman/werkzeugkasten/tui/styles"
	"charm.land/lipgloss/v2"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type MainModel struct {
	CurrentView      string
	ToolsListView    list.Model
	SelectedTool     item.Tool
	DetailView       Output
	ReleasesListView list.Model
	ProcessingModel  Output
	ToolData         tool.ToolData
	config           cli.CliConfig
	version          string
}

type Output struct {
	ItemName string
	ItemTag  string
	ViewPort viewport.Model
	Content  string
	Help     help.Model
}

func InitialModel(toolData tool.ToolData, cfg cli.CliConfig) *MainModel {
	items := make([]list.Item, 0, len(toolData.Tools))
	sortedTools := tool.SortTools(toolData)

	for _, tool := range sortedTools {
		identifier := toolData.Tools[tool].Identifier
		description := toolData.Tools[tool].Description
		items = append(items, item.NewItem(tool, identifier, description))
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Available Tools"
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return keys.ToolsKeys.ShortHelp()
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return keys.ToolsKeys.FullHelp()
	}

	view := viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))

	return &MainModel{
		config:           cfg,
		CurrentView:      "tools",
		ToolsListView:    l,
		ToolData:         toolData,
		DetailView:       Output{ViewPort: view, Help: help.New()},
		ReleasesListView: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		ProcessingModel:  Output{ViewPort: view, Help: help.New()},
		version:          cli.Version,
	}
}

func (m MainModel) headerView() string {
	var title, line string
	if m.CurrentView == "detail" {
		title = styles.TitleStyle.Render("README of", m.DetailView.ItemName)
		line = strings.Repeat("─", max(0, m.DetailView.ViewPort.Width()-lipgloss.Width(title)))
	}
	if m.CurrentView == "processing" {
		title = styles.TitleStyle.Render("Installing", m.ProcessingModel.ItemName)
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m MainModel) footerView() string {
	var info, line string
	if m.CurrentView == "detail" {
		info = styles.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.DetailView.ViewPort.ScrollPercent()*100))
		line = strings.Repeat("─", max(0, m.DetailView.ViewPort.Width()-lipgloss.Width(info)))
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *MainModel) showVersion() string {
	return cli.Version
}

func (m MainModel) Init() tea.Cmd {
	return nil
}
