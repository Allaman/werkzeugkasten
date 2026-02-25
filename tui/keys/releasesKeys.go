package keys

import "charm.land/bubbles/v2/key"

type releasesKeyMap struct {
	Install key.Binding
	Esc     key.Binding
}

var ReleasesKeys = releasesKeyMap{
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install release"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Esc"),
	),
}

func (k releasesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Install, k.Esc}
}

func (k releasesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Install, k.Esc},
	}
}
