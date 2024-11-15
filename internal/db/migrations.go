package db

const (
	addArchivedToGoals = `ALTER TABLE goals ADD COLUMN is_archived BOOLEAN;`
)

var migrations = map[int]string{
	1: addArchivedToGoals,
}
