//go:build mage
// +build mage

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/magefile/mage/mg"
)

// SeedWord represents a word from seed JSON
type SeedWord struct {
	Kanji   string `json:"kanji"`
	Romaji  string `json:"romaji"`
	English string `json:"english"`
}

// SeedConfig represents configuration for seeding a group of words
type SeedConfig struct {
	File      string
	GroupName string
}

// DB initializes the database connection
func DB() (*sql.DB, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	
	return db, nil
}

// getDBPath returns the appropriate database path based on environment
func getDBPath() string {
	if os.Getenv("TEST_DB") == "true" {
		fmt.Println("Using test database: words.test.db")
		return "words.test.db"
	}
	fmt.Println("Using production database: words.db")
	return "words.db"
}

// Init creates a new database file if it doesn't exist
func Init() error {
	dbPath := getDBPath()
	if _, err := os.Stat(dbPath); err == nil {
		fmt.Printf("Database already exists at %s\n", dbPath)
		return nil
	}

	file, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database file: %w", err)
	}
	file.Close()
	fmt.Printf("Created new database at %s\n", dbPath)
	return nil
}

// Migrate runs all database migrations
func Migrate() error {
	fmt.Println("Running database migrations...")
	
	db, err := DB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Run migrations
	migrations, err := filepath.Glob("db/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to find migration files: %w", err)
	}

	for _, migration := range migrations {
		fmt.Printf("Running migration: %s\n", filepath.Base(migration))
		sqlBytes, err := os.ReadFile(migration)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migration, err)
		}

		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration, err)
		}
	}

	fmt.Println("Migrations completed successfully")
	return nil
}

// Seed imports seed data into the database
func Seed() error {
	fmt.Println("Seeding database...")
	
	db, err := DB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Define seed configurations
	seedConfigs := []SeedConfig{
		{File: "db/seeds/basic_greetings.json", GroupName: "Basic Greetings"},
		{File: "db/seeds/numbers.json", GroupName: "Numbers"},
		{File: "db/seeds/colors.json", GroupName: "Colors"},
		{File: "db/seeds/family.json", GroupName: "Family"},
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for _, config := range seedConfigs {
		fmt.Printf("Processing seed file: %s\n", config.File)
		
		// Read the seed file
		data, err := os.ReadFile(config.File)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to read seed file %s: %w", config.File, err)
		}

		// Parse the JSON data
		var words []SeedWord
		if err := json.Unmarshal(data, &words); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to parse seed file %s: %w", config.File, err)
		}

		// Insert group if it doesn't exist
		var groupID int64
		err = tx.QueryRow("INSERT OR IGNORE INTO groups (name) VALUES (?) RETURNING id", config.GroupName).Scan(&groupID)
		if err == sql.ErrNoRows {
			err = tx.QueryRow("SELECT id FROM groups WHERE name = ?", config.GroupName).Scan(&groupID)
		}
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get/create group %s: %w", config.GroupName, err)
		}

		// Insert words
		for _, word := range words {
			var wordID int64
			err := tx.QueryRow(`
				INSERT OR IGNORE INTO words (japanese, romaji, english) 
				VALUES (?, ?, ?) 
				RETURNING id
			`, word.Kanji, word.Romaji, word.English).Scan(&wordID)
			
			if err != nil && err != sql.ErrNoRows {
				tx.Rollback()
				return fmt.Errorf("failed to insert word %v: %w", word, err)
			}

			// If word already exists, get its ID
			if err == sql.ErrNoRows {
				err = tx.QueryRow("SELECT id FROM words WHERE japanese = ?", word.Kanji).Scan(&wordID)
				if err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to get word ID: %w", err)
				}
			}

			// Associate word with group
			_, err = tx.Exec(`
				INSERT OR IGNORE INTO words_groups (word_id, group_id) 
				VALUES (?, ?)
			`, wordID, groupID)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to associate word with group: %w", err)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Println("Database seeded successfully!")
	return nil
}

// Reset drops all tables and recreates them
func Reset() error {
	fmt.Println("Resetting database...")
	
	dbPath := getDBPath()
	
	// Remove existing database
	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove database: %w", err)
	}

	// Reinitialize database
	if err := Init(); err != nil {
		return err
	}
	if err := Migrate(); err != nil {
		return err
	}
	if err := Seed(); err != nil {
		return err
	}

	fmt.Println("Database reset successfully!")
	return nil
}

// TestDB runs all database tests
func TestDB() error {
	mg.Deps(Reset)
	fmt.Println("Running database tests...")
	// Add your test commands here
	return nil
}
