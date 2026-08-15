package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping DB:", err)
	}

	log.Println("postgress connected")
	DB=db
}


func MigrateDB() {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
	id SERIAL PRIMARY KEY,
	code VERCHAR(10) UNIQUE NOT NULL,
	original TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
	);
	`

	if _, err := DB.Exec(query); err != nil {
		log.Fatal("Migration faied:", err)
	}
	log.Println("URLs table ready")
}