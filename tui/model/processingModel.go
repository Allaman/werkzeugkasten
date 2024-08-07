package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/allaman/werkzeugkasten/tui/item"
	"github.com/allaman/werkzeugkasten/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m ProcessingModel) headerView() string {
	title := styles.TitleStyle.Render("Installing", m.ItemName)
	line := strings.Repeat("─", max(0, m.DetailView.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m ProcessingModel) footerView() string {
	info := styles.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.DetailView.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.DetailView.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *MainModel) processSelectedItem(i item.Item) tea.Cmd {
	return func() tea.Msg {
		output := "Starting processing...\n"
		output += "Processing: " + i.Title() + "\n"
		time.Sleep(1000 * time.Millisecond)
		output += "Processed: " + i.Title() + "\n"
		output += "Processing complete.\n"
		return processUpdateMsg(output)
	}
}

type processUpdateMsg string
type tickMsg struct{}