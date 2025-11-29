package hierarchy

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	showAllTree key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		showAllTree: key.NewBinding(
			key.WithKeys("a", "ф"),
			key.WithHelp("a", "Show full tree"),
		),
	}
}
