package goallist

import (
	"fmt"
	"hinoki-cli/internal/dates"
	"hinoki-cli/internal/theme"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GoalItemDelegate struct {
	keys listKeyMap
}

var (
	doneItemStyle     = lipgloss.NewStyle().Foreground(theme.TextDisabled())
	selectedItemStyle = lipgloss.NewStyle().Foreground(theme.TextSelected())
	parentStyle       = lipgloss.NewStyle().Foreground(theme.TextMuted())
)

func (d GoalItemDelegate) Height() int {
	return 2
}

func (d GoalItemDelegate) Spacing() int { return 1 }
func (d GoalItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	item, ok := m.SelectedItem().(GoalItem)

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
	i, ok := listItem.(GoalItem)
	if !ok {
		return
	}

	itemStyle := lipgloss.NewStyle().Foreground(theme.TextPrimary())

	checkmark := " "
	if i.IsDone {
		checkmark = "x"
		itemStyle = doneItemStyle
	}

	dateTime := ""

	if i.mode == Subgoal && i.Date != nil && i.Timeframe != nil {
		dateTime = "\n    " + dates.DateString(*i.Date, *i.Timeframe)
	}

	str := fmt.Sprintf("[%s] %s%s", checkmark, i.Title, parentStyle.Render(dateTime))

	if i.ParentId != nil && i.ParentTitle != nil {
		str = fmt.Sprintf("%s\n    %s", str, parentStyle.Render(*i.ParentTitle))
	}

	if index == m.Index() {
		itemStyle = selectedItemStyle
	}

	wrapped := lipgloss.NewStyle().
		Width(m.Width()).
		Render(str)

	fmt.Fprint(w, itemStyle.Render(wrapped))
}
