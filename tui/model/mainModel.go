package model

import (
	"fmt"
	"strings"

	"github.com/allaman/werkzeugkasten/cli"
	"github.com/allaman/werkzeugkasten/tool"
	"github.com/allaman/werkzeugkasten/tui/item"
	"github.com/allaman/werkzeugkasten/tui/keys"
	"github.com/allaman/werkzeugkasten/tui/styles"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	CurrentView     string
	List            list.Model
	DetailView      Output
	ProcessingModel Output
	ToolData        tool.ToolData
	config          cli.CliConfig
	version         string
}

type Output struct {
	ItemName string
	ViewPort viewport.Model
	Content  string
	Help     help.Model
}

func InitialModel(toolData tool.ToolData) *MainModel {
	items := make([]list.Item, 0, len(toolData.Tools))
	sortedTools := tool.SortTools(toolData)

	for _, tool := range sortedTools {
		identifier := toolData.Tools[tool].Identifier
		description := toolData.Tools[tool].Description
		items = append(items, item.NewItem(tool, identifier, description))
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Werkzeugkasten"
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.Keys.Install,
			keys.Keys.Describe,
		}
	}

	view := viewport.New(80, 20)

	s := spinner.New()
	s.Spinner = spinner.Dot

	return &MainModel{
		CurrentView:     "list",
		List:            l,
		ToolData:        toolData,
		DetailView:      Output{ViewPort: view, Help: help.New()},
		ProcessingModel: Output{ViewPort: view},
		version:         cli.Version,
	}
}

func (m MainModel) headerView() string {
	var title string
	if m.CurrentView == "detail" {
		title = styles.TitleStyle.Render("README of", m.DetailView.ItemName)
	}
	if m.CurrentView == "processing" {
		title = styles.TitleStyle.Render("Installing", m.ProcessingModel.ItemName)
	}
	line := strings.Repeat("─", max(0, m.DetailView.ViewPort.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m MainModel) footerView() string {
	info := styles.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.DetailView.ViewPort.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.DetailView.ViewPort.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m MainModel) Init() tea.Cmd {
	return tea.SetWindowTitle("Werkzeugkasten")
}
