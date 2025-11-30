package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetBackupDir reads the backup directory from ~/.hinoki.rc
// Creates the config file with default settings if it doesn't exist
// Returns the backup directory path, or an error if unable to read or create
func GetBackupDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".hinoki.rc")

	// Check if config file exists, if not create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultBackupDir := filepath.Join(homeDir, "Documents", "hinoki-backups")
		configContent := fmt.Sprintf("# Hinoki Planner Configuration\n# Backup directory for database backups\nbackup_dir=%s\n", defaultBackupDir)

		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			return "", fmt.Errorf("failed to create config file: %w", err)
		}

		// Return the default directory we just created
		return defaultBackupDir, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		if strings.HasPrefix(line, "backup_dir=") {
			dir := strings.TrimPrefix(line, "backup_dir=")
			dir = strings.TrimSpace(dir)
			// Expand ~ in path
			if strings.HasPrefix(dir, "~/") {
				dir = filepath.Join(homeDir, dir[2:])
			}
			return dir, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading config file: %w", err)
	}

	// If backup_dir not found in existing config, add it
	defaultBackupDir := filepath.Join(homeDir, "Documents", "hinoki-backups")
	configContent := fmt.Sprintf("\n# Backup directory for database backups\nbackup_dir=%s\n", defaultBackupDir)

	file.Close() // Close the read file

	// Append to existing config file
	file, err = os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to append to config file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(configContent); err != nil {
		return "", fmt.Errorf("failed to write to config file: %w", err)
	}

	return defaultBackupDir, nil
}

// GetBackupDirOrDefault returns the backup directory from config, or a default location
func GetBackupDirOrDefault() (string, error) {
	dir, err := GetBackupDir()
	if err != nil {
		// Return default backup directory in user's Documents folder
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		defaultDir := filepath.Join(homeDir, "Documents", "hinoki-backups")
		return defaultDir, nil
	}
	return dir, nil
}
