package hierarchy

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	showAllTree   key.Binding
	cursorUp      key.Binding
	cursorDown    key.Binding
	openDetails   key.Binding
	openTimeframe key.Binding
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
		openDetails: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Open goal details"),
		),
		openTimeframe: key.NewBinding(
			key.WithKeys("o", "щ"),
			key.WithHelp("o", "Open timeframe"),
		),
	}
}
