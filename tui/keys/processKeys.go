package keys

import "github.com/charmbracelet/bubbles/key"

type processKeyMap struct {
	Esc  key.Binding
	Quit key.Binding
}

var ProcessingKeys = processKeyMap{
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Esc"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
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
