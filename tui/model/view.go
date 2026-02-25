package model

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"github.com/allaman/werkzeugkasten/tui/keys"
)

func (m *MainModel) View() tea.View {
	var content string
	switch m.CurrentView {
	case "tools":
		content = m.ToolsListView.View()
	case "detail":
		helpView := m.DetailView.Help.View(keys.DetailKeys)
		content = fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.DetailView.ViewPort.View(), m.footerView(), helpView)
	case "releases":
		content = m.ReleasesListView.View()
	case "processing":
		helpView := m.ProcessingModel.Help.View(keys.ProcessingKeys)
		content = fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.ProcessingModel.ViewPort.View(), m.footerView(), helpView)
	case "version":
		content = fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.showVersion(), m.footerView())
	default:
		content = "Unknown view"
	}

	view := tea.NewView(content)
	view.AltScreen = true
	view.WindowTitle = "Werkzeugkasten"
	return view
}
