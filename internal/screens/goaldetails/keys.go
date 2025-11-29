package goaldetails

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	goBack   key.Binding
	openGoal key.Binding
}

func NewListKeyMap() listKeyMap {
	return listKeyMap{
		goBack: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Escape", "Go back"),
		),
		openGoal: key.NewBinding(
			key.WithKeys("o", "щ"),
			key.WithHelp("o", "Open goal"),
		),
	}
}
