package screens

import (
	"hinoki-cli/internal/goal"
)

type Navigation interface {
	Push(screen Screen)
	Pop() Screen
	Replace(screen Screen)
	Top() Screen
}

type NavigationState struct {
	stack []Screen
}

type GoBack struct{}
type OpenTimeframeScreen struct{}
type OpenGoalDetailsScreen struct {
	Goal *goal.Goal
}

func (m *NavigationState) Push(screen Screen) {
	m.stack = append(m.stack, screen)
}

func (m *NavigationState) Pop() Screen {
	lastIndex := len(m.stack) - 1
	screen := m.stack[lastIndex]
	m.stack = m.stack[:lastIndex]

	return screen
}

func (m *NavigationState) Replace(screen Screen) {
	m.stack = append(m.stack[:len(m.stack)-1], screen)
}

func (m *NavigationState) Top() Screen {
	if len(m.stack) == 0 {
		return nil
	}

	return m.stack[len(m.stack)-1]
}
