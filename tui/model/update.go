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
		m.DetailView.ViewPort.Width = msg.Width - 4
		m.DetailView.ViewPort.Height = msg.Height - 4
		m.ProcessingModel.ViewPort.Width = msg.Width - 4
		m.ProcessingModel.ViewPort.Height = msg.Height - 4

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

			case key.Matches(msg, keys.Keys.Version):
				m.CurrentView = "version"
				return m, nil
			}

		case "detail":
			switch {
			case key.Matches(msg, keys.ViewPortKeys.Down):
				m.DetailView.ViewPort.LineDown(1)
			case key.Matches(msg, keys.ViewPortKeys.Up):
				m.DetailView.ViewPort.LineUp(1)
			case key.Matches(msg, keys.ViewPortKeys.HalfPageDown):
				m.DetailView.ViewPort.HalfViewDown()
			case key.Matches(msg, keys.ViewPortKeys.HalfPageUp):
				m.DetailView.ViewPort.HalfViewUp()
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

		case "version":
			if msg.String() == "esc" {
				m.CurrentView = "list"
				return m, nil
			}
		}

	case fetchReadmeSuccessMsg:
		m.DetailView.ViewPort.SetContent(string(msg))
		m.DetailView.ViewPort.GotoTop()
		return m, nil

	case fetchReadmeErrMsg:
		m.DetailView.ViewPort.SetContent(msg.err.Error())
		return m, nil

	case processSuccessMsg:
		m.ProcessingModel.ViewPort.SetContent(string(msg))
		m.ProcessingModel.ViewPort.GotoTop()
		return m, nil

	case processErrMsg:
		m.ProcessingModel.ViewPort.SetContent(msg.err.Error())
		return m, nil

	}

	var cmd tea.Cmd
	// TODO:
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}
