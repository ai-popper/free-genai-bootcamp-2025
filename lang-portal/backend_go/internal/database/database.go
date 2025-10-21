package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	_, err = DB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return DB.Ping()
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// RunMigrations executes all SQL migration files in order
func RunMigrations(migrationsDir string) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort files to ensure they run in order
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, filename := range sqlFiles {
		fmt.Printf("Running migration: %s\n", filename)
		content, err := os.ReadFile(filepath.Join(migrationsDir, filename))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		_, err = DB.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}
	}

	fmt.Println("All migrations completed successfully")
	return nil
}

// ResetHistory deletes all study sessions and review items
func ResetHistory() error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM word_review_items")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM study_sessions")
	if err != nil {
		return err
	}

	return tx.Commit()
}

// FullReset drops all tables and recreates them
func FullReset(migrationsDir string) error {
	// Drop all tables
	tables := []string{
		"word_review_items",
		"study_sessions",
		"study_activities",
		"words_groups",
		"groups",
		"words",
	}

	for _, table := range tables {
		_, err := DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	// Re-run migrations
	return RunMigrations(migrationsDir)
}
