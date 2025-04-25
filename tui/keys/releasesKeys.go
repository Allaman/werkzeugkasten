package keys

import "github.com/charmbracelet/bubbles/key"

type releasesKeyMap struct {
	Down key.Binding
	Up   key.Binding
	Help key.Binding
	Esc  key.Binding
	Quit key.Binding
}

var ReleasesKeys = releasesKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "scroll down"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "scroll up"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Esc"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

func (k releasesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Down, k.Up, k.Esc}
}

func (k releasesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Down, k.Up, k.Esc, k.Help, k.Quit},
	}
}
