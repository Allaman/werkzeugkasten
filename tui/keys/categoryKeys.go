package keys

import "charm.land/bubbles/v2/key"

type categoryKeyMap struct {
	Select key.Binding
	Esc    key.Binding
}

var CategoryKeys = categoryKeyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select category"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}

func (k categoryKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Esc}
}

func (k categoryKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Esc},
	}
}
