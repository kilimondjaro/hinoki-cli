package db

import (
	"fmt"
	"hinoki-cli/internal/config"
	"io"
	"os"
	"path/filepath"
	"time"
)

// CreateBackup creates a backup of the database file in the configured backup directory
// Returns the full path to the created backup file
func CreateBackup() (string, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return "", fmt.Errorf("failed to get DB path: %w", err)
	}

	// Get backup directory from config (creates config file if it doesn't exist)
	backupDir, err := config.GetBackupDir()
	if err != nil {
		return "", fmt.Errorf("failed to get backup directory: %w", err)
	}

	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("hinoki_backup_%s.db", timestamp)
	backupPath := filepath.Join(backupDir, backupName)

	// Copy database file
	source, err := os.Open(dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open source DB: %w", err)
	}
	defer source.Close()

	dest, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		os.Remove(backupPath) // Clean up on error
		return "", fmt.Errorf("failed to copy DB: %w", err)
	}

	return backupPath, nil
}
