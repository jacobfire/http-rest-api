package model_test

import (
	"github.com/jacobfire/http-rest-api/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}

// Hi, Im Yakov and have been working with IT technologies. I am a web magician and helped a lot of companies to start their offline businesses to start it in the web.
//It's easy and profitable. And also can bring your business to new level
func TestUser_Validate(t *testing.T) {
	//u := model.TestUser(t)
	//assert.NoError(t, u.Validate())

	testCases := []struct {
		name string
		u func() *model.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "invalid email",
			u: func() *model.User {
				user := model.TestUser(t)
				user.Email = "invalid"

				return user
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *model.User {
				user := model.TestUser(t)
				user.Password = ""

				return user
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *model.User {
				user := model.TestUser(t)
				user.Password = "pas"

				return user
			},
			isValid: false,
		},
		{
			name: "with encrypted password",
			u: func() *model.User {
				user := model.TestUser(t)
				user.Password = ""
				user.EncryptedPassword = "encryptedPassword"

				return user
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}