package goallist

import (
	"database/sql"
	"fmt"
	"hinoki-cli/internal/dates"
	"hinoki-cli/internal/db"
	"hinoki-cli/internal/goal"
	"time"
)

func getGoalsByDate(timeframe goal.Timeframe, date time.Time) ([]goal.Goal, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `
		SELECT id, title, created_at, updated_at, is_done, timeframe, date
		FROM goals
	`
	orderByQuery := `
		ORDER BY is_done ASC, date DESC;
	`

	filterArchivedQuery := `AND is_archived IS NOT true`

	composeQuery := func(query string) string {
		return baseQuery + query + filterArchivedQuery + orderByQuery
	}

	switch timeframe {
	case goal.Day:
		rows, err = db.QueryDB(
			composeQuery(`WHERE timeframe = ? AND DATE(date) = ?`),
			string(timeframe),
			dates.TimeframeDateString(date),
		)
	case goal.Week:
		rows, err = db.QueryDB(
			composeQuery(`WHERE timeframe = ? AND DATE(date) >= ? AND DATE(date) <= ?`),
			string(timeframe),
			dates.TimeframeDateString(dates.StartOfWeek(date)),
			dates.TimeframeDateString(dates.EndOfWeek(date)),
		)
	case goal.Month:
		rows, err = db.QueryDB(
			composeQuery(`WHERE timeframe = ? AND DATE(date) LIKE ?`),
			string(timeframe),
			date.Format("2006-01%"),
		)
	case goal.Quarter:
		rows, err = db.QueryDB(
			composeQuery(`WHERE timeframe = ? AND DATE(date) > ? AND DATE(date) < ?`),
			string(timeframe),
			dates.TimeframeDateString(dates.StartOfQuarter(date)),
			dates.TimeframeDateString(dates.EndOfQuarter(date)),
		)
	case goal.Year:
		rows, err = db.QueryDB(
			composeQuery(`WHERE timeframe = ? AND DATE(date) LIKE ?`),
			string(timeframe),
			date.Format("2006%"),
		)
	case goal.Life:
		rows, err = db.QueryDB(
			composeQuery(`WHERE timeframe = ?`),
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
		if err := rows.Scan(&goal.ID, &goal.Title, &goal.CreatedAt, &goal.UpdatedAt, &goal.IsDone, &goal.Timeframe, &goal.Date); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		goals = append(goals, goal)
	}
	return goals, rows.Err()
}

func addGoal(goal goal.Goal) error {
	_, err := db.ExecQuery("INSERT INTO goals (id, title, is_done, timeframe, date) VALUES (?, ?, ?, ?, ?)", goal.ID, goal.Title, goal.IsDone, goal.Timeframe, goal.Date)

	return err
}

func updateGoal(goal goal.Goal) error {
	_, err := db.ExecQuery("UPDATE goals SET title = ?, is_done = ?, timeframe = ?, date = ?, is_archived = ? WHERE id = ?", goal.Title, goal.IsDone, goal.Timeframe, goal.Date, goal.IsArchived, goal.ID)

	return err
}
