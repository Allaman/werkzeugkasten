package keys

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Install  key.Binding
	Releases key.Binding
	Describe key.Binding
	Browse   key.Binding
	Quit     key.Binding
	Esc      key.Binding
	Version  key.Binding
}

var ToolsKeys = KeyMap{
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install"),
	),
	Describe: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "describe"),
	),
	Releases: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "releases"),
	),
	Browse: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "browse"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Esc"),
	),
	Version: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "version"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Install}
}

func (k KeyMap) FullHelp() []key.Binding {
	return []key.Binding{k.Install, k.Describe, k.Releases, k.Browse, k.Esc, k.Version}
}
