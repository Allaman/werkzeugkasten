package keys

import "github.com/charmbracelet/bubbles/key"

type processKeyMap struct {
	Down key.Binding
	Up   key.Binding
	Esc  key.Binding
	Quit key.Binding
}

var ProcessingKeys = detailsKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Esc"),
	),
}

func (k processKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Esc, k.Quit}
}

func (k processKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Esc, k.Quit},
	}
}
