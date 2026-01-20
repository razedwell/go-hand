package postgres

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

func BuildDBUrl(user, password, host, port, dbname, sslmode string) string {
	return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode
}

func NewClient(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.New("failed to open database connection")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.New("failed to ping database")
	}

	db.SetMaxOpenConns(25)                 // Max total connections
	db.SetMaxIdleConns(25)                 // Max connections sitting doing nothing
	db.SetConnMaxLifetime(5 * time.Minute) // Periodically kill old connections
	return db, nil
}
