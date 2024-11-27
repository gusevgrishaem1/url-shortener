package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/gusevgrishaem1/url-shortener/internal/shortener/model"
)

// PostgresStorage реализация интерфейса Storage для PostgreSQL.
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage создает новое соединение с базой данных PostgreSQL.
func NewPostgresStorage(connStr string, maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	sqlFile := "./db_scripts/001_create_tables.sql"
	script, err := os.ReadFile(sqlFile)
	if err != nil {
		return nil, err
	}

	// Выполнение SQL-скрипта
	_, err = db.ExecContext(context.TODO(), string(script))
	if err != nil {
		return nil, err
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	return &PostgresStorage{db: db}, nil
}

// Save сохраняет оригинальный и сокращенный URL в базе данных.
func (s *PostgresStorage) Save(url model.URL) error {
	query := `INSERT INTO urls (original_url, short_url) VALUES ($1, $2)`
	_, err := s.db.Exec(query, url.Original, url.Short)
	return err
}

// Get получает оригинальный URL по сокращенному URL из базы данных.
func (s *PostgresStorage) Get(short string) (string, bool) {
	query := `SELECT original_url, create_ts FROM urls WHERE short_url = $1`
	row := s.db.QueryRow(query, short)

	var original string
	var createTS time.Time
	err := row.Scan(&original, &createTS)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false
		}
		log.Println("Error querying original URL:", err)
		return "", false
	}

	return original, true
}
