package model

import (
	"fmt"
	"strings"

	"github.com/allaman/werkzeugkasten/tui/styles"
	"github.com/allaman/werkzeugkasten/tui/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/glamour"
)

func (m DetailModel) headerView() string {
	title := styles.TitleStyle.Render("README of", m.ItemName)
	line := strings.Repeat("─", max(0, m.DetailView.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m DetailModel) footerView() string {
	info := styles.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.DetailView.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.DetailView.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func fetchReadmeCmd(url string) tea.Cmd {
	return func() tea.Msg {
		content, err := utils.FetchReadme(url)
		if err != nil {
			fmt.Println("Error fetching README:", err)
			return fetchReadmeErrMsg(err)
		}

		renderedContent, err := glamour.Render(content, "dark")
		if err != nil {
			fmt.Println("Error rendering content:", err)
			return fetchReadmeErrMsg(err)
		}

		return fetchReadmeSuccessMsg(renderedContent)
	}
}

type fetchReadmeSuccessMsg string
type fetchReadmeErrMsg error
