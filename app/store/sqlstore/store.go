package sqlstore

import (
	"database/sql"
	"github.com/jacobfire/http-rest-api/app/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
	userRepository *UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User returns new created rep
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}