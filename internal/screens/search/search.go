package search

import (
	"strings"

	"hinoki-cli/internal/goal"
	"hinoki-cli/internal/goallist"
	"hinoki-cli/internal/screens"
	"hinoki-cli/internal/theme"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchScreen struct {
	searchInput textinput.Model
	searchList  list.Model
	keys        keyMap

	width, height int
}

var (
	actionInputLightStyle = lipgloss.NewStyle().MarginBottom(1).Foreground(theme.TextSecondary())
	actionInputDarkStyle  = lipgloss.NewStyle().MarginBottom(1).Foreground(theme.TextSecondary())
)

const (
	maxWidth = 130
)

type searchGoalsResult struct {
	goals []goal.Goal
}

func NewSearchScreen() screens.Screen {
	searchInput := textinput.New()
	searchInput.Prompt = "Search: "
	searchInput.Placeholder = "Type to find goals..."
	searchInput.CharLimit = 256
	searchInput.Focus()

	searchList := list.New([]list.Item{}, newSearchItemDelegate(), 0, 0)
	searchList.SetShowHelp(false)
	searchList.SetShowStatusBar(false)
	searchList.SetShowTitle(false)
	searchList.SetFilteringEnabled(false)
	searchList.SetShowPagination(false)
	searchList.DisableQuitKeybindings()

	searchList.KeyMap.CursorUp = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	)
	searchList.KeyMap.CursorDown = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	)

	return &SearchScreen{
		searchInput: searchInput,
		searchList:  searchList,
		keys:        newKeyMap(),
	}
}

func (m *SearchScreen) Init() tea.Cmd {
	return nil
}

func (m *SearchScreen) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle Enter and Esc first, before list can consume them
		if msg.Type == tea.KeyEnter || msg.Type == tea.KeyEsc {
			cmd := m.handleKeyMsg(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return tea.Batch(cmds...)
		}

		// For navigation keys, update both handler and list
		navKeys := map[string]bool{
			"up":     true,
			"down":   true,
			"pgup":   true,
			"pgdown": true,
		}

		if navKeys[msg.String()] && len(m.searchList.Items()) > 0 {
			// Navigation keys: update list
			var listCmd tea.Cmd
			m.searchList, listCmd = m.searchList.Update(msg)
			if listCmd != nil {
				cmds = append(cmds, listCmd)
			}
			return tea.Batch(cmds...)
		}

		// All other keys go through handleKeyMsg (which sends to input)
		cmd := m.handleKeyMsg(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case searchGoalsResult:
		m.handleSearchGoals(msg)
	case error:
		// swallow errors in UI loop
	}

	return tea.Batch(cmds...)
}

func (m *SearchScreen) View() string {
	inputView := lipgloss.NewStyle().
		SetString(m.searchInput.View()).
		Render()

	inputStyle := actionInputLightStyle
	if lipgloss.HasDarkBackground() {
		inputStyle = actionInputDarkStyle
	}

	inputView = inputStyle.Render(inputView)

	inputHeight := lipgloss.Height(inputView)
	listHeight := m.height - inputHeight - 2 // Account for top padding
	if listHeight < 3 {
		listHeight = 3
	}

	style := lipgloss.NewStyle().PaddingLeft(2).PaddingTop(2)
	horizontalPadding := (m.width - maxWidth) / 2

	if m.width > maxWidth {
		style = style.PaddingLeft(horizontalPadding).PaddingRight(horizontalPadding)
	}

	contentWidth := min(m.width, maxWidth)
	m.searchList.SetSize(contentWidth, listHeight)

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		inputView,
		m.searchList.View(),
	)

	return style.
		SetString(body).
		Render()
}

func (m *SearchScreen) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *SearchScreen) Refresh() tea.Cmd {
	return nil
}

func (m *SearchScreen) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyEsc:
		return func() tea.Msg {
			return screens.GoBack{}
		}
	case tea.KeyEnter:
		// If list has items, open selected goal
		if len(m.searchList.Items()) > 0 {
			return m.openSelectedGoal()
		}
		// Otherwise, let input handle it
		var inputCmd tea.Cmd
		m.searchInput, inputCmd = m.searchInput.Update(msg)
		return inputCmd
	}

	// All keys (including j/k) go to the input
	// Navigation keys are handled in Update method before this
	prevValue := m.searchInput.Value()
	var inputCmd tea.Cmd
	m.searchInput, inputCmd = m.searchInput.Update(msg)

	if prevValue == m.searchInput.Value() {
		return inputCmd
	}

	searchCmd := m.searchGoalsCmd(m.searchInput.Value())
	return tea.Batch(inputCmd, searchCmd)
}

func (m *SearchScreen) openSelectedGoal() tea.Cmd {
	item, ok := m.searchList.SelectedItem().(searchItem)
	if !ok {
		return nil
	}

	goal := item.goal
	if goal.Timeframe == nil || goal.Date == nil {
		return nil
	}

	return func() tea.Msg {
		return screens.OpenTimeframeScreenWithGoal{
			Timeframe: *goal.Timeframe,
			Date:      *goal.Date,
			GoalID:    goal.ID,
		}
	}
}

func (m *SearchScreen) searchGoalsCmd(term string) tea.Cmd {
	trimmed := strings.TrimSpace(term)
	if trimmed == "" {
		m.searchList.SetItems([]list.Item{})
		return nil
	}

	return func() tea.Msg {
		goals, err := goallist.SearchGoals(trimmed, 50)
		if err != nil {
			return err
		}
		return searchGoalsResult{goals: goals}
	}
}

func (m *SearchScreen) handleSearchGoals(msg searchGoalsResult) {
	items := make([]list.Item, 0, len(msg.goals))
	for _, g := range msg.goals {
		items = append(items, searchItem{goal: g})
	}
	m.searchList.SetItems(items)
}
