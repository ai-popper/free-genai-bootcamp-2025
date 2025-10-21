//go:build mage
// +build mage

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lang-portal/backend/internal/database"
	"github.com/magefile/mage/mg"
)

// SeedWord represents a word from seed JSON
type SeedWord struct {
	Kanji   string `json:"kanji"`
	Romaji  string `json:"romaji"`
	English string `json:"english"`
}

// SeedConfig maps seed files to group names
type SeedConfig struct {
	File      string
	GroupName string
}

// Init initializes the database file
func Init() error {
	fmt.Println("Initializing database...")
	
	dbPath := "words.db"
	
	// Check if database already exists
	if _, err := os.Stat(dbPath); err == nil {
		fmt.Println("Database already exists at", dbPath)
		return nil
	}
	
	// Create empty database file
	file, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database file: %w", err)
	}
	file.Close()
	
	fmt.Println("Database file created:", dbPath)
	return nil
}

// Migrate runs all database migrations
func Migrate() error {
	fmt.Println("Running database migrations...")
	
	dbPath := "words.db"
	err := database.InitDB(dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.CloseDB()
	
	migrationsDir := filepath.Join("db", "migrations")
	err = database.RunMigrations(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	
	fmt.Println("Migrations completed successfully")
	return nil
}

// Seed imports seed data into the database
func Seed() error {
	fmt.Println("Seeding database...")
	
	dbPath := "words.db"
	err := database.InitDB(dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.CloseDB()
	
	// Define seed configurations
	seedConfigs := []SeedConfig{
		{File: "db/seeds/basic_greetings.json", GroupName: "Basic Greetings"},
		{File: "db/seeds/numbers.json", GroupName: "Numbers"},
		{File: "db/seeds/colors.json", GroupName: "Colors"},
		{File: "db/seeds/family.json", GroupName: "Family"},
		{File: "db/seeds/food.json", GroupName: "Food"},
	}
	
	for _, config := range seedConfigs {
		if _, err := os.Stat(config.File); os.IsNotExist(err) {
			fmt.Printf("Skipping %s (file not found)\n", config.File)
			continue
		}
		
		err := seedFile(config.File, config.GroupName)
		if err != nil {
			fmt.Printf("Warning: Failed to seed %s: %v\n", config.File, err)
			continue
		}
		
		fmt.Printf("Seeded: %s -> %s\n", config.File, config.GroupName)
	}
	
	fmt.Println("Seeding completed")
	return nil
}

// Setup runs Init, Migrate, and Seed in sequence
func Setup() error {
	mg.Deps(Init)
	mg.Deps(Migrate)
	mg.Deps(Seed)
	return nil
}

// seedFile reads a JSON seed file and imports words into a group
func seedFile(filePath, groupName string) error {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	
	// Parse JSON
	var words []SeedWord
	err = json.Unmarshal(data, &words)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	// Create or get group
	groupID, err := database.InsertGroup(groupName)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}
	
	// Insert words and link to group
	for _, word := range words {
		wordID, err := database.InsertWord(word.Kanji, word.Romaji, word.English, "")
		if err != nil {
			return fmt.Errorf("failed to insert word %s: %w", word.Kanji, err)
		}
		
		err = database.LinkWordToGroup(wordID, groupID)
		if err != nil {
			return fmt.Errorf("failed to link word to group: %w", err)
		}
	}
	
	return nil
}

// Clean removes the database file
func Clean() error {
	fmt.Println("Cleaning database...")
	
	dbPath := "words.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("No database file to clean")
		return nil
	}
	
	err := os.Remove(dbPath)
	if err != nil {
		return fmt.Errorf("failed to remove database: %w", err)
	}
	
	fmt.Println("Database cleaned")
	return nil
}

// Reset cleans and sets up the database from scratch
func Reset() error {
	mg.Deps(Clean)
	mg.Deps(Setup)
	return nil
}
