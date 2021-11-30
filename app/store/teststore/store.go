package teststore

import (
	"github.com/jacobfire/http-rest-api/app/model"
	"github.com/jacobfire/http-rest-api/app/store"
)

type Store struct {
	userRepository *UserRepository
}

func New() *Store {
	return &Store{}
}

// User returns new created rep
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[string]*model.User),
	}

	return s.userRepository

}