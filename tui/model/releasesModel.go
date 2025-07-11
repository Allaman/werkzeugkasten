package model

import (
	"log/slog"
	"os"

	"github.com/allaman/werkzeugkasten/tui/item"

	tea "github.com/charmbracelet/bubbletea"
)

func fetchReleasesCmd(identifier string) tea.Cmd {
	return func() tea.Msg {
		token := os.Getenv("EGET_GITHUB_TOKEN")
		releases, err := item.FetchReleases(identifier, token)
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
