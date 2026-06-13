package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	DB *sql.DB
}

func NewSQLite(path string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	s := &SQLite{
		DB: db,
	}

	if err := s.migrate(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SQLite) Close() error {
	return s.DB.Close()
}

func (s *SQLite) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS attachments (
		id TEXT PRIMARY KEY,
		ticket_id TEXT,
		file_name TEXT,
		status TEXT
	);
	`

	_, err := s.DB.Exec(query)
	return err
}
