package model

import (
	"fmt"
	"log/slog"

	"github.com/allaman/werkzeugkasten/tui/utils"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/glamour"
)

func fetchReleasesCmd(identifier string) tea.Cmd {
	return func() tea.Msg {
		// TODO: Pass GitHub token
		releases, err := utils.FetchReleases(identifier, "")
		if err != nil {
			slog.Debug("error fetching releases", "error", err)
			return fetchReleasesErrMsg{err: err}
		}

		content := "| Release | Published |\n| --- | --- |\n"
		for _, release := range releases {
			formattedTime := release.PublishedAt.Format("2006-01-02")
			content += fmt.Sprintf("| %s | %s |\n", release.Name, formattedTime)
		}

		renderedContent, err := glamour.Render(content, "dark")
		if err != nil {
			slog.Debug("error rendering content", "error", err)
			return fetchReleasesErrMsg{err: err}
		}

		return fetchReleasesSuccessMsg(renderedContent)
	}
}

type fetchReleasesSuccessMsg string
type fetchReleasesErrMsg struct {
	err error
}
