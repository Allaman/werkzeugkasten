package model

import (
	"fmt"

	"github.com/allaman/werkzeugkasten/tui/keys"
)

func (m *MainModel) View() string {
	switch m.CurrentView {
	case "list":
		return m.List.View()
	case "detail":
		helpView := m.DetailView.Help.View(keys.ViewPortKeys)
		return fmt.Sprintf("%s\n%s\n%s\n%s", m.DetailView.headerView(), m.DetailView.DetailView.View(), m.DetailView.footerView(), helpView)
	case "processing":
		// m.ProcessingModel.DetailView.SetContent(m.ProcessingModel.Spinner.View())
		return fmt.Sprintf("%s\n%s\n%s", m.ProcessingModel.headerView(), m.ProcessingModel.DetailView.View(), m.ProcessingModel.footerView())
	default:
		return "Unknown view"
	}
}
