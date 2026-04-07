package database

import (
	"database/sql"
	"fmt"
	"log"

	"RestApiGo/internal/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(cfg *config.Config) error {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Printf("✅ Connected to PostgreSQL at %s:%s", cfg.DBHost, cfg.DBPort)
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("✅ Database connection closed")
	}
}