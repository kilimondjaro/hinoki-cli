package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

const (
	migrationsDir = "./internal/db/migrations"
)

var (
	instance *sql.DB
	once     sync.Once
	mu       sync.Mutex
)

// •	macOS: ~/Library/Application Support/<your-app-name>/local.db
// •	Linux: ~/.config/<your-app-name>/local.db
// •	Windows: %AppData%\<your-app-name>\local.db
func getDBPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	hinokiCliDir := filepath.Join(configDir, "hinoki-cli")

	if _, err := os.Stat(hinokiCliDir); os.IsNotExist(err) {
		err := os.Mkdir(hinokiCliDir, 0700)
		return "", err
	}

	return filepath.Join(configDir, "hinoki-cli", "hinoki.db"), nil
}

func InitDB() *sql.DB {
	once.Do(func() {
		path, err := getDBPath()
		if err != nil {
			panic(err)
		}

		inst, err := sql.Open("sqlite3", path)
		instance = inst

		if err != nil {
			panic(err)
		}

		err = createSchemaVersionTable()
		if err != nil {
			panic(err)
		}

		err = applyMigrations(inst)
		if err != nil {
			panic(err)
		}
	})

	return instance
}

func CloseDB() {
	instance.Close()
}

func ExecQuery(query string, args ...interface{}) (sql.Result, error) {
	mu.Lock()
	defer mu.Unlock()
	return instance.Exec(query, args...)
}

func QueryDB(query string, args ...interface{}) (*sql.Rows, error) {
	mu.Lock()
	defer mu.Unlock()
	return instance.Query(query, args...)
}

func createSchemaVersionTable() error {
	_, err := ExecQuery(`
		CREATE TABLE IF NOT EXISTS schema_version (
    		version INTEGER PRIMARY KEY,
    		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)

	return err
}

func applyMigrations(db *sql.DB) error {
	// Retrieve current schema version
	var currentVersion int
	err := db.QueryRow("SELECT version FROM schema_version ORDER BY version DESC LIMIT 1").Scan(&currentVersion)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Migration scripts: For each version, apply the corresponding migrationsArr
	var migrationsArr []struct {
		version int
		query   string
	}

	for version, query := range migrations {
		if version > currentVersion {
			migrationsArr = append(migrationsArr, struct {
				version int
				query   string
			}{version: version, query: query})
		}
	}

	// Apply migrationsArr in order
	sort.Slice(migrationsArr, func(i, j int) bool {
		return migrationsArr[i].version < migrationsArr[j].version
	})

	for _, migration := range migrationsArr {
		_, err = db.Exec(migration.query)
		if err != nil {
			return fmt.Errorf("failed to apply migrationsArr version %d: %w", migration.version, err)
		}

		// Update the schema_version table
		_, err = db.Exec("INSERT INTO schema_version (version) VALUES (?)", migration.version)
		if err != nil {
			return fmt.Errorf("failed to update schema_version table: %w", err)
		}

		fmt.Println("Applied migrationsArr version:", migration.version)
	}
	return nil
}
