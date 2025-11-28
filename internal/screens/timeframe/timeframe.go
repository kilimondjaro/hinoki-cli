package timeframe

import (
	"time"

	"hinoki-cli/internal/dates"
	"hinoki-cli/internal/goal"
	"hinoki-cli/internal/goallist"
	"hinoki-cli/internal/screens"
	"hinoki-cli/internal/theme"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	Normal = iota
	GotoDate
)

type State int

type TimeframeScreen struct {
	list        goallist.GoalList
	keys        listKeyMap
	actionInput textinput.Model
	state       State

	date      time.Time
	timeframe goal.Timeframe

	width, height int
}

var (
	actionInputLightStyle = lipgloss.NewStyle().MarginBottom(1).Foreground(theme.TextSecondary())
	actionInputDarkStyle  = lipgloss.NewStyle().MarginBottom(1).Foreground(theme.TextSecondary())
)

const (
	maxWidth = 130
)

type GoalsResult struct {
	goals []goal.Goal
}

type AddGoalSuccess struct{}
type UpdateGoalSuccess struct{}

func NewTimeframeScreen() screens.Screen {
	keys := NewListKeyMap()

	actionInput := textinput.New()
	actionInput.Focus()

	timeframe := goal.Day
	date := time.Now()

	goalList := goallist.NewGoalList(&timeframe, &date)

	return &TimeframeScreen{
		keys:        keys,
		actionInput: actionInput,
		list:        goalList,
		timeframe:   timeframe,
		date:        date,
	}
}

func (m *TimeframeScreen) SetTimeframeAndDate(timeframe goal.Timeframe, date time.Time) {
	m.timeframe = timeframe
	m.date = date
	m.list.SetDate(timeframe, date)
}

func (m *TimeframeScreen) SetSelectedGoalID(goalID string) {
	m.list.SetGoalIDToSelect(goalID)
}

func (m *TimeframeScreen) Init() tea.Cmd {
	return m.list.Init()
}

func (m *TimeframeScreen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd := m.handleKeyMsg(msg)

		if cmd != nil {
			return cmd
		}
	case error:
		// swallow errors in UI loop, they will be logged by Bubble Tea
	}

	if m.state == Normal {
		cmds = append(cmds, m.list.Update(msg))
	}

	return tea.Batch(cmds...)
}

func (m *TimeframeScreen) View() string {
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
	if m.state == GotoDate {
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

	headerHeight := lipgloss.Height(header)
	actionInputHeight := 0

	if m.state != Normal {
		actionInputHeight = lipgloss.Height(actionInput)
	}

	listHeight := m.height - headerHeight - actionInputHeight

	style := lipgloss.NewStyle().PaddingLeft(2)
	horizontalPadding := (m.width - maxWidth) / 2

	if m.width > maxWidth {
		style = style.PaddingLeft(horizontalPadding).PaddingRight(horizontalPadding)
	}

	contentWidth := min(m.width, maxWidth)

	var body string
	switch m.state {
	case GotoDate:
		m.list.SetSize(contentWidth, listHeight)
		body = m.list.View()
	case Normal:
		m.list.SetSize(contentWidth, listHeight)
		body = m.list.View()
	}

	view := lipgloss.JoinVertical(lipgloss.Left, header, body)

	if m.state == GotoDate {
		view = lipgloss.JoinVertical(lipgloss.Left, view, actionInput)
	}

	return style.
		SetString(view).
		Render()
}

func (m *TimeframeScreen) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *TimeframeScreen) Refresh() tea.Cmd {
	m.list.SetDate(m.timeframe, m.date)
	return m.list.RefreshData()
}

func (m *TimeframeScreen) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	if m.list.IsInActiveState() {
		return nil
	}

	switch m.state {
	case Normal:
		cmds = append(cmds, m.handleKeyMsgInNormalState(msg))
	case GotoDate:
		cmds = append(cmds, m.handleKeyMsgInGotoDateState(msg))
	}

	return tea.Batch(cmds...)
}

func (m *TimeframeScreen) handleKeyMsgInNormalState(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keys.dayTimeslice):
		m.timeframe = goal.Day
		m.date = time.Now()
		return m.Refresh()
	case key.Matches(msg, m.keys.weekTimeslice):
		m.timeframe = goal.Week
		m.date = time.Now()
		return m.Refresh()
	case key.Matches(msg, m.keys.monthTimeslice):
		m.timeframe = goal.Month
		m.date = time.Now()
		return m.Refresh()
	case key.Matches(msg, m.keys.quarterTimeslice):
		m.timeframe = goal.Quarter
		m.date = time.Now()
		return m.Refresh()
	case key.Matches(msg, m.keys.yearTimeslice):
		m.date = time.Now()
		m.timeframe = goal.Year
		return m.Refresh()
	case key.Matches(msg, m.keys.currentPeriod):
		m.date = time.Now()
		m.timeframe = goal.Day
		return m.Refresh()
	case key.Matches(msg, m.keys.lifeTimeslice):
		m.date = time.Now()
		m.timeframe = goal.Life
		return m.Refresh()
	case key.Matches(msg, m.keys.previousPeriod):
		m.date = dates.ChangePeriod(m.date, m.timeframe, -1)
		return m.Refresh()
	case key.Matches(msg, m.keys.nextPeriod):
		m.date = dates.ChangePeriod(m.date, m.timeframe, 1)
		return m.Refresh()
	case key.Matches(msg, m.keys.gotoPeriod):
		m.actionInput.Placeholder = ""
		m.actionInput.Prompt = "Jump to date: "
		m.state = GotoDate
	case key.Matches(msg, m.keys.searchGoals):
		return func() tea.Msg {
			return screens.OpenSearchScreen{}
		}
	case key.Matches(msg, m.keys.goToParent):
		selectedGoal := m.list.GetSelectedGoal()
		if selectedGoal == nil || selectedGoal.ParentId == nil {
			return nil
		}
		return m.goToParentGoalCmd(*selectedGoal.ParentId)
	}
	return nil
}

func (m *TimeframeScreen) handleKeyMsgInGotoDateState(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEsc:
		m.state = Normal
		m.actionInput.SetValue("")
		return nil
	case tea.KeyEnter:
		m.state = Normal
		date, timeframe, err := dates.ParseDate(time.Now(), m.actionInput.Value())
		m.actionInput.SetValue("")
		if err == nil {
			m.timeframe = timeframe
			m.date = date
		}
		return m.Refresh()
	}

	m.actionInput, cmd = m.actionInput.Update(msg)
	return cmd
}

func (m *TimeframeScreen) goToParentGoalCmd(parentID string) tea.Cmd {
	return func() tea.Msg {
		parentGoal, err := goallist.GetGoalByID(parentID)
		if err != nil || parentGoal == nil {
			return nil
		}

		if parentGoal.Timeframe == nil || parentGoal.Date == nil {
			return nil
		}

		return screens.OpenTimeframeScreenWithGoal{
			Timeframe: *parentGoal.Timeframe,
			Date:      *parentGoal.Date,
			GoalID:    parentGoal.ID,
		}
	}
}
