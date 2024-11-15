package goallist

type GoalItem struct {
	id     string
	title  string
	isDone bool
}

func (i GoalItem) FilterValue() string {
	return i.title
}
