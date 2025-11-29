package timeframe

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	dayTimeslice     key.Binding
	weekTimeslice    key.Binding
	monthTimeslice   key.Binding
	quarterTimeslice key.Binding
	yearTimeslice    key.Binding
	lifeTimeslice    key.Binding
	nextPeriod       key.Binding
	previousPeriod   key.Binding
	currentPeriod    key.Binding
	gotoPeriod       key.Binding
	searchGoals      key.Binding
	goToParent       key.Binding
	unlinkParent     key.Binding
}

func NewListKeyMap() listKeyMap {
	return listKeyMap{
		dayTimeslice: key.NewBinding(
			key.WithKeys("d", "в"),
			key.WithHelp("d", "Day timeframe"),
		),
		weekTimeslice: key.NewBinding(
			key.WithKeys("w", "ц"),
			key.WithHelp("w", "Week timeframe"),
		),
		monthTimeslice: key.NewBinding(
			key.WithKeys("m", "ь"),
			key.WithHelp("m", "Month timeframe"),
		),
		quarterTimeslice: key.NewBinding(
			key.WithKeys("q", "й"),
			key.WithHelp("q", "Quarter timeframe"),
		),
		yearTimeslice: key.NewBinding(
			key.WithKeys("y", "н"),
			key.WithHelp("y", "Year timeframe"),
		),
		lifeTimeslice: key.NewBinding(
			key.WithKeys("L", "Д"),
			key.WithHelp("L", "Life timeframe"),
		),
		nextPeriod: key.NewBinding(
			key.WithKeys("right", "l", "д"),
			key.WithHelp("->", "Next period"),
		),
		previousPeriod: key.NewBinding(
			key.WithKeys("left", "h", "р"),
			key.WithHelp("<-", "Previous period"),
		),
		currentPeriod: key.NewBinding(
			key.WithKeys("t", "е"),
			key.WithHelp("t", "Current period"),
		),
		gotoPeriod: key.NewBinding(
			key.WithKeys("g", "п"),
			key.WithHelp("g", "Go to period"),
		),
		searchGoals: key.NewBinding(
			key.WithKeys("f", "/"),
			key.WithHelp("f", "Search goals"),
		),
		goToParent: key.NewBinding(
			key.WithKeys("p", "з"),
			key.WithHelp("p", "Go to parent goal"),
		),
		unlinkParent: key.NewBinding(
			key.WithKeys("u", "г"),
			key.WithHelp("u", "Unlink from parent"),
		),
	}
}
