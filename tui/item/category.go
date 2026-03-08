package item

import "fmt"

type Category struct {
	Name  string
	Count int
}

func (c Category) Title() string       { return c.Name }
func (c Category) Description() string { return fmt.Sprintf("%d tools", c.Count) }
func (c Category) FilterValue() string { return c.Name }
