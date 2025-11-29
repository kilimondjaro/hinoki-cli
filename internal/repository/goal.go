package repository

import (
	"database/sql"
	"fmt"
	"hinoki-cli/internal/dates"
	"hinoki-cli/internal/db"
	"hinoki-cli/internal/goal"
	"strings"
	"time"
)

// GetGoalsByParent retrieves all goals that have the specified parent ID
func GetGoalsByParent(parentId string) ([]goal.Goal, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `
		SELECT id, parent_id, title, created_at, updated_at, is_done, timeframe, date
		FROM goals
	`

	orderByQuery := `
		ORDER BY is_done ASC, created_at ASC;
	`

	filterArchivedQuery := `AND is_archived IS NOT true`

	composeQuery := func(query string) string {
		return baseQuery + query + filterArchivedQuery + orderByQuery
	}

	rows, err = db.QueryDB(
		composeQuery(`WHERE parent_id = ?`),
		parentId,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var goals []goal.Goal

	for rows.Next() {
		var goal goal.Goal
		if err := rows.Scan(&goal.ID, &goal.ParentId, &goal.Title, &goal.CreatedAt, &goal.UpdatedAt, &goal.IsDone, &goal.Timeframe, &goal.Date); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		goals = append(goals, goal)
	}
	return goals, rows.Err()
}

// GetGoalsByDate retrieves all goals for a specific timeframe and date
func GetGoalsByDate(timeframe goal.Timeframe, date time.Time) ([]goal.Goal, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `
		SELECT g.id, g.title, g.created_at, g.updated_at, g.is_done, g.timeframe, g.date, p.id, p.title
		FROM goals g
		LEFT JOIN goals p ON g.parent_id = p.id
	`
	orderByQuery := `
		ORDER BY g.is_done ASC, g.created_at ASC;
	`

	filterArchivedQuery := `AND g.is_archived IS NOT true`

	composeQuery := func(query string) string {
		return baseQuery + query + filterArchivedQuery + orderByQuery
	}

	switch timeframe {
	case goal.Day:
		rows, err = db.QueryDB(
			composeQuery(`WHERE g.timeframe = ? AND DATE(g.date) = ?`),
			string(timeframe),
			dates.TimeframeDateString(date),
		)
	case goal.Week:
		rows, err = db.QueryDB(
			composeQuery(`WHERE g.timeframe = ? AND DATE(g.date) >= ? AND DATE(g.date) <= ?`),
			string(timeframe),
			dates.TimeframeDateString(dates.StartOfWeek(date)),
			dates.TimeframeDateString(dates.EndOfWeek(date)),
		)
	case goal.Month:
		rows, err = db.QueryDB(
			composeQuery(`WHERE g.timeframe = ? AND DATE(g.date) LIKE ?`),
			string(timeframe),
			date.Format("2006-01%"),
		)
	case goal.Quarter:
		rows, err = db.QueryDB(
			composeQuery(`WHERE g.timeframe = ? AND DATE(g.date) > ? AND DATE(g.date) < ?`),
			string(timeframe),
			dates.TimeframeDateString(dates.StartOfQuarter(date)),
			dates.TimeframeDateString(dates.EndOfQuarter(date)),
		)
	case goal.Year:
		rows, err = db.QueryDB(
			composeQuery(`WHERE g.timeframe = ? AND DATE(g.date) LIKE ?`),
			string(timeframe),
			date.Format("2006%"),
		)
	case goal.Life:
		rows, err = db.QueryDB(
			composeQuery(`WHERE g.timeframe = ?`),
			string(timeframe),
		)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var goals []goal.Goal

	for rows.Next() {
		var goal goal.Goal

		if err := rows.Scan(&goal.ID, &goal.Title, &goal.CreatedAt, &goal.UpdatedAt, &goal.IsDone, &goal.Timeframe, &goal.Date, &goal.ParentId, &goal.ParentTitle); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		goals = append(goals, goal)
	}
	return goals, rows.Err()
}

// GetGoalByID retrieves a single goal by its ID
func GetGoalByID(goalID string) (*goal.Goal, error) {
	query := `
		SELECT id, parent_id, title, created_at, updated_at, is_done, timeframe, date, COALESCE(is_archived, 0) as is_archived
		FROM goals
		WHERE id = ? AND COALESCE(is_archived, 0) = 0
	`

	row := db.QueryRowDB(query, goalID)

	var g goal.Goal
	err := row.Scan(&g.ID, &g.ParentId, &g.Title, &g.CreatedAt, &g.UpdatedAt, &g.IsDone, &g.Timeframe, &g.Date, &g.IsArchived)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

// AddGoal creates a new goal in the database
func AddGoal(goal goal.Goal) error {
	_, err := db.ExecQuery("INSERT INTO goals (id, parent_id, title, is_done, timeframe, date) VALUES (?, ?, ?, ?, ?, ?)", goal.ID, goal.ParentId, goal.Title, goal.IsDone, goal.Timeframe, goal.Date)

	return err
}

// UpdateGoal updates an existing goal in the database
func UpdateGoal(goal goal.Goal) error {
	_, err := db.ExecQuery("UPDATE goals SET title = ?, is_done = ?, timeframe = ?, date = ?, is_archived = ?, parent_id = ? WHERE id = ?", goal.Title, goal.IsDone, goal.Timeframe, goal.Date, goal.IsArchived, goal.ParentId, goal.ID)

	return err
}

// SearchGoals searches for goals matching the given search term
func SearchGoals(term string, limit int) ([]goal.Goal, error) {
	if limit <= 0 {
		limit = 20
	}

	trimmed := strings.TrimSpace(term)
	if trimmed == "" {
		return []goal.Goal{}, nil
	}

	query := `
		SELECT g.id, g.title, g.created_at, g.updated_at, g.is_done, g.timeframe, g.date, g.parent_id, p.title
		FROM goals g
		LEFT JOIN goals p ON g.parent_id = p.id
		WHERE g.is_archived IS NOT true AND LOWER(g.title) LIKE LOWER(?)
		ORDER BY g.date IS NULL, g.date DESC, g.updated_at DESC
		LIMIT ?
	`

	rows, err := db.QueryDB(query, "%"+trimmed+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []goal.Goal

	for rows.Next() {
		var g goal.Goal
		if err := rows.Scan(&g.ID, &g.Title, &g.CreatedAt, &g.UpdatedAt, &g.IsDone, &g.Timeframe, &g.Date, &g.ParentId, &g.ParentTitle); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		goals = append(goals, g)
	}

	return goals, rows.Err()
}

// GetAncestorChain retrieves all ancestors of a goal, from the goal itself up to the root parent
// Returns goals in order from root (topmost parent) to the goal itself
func GetAncestorChain(goalID string) ([]goal.Goal, error) {
	var chain []goal.Goal
	currentID := goalID
	visited := make(map[string]bool) // Prevent infinite loops

	for currentID != "" {
		// Prevent infinite loops
		if visited[currentID] {
			break
		}
		visited[currentID] = true

		g, err := GetGoalByID(currentID)
		if err != nil || g == nil {
			break
		}

		// Prepend to chain (we want root first, goal last)
		chain = append([]goal.Goal{*g}, chain...)

		// Move to parent
		if g.ParentId == nil || *g.ParentId == "" {
			break
		}
		currentID = *g.ParentId
	}

	return chain, nil
}

// GetOverdueGoals retrieves all undone goals that are overdue
// A goal is overdue if it has a date and timeframe, is not done, and the period has passed
func GetOverdueGoals() ([]goal.Goal, error) {
	today := dates.DateWithoutTime(time.Now())

	baseQuery := `
		SELECT g.id, g.title, g.created_at, g.updated_at, g.is_done, g.timeframe, g.date, g.parent_id, p.title
		FROM goals g
		LEFT JOIN goals p ON g.parent_id = p.id
		WHERE g.is_archived IS NOT true 
		AND g.is_done = 0
		AND g.timeframe IS NOT NULL
		AND g.date IS NOT NULL
		AND (
			(g.timeframe = 'day' AND DATE(g.date) < ?)
			OR (g.timeframe = 'week' AND DATE(g.date) < ?)
			OR (g.timeframe = 'month' AND DATE(g.date) < ?)
			OR (g.timeframe = 'quarter' AND DATE(g.date) < ?)
			OR (g.timeframe = 'year' AND DATE(g.date) < ?)
		)
		ORDER BY g.date DESC, g.created_at DESC
	`

	// Initial SQL filter to get potential overdue goals
	// The dates.IsOverdue function will do the final accurate check
	todayStr := dates.TimeframeDateString(today)

	rows, err := db.QueryDB(baseQuery, todayStr, todayStr, todayStr, todayStr, todayStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []goal.Goal

	for rows.Next() {
		var g goal.Goal
		if err := rows.Scan(&g.ID, &g.Title, &g.CreatedAt, &g.UpdatedAt, &g.IsDone, &g.Timeframe, &g.Date, &g.ParentId, &g.ParentTitle); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		// Additional check: filter out goals that aren't actually overdue
		// This handles edge cases for week/month/quarter/year timeframes
		if dates.IsOverdue(g.Date, g.Timeframe) {
			goals = append(goals, g)
		}
	}

	return goals, rows.Err()
}
