package item

type Tool struct {
	title       string
	identifier  string
	description string
}

type Item Tool

func (i Item) Title() string       { return i.title }
func (i Item) Identifier() string  { return i.identifier }
func (i Item) Description() string { return i.description }

func NewItem(title, identifier, description string) Item {
	return Item{
		title:       title,
		description: description,
		identifier:  identifier,
	}
}

func (i Item) FilterValue() string {
	return i.title + " " + i.description
}
