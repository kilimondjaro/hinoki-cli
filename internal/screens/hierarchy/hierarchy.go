package hierarchy

import (
	"fmt"
	"hinoki-cli/internal/dates"
	"hinoki-cli/internal/goal"
	"hinoki-cli/internal/repository"
	"hinoki-cli/internal/screens"
	"hinoki-cli/internal/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HierarchyScreen struct {
	goal      *goal.Goal
	ancestors []goal.Goal
	keys      keyMap

	width, height int
}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.TextPrimary()).
			MarginBottom(2).
			PaddingTop(2)

	rootStyle = lipgloss.NewStyle().
			Foreground(theme.TextPrimary()).
			Bold(true)

	ancestorStyle = lipgloss.NewStyle().
			Foreground(theme.TextSecondary())

	currentStyle = lipgloss.NewStyle().
			Foreground(theme.TextSelected()).
			Bold(true)

	metaStyle = lipgloss.NewStyle().
			Foreground(theme.TextMuted())

	treeCharStyle = lipgloss.NewStyle().
			Foreground(theme.TextMuted())
)

const (
	maxWidth = 130
)

type AncestorChainResult struct {
	ancestors []goal.Goal
}

func NewHierarchyScreen(goal *goal.Goal) screens.Screen {
	return &HierarchyScreen{
		goal: goal,
		keys: newKeyMap(),
	}
}

func (m *HierarchyScreen) Init() tea.Cmd {
	return m.getAncestorChainCmd()
}

func (m *HierarchyScreen) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case AncestorChainResult:
		m.ancestors = msg.ancestors
	case error:
		// swallow errors in UI loop
	}

	return nil
}

func (m *HierarchyScreen) View() string {
	header := headerStyle.Render("Goal Hierarchy")

	headerHeight := lipgloss.Height(header)
	contentHeight := m.height - headerHeight - 2 // Account for padding

	style := lipgloss.NewStyle().PaddingLeft(2).PaddingTop(1)
	horizontalPadding := (m.width - maxWidth) / 2

	if m.width > maxWidth {
		style = style.PaddingLeft(horizontalPadding).PaddingRight(horizontalPadding)
	}

	contentWidth := min(m.width, maxWidth)
	treeView := m.renderTree(contentWidth, contentHeight)

	view := lipgloss.JoinVertical(lipgloss.Left, header, treeView)

	return style.
		SetString(view).
		Render()
}

func (m *HierarchyScreen) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *HierarchyScreen) Refresh() tea.Cmd {
	return m.getAncestorChainCmd()
}

func (m *HierarchyScreen) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyEsc:
		return func() tea.Msg {
			return screens.GoBack{}
		}
	}

	return nil
}

func (m *HierarchyScreen) getAncestorChainCmd() tea.Cmd {
	return func() tea.Msg {
		ancestors, err := repository.GetAncestorChain(m.goal.ID)
		if err != nil {
			return err
		}
		return AncestorChainResult{ancestors: ancestors}
	}
}

func (m *HierarchyScreen) renderTree(width, height int) string {
	if len(m.ancestors) == 0 {
		return metaStyle.Render("Loading...")
	}

	var lines []string

	// Render the tree from root to current goal
	for i, g := range m.ancestors {
		isLast := i == len(m.ancestors)-1
		isCurrent := i == len(m.ancestors)-1
		isRoot := i == 0

		// Build indentation and tree characters
		var prefix string
		if isRoot {
			// Root - show with special indicator
			prefix = treeCharStyle.Render("┌─ ")
		} else if isLast {
			// Current goal - use └─
			prefix = treeCharStyle.Render("└─ ")
		} else {
			// Intermediate ancestor - use ├─
			prefix = treeCharStyle.Render("├─ ")
		}

		// Determine style
		var goalStyle lipgloss.Style
		if isRoot {
			goalStyle = rootStyle
		} else if isCurrent {
			goalStyle = currentStyle
		} else {
			goalStyle = ancestorStyle
		}

		// Add status indicator with appropriate styling
		var status string
		if g.IsDone {
			status = metaStyle.Render("✓")
		} else {
			status = metaStyle.Render("○")
		}

		// Build the goal line
		goalLine := goalStyle.Render(g.Title)

		// Add metadata (timeframe and date)
		meta := m.formatMeta(&g)
		if meta != "" {
			goalLine = fmt.Sprintf("%s %s", goalLine, metaStyle.Render(meta))
		}

		// Combine prefix, status, and goal line
		fullLine := fmt.Sprintf("%s%s %s", prefix, status, goalLine)

		if !isLast && !isRoot {
			lines = append(lines, treeCharStyle.Render("│ "))
		}

		lines = append(lines, fullLine)

		// Add vertical connector for non-last items (except root)
		// Only add if there are more items after this one
		// Use a single space after the connector to maintain consistent spacing
		if !isLast && !isRoot {
			lines = append(lines, treeCharStyle.Render("│ "))
		}
	}

	// Join all lines with no extra spacing
	treeContent := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Wrap in a container with proper width and no extra padding
	return lipgloss.NewStyle().
		Width(width).
		MaxHeight(height).
		Render(treeContent)
}

func (m *HierarchyScreen) formatMeta(g *goal.Goal) string {
	if g.Timeframe == nil {
		return ""
	}

	timeframe := g.Timeframe.String()
	date := ""
	if g.Date != nil {
		date = dates.DateString(*g.Date, *g.Timeframe)
	}

	if date != "" {
		return fmt.Sprintf("• %s • %s", timeframe, date)
	}
	return fmt.Sprintf("• %s", timeframe)
}
