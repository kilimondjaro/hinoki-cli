package goallist

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hinoki-cli/internal/goal"
	"io"
)

type GoalItemDelegate struct {
	keys listKeyMap
}

var (
	doneItemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	itemDarkStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	itemLightStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000"))
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
)

func (d GoalItemDelegate) Height() int  { return 1 }
func (d GoalItemDelegate) Spacing() int { return 1 }
func (d GoalItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	item, ok := m.SelectedItem().(goal.Goal)

	if !ok {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys.markGoalDone):
			item.IsDone = true
		}
	}
	return nil
}

func (d GoalItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(goal.Goal)
	if !ok {
		return
	}

	fn := itemLightStyle.Render
	if lipgloss.HasDarkBackground() {
		fn = itemDarkStyle.Render
	}

	checkmark := " "
	if i.IsDone {
		checkmark = "x"
		fn = doneItemStyle.Render
	}

	str := fmt.Sprintf("[%s] %s", checkmark, i.Title)

	if index == m.Index() {
		fn = selectedItemStyle.Render
	}

	fmt.Fprint(w, fn(str))
}
