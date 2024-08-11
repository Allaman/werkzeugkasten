package model

import (
	"github.com/allaman/werkzeugkasten/tool"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *MainModel) processSelectedItem() tea.Cmd {
	return func() tea.Msg {
		tool.InstallEget(m.config.DownloadDir)
		err := tool.DownloadToolWithEget(m.config.DownloadDir, m.ToolData.Tools[m.ProcessingModel.ItemName])
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
