package model

import (
	"fmt"

	"github.com/allaman/werkzeugkasten/tui/item"
	"github.com/allaman/werkzeugkasten/tui/keys"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		m.List.SetHeight(msg.Height)
		m.DetailView.DetailView.Width = msg.Width - 4
		m.DetailView.DetailView.Height = msg.Height - 4
		m.ProcessingModel.DetailView.Width = msg.Width - 4
		m.ProcessingModel.DetailView.Height = msg.Height - 4
		// m.ProcessingModel.Width = msg.Width
		// m.ProcessingModel.Height = msg.Height
		// headerHeight := lipgloss.Height(m.headerView())
		// footerHeight := lipgloss.Height(m.footerView())
		// verticalMarginHeight := headerHeight + footerHeight
		// m.DetailView.Height = msg.Height - verticalMarginHeight
		// m.DetailView.Width = (msg.Width)
		// m.DetailView.YPosition = headerHeight
		// m.DetailView.YPosition = headerHeight + 1

	case tea.KeyMsg:
		if m.List.FilterState() == list.Filtering {
			break
		}

		switch m.CurrentView {

		case "list":
			switch {

			case key.Matches(msg, keys.Keys.Install):
				selectedItem, ok := m.List.SelectedItem().(item.Item)
				if ok {
					m.CurrentView = "processing"
					m.ProcessingModel.ItemName = selectedItem.Title()
					return m, m.processSelectedItem()
				}

			case key.Matches(msg, keys.Keys.Describe):
				selectedItem, ok := m.List.SelectedItem().(item.Item)
				if ok {
					m.CurrentView = "detail"
					m.DetailView.ItemName = selectedItem.Title()
					return m, fetchReadmeCmd(fmt.Sprintf("https://raw.githubusercontent.com/%s/main/README.md", selectedItem.Identifier()))
				}
			}

		case "detail":
			switch {
			case key.Matches(msg, keys.ViewPortKeys.Down):
				m.DetailView.DetailView.LineDown(1)
			case key.Matches(msg, keys.ViewPortKeys.Up):
				m.DetailView.DetailView.LineUp(1)
			case key.Matches(msg, keys.ViewPortKeys.HalfPageDown):
				m.DetailView.DetailView.HalfViewDown()
			case key.Matches(msg, keys.ViewPortKeys.HalfPageUp):
				m.DetailView.DetailView.HalfViewUp()
			case key.Matches(msg, keys.ViewPortKeys.Help):
				m.DetailView.Help.ShowAll = !m.DetailView.Help.ShowAll
			case key.Matches(msg, keys.ViewPortKeys.Install):
				selectedItem, ok := m.List.SelectedItem().(item.Item)
				if ok {
					m.CurrentView = "processing"
					m.ProcessingModel.ItemName = selectedItem.Title()
					return m, m.processSelectedItem()
				}
			case key.Matches(msg, keys.ViewPortKeys.Esc):
				m.CurrentView = "list"
				return m, nil
			}

		case "processing":
			if msg.String() == "esc" {
				m.CurrentView = "list"
				return m, nil
			}
		}

	case fetchReadmeSuccessMsg:
		m.DetailView.DetailView.SetContent(string(msg))
		m.DetailView.DetailView.GotoTop()
		return m, nil

	case fetchReadmeErrMsg:
		m.DetailView.DetailView.SetContent(msg.err.Error())
		return m, nil

	case processSuccessMsg:
		m.ProcessingModel.DetailView.SetContent(string(msg))
		m.ProcessingModel.DetailView.GotoTop()
		return m, nil

	case processErrMsg:
		m.ProcessingModel.DetailView.SetContent(msg.err.Error())
		return m, nil

	}

	var cmd tea.Cmd
	// TODO:
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}
