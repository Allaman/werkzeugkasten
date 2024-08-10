package model

import (
	"github.com/allaman/werkzeugkasten/cli"
	"github.com/allaman/werkzeugkasten/tool"
	"github.com/allaman/werkzeugkasten/tui/item"
	"github.com/allaman/werkzeugkasten/tui/keys"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	CurrentView     string
	List            list.Model
	DetailView      DetailModel
	ProcessingModel ProcessingModel
	ToolData        tool.ToolData
	config          cli.CliConfig
}

type ProcessingModel struct {
	ItemName   string
	DetailView viewport.Model
}

type DetailModel struct {
	Content    string
	Help       help.Model
	DetailView viewport.Model
	ItemName   string
	// Tool       tool.Tool
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

	detailView := viewport.New(80, 20) // Start with a reasonable size

	processingView := viewport.New(80, 20) // Start with a reasonable size

	return &MainModel{
		CurrentView:     "list",
		List:            l,
		ToolData:        toolData,
		DetailView:      DetailModel{DetailView: detailView, Help: help.New()},
		ProcessingModel: ProcessingModel{DetailView: processingView},
	}
}

func (m MainModel) Init() tea.Cmd {
	return tea.SetWindowTitle("Werkzeugkasten")
}
