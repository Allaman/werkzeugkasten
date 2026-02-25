package model

import (
	"github.com/allaman/werkzeugkasten/tui/item"
	tea "charm.land/bubbletea/v2"
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
