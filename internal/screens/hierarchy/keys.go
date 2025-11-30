package hierarchy

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	showAllTree key.Binding
	cursorUp    key.Binding
	cursorDown  key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		showAllTree: key.NewBinding(
			key.WithKeys("a", "ф"),
			key.WithHelp("a", "Show full tree"),
		),
		cursorUp: key.NewBinding(
			key.WithKeys("up", "k", "л"),
			key.WithHelp("↑/k", "up"),
		),
		cursorDown: key.NewBinding(
			key.WithKeys("down", "j", "о"),
			key.WithHelp("↓/j", "down"),
		),
	}
}
