package search

import (
	"fmt"
	"hinoki-cli/internal/dates"
	"hinoki-cli/internal/goal"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchItem struct {
	goal goal.Goal
}

func (i searchItem) FilterValue() string {
	return i.goal.Title
}

type searchItemDelegate struct{}

var (
	searchTitleStyle     = lipgloss.NewStyle()
	searchMetaStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	searchSelectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	searchUnselectedDark = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	searchUnselectedLite = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000"))
)

func newSearchItemDelegate() list.ItemDelegate {
	return searchItemDelegate{}
}

func (d searchItemDelegate) Height() int { return 2 }

func (d searchItemDelegate) Spacing() int { return 1 }

func (d searchItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d searchItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(searchItem)
	if !ok {
		return
	}

	render := searchUnselectedLite.Render
	if lipgloss.HasDarkBackground() {
		render = searchUnselectedDark.Render
	}

	meta := d.metaLine(item.goal)

	line := item.goal.Title
	if meta != "" {
		line = fmt.Sprintf("%s\n%s", line, searchMetaStyle.Render(meta))
	}

	if index == m.Index() {
		render = searchSelectedStyle.Render
	}

	fmt.Fprint(w, render(lipgloss.NewStyle().Width(m.Width()).Render(line)))
}

func (d searchItemDelegate) metaLine(goal goal.Goal) string {
	if goal.Timeframe == nil {
		return ""
	}

	timeframe := goal.Timeframe.String()
	date := ""
	if goal.Date != nil {
		date = dates.DateString(*goal.Date, *goal.Timeframe)
	}

	meta := timeframe
	if date != "" {
		meta = fmt.Sprintf("%s • %s", meta, date)
	}

	if goal.ParentTitle != nil {
		meta = fmt.Sprintf("%s • Parent: %s", meta, *goal.ParentTitle)
	}

	return meta
}
