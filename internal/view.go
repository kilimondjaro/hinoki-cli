package internal

import "github.com/charmbracelet/lipgloss"

func startupView(width int, height int) string {
	return lipgloss.NewStyle().Width(width).Height(height).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render("Hinoki Planner🌲")
}
