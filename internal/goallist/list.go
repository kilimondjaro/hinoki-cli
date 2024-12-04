package goallist

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"hinoki-cli/internal/dates"
	"hinoki-cli/internal/goal"
	"time"
)

const (
	Initial = iota
	Normal
	NewGoalInProgress
	GoalEditing
	GoalEditDate
)

type listState int

type GoalList struct {
	list        list.Model
	timeframe   goal.Timeframe
	keys        listKeyMap
	state       listState
	actionInput textinput.Model
	date        time.Time

	width, height int
}

var (
	actionInputLightStyle = lipgloss.NewStyle().MarginBottom(1).Foreground(lipgloss.Color("#666666"))
	actionInputDarkStyle  = lipgloss.NewStyle().MarginBottom(1).Foreground(lipgloss.Color("#cccccc"))
)

const (
	maxWidth = 130
)

type GoalsResult struct {
	goals []goal.Goal
}

type AddGoalSuccess struct{}
type UpdateGoalSuccess struct{}

func NewGoalList(width int, height int, timeframe goal.Timeframe, date time.Time) GoalList {
	keys := NewListKeyMap()

	l := list.New([]list.Item{}, GoalItemDelegate{keys: keys}, width, height)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetStatusBarItemName("goal", "goals")

	actionInput := textinput.New()
	actionInput.Focus()

	return GoalList{list: l, keys: keys, state: Initial, actionInput: actionInput, timeframe: timeframe, date: date}
}

func (m *GoalList) Init() tea.Cmd {
	createGoalsTable()
	return m.getGoalsCmd()
}

func (m *GoalList) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case AddGoalSuccess, UpdateGoalSuccess:
		cmds = append(cmds, m.getGoalsCmd())
	case GoalsResult:
		m.handleGoalResult(msg)
	case tea.KeyMsg:
		cmds = append(cmds, m.handleKeyMsg(msg))
	}

	if m.state == Normal {
		m.list, _ = m.list.Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *GoalList) View() string {
	var actionInput string

	if m.state == NewGoalInProgress || m.state == GoalEditing || m.state == GoalEditDate {
		actionInput = lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.
				NewStyle().
				SetString(m.actionInput.View()).
				Render(),
		)

		actionInputStyle := actionInputLightStyle

		if lipgloss.HasDarkBackground() {
			actionInputStyle = actionInputDarkStyle
		}
		actionInput = actionInputStyle.Render(actionInput)
	}

	actionInputHeight := lipgloss.Height(actionInput)

	listHeight := m.height - actionInputHeight
	m.list.SetSize(maxWidth, listHeight)

	return lipgloss.NewStyle().
		SetString(lipgloss.JoinVertical(lipgloss.Left, m.list.View(), actionInput)).
		Render()
}

func (m *GoalList) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *GoalList) RefreshData() func() tea.Msg {
	return m.getGoalsCmd()
}

func (m *GoalList) SetDate(timeframe goal.Timeframe, date time.Time) {
	m.timeframe = timeframe
	m.date = date
}

func (m *GoalList) IsInActiveState() bool {
	return m.state != Normal
}

func (m *GoalList) getGoalsCmd() func() tea.Msg {
	return func() tea.Msg {
		goals, err := getGoalsByDate(m.timeframe, m.date)

		if err != nil {
			return err
		}

		return GoalsResult{goals: goals}
	}
}

func (m *GoalList) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	switch m.state {
	case Initial:
		cmds = nil
	case NewGoalInProgress, GoalEditing, GoalEditDate:
		cmds = append(cmds, m.handleActionInputKeyMsg(msg))
	case Normal:
		cmds = append(cmds, m.handleKeyMsgInNormalState(msg))
	}
	return tea.Batch(cmds...)
}

func (m *GoalList) handleKeyMsgInNormalState(msg tea.KeyMsg) tea.Cmd {
	item, _ := m.list.SelectedItem().(goal.Goal)

	switch {
	case key.Matches(msg, m.keys.markGoalDone):
		if len(m.list.Items()) == 0 {
			return nil
		}
		item.IsDone = !item.IsDone
		return m.updateGoalCmd(item)
	case key.Matches(msg, m.keys.createGoal):
		m.actionInput.Prompt = "[ ] "
		m.actionInput.Placeholder = "New goal..."
		m.state = NewGoalInProgress
	case key.Matches(msg, m.keys.editGoal):
		if len(m.list.Items()) == 0 {
			return nil
		}

		m.actionInput.Placeholder = ""
		m.actionInput.SetValue(item.Title)
		m.actionInput.Prompt = "Edit: "
		m.state = GoalEditing
	case key.Matches(msg, m.keys.reloadGoals):
		return m.getGoalsCmd()
	case key.Matches(msg, m.keys.archiveGoal):
		item.IsArchived = true
		return m.updateGoalCmd(item)
	case key.Matches(msg, m.keys.changeDate):
		if len(m.list.Items()) == 0 {
			return nil
		}

		m.actionInput.Placeholder = ""
		m.actionInput.Prompt = "Change date: "
		m.state = GoalEditDate
	}

	return nil
}

func (m *GoalList) handleActionInputKeyMsg(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	item, _ := m.list.SelectedItem().(goal.Goal)

	switch msg.Type {
	case tea.KeyEsc:
		m.state = Normal
		m.actionInput.SetValue("")
	case tea.KeyEnter:
		switch m.state {
		case GoalEditing:
			item.Title = m.actionInput.Value()
			m.actionInput.SetValue("")
			cmd = m.updateGoalCmd(item)
		case GoalEditDate:
			date, timeframe, err := dates.ParseDate(time.Now(), m.actionInput.Value())
			if err != nil {
				return nil
			}

			item.Date = date
			item.Timeframe = timeframe

			return m.updateGoalCmd(item)
		case NewGoalInProgress:
			goal := goal.Goal{ID: uuid.New().String(), Title: m.actionInput.Value(), Date: m.date, Timeframe: m.timeframe}
			m.actionInput.SetValue("")

			return m.addGoalCmd(goal)
		}
		m.state = Normal
	default:
		m.actionInput, cmd = m.actionInput.Update(msg)

	}

	return cmd
}

func (m *GoalList) handleGoalResult(msg GoalsResult) {
	m.state = Normal

	var items []list.Item

	for _, goal := range msg.goals {
		items = append(items, goal)
	}
	m.list.SetItems(items)
}

func (m *GoalList) addGoalCmd(goal goal.Goal) func() tea.Msg {
	return func() tea.Msg {
		err := addGoal(goal)
		if err != nil {
			return err
		}

		return AddGoalSuccess{}
	}
}

func (m *GoalList) updateGoalCmd(goal goal.Goal) func() tea.Msg {
	return func() tea.Msg {
		err := updateGoal(goal)
		if err != nil {
			return err
		}

		return UpdateGoalSuccess{}
	}
}
