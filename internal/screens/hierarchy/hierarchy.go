package hierarchy

import (
	"fmt"
	"hinoki-cli/internal/dates"
	"hinoki-cli/internal/goal"
	"hinoki-cli/internal/repository"
	"hinoki-cli/internal/screens"
	"hinoki-cli/internal/theme"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HierarchyScreen struct {
	goal      *goal.Goal
	ancestors []goal.Goal
	keys      keyMap
	showAll   bool // Toggle for showing full tree vs just ancestors

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
	return m.getAncestorChainCmd()
}

func (m *HierarchyScreen) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case AncestorChainResult:
		m.ancestors = msg.ancestors
	case FullTreeResult:
		// Convert tree nodes back to ancestors for rendering
		m.ancestors = m.treeNodesToGoals(msg.treeNodes)
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
		if m.showAll {
			return m.getFullTreeCmd()
		}
		return m.getAncestorChainCmd()
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

	if m.showAll {
		return m.renderFullTree(width, height)
	}

	return m.renderAncestorChain(width, height)
}

func (m *HierarchyScreen) renderAncestorChain(width, height int) string {
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

func (m *HierarchyScreen) renderFullTree(width, height int) string {
	// Build tree structure
	treeNodes := m.buildFullTree()

	// Render tree nodes recursively
	// Start with empty prefix and mark as root level
	var lines []string
	m.renderTreeNode(treeNodes, &lines, "", false, true)

	// Join all lines
	treeContent := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Wrap in a container with proper width
	return lipgloss.NewStyle().
		Width(width).
		MaxHeight(height).
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

func (m *HierarchyScreen) renderTreeNode(nodes []TreeNode, lines *[]string, prefix string, _ bool, isRoot bool) {
	for i, node := range nodes {
		isNodeLast := i == len(nodes)-1

		// Determine tree characters for this node
		// Only use ┌─ for the very first root node
		var connector string
		if isRoot && i == 0 && len(nodes) == 1 {
			connector = treeCharStyle.Render("┌─ ")
		} else if isNodeLast {
			connector = treeCharStyle.Render("└─ ")
		} else {
			connector = treeCharStyle.Render("├─ ")
		}

		// Determine style
		var goalStyle lipgloss.Style
		if node.depth == 0 {
			goalStyle = rootStyle
		} else if node.isCurrent {
			goalStyle = currentStyle
		} else {
			goalStyle = ancestorStyle
		}

		// Add status indicator
		var status string
		if node.goal.IsDone {
			status = metaStyle.Render("✓")
		} else {
			status = metaStyle.Render("○")
		}

		// Build the goal line
		goalLine := goalStyle.Render(node.goal.Title)

		// Add metadata
		meta := m.formatMeta(&node.goal)
		if meta != "" {
			goalLine = fmt.Sprintf("%s %s", goalLine, metaStyle.Render(meta))
		}

		// Combine everything
		fullLine := fmt.Sprintf("%s%s%s %s", prefix, connector, status, goalLine)
		*lines = append(*lines, fullLine)

		// Render children if any
		if len(node.children) > 0 {
			// Build prefix for children - properly indent based on parent's position
			childPrefix := prefix
			if !isRoot {
				// If this node is last, use spaces; otherwise use vertical line
				if isNodeLast {
					childPrefix += treeCharStyle.Render("   ")
				} else {
					childPrefix += treeCharStyle.Render("│  ")
				}
			} else {
				// For root level, add spacing for children
				childPrefix += treeCharStyle.Render("   ")
			}

			// Recursively render children with proper indentation
			// Pass isNodeLast so children know if their parent is the last sibling
			m.renderTreeNode(node.children, lines, childPrefix, isNodeLast, false)
		}

		// Add vertical connector after node if it's not the last and has siblings after
		// Only add if node has no children (to avoid double connectors)
		if !isNodeLast && !isRoot && len(node.children) == 0 {
			*lines = append(*lines, prefix+treeCharStyle.Render("│ "))
		}
	}
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
