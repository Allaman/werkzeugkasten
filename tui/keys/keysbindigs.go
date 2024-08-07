package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Down         key.Binding
	Up           key.Binding
	Install      key.Binding
	Describe     key.Binding
	Quit         key.Binding
	HalfPageDown key.Binding
	HalfPageUp   key.Binding
	Esc          key.Binding
}

var Keys = KeyMap{
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install"),
	),
	Describe: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "describe"),
	),
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
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Install, k.Describe, k.Esc, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Install, k.Describe, k.Esc, k.Quit},
	}
}
