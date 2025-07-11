package model

import (
	"github.com/allaman/werkzeugkasten/tui/item"
	tea "github.com/charmbracelet/bubbletea"
)

func openBrowserCmd(identifier string) tea.Cmd {
	return func() tea.Msg {
		err := item.OpenInBrowser(identifier)
		if err != nil {
			return browserErrMsg{err: err}
		}
		return browserSuccessMsg{}
	}
}

type browserSuccessMsg struct{}
type browserErrMsg struct {
	err error
}
