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
	u := model.TestUser(t)
	assert.NoError(t, u.Validate())
}