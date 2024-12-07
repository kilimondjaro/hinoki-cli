package goaldetails

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	goBack key.Binding
}

func NewListKeyMap() listKeyMap {
	return listKeyMap{
		goBack: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Escape", "Go back"),
		),
	}
}
