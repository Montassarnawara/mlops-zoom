// hotel-sync/config/postgres.go
package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func ConnectPostgres() (*sql.DB, error) {
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5433",
		User:     "myuser",
		Password: "mypassword",
		DBName:   "hotel_db",
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
