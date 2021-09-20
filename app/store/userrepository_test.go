package store_test

import (
	"fmt"
	"github.com/jacobfire/http-rest-api/app/model"
	"github.com/jacobfire/http-rest-api/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s, _ := store.TestStore(t, databaseURL)
	//defer teardown("users")

	u, err := s.User().Create(&model.User{
		Email: "jack@gmail.com",
	})

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("users")

	email := "test@test.com"

	u, err := s.User().FindByEmail(email)

	assert.NoError(t, err)
	assert.Nil(t, u)

	s.User().Create(&model.User {
		Email: "test@test.com",
	})

	u, err = s.User().FindByEmail(email)
	fmt.Println(err)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}