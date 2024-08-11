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
		return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.DetailView.ViewPort.View(), m.footerView(), helpView)
	case "processing":
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.ProcessingModel.ViewPort.View(), m.footerView())
	case "version":
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.showVersion(), m.footerView())
	default:
		return "Unknown view"
	}
}
