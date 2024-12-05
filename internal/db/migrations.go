package db

const (
	createGoalsTable = `
	CREATE TABLE IF NOT EXISTS goals (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_done BOOLEAN NOT NULL,
		timeframe TEXT CHECK(timeframe IN ('day', 'week', 'month', 'quarter', 'year', 'life')),
	   	date DATETIME                              
	)`
	addArchivedToGoals = `ALTER TABLE goals ADD COLUMN is_archived BOOLEAN;`
)

var migrations = map[int]string{
	1: createGoalsTable,
	2: addArchivedToGoals,
}
