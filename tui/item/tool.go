package item

type Tool struct {
	name        string
	identifier  string
	description string
}

func (i Tool) Title() string       { return i.name }
func (i Tool) Identifier() string  { return i.identifier }
func (i Tool) Description() string { return i.description }

func NewItem(name, identifier, description string) Tool {
	return Tool{
		name:        name,
		description: description,
		identifier:  identifier,
	}
}

func (i Tool) FilterValue() string {
	return i.name + " " + i.description
}
