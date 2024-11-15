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
	Loading
	Normal
	NewGoalInProgress
	GoalEditing
	GotoDate
	GoalEditDate
)

type listState int

type GoalList struct {
	list      list.Model
	timeframe goal.Timeframe
	keys      listKeyMap
	state     listState
	goalInput textinput.Model
	date      time.Time

	width, height int
}

var (
	actionInputLightStyle = lipgloss.NewStyle().MarginBottom(1).Foreground(lipgloss.Color("#666666"))
	actionInputDarkStyle  = lipgloss.NewStyle().MarginBottom(1).Foreground(lipgloss.Color("#cccccc"))
)

type GoalsResult struct {
	goals []goal.Goal
}

type AddGoalSuccess struct{}
type UpdateGoalSuccess struct{}

func NewGoalList(width int, height int) GoalList {
	keys := NewListKeyMap()

	l := list.New([]list.Item{}, GoalItemDelegate{keys: keys}, width, height)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetStatusBarItemName("goal", "goals")

	goalInput := textinput.New()
	goalInput.Focus()

	return GoalList{list: l, keys: keys, state: Initial, goalInput: goalInput, timeframe: goal.Day, date: time.Now()}
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
	slice := lipgloss.NewStyle().
		SetString(m.timeframe.String()).
		Underline(true).
		MarginBottom(1).
		Render()

	date := lipgloss.
		NewStyle().
		SetString(dates.DateString(m.date, m.timeframe)).
		Height(1).
		Bold(true).
		Render()

	header := lipgloss.NewStyle().MarginBottom(2).PaddingTop(2).Render(lipgloss.JoinVertical(lipgloss.Left, slice, date))

	var actionInput string
	if m.state == NewGoalInProgress || m.state == GoalEditing || m.state == GotoDate || m.state == GoalEditDate {
		actionInput = lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.
				NewStyle().
				SetString(m.goalInput.View()).
				Render(),
		)

		actionInputStyle := actionInputLightStyle

		if lipgloss.HasDarkBackground() {
			actionInputStyle = actionInputDarkStyle
		}
		actionInput = actionInputStyle.Render(actionInput)
	}

	headerHeight := lipgloss.Height(header)
	actionInputHeight := lipgloss.Height(actionInput)
	listHeight := m.height - headerHeight - actionInputHeight
	m.list.SetSize(m.width, listHeight)

	return lipgloss.NewStyle().
		SetString(lipgloss.JoinVertical(lipgloss.Left, header, m.list.View(), actionInput)).
		PaddingLeft(2).
		Render()
}

func (m *GoalList) SetSize(width, height int) {
	m.width = width
	m.height = height
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

	switch {
	case m.state == NewGoalInProgress:
		cmds = append(cmds, m.handleKeyMsgInNewGoalInProgressState(msg))
	case m.state == Normal:
		cmds = append(cmds, m.handleKeyMsgInNormalState(msg))
	case m.state == GoalEditing:
		cmds = append(cmds, m.handleKeyMsgInGoalEditingState(msg))
	case m.state == GotoDate:
		cmds = append(cmds, m.handleKeyMsgInGotoDateState(msg))
	case m.state == GoalEditDate:
		cmds = append(cmds, m.handleKeyMsgInGoalEditDateState(msg))
	}
	return tea.Batch(cmds...)
}

func (m *GoalList) handleKeyMsgInNormalState(msg tea.KeyMsg) tea.Cmd {
	item, _ := m.list.SelectedItem().(goal.Goal)

	switch {
	case key.Matches(msg, m.keys.dayTimeslice):
		m.timeframe = goal.Day
		m.date = time.Now()
		return m.getGoalsCmd()
	case key.Matches(msg, m.keys.weekTimeslice):
		m.timeframe = goal.Week
		m.date = time.Now()
		return m.getGoalsCmd()
	case key.Matches(msg, m.keys.monthTimeslice):
		m.timeframe = goal.Month
		m.date = time.Now()
		return m.getGoalsCmd()
	case key.Matches(msg, m.keys.quarterTimeslice):
		m.timeframe = goal.Quarter
		m.date = time.Now()
		return m.getGoalsCmd()
	case key.Matches(msg, m.keys.yearTimeslice):
		m.date = time.Now()
		m.timeframe = goal.Year
		return m.getGoalsCmd()
	case key.Matches(msg, m.keys.currentPeriod):
		m.date = time.Now()
		m.timeframe = goal.Day
		return m.getGoalsCmd()
	case key.Matches(msg, m.keys.lifeTimeslice):
		m.date = time.Now()
		m.timeframe = goal.Life
		return m.getGoalsCmd()
	case key.Matches(msg, m.keys.markGoalDone):
		if len(m.list.Items()) == 0 {
			return nil
		}
		item.IsDone = !item.IsDone
		return m.updateGoalCmd(item)
	case key.Matches(msg, m.keys.createGoal):
		m.goalInput.Prompt = "[ ] "
		m.goalInput.Placeholder = "New goal..."
		m.state = NewGoalInProgress
	case key.Matches(msg, m.keys.editGoal):
		if len(m.list.Items()) == 0 {
			return nil
		}

		m.goalInput.Placeholder = ""
		m.goalInput.SetValue(item.Title)
		m.goalInput.Prompt = "Edit: "
		m.state = GoalEditing
	case key.Matches(msg, m.keys.previousPeriod):
		m.date = dates.ChangePeriod(m.date, m.timeframe, -1)
		return m.getGoalsCmd()
	case key.Matches(msg, m.keys.nextPeriod):
		m.date = dates.ChangePeriod(m.date, m.timeframe, 1)
		return m.getGoalsCmd()
	case key.Matches(msg, m.keys.reloadGoals):
		return m.getGoalsCmd()
	case key.Matches(msg, m.keys.archiveGoal):
		item.IsArchived = true
		return m.updateGoalCmd(item)
	case key.Matches(msg, m.keys.gotoPeriod):
		m.goalInput.Placeholder = ""
		m.goalInput.Prompt = "Jump to date: "
		m.state = GotoDate
	case key.Matches(msg, m.keys.changeDate):
		if len(m.list.Items()) == 0 {
			return nil
		}

		m.goalInput.Placeholder = ""
		m.goalInput.Prompt = "Change date: "
		m.state = GoalEditDate
	}

	return nil
}

func (m *GoalList) handleKeyMsgInNewGoalInProgressState(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEsc:
		m.state = Normal
		m.goalInput.SetValue("")
		return nil
	case tea.KeyEnter:
		m.state = Normal

		goal := goal.Goal{ID: uuid.New().String(), Title: m.goalInput.Value(), Date: m.date, Timeframe: m.timeframe}
		m.goalInput.SetValue("")

		return m.addGoalCmd(goal)
	}

	m.goalInput, cmd = m.goalInput.Update(msg)
	return cmd
}

func (m *GoalList) handleKeyMsgInGoalEditingState(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	item, _ := m.list.SelectedItem().(goal.Goal)

	switch msg.Type {
	case tea.KeyEsc:
		m.state = Normal
		m.goalInput.SetValue("")
		return nil
	case tea.KeyEnter:
		m.state = Normal
		item.Title = m.goalInput.Value()
		m.goalInput.SetValue("")
		return m.updateGoalCmd(item)
	}

	m.goalInput, cmd = m.goalInput.Update(msg)
	return cmd
}

func (m *GoalList) handleKeyMsgInGotoDateState(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEsc:
		m.state = Normal
		m.goalInput.SetValue("")
		return nil
	case tea.KeyEnter:
		m.state = Normal
		date, timeframe, err := dates.ParseDate(time.Now(), m.goalInput.Value())
		m.goalInput.SetValue("")
		if err == nil {
			m.timeframe = timeframe
			m.date = date
		}
		return m.getGoalsCmd()
	}

	m.goalInput, cmd = m.goalInput.Update(msg)
	return cmd
}

func (m *GoalList) handleKeyMsgInGoalEditDateState(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	item, _ := m.list.SelectedItem().(goal.Goal)

	switch msg.Type {
	case tea.KeyEsc:
		m.state = Normal
		m.goalInput.SetValue("")
		return nil
	case tea.KeyEnter:
		m.state = Normal
		date, timeframe, err := dates.ParseDate(time.Now(), m.goalInput.Value())
		m.goalInput.SetValue("")

		if err != nil {
			return nil
		}

		item.Date = date
		item.Timeframe = timeframe

		return m.updateGoalCmd(item)
	}

	m.goalInput, cmd = m.goalInput.Update(msg)
	return cmd
}

func (m *GoalList) handleTimeframeChange() {

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
