package theme

import "github.com/charmbracelet/lipgloss"

// Color palette - unified design system colors
const (
	// Primary colors
	ColorBlack = "#000000"
	ColorWhite = "#ffffff"

	// Grayscale
	ColorGrayLight  = "#cccccc" // Light gray for secondary text
	ColorGrayMedium = "#888888" // Medium gray for metadata
	ColorGrayDark   = "#666666" // Dark gray for disabled/done items

	// Accent colors
	ColorAccent = "170" // Pink/magenta for selected items
)

// Semantic color functions - adapt to light/dark background
func TextPrimary() lipgloss.Color {
	if lipgloss.HasDarkBackground() {
		return lipgloss.Color(ColorWhite)
	}
	return lipgloss.Color(ColorBlack)
}

func TextSecondary() lipgloss.Color {
	if lipgloss.HasDarkBackground() {
		return lipgloss.Color(ColorGrayLight)
	}
	return lipgloss.Color(ColorGrayDark)
}

func TextMuted() lipgloss.Color {
	return lipgloss.Color(ColorGrayMedium)
}

func TextDisabled() lipgloss.Color {
	return lipgloss.Color(ColorGrayDark)
}

func TextSelected() lipgloss.Color {
	return lipgloss.Color(ColorAccent)
}

// Direct color access (for cases where semantic doesn't fit)
func Color(color string) lipgloss.Color {
	return lipgloss.Color(color)
}
