package hierarchy

import (
	"fmt"
	"hinoki-cli/internal/dates"
	"hinoki-cli/internal/goal"
	"hinoki-cli/internal/repository"
	"hinoki-cli/internal/screens"
	"hinoki-cli/internal/theme"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HierarchyScreen struct {
	goal      *goal.Goal
	ancestors []goal.Goal
	keys      keyMap
	showAll   bool // Toggle for showing full tree vs just ancestors

	// Navigation state
	cursor         int        // Index of currently selected item in flattened list
	scrollOffset   int        // Number of lines scrolled up
	flattenedItems []TreeItem // Flattened list of all tree items for navigation

	width, height int
}

// TreeItem represents a single item in the flattened tree for navigation
type TreeItem struct {
	goal      goal.Goal
	prefix    string // Tree characters (├─, └─, etc.)
	status    string // Status indicator (✓, ○)
	style     lipgloss.Style
	isCurrent bool
	depth     int
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

type FullTreeResult struct {
	treeNodes []TreeNode
}

type TreeNode struct {
	goal      goal.Goal
	children  []TreeNode
	depth     int
	isCurrent bool
	isLast    bool
}

func NewHierarchyScreen(goal *goal.Goal) screens.Screen {
	return &HierarchyScreen{
		goal: goal,
		keys: newKeyMap(),
	}
}

func (m *HierarchyScreen) Init() tea.Cmd {
	m.cursor = 0
	m.scrollOffset = 0
	return m.getAncestorChainCmd()
}

func (m *HierarchyScreen) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case AncestorChainResult:
		m.ancestors = msg.ancestors
		m.updateFlattenedItems()
	case FullTreeResult:
		// Convert tree nodes back to ancestors for rendering
		m.ancestors = m.treeNodesToGoals(msg.treeNodes)
		m.updateFlattenedItems()
	case error:
		// swallow errors in UI loop
	}

	return nil
}

func (m *HierarchyScreen) View() string {
	headerText := "Goal Hierarchy"
	if m.showAll {
		headerText += " (Full Tree)"
	}
	header := headerStyle.Render(headerText)

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
	// Adjust scroll offset when window size changes to ensure cursor is still visible
	m.adjustScrollOffset()
}

func (m *HierarchyScreen) Refresh() tea.Cmd {
	return m.getAncestorChainCmd()
}

func (m *HierarchyScreen) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch {
	case msg.Type == tea.KeyEsc:
		return func() tea.Msg {
			return screens.GoBack{}
		}
	case key.Matches(msg, m.keys.showAllTree):
		m.showAll = !m.showAll
		m.cursor = 0
		m.scrollOffset = 0
		if m.showAll {
			return m.getFullTreeCmd()
		}
		return m.getAncestorChainCmd()
	case key.Matches(msg, m.keys.cursorUp):
		m.moveCursor(-1)
	case key.Matches(msg, m.keys.cursorDown):
		m.moveCursor(1)
	case key.Matches(msg, m.keys.openDetails):
		return m.openGoalDetailsCmd()
	case key.Matches(msg, m.keys.openTimeframe):
		return m.openTimeframeCmd()
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

	if len(m.flattenedItems) == 0 {
		return metaStyle.Render("Loading...")
	}

	return m.renderFlattenedTree(width, height)
}

// renderFlattenedTree renders the visible portion of the flattened tree
// height parameter is the available content height (already accounting for header)
func (m *HierarchyScreen) renderFlattenedTree(width, height int) string {
	if len(m.flattenedItems) == 0 {
		return metaStyle.Render("Loading...")
	}

	// height is already the content height, so we can use it directly
	// Account for any padding/margins (subtract a small buffer)
	maxVisibleItems := max(1, height-1)

	// Calculate visible range
	startIdx := m.scrollOffset
	endIdx := min(startIdx+maxVisibleItems, len(m.flattenedItems))

	// Ensure the cursor is always included in the visible range
	if m.cursor >= startIdx && m.cursor < len(m.flattenedItems) {
		// If cursor is beyond endIdx, extend endIdx to include it
		if m.cursor >= endIdx {
			endIdx = m.cursor + 1
		}
		// Ensure we don't go beyond the list
		endIdx = min(endIdx, len(m.flattenedItems))
	}

	var lines []string
	for i := startIdx; i < endIdx; i++ {
		item := m.flattenedItems[i]
		isSelected := i == m.cursor

		// Skip empty connector lines when rendering
		if item.goal.ID == "" && item.status == "" {
			// This is a connector line
			lines = append(lines, treeCharStyle.Render(item.prefix))
			continue
		}

		// Determine if this item should be highlighted (selected)
		var itemStyle lipgloss.Style = item.style
		if isSelected {
			// Highlight selected item - invert colors for visibility
			if lipgloss.HasDarkBackground() {
				itemStyle = item.style.Copy().Background(theme.TextSelected()).Foreground(theme.TextPrimary())
			} else {
				itemStyle = item.style.Copy().Foreground(theme.TextSelected()).Bold(true)
			}
		}

		// Render prefix
		prefix := treeCharStyle.Render(item.prefix)

		// Render status
		var status string
		if item.status != "" {
			status = metaStyle.Render(item.status) + " "
		}

		// Render goal line
		goalLine := itemStyle.Render(item.goal.Title)

		// Add metadata if this is a real goal (not a connector)
		if item.goal.ID != "" {
			meta := m.formatMeta(&item.goal)
			if meta != "" {
				goalLine = fmt.Sprintf("%s %s", goalLine, metaStyle.Render(meta))
			}
		}

		// Combine everything
		fullLine := fmt.Sprintf("%s%s%s", prefix, status, goalLine)
		lines = append(lines, fullLine)
	}

	// Join all lines
	treeContent := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Wrap in a container with proper width
	// Don't use MaxHeight here as it might cut off content - height is already constrained
	return lipgloss.NewStyle().
		Width(width).
		Render(treeContent)
}

func (m *HierarchyScreen) buildFullTree() []TreeNode {
	// Build tree structure from ancestors
	// Ancestors form a chain, not siblings
	if len(m.ancestors) == 0 {
		return []TreeNode{}
	}

	// Build the chain recursively, starting from root
	return m.buildTreeNodeChain(0)
}

func (m *HierarchyScreen) buildTreeNodeChain(index int) []TreeNode {
	if index >= len(m.ancestors) {
		return []TreeNode{}
	}

	g := m.ancestors[index]
	isCurrent := index == len(m.ancestors)-1

	// Get all children for this goal
	allChildren, _ := repository.GetGoalsByParent(g.ID)

	// Separate children into: ancestors in the chain vs other siblings
	var chainChildren []TreeNode
	var siblingChildren []TreeNode

	for _, child := range allChildren {
		// Check if this child is the next ancestor in the chain
		isNextAncestor := index+1 < len(m.ancestors) && child.ID == m.ancestors[index+1].ID

		if isNextAncestor {
			// This is the next ancestor in the chain - build it recursively
			nextChain := m.buildTreeNodeChain(index + 1)
			chainChildren = append(chainChildren, nextChain...)
		} else {
			// This is a sibling - get its children recursively
			siblingChildNodes := m.getChildrenRecursive(child.ID, index+1)
			siblingChildren = append(siblingChildren, TreeNode{
				goal:      child,
				children:  siblingChildNodes,
				depth:     index + 1,
				isCurrent: child.ID == m.goal.ID,
				isLast:    false, // Will be set when we combine
			})
		}
	}

	// Combine chain children and sibling children
	allChildNodes := append(chainChildren, siblingChildren...)

	// Mark the last child as last
	if len(allChildNodes) > 0 {
		allChildNodes[len(allChildNodes)-1].isLast = true
	}

	// Return this node with its children
	return []TreeNode{{
		goal:      g,
		children:  allChildNodes,
		depth:     index,
		isCurrent: isCurrent,
		isLast:    true, // Only one node at root level
	}}
}

func (m *HierarchyScreen) getChildrenRecursive(goalID string, depth int) []TreeNode {
	children, _ := repository.GetGoalsByParent(goalID)
	if len(children) == 0 {
		return []TreeNode{}
	}

	var childNodes []TreeNode
	for i, child := range children {
		childChildren := m.getChildrenRecursive(child.ID, depth+1)
		childNodes = append(childNodes, TreeNode{
			goal:      child,
			children:  childChildren,
			depth:     depth,
			isCurrent: child.ID == m.goal.ID,
			isLast:    i == len(children)-1,
		})
	}

	return childNodes
}

func (m *HierarchyScreen) treeNodesToGoals(nodes []TreeNode) []goal.Goal {
	// Flatten tree nodes back to goals for simple ancestor view
	var goals []goal.Goal
	for _, node := range nodes {
		goals = append(goals, node.goal)
	}
	return goals
}

func (m *HierarchyScreen) getFullTreeCmd() tea.Cmd {
	return func() tea.Msg {
		ancestors, err := repository.GetAncestorChain(m.goal.ID)
		if err != nil {
			return err
		}

		// Build full tree structure
		var treeNodes []TreeNode
		for i, g := range ancestors {
			isCurrent := i == len(ancestors)-1

			// Get children for this goal
			children, _ := repository.GetGoalsByParent(g.ID)

			// Convert children to tree nodes
			childNodes := make([]TreeNode, len(children))
			for j, child := range children {
				childNodes[j] = TreeNode{
					goal:      child,
					children:  []TreeNode{},
					depth:     i + 1,
					isCurrent: child.ID == m.goal.ID,
					isLast:    j == len(children)-1,
				}
			}

			treeNodes = append(treeNodes, TreeNode{
				goal:      g,
				children:  childNodes,
				depth:     i,
				isCurrent: isCurrent,
				isLast:    i == len(ancestors)-1,
			})
		}

		return FullTreeResult{treeNodes: treeNodes}
	}
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

// updateFlattenedItems creates a flattened list of all tree items for navigation
func (m *HierarchyScreen) updateFlattenedItems() {
	m.flattenedItems = []TreeItem{}

	if len(m.ancestors) == 0 {
		return
	}

	if m.showAll {
		treeNodes := m.buildFullTree()
		m.flattenTreeNodes(treeNodes, "", false, true)
	} else {
		m.flattenAncestorChain()
	}

	// Set cursor to current goal if it exists, otherwise keep current position
	if m.goal != nil {
		for i, item := range m.flattenedItems {
			if item.goal.ID == m.goal.ID {
				m.cursor = i
				break
			}
		}
	}

	// Ensure cursor is within bounds and points to a valid item (not a connector)
	if m.cursor >= len(m.flattenedItems) {
		m.cursor = max(0, len(m.flattenedItems)-1)
	}

	// If cursor is on a connector line, find the nearest valid item
	if m.cursor < len(m.flattenedItems) && m.flattenedItems[m.cursor].goal.ID == "" {
		// Try to find a valid item nearby
		found := false
		// First try going forward
		for i := m.cursor; i < len(m.flattenedItems); i++ {
			if m.flattenedItems[i].goal.ID != "" {
				m.cursor = i
				found = true
				break
			}
		}
		// If not found, go backward
		if !found {
			for i := m.cursor; i >= 0; i-- {
				if m.flattenedItems[i].goal.ID != "" {
					m.cursor = i
					break
				}
			}
		}
	}

	// Adjust scroll to show cursor
	m.adjustScrollOffset()
}

// flattenAncestorChain flattens the ancestor chain into navigable items
func (m *HierarchyScreen) flattenAncestorChain() {
	for i, g := range m.ancestors {
		isLast := i == len(m.ancestors)-1
		isCurrent := i == len(m.ancestors)-1
		isRoot := i == 0

		var prefix string
		if isRoot {
			prefix = "┌─ "
		} else if isLast {
			prefix = "└─ "
		} else {
			prefix = "├─ "
		}

		var goalStyle lipgloss.Style
		if isRoot {
			goalStyle = rootStyle
		} else if isCurrent {
			goalStyle = currentStyle
		} else {
			goalStyle = ancestorStyle
		}

		var status string
		if g.IsDone {
			status = "✓"
		} else {
			status = "○"
		}

		m.flattenedItems = append(m.flattenedItems, TreeItem{
			goal:      g,
			prefix:    prefix,
			status:    status,
			style:     goalStyle,
			isCurrent: isCurrent,
			depth:     i,
		})

		// Add vertical connector lines for non-last items (except root)
		if !isLast {
			m.flattenedItems = append(m.flattenedItems, TreeItem{
				goal:      goal.Goal{}, // Empty goal for connector lines
				prefix:    "│ ",
				status:    "",
				style:     treeCharStyle,
				isCurrent: false,
				depth:     i,
			})
		}
	}
}

// flattenTreeNodes recursively flattens tree nodes into navigable items
func (m *HierarchyScreen) flattenTreeNodes(nodes []TreeNode, prefix string, _ bool, isRoot bool) {
	for i, node := range nodes {
		isNodeLast := i == len(nodes)-1

		var connector string
		if isRoot && i == 0 && len(nodes) == 1 {
			connector = "   "
		} else if isNodeLast {
			connector = "└─ "
		} else {
			connector = "├─ "
		}

		var goalStyle lipgloss.Style
		if node.depth == 0 {
			goalStyle = rootStyle
		} else if node.isCurrent {
			goalStyle = currentStyle
		} else {
			goalStyle = ancestorStyle
		}

		var status string
		if node.goal.IsDone {
			status = "✓"
		} else {
			status = "○"
		}

		m.flattenedItems = append(m.flattenedItems, TreeItem{
			goal:      node.goal,
			prefix:    prefix + connector,
			status:    status,
			style:     goalStyle,
			isCurrent: node.isCurrent,
			depth:     node.depth,
		})

		// Render children if any
		if len(node.children) > 0 {
			childPrefix := prefix
			if !isRoot {
				if isNodeLast {
					childPrefix += "   "
				} else {
					childPrefix += "│  "
				}
			} else {
				childPrefix += "   "
			}
			m.flattenTreeNodes(node.children, childPrefix, isNodeLast, false)
		}

		// Add vertical connector after node if needed
		if !isNodeLast && !isRoot && len(node.children) == 0 {
			m.flattenedItems = append(m.flattenedItems, TreeItem{
				goal:      goal.Goal{},
				prefix:    prefix + "│ ",
				status:    "",
				style:     treeCharStyle,
				isCurrent: false,
				depth:     node.depth,
			})
		}
	}
}

// moveCursor moves the cursor up or down and adjusts scroll offset
// It skips connector lines (items with empty goal.ID) to ensure cursor always points to a valid goal
func (m *HierarchyScreen) moveCursor(delta int) {
	if len(m.flattenedItems) == 0 {
		return
	}

	// Start from the next position
	startPos := m.cursor + delta

	// Search for the next valid item in the direction of movement
	if delta < 0 {
		// Moving up: search backwards from start position
		for i := startPos; i >= 0; i-- {
			if i < len(m.flattenedItems) && m.flattenedItems[i].goal.ID != "" {
				m.cursor = i
				m.adjustScrollOffset()
				return
			}
		}
		// If nothing found going up, stay at current position if it's valid
		if m.cursor < len(m.flattenedItems) && m.flattenedItems[m.cursor].goal.ID != "" {
			m.adjustScrollOffset()
			return
		}
	} else {
		// Moving down: search forwards from start position
		for i := startPos; i < len(m.flattenedItems); i++ {
			if m.flattenedItems[i].goal.ID != "" {
				m.cursor = i
				m.adjustScrollOffset()
				return
			}
		}
		// If nothing found going down, stay at current position if it's valid
		if m.cursor < len(m.flattenedItems) && m.flattenedItems[m.cursor].goal.ID != "" {
			m.adjustScrollOffset()
			return
		}
	}

	// Fallback: find any valid item
	for i := 0; i < len(m.flattenedItems); i++ {
		if m.flattenedItems[i].goal.ID != "" {
			m.cursor = i
			break
		}
	}

	m.adjustScrollOffset()
}

// adjustScrollOffset ensures the cursor is visible in the viewport
func (m *HierarchyScreen) adjustScrollOffset() {
	if len(m.flattenedItems) == 0 {
		return
	}

	headerHeight := lipgloss.Height(headerStyle.Render("Goal Hierarchy"))
	availableHeight := m.height - headerHeight - 2 // Account for padding

	// Calculate how many items can fit in the viewport
	// We'll estimate 1 line per item (this is approximate)
	// Use the same calculation as in renderFlattenedTree
	maxVisibleItems := max(1, availableHeight-1)

	// Calculate the maximum scroll offset (to show the bottom of the list)
	// We want to show the last maxVisibleItems items, so maxScrollOffset should allow that
	maxScrollOffset := max(0, len(m.flattenedItems)-maxVisibleItems)

	// Priority 1: Ensure cursor is visible
	// If cursor is above the visible area, scroll up to show it
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}

	// If cursor is below the visible area, scroll down to show it
	// The cursor should be visible, so it should be within [scrollOffset, scrollOffset+maxVisibleItems-1]
	// We use > instead of >= to ensure we scroll when cursor is exactly at the boundary
	if m.cursor > m.scrollOffset+maxVisibleItems-1 {
		m.scrollOffset = m.cursor - maxVisibleItems + 1
	}

	// Priority 2: Clamp to valid bounds, but ensure cursor remains visible
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}

	// Special handling for when cursor is at the last item
	// Always ensure the last item is visible when cursor is on it
	lastItemIndex := len(m.flattenedItems) - 1
	if m.cursor == lastItemIndex {
		// Cursor is at the very last item - ensure it's visible
		if len(m.flattenedItems) > maxVisibleItems {
			// Set scroll to show the last maxVisibleItems items
			m.scrollOffset = maxScrollOffset
			// Double-check: if cursor is still not visible, adjust
			if m.cursor > m.scrollOffset+maxVisibleItems-1 {
				m.scrollOffset = m.cursor - maxVisibleItems + 1
			}
		} else {
			// List is shorter than viewport, start from beginning
			m.scrollOffset = 0
		}
	} else if m.cursor >= lastItemIndex-maxVisibleItems+1 && m.cursor < lastItemIndex {
		// Cursor is near the last item but not at it - show the bottom
		if len(m.flattenedItems) > maxVisibleItems {
			m.scrollOffset = maxScrollOffset
		} else {
			m.scrollOffset = 0
		}
		// But ensure cursor is still visible
		if m.cursor < m.scrollOffset {
			m.scrollOffset = m.cursor
		}
	} else if m.scrollOffset > maxScrollOffset {
		// Normal case: clamp to maxScrollOffset if cursor would still be visible
		cursorWouldBeVisible := m.cursor <= maxScrollOffset+maxVisibleItems-1
		if cursorWouldBeVisible {
			m.scrollOffset = maxScrollOffset
		}
	}

	// Final safety check: ensure cursor is actually visible with current scrollOffset
	// If not, force it to be visible (this should rarely happen, but acts as a safeguard)
	if m.cursor < m.scrollOffset || m.cursor > m.scrollOffset+maxVisibleItems-1 {
		// Force cursor to be visible
		if m.cursor < m.scrollOffset {
			m.scrollOffset = m.cursor
		} else {
			m.scrollOffset = m.cursor - maxVisibleItems + 1
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}
		}
	}
}

// getSelectedGoal returns the goal at the current cursor position
func (m *HierarchyScreen) getSelectedGoal() *goal.Goal {
	if len(m.flattenedItems) == 0 || m.cursor < 0 || m.cursor >= len(m.flattenedItems) {
		return nil
	}

	item := m.flattenedItems[m.cursor]
	// Skip connector lines (items with empty goal.ID)
	if item.goal.ID == "" {
		return nil
	}

	return &item.goal
}

// openGoalDetailsCmd opens the goal details screen for the selected goal
func (m *HierarchyScreen) openGoalDetailsCmd() tea.Cmd {
	selectedGoal := m.getSelectedGoal()
	if selectedGoal == nil {
		return nil
	}

	return func() tea.Msg {
		return screens.OpenGoalDetailsScreen{
			Goal: selectedGoal,
		}
	}
}

// openTimeframeCmd opens the timeframe screen for the selected goal if it has a timeframe
func (m *HierarchyScreen) openTimeframeCmd() tea.Cmd {
	selectedGoal := m.getSelectedGoal()
	if selectedGoal == nil {
		return nil
	}

	// Check if the goal has a timeframe
	if selectedGoal.Timeframe == nil {
		return nil
	}

	return func() tea.Msg {
		// Life goals don't have dates, so use current time
		var date time.Time
		if selectedGoal.Date != nil {
			date = *selectedGoal.Date
		} else if *selectedGoal.Timeframe == goal.Life {
			date = time.Now()
		} else {
			// Other timeframes require a date
			return nil
		}

		return screens.OpenTimeframeScreenWithGoal{
			Timeframe: *selectedGoal.Timeframe,
			Date:      date,
			GoalID:    selectedGoal.ID,
		}
	}
}
