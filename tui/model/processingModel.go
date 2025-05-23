package model

import (
	"github.com/allaman/werkzeugkasten/tool"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *MainModel) processSelectedItem() tea.Cmd {
	return func() tea.Msg {
		tool.InstallEget(m.config.DownloadDir)
		item := m.ToolData.Tools[m.ProcessingModel.ItemName]
		if m.ProcessingModel.ItemTag != "" {
			item.Tag = m.ProcessingModel.ItemTag
		}
		err := tool.DownloadToolWithEget(m.config.DownloadDir, item)
		if err != nil {
			return processErrMsg{err: err}
		}
		return processSuccessMsg("Install complete.\n")
	}
}

type processSuccessMsg string
type processErrMsg struct {
	err error
}
