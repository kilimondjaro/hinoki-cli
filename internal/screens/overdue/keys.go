package overdue

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	openGoal     key.Binding
	assignParent key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		openGoal: key.NewBinding(
			key.WithKeys("o", "щ"),
			key.WithHelp("o", "Open goal in timeframe"),
		),
		assignParent: key.NewBinding(
			key.WithKeys("p", "з"),
			key.WithHelp("p", "Assign parent"),
		),
	}
}
