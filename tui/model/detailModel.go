package model

import (
	"log/slog"

	"github.com/allaman/werkzeugkasten/tui/utils"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/glamour"
)

func fetchReadmeCmd(url string) tea.Cmd {
	return func() tea.Msg {
		content, err := utils.FetchReadme(url)
		if err != nil {
			slog.Debug("error fetching README", "error", err)
			return fetchReadmeErrMsg{err: err}
		}

		renderedContent, err := glamour.Render(content, "dark")
		if err != nil {
			slog.Debug("error rendering content", "error", err)
			return fetchReadmeErrMsg{err: err}
		}

		return fetchReadmeSuccessMsg(renderedContent)
	}
}

type fetchReadmeSuccessMsg string
type fetchReadmeErrMsg struct {
	err error
}
