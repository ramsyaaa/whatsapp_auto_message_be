package main

import (
	"go_whatsapp/config"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables instead")
	}

	// Connect to database
	db, _ := config.Connect()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}

	// Get all SQL files in the migrations directory
	files, err := ioutil.ReadDir("./migrations")
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Execute each SQL file
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			log.Printf("Running migration: %s", file.Name())

			// Read SQL file
			sqlFile := filepath.Join("./migrations", file.Name())
			sqlBytes, err := ioutil.ReadFile(sqlFile)
			if err != nil {
				log.Fatalf("Failed to read SQL file %s: %v", sqlFile, err)
			}

			// Split SQL file into statements
			sqlContent := string(sqlBytes)
			statements := strings.Split(sqlContent, ";")

			// Execute each statement
			for _, statement := range statements {
				statement = strings.TrimSpace(statement)
				if statement == "" {
					continue
				}

				_, err = sqlDB.Exec(statement)
				if err != nil {
					log.Fatalf("Failed to execute SQL statement: %v", err)
				}
			}

			log.Printf("Migration %s completed successfully", file.Name())
		}
	}

	log.Println("All migrations completed successfully")
}
