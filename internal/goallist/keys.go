package goallist

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	markGoalDone key.Binding
	createGoal   key.Binding
	editGoal     key.Binding
	reloadGoals  key.Binding
	archiveGoal  key.Binding
	changeDate   key.Binding
}

func NewListKeyMap() listKeyMap {
	return listKeyMap{
		reloadGoals: key.NewBinding(
			key.WithKeys("r", "к"),
			key.WithHelp("r", "Reload goals"),
		),
		markGoalDone: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Mark goal done"),
		),
		createGoal: key.NewBinding(
			key.WithKeys("n", "т"),
			key.WithHelp("n", "Create new goal"),
		),
		editGoal: key.NewBinding(
			key.WithKeys("e", "у"),
			key.WithHelp("e", "Edit goal"),
		),
		archiveGoal: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("Backspace", "Archive goal"),
		),
		changeDate: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "Change date"),
		),
	}
}
