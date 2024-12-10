package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

const title = "Hinoki Planner🌲"

type StartupTitleAnimationTickMsg struct{}

type StartupModel struct {
	width, height int
	delay         time.Duration
	letterCounter int
}

func NewStartupScreen(delay time.Duration) Screen {
	return &StartupModel{delay: delay}
}

func (m *StartupModel) Init() tea.Cmd {
	return titleAnimationCmd(calcTickDuration(m.delay, len(title)))
}

func (m *StartupModel) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case StartupTitleAnimationTickMsg:
		if m.letterCounter >= len(title) {
			return nil
		}
		m.letterCounter++
		return titleAnimationCmd(calcTickDuration(m.delay, len(title)))
	}
	return nil
}

func (m *StartupModel) View() string {
	str := ""
	for i, s := range title {
		if i >= m.letterCounter {
			break
		}
		str += string(s)
	}
	return lipgloss.NewStyle().Width(m.width).Height(m.height).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render(str)
}

func (m *StartupModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *StartupModel) Refresh() tea.Cmd {
	return nil
}

func calcTickDuration(d time.Duration, strLen int) time.Duration {
	return time.Duration(float64(d) * 0.5 / float64(strLen))
}

func titleAnimationCmd(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return StartupTitleAnimationTickMsg{}
	})
}
