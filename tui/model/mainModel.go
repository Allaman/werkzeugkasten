package model

import (
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
}

func InitialModel() *MainModel {
	items := []list.Item{
		item.NewItem("Item 1", "allaman/werkzeugkasten", "Description for Item 1"),
		item.NewItem("Item 2", "allaman/werkzeugkasten", "Description for Item 3"),
		item.NewItem("Item 3", "allaman/werkzeugkasten", "Description for Item 3"),
		item.NewItem("Item 4", "allaman/werkzeugkasten", "Description for Item 4"),
		item.NewItem("Item 5", "allaman/werkzeugkasten", "Description for Item 5"),
		item.NewItem("Item 6", "allaman/werkzeugkasten", "Description for Item 6"),
		item.NewItem("Item 7", "allaman/werkzeugkasten", "Description for Item 7"),
		item.NewItem("Item 8", "allaman/werkzeugkasten", "Description for Item 8"),
		item.NewItem("Item 9", "allaman/werkzeugkasten", "Description for Item 9"),
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
		DetailView:      DetailModel{DetailView: detailView, Help: help.New()},
		ProcessingModel: ProcessingModel{DetailView: processingView},
	}
}

func (m MainModel) Init() tea.Cmd {
	return tea.SetWindowTitle("Werkzeugkasten")
}
