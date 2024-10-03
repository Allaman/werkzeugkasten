package item

type Tool struct {
	name        string
	identifier  string
	description string
}

type Item Tool

func (i Item) Title() string       { return i.name }
func (i Item) Identifier() string  { return i.identifier }
func (i Item) Description() string { return i.description }

func NewItem(name, identifier, description string) Item {
	return Item{
		name:        name,
		description: description,
		identifier:  identifier,
	}
}

func (i Item) FilterValue() string {
	return i.name + " " + i.description
}
