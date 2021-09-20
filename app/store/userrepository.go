package store

import (
	"database/sql"
	"github.com/jacobfire/http-rest-api/app/model"
)

// UserRepository for managing models
type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	if err := r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}

	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}