package model

import (
	"log/slog"

	"github.com/allaman/werkzeugkasten/tui/item"

	tea "github.com/charmbracelet/bubbletea"
)

func fetchReleasesCmd(identifier string) tea.Cmd {
	return func() tea.Msg {
		// TODO: Pass GitHub token
		releases, err := item.FetchReleases(identifier, "")
		if err != nil {
			slog.Debug("error fetching releases", "error", err)
			return fetchReleasesErrMsg{err: err}
		}
		return fetchReleasesSuccessMsg(releases)
	}
}

type fetchReleasesSuccessMsg []item.FetchRelease
type fetchReleasesErrMsg struct {
	err error
}
