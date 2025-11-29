package overdue

import (
	"time"

	"hinoki-cli/internal/goal"
	"hinoki-cli/internal/goallist"
	"hinoki-cli/internal/repository"
	"hinoki-cli/internal/screens"
	"hinoki-cli/internal/theme"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type OverdueScreen struct {
	list goallist.GoalList
	keys keyMap

	width, height int
}

var (
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.TextPrimary()).
		MarginBottom(2).
		PaddingTop(2)
)

const (
	maxWidth = 130
)

func NewOverdueScreen() screens.Screen {
	goalList := goallist.NewGoalList(nil, nil)
	goalList.SetDisplayMode(goallist.Overdue)

	return &OverdueScreen{
		list: goalList,
		keys: newKeyMap(),
	}
}

func (m *OverdueScreen) Init() tea.Cmd {
	return m.getOverdueGoalsCmd()
}

func (m *OverdueScreen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.list.IsInActiveState() {
			cmd := m.handleKeyMsg(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	case goallist.AddGoalSuccess, goallist.UpdateGoalSuccess:
		cmds = append(cmds, m.getOverdueGoalsCmd())
	case error:
		// swallow errors in UI loop
	}

	// Always pass messages to the list so it can handle GoalsResult and other messages
	cmds = append(cmds, m.list.Update(msg))

	return tea.Batch(cmds...)
}

func (m *OverdueScreen) View() string {
	header := headerStyle.Render("Overdue Goals")

	headerHeight := lipgloss.Height(header)
	listHeight := m.height - headerHeight

	style := lipgloss.NewStyle().PaddingLeft(2)
	horizontalPadding := (m.width - maxWidth) / 2

	if m.width > maxWidth {
		style = style.PaddingLeft(horizontalPadding).PaddingRight(horizontalPadding)
	}

	contentWidth := min(m.width, maxWidth)
	m.list.SetSize(contentWidth, listHeight)

	body := m.list.View()

	view := lipgloss.JoinVertical(lipgloss.Left, header, body)

	return style.
		SetString(view).
		Render()
}

func (m *OverdueScreen) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *OverdueScreen) Refresh() tea.Cmd {
	return m.getOverdueGoalsCmd()
}

func (m *OverdueScreen) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	if m.list.IsInActiveState() {
		return nil
	}

	switch {
	case msg.Type == tea.KeyEsc:
		return func() tea.Msg {
			return screens.GoBack{}
		}
	case key.Matches(msg, m.keys.openGoal):
		selectedGoal := m.list.GetSelectedGoal()
		if selectedGoal == nil {
			return nil
		}
		return m.openGoalInTimeframeCmd(selectedGoal)
	case key.Matches(msg, m.keys.assignParent):
		selectedGoal := m.list.GetSelectedGoal()
		if selectedGoal == nil {
			return nil
		}
		// If goal has a parent, go to parent. Otherwise, open search to assign parent.
		if selectedGoal.ParentId != nil {
			return m.goToParentGoalCmd(*selectedGoal.ParentId)
		}
		// Open search screen for parent assignment
		return func() tea.Msg {
			return screens.OpenSearchScreenForParent{
				GoalID: selectedGoal.ID,
			}
		}
	}

	return nil
}

func (m *OverdueScreen) openGoalInTimeframeCmd(selectedGoal *goal.Goal) tea.Cmd {
	return func() tea.Msg {
		// Check if the goal has a timeframe and date
		if selectedGoal.Timeframe == nil {
			return nil
		}

		// Life goals don't have dates, so use current time
		var date time.Time
		if selectedGoal.Date != nil {
			date = *selectedGoal.Date
		} else if *selectedGoal.Timeframe == goal.Life {
			date = time.Now()
		} else {
			// Other timeframes require a date
			return nil
		}

		return screens.OpenTimeframeScreenWithGoal{
			Timeframe: *selectedGoal.Timeframe,
			Date:      date,
			GoalID:    selectedGoal.ID,
		}
	}
}

func (m *OverdueScreen) goToParentGoalCmd(parentID string) tea.Cmd {
	return func() tea.Msg {
		parentGoal, err := repository.GetGoalByID(parentID)
		if err != nil || parentGoal == nil {
			return nil
		}

		if parentGoal.Timeframe == nil {
			return nil
		}

		// Life goals don't have dates, so use current time
		var date time.Time
		if parentGoal.Date != nil {
			date = *parentGoal.Date
		} else if *parentGoal.Timeframe == goal.Life {
			date = time.Now()
		} else {
			// Other timeframes require a date
			return nil
		}

		return screens.OpenTimeframeScreenWithGoal{
			Timeframe: *parentGoal.Timeframe,
			Date:      date,
			GoalID:    parentGoal.ID,
		}
	}
}

func (m *OverdueScreen) getOverdueGoalsCmd() tea.Cmd {
	return func() tea.Msg {
		goals, err := repository.GetOverdueGoals()
		if err != nil {
			return err
		}
		return goallist.GoalsResult{Goals: goals}
	}
}
