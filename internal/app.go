package internal

import (
	"fmt"
	"hinoki-cli/internal/db"
	"hinoki-cli/internal/screens"
	"hinoki-cli/internal/screens/goaldetails"
	"hinoki-cli/internal/screens/search"
	"hinoki-cli/internal/screens/timeframe"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const startupDelay = time.Second

type model struct {
	navigation screens.Navigation

	width  int
	height int
}

type AppLaunchStart struct{}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return AppLaunchStart{}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	currentScreen := m.navigation.Top()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		fmt.Print("\033]2;Hinoki Planner\a")
		m.width = msg.Width
		m.height = msg.Height
		if currentScreen != nil {
			currentScreen.SetSize(m.width, m.height)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	cmds = append(cmds, m.handleNavigation(msg))

	if currentScreen != nil {
		cmds = append(cmds, m.navigation.Top().Update(msg))
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	currentScreen := m.navigation.Top()
	if currentScreen == nil {
		return ""
	}

	return currentScreen.View()
}

func (m model) handleNavigation(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case AppLaunchStart:
		startupScreen := screens.NewStartupScreen(startupDelay)
		startupScreen.SetSize(m.width, m.height)
		cmds = append(cmds,
			startupScreen.Init(),
			startupDelayCmd(startupDelay),
		)
		m.navigation.Push(startupScreen)
	case screens.OpenTimeframeScreen:
		timeframeScreen := timeframe.NewTimeframeScreen()
		timeframeScreen.SetSize(m.width, m.height)
		cmds = append(cmds, timeframeScreen.Init())
		m.navigation.Replace(timeframeScreen)
	case screens.OpenTimeframeScreenWithGoal:
		timeframeScreen := timeframe.NewTimeframeScreen()
		if ts, ok := timeframeScreen.(*timeframe.TimeframeScreen); ok {
			ts.SetTimeframeAndDate(msg.Timeframe, msg.Date)
			if msg.GoalID != "" {
				ts.SetSelectedGoalID(msg.GoalID)
			}
		}
		timeframeScreen.SetSize(m.width, m.height)
		cmds = append(cmds, timeframeScreen.Init(), timeframeScreen.Refresh())
		m.navigation.Replace(timeframeScreen)
	case screens.OpenSearchScreen:
		searchScreen := search.NewSearchScreen()
		searchScreen.SetSize(m.width, m.height)
		cmds = append(cmds, searchScreen.Init())
		m.navigation.Push(searchScreen)
	case screens.OpenGoalDetailsScreen:
		goalDetailsScreen := goaldetails.NewGoalDetailsScreen(msg.Goal)
		goalDetailsScreen.SetSize(m.width, m.height)
		cmds = append(cmds, goalDetailsScreen.Init())
		m.navigation.Push(goalDetailsScreen)
	case screens.GoBack:
		m.navigation.Pop()
		cmds = append(cmds, m.navigation.Top().Refresh())
	}

	return tea.Batch(cmds...)
}

func CreateApp() {
	defer db.CloseDB()

	p := tea.NewProgram(model{navigation: &screens.NavigationState{}}, tea.WithAltScreen())

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
			return screens.OpenTimeframeScreen{}
		}
		return nil
	})

	return tea.Batch(tickCmd, dbCmd)
}
