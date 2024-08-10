package model

import (
	"github.com/allaman/werkzeugkasten/tool"

	tea "github.com/charmbracelet/bubbletea"
)

// func (m MainModel) headerView() string {
// 	title := styles.TitleStyle.Render("Installing", m.ItemName)
// 	line := strings.Repeat("─", max(0, m.DetailView.Width-lipgloss.Width(title)))
// 	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
// }
//
// func (m MainModel) footerView() string {
// 	info := styles.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.DetailView.ScrollPercent()*100))
// 	line := strings.Repeat("─", max(0, m.DetailView.Width-lipgloss.Width(info)))
// 	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
// }

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
