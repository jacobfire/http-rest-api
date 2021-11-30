package sqlstore_test

import (
	"github.com/jacobfire/http-rest-api/app/model"
	"github.com/jacobfire/http-rest-api/app/store"
	"github.com/jacobfire/http-rest-api/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("users")

	email := "test@test.com"

	s := sqlstore.New(db)

	u, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(&model.User {
		Email: "test@test.com",
	})

	u = model.TestUser(t)
	u.Email = email
	err = s.User().Create(u)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)}
