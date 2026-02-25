package keys

import "charm.land/bubbles/v2/key"

type detailsKeyMap struct {
	Down         key.Binding
	Up           key.Binding
	HalfPageDown key.Binding
	HalfPageUp   key.Binding
	Help         key.Binding
	Install      key.Binding
	Esc          key.Binding
	Quit         key.Binding
}

var DetailKeys = detailsKeyMap{
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
	HalfPageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	HalfPageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install"),
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

func (k detailsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Down, k.Up, k.Install, k.Help, k.Esc, k.Quit}
}

func (k detailsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Down, k.Up, k.Esc, k.Help, k.Quit},
		{k.HalfPageUp, k.HalfPageDown},
	}
}
