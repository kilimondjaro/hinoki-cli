package goallist

import "hinoki-cli/internal/goal"

const (
	Timeframe = iota
	Subgoal
)

type GoalItem struct {
	goal.Goal
	mode int
}

func (i GoalItem) FilterValue() string {
	return i.Goal.Title
}
