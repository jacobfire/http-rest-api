package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	config *Config
	db *sql.DB
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

// Open connection and ping db
func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	s.db = db
	return nil
}

// Close connection
func (s *Store) Close() {
	s.db.Close()
}