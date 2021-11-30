package store

import "github.com/jacobfire/http-rest-api/app/model"

// UserRepository ....
type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(string) (*model.User, error)
}
