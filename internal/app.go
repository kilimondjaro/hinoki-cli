package internal

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"hinoki-cli/internal/db"
	"hinoki-cli/internal/goallist"
	"hinoki-cli/internal/screens"
	"hinoki-cli/internal/screens/goaldetails"
	"hinoki-cli/internal/screens/timeframe"
	"log"
	"time"
)

type State int

const (
	StartupView = iota
	TimeframeView
	GoalsDetailsView
)

const startupDelay = time.Second

type model struct {
	state        State
	activeScreen screens.Screen
	startup      StartupModel

	width  int
	height int
}

type AppLaunchStart struct{}
type AppLaunchFinish struct{}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return AppLaunchStart{}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		fmt.Print("\033]2;Hinoki Planner\a")
		m.width = msg.Width
		m.height = msg.Height
		m.startup.SetSize(msg.Width, msg.Height)
		m.activeScreen.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "Q", "Й":
			return m, tea.Quit
		}
	case AppLaunchStart:
		m.state = StartupView
		cmds = append(cmds, m.startup.Init())
		cmds = append(cmds, startupDelayCmd(startupDelay))
	case AppLaunchFinish, goaldetails.OpenTimeframeScreen:
		m.state = TimeframeView
		m.activeScreen = timeframe.NewTimeframeScreen()
		m.activeScreen.SetSize(m.width, m.height)
		cmds = append(cmds, m.activeScreen.Init())
	case goallist.OpenGoalDetails:
		m.state = GoalsDetailsView
		m.activeScreen = goaldetails.NewGoalDetailsScreen(msg.Goal)
		m.activeScreen.SetSize(m.width, m.height)
		cmds = append(cmds, m.activeScreen.Init())
	}

	switch m.state {
	case StartupView:
		cmds = append(cmds, m.startup.Update(msg))
	case TimeframeView, GoalsDetailsView:
		cmds = append(cmds, m.activeScreen.Update(msg))
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state {
	case StartupView:
		return m.startup.View()
	case TimeframeView, GoalsDetailsView:
		return m.activeScreen.View()
	}
	return ""
}

func CreateApp() {
	defer db.CloseDB()

	p := tea.NewProgram(model{activeScreen: timeframe.NewTimeframeScreen(), startup: StartupModel{delay: startupDelay}}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func startupDelayCmd(duration time.Duration) tea.Cmd {
	ch := make(chan int)

	dbCmd := func() tea.Msg {
		db.InitDB()

		ch <- 1

		return nil
	}

	tickCmd := tea.Tick(duration, func(t time.Time) tea.Msg {

		dbInitRes := <-ch

		if dbInitRes > 0 {
			return AppLaunchFinish{}
		}
		return nil
	})

	return tea.Batch(tickCmd, dbCmd)
}
