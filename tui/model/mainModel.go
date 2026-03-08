package model

import (
	"fmt"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/allaman/werkzeugkasten/cli"
	"github.com/allaman/werkzeugkasten/tool"
	"github.com/allaman/werkzeugkasten/tui/item"
	"github.com/allaman/werkzeugkasten/tui/keys"
	"github.com/allaman/werkzeugkasten/tui/styles"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type MainModel struct {
	CurrentView        string
	ToolsListView      list.Model
	SelectedTool       item.Tool
	DetailView         Output
	ReleasesListView   list.Model
	CategoriesListView list.Model
	ActiveCategory     string
	ProcessingModel    Output
	ToolData           tool.ToolData
	config             cli.CLIConfig
	version            string
}

type Output struct {
	ItemName string
	ItemTag  string
	ViewPort viewport.Model
	Content  string
	Help     help.Model
}

func buildToolsList(toolData tool.ToolData, category string, width, height int) list.Model {
	filteredData := toolData
	if category != "" && category != "All" {
		filteredData = tool.GetToolsByCategory(category, toolData)
	}

	sortedTools := tool.SortTools(filteredData)
	items := make([]list.Item, 0, len(filteredData.Tools))
	for _, t := range sortedTools {
		identifier := filteredData.Tools[t].Identifier
		description := filteredData.Tools[t].Description
		items = append(items, item.NewItem(t, identifier, description))
	}

	l := list.New(items, list.NewDefaultDelegate(), width, height)
	l.Title = "Available Tools"
	if category != "" && category != "All" {
		l.Title = fmt.Sprintf("Available Tools [%s]", category)
	}
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return keys.ToolsKeys.ShortHelp()
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return keys.ToolsKeys.FullHelp()
	}
	return l
}

func buildCategoriesList(toolData tool.ToolData, width, height int) list.Model {
	categories := tool.GetCategories(toolData)
	sortedCategories := make([]string, 0, len(categories))
	for k := range categories {
		sortedCategories = append(sortedCategories, k)
	}
	slices.Sort(sortedCategories)

	items := make([]list.Item, 0, len(categories)+1)
	items = append(items, item.Category{Name: "All", Count: len(toolData.Tools)})
	for _, c := range sortedCategories {
		items = append(items, item.Category{Name: c, Count: categories[c]})
	}

	l := list.New(items, list.NewDefaultDelegate(), width, height)
	l.Title = "Categories"
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return keys.CategoryKeys.ShortHelp()
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return keys.CategoryKeys.ShortHelp()
	}
	return l
}

func InitialModel(toolData tool.ToolData, cfg cli.CLIConfig) *MainModel {
	l := buildToolsList(toolData, "", 0, 0)
	cl := buildCategoriesList(toolData, 0, 0)

	view := viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))

	return &MainModel{
		config:             cfg,
		CurrentView:        "tools",
		ToolsListView:      l,
		CategoriesListView: cl,
		ToolData:           toolData,
		DetailView:         Output{ViewPort: view, Help: help.New()},
		ReleasesListView:   list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		ProcessingModel:    Output{ViewPort: view, Help: help.New()},
		version:            cli.Version,
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
