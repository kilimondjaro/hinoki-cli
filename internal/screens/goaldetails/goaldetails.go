package goaldetails

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hinoki-cli/internal/goal"
	"hinoki-cli/internal/goallist"
	"hinoki-cli/internal/screens"
)

const (
	Normal = iota
	GotoDate
)

type State int

type GoalDetailsScreen struct {
	list        goallist.GoalList
	keys        listKeyMap
	actionInput textinput.Model
	state       State
	goal        *goal.Goal

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
type OpenTimeframeScreen struct{}

func NewGoalDetailsScreen(goal *goal.Goal) screens.Screen {
	keys := NewListKeyMap()

	actionInput := textinput.New()
	actionInput.Focus()

	goalList := goallist.NewSubgoalsList(goal)

	return &GoalDetailsScreen{keys: keys, actionInput: actionInput, list: goalList, goal: goal}
}

func (m *GoalDetailsScreen) Init() tea.Cmd {
	return m.list.Init()
}

func (m *GoalDetailsScreen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd := m.handleKeyMsg(msg)

		if cmd != nil {
			return cmd
		}
	}

	if m.state == Normal {
		cmds = append(cmds, m.list.Update(msg))
	}

	return tea.Batch(cmds...)
}

func (m *GoalDetailsScreen) View() string {

	header := lipgloss.NewStyle().MarginBottom(2).PaddingTop(2).Render(lipgloss.JoinVertical(lipgloss.Left, m.goal.Title))

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

	m.list.SetSize(maxWidth, listHeight)

	view := lipgloss.JoinVertical(lipgloss.Left, header, m.list.View())

	if m.state != Normal {
		view = lipgloss.JoinVertical(lipgloss.Left, view, actionInput)
	}

	return style.
		SetString(view).
		Render()
}

func (m *GoalDetailsScreen) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *GoalDetailsScreen) refreshData() func() tea.Msg {
	return m.list.RefreshData()
}

func (m *GoalDetailsScreen) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	if m.list.IsInActiveState() {
		return nil
	}

	switch {
	case m.state == Normal:
		cmds = append(cmds, m.handleKeyMsgInNormalState(msg))
	case m.state == GotoDate:
		cmds = append(cmds, m.handleKeyMsgInGotoDateState(msg))
	}

	return tea.Batch(cmds...)
}

func (m *GoalDetailsScreen) handleKeyMsgInNormalState(msg tea.KeyMsg) tea.Cmd {
	switch {
	case msg.String() == "esc":
		return func() tea.Msg {
			return OpenTimeframeScreen{}
		}
	}
	return nil
}

func (m *GoalDetailsScreen) handleKeyMsgInGotoDateState(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEsc:
		m.state = Normal
		m.actionInput.SetValue("")
		return nil
	case tea.KeyEnter:
		m.state = Normal
		//date, timeframe, err := dates.ParseDate(time.Now(), m.actionInput.Value())
		//m.actionInput.SetValue("")
		//if err == nil {
		//	m.timeframe = timeframe
		//	m.date = date
		//}
		//return m.refreshData()
	}

	m.actionInput, cmd = m.actionInput.Update(msg)
	return cmd
}
