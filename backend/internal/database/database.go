package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"lms/internal/config"
)

type DB struct {
	*sql.DB
	Driver string
}

func Connect(cfg *config.Config) (*DB, error) {
	var db *sql.DB
	var err error
	driver := cfg.DBDriver

	if driver == "postgres" {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=3",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
		)
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
		}

		if err != nil {
			log.Printf("[DB WARNING] PostgreSQL connection to %s:%s failed: %v. Falling back to SQLite for seamless local execution.", cfg.DBHost, cfg.DBPort, err)
			driver = "sqlite"
		}
	}

	if driver == "sqlite" {
		db, err = sql.Open("sqlite", cfg.DBSQLitePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open SQLite database: %w", err)
		}
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
		}
		// Enable foreign keys and WAL mode for SQLite
		_, _ = db.Exec("PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;")
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("[DB] Connected successfully using driver: %s", driver)
	return &DB{DB: db, Driver: driver}, nil
}
