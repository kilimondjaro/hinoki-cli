package goallist

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	markGoalDone     key.Binding
	createGoal       key.Binding
	editGoal         key.Binding
	dayTimeslice     key.Binding
	weekTimeslice    key.Binding
	monthTimeslice   key.Binding
	quarterTimeslice key.Binding
	yearTimeslice    key.Binding
	lifeTimeslice    key.Binding
	nextPeriod       key.Binding
	previousPeriod   key.Binding
	reloadGoals      key.Binding
	archiveGoal      key.Binding
	currentPeriod    key.Binding
	gotoPeriod       key.Binding
	changeDate       key.Binding
}

func NewListKeyMap() listKeyMap {
	return listKeyMap{
		reloadGoals: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "Reload goals"),
		),
		markGoalDone: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Mark goal done"),
		),
		createGoal: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "Create new goal"),
		),
		editGoal: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "Edit goal"),
		),
		dayTimeslice: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "Day timeframe"),
		),
		weekTimeslice: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "Week timeframe"),
		),
		monthTimeslice: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "Month timeframe"),
		),
		quarterTimeslice: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "Quarter timeframe"),
		),
		yearTimeslice: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "Year timeframe"),
		),
		lifeTimeslice: key.NewBinding(
			key.WithKeys("L"),
			key.WithHelp("L", "Life timeframe"),
		),
		nextPeriod: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("->", "Next period"),
		),
		previousPeriod: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("<-", "Previous period"),
		),
		archiveGoal: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("Backspace", "Archive goal"),
		),
		currentPeriod: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "Current period"),
		),
		gotoPeriod: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "Go to period"),
		),
		changeDate: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "Change date"),
		),
	}
}
