package main

import (
	"log"

	"lms/internal/config"
	"lms/internal/database"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := database.SeedData(db); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	log.Println("Database seeding completed successfully.")
}
