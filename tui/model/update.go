package model

import (
	"fmt"
	"time"

	"github.com/allaman/werkzeugkasten/tui/item"
	"github.com/allaman/werkzeugkasten/tui/keys"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ToolsListView.SetWidth(msg.Width)
		m.ToolsListView.SetHeight(msg.Height)
		m.DetailView.ViewPort.SetWidth(msg.Width - 4)
		m.DetailView.ViewPort.SetHeight(msg.Height - 4)
		m.ProcessingModel.ViewPort.SetWidth(msg.Width - 4)
		m.ProcessingModel.ViewPort.SetHeight(msg.Height - 4)

	case tea.KeyPressMsg:
		if m.ToolsListView.FilterState() == list.Filtering {
			break
		}

		switch m.CurrentView {

		case "tools":
			switch {

			case key.Matches(msg, keys.ToolsKeys.Install):
				selectedItem, ok := m.ToolsListView.SelectedItem().(item.Tool)
				if ok {
					m.CurrentView = "processing"
					m.ProcessingModel.ItemName = selectedItem.Title()
					return m, m.processSelectedItem()
				}

			case key.Matches(msg, keys.ToolsKeys.Describe):
				selectedItem, ok := m.ToolsListView.SelectedItem().(item.Tool)
				if ok {
					m.CurrentView = "detail"
					m.DetailView.ItemName = selectedItem.Title()
					return m, fetchReadmeCmd(fmt.Sprintf("https://raw.githubusercontent.com/%s/%%s/%%s", selectedItem.Identifier()))
				}

			case key.Matches(msg, keys.ToolsKeys.Releases):
				selectedItem, ok := m.ToolsListView.SelectedItem().(item.Tool)
				if ok {
					m.CurrentView = "releases"
					m.SelectedTool = selectedItem
					m.ReleasesListView.Title = fmt.Sprintf("Releases of %s (max. last 100)", selectedItem.Title())
					m.ReleasesListView.AdditionalFullHelpKeys = func() []key.Binding {
						return []key.Binding{
							keys.ReleasesKeys.Install,
							keys.ReleasesKeys.Esc,
						}
					}
					return m, fetchReleasesCmd(selectedItem.Identifier())
				}

			case key.Matches(msg, keys.ToolsKeys.Browse):
				selectedItem, ok := m.ToolsListView.SelectedItem().(item.Tool)
				if ok {
					return m, openBrowserCmd(selectedItem.Identifier())
				}

			case key.Matches(msg, keys.ToolsKeys.Version):
				m.CurrentView = "version"
				return m, nil
			}

		case "detail":
			switch {
			case key.Matches(msg, keys.DetailKeys.Down):
				m.DetailView.ViewPort.ScrollDown(1)
			case key.Matches(msg, keys.DetailKeys.Up):
				m.DetailView.ViewPort.ScrollUp(1)
			case key.Matches(msg, keys.DetailKeys.HalfPageDown):
				m.DetailView.ViewPort.HalfPageDown()
			case key.Matches(msg, keys.DetailKeys.HalfPageUp):
				m.DetailView.ViewPort.HalfPageUp()
			case key.Matches(msg, keys.DetailKeys.Help):
				m.DetailView.Help.ShowAll = !m.DetailView.Help.ShowAll
			case key.Matches(msg, keys.DetailKeys.Install):
				selectedItem, ok := m.ToolsListView.SelectedItem().(item.Tool)
				if ok {
					m.CurrentView = "processing"
					m.ProcessingModel.ItemName = selectedItem.Title()
					return m, m.processSelectedItem()
				}
			case key.Matches(msg, keys.DetailKeys.Esc):
				m.CurrentView = "tools"
				return m, nil
			}

		case "releases":
			switch {
			case key.Matches(msg, keys.ReleasesKeys.Install):
				selectedItem, ok := m.ReleasesListView.SelectedItem().(item.Release)
				if ok {
					m.CurrentView = "processing"
					m.ProcessingModel.ItemName = m.SelectedTool.Title()
					m.ProcessingModel.ItemTag = selectedItem.Tag
					return m, m.processSelectedItem()
				}
			case key.Matches(msg, keys.ReleasesKeys.Esc):
				m.CurrentView = "tools"
				return m, nil
			}

		case "processing":
			switch {
			case key.Matches(msg, keys.ProcessingKeys.Esc):
				m.CurrentView = "tools"
				return m, nil
			case key.Matches(msg, keys.ProcessingKeys.Quit):
				return m, tea.Quit
			}

		case "version":
			if msg.String() == "esc" {
				m.CurrentView = "tools"
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

	case fetchReleasesSuccessMsg:
		releases := []item.FetchRelease(msg)
		items := make([]list.Item, 0, len(releases))

		for _, release := range releases {
			items = append(items, item.Release{Tag: release.TagName, PublishedAt: release.PublishedAt})
		}
		// Update the existing list with the new items
		cmd := m.ReleasesListView.SetItems(items)

		// Make sure dimensions are correct
		m.ReleasesListView.SetWidth(m.ToolsListView.Width())
		m.ReleasesListView.SetHeight(m.ToolsListView.Height())

		// Return the command from SetItems to ensure proper updates
		return m, cmd

	case fetchReleasesErrMsg:
		errorItem := item.NewRelease("Error: "+msg.err.Error(), time.Now())
		m.ReleasesListView.SetItems([]list.Item{errorItem})
		return m, nil

	case processSuccessMsg:
		m.ProcessingModel.ViewPort.SetContent(string(msg))
		m.ProcessingModel.ViewPort.GotoTop()
		return m, nil

	case processErrMsg:
		m.ProcessingModel.ViewPort.SetContent(msg.err.Error())
		return m, nil

	case browserSuccessMsg:
		return m, nil

	case browserErrMsg:
		// TODO: handle error
		return m, nil

	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.CurrentView {
	case "tools":
		newToolsListView, toolsCmd := m.ToolsListView.Update(msg)
		m.ToolsListView = newToolsListView
		if toolsCmd != nil {
			cmds = append(cmds, toolsCmd)
		}
	case "releases":
		newReleasesListView, releasesCmd := m.ReleasesListView.Update(msg)
		m.ReleasesListView = newReleasesListView
		if releasesCmd != nil {
			cmds = append(cmds, releasesCmd)
		}
	}

	if len(cmds) > 0 {
		cmd = tea.Batch(cmds...)
	}
	return m, cmd
}
