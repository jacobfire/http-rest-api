package teststore_test

import (
	"fmt"
	"github.com/jacobfire/http-rest-api/app/model"
	"github.com/jacobfire/http-rest-api/app/store"
	"github.com/jacobfire/http-rest-api/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()
	email := "test@test.com"
	u := model.TestUser(t)
	u, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	email = "test@test.com"
	u = model.TestUser(t)
	u.Email = email
	err = s.User().Create(u)

	u, err = s.User().FindByEmail(email)
	fmt.Println(err)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
