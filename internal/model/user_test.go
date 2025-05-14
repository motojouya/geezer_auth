package model_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	// "time"
	"github.com/motojouya/geezer_auth/internal/model"
)

func TestCreateUser(t *testing.T) {
	var exposeId = "TestExposeId"
	var emailId = "test@gmail.com"
	var name = "TestName"
	var botFlag = false

	var user = model.CreateUser(exposeId, emailId, name, botFlag)

	assert.Equal(t, exposeId, user.ExposeId)
	assert.Equal(t, emailId, user.ExposeEmailId)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, botFlag, user.BotFlag)

	t.Logf("user: %+v", user)
	t.Logf("user.ExposeId: %s", user.ExposeId)
	t.Logf("user.ExposeEmailId: %s", user.ExposeEmailId)
	t.Logf("user.Name: %s", user.Name)
	t.Logf("user.BotFlag: %t", user.BotFlag)
}

func TestNewUser(t *testing.T) {
	var exposeId = "TestExposeId"
	var emailId = "test@gmail.com"
	var name = "TestName"
	var botFlag = false

	var user = model.CreateUser(exposeId, emailId, name, botFlag)

	assert.Equal(t, exposeId, user.ExposeId)
	assert.Equal(t, emailId, user.ExposeEmailId)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, botFlag, user.BotFlag)

	t.Logf("user: %+v", user)
	t.Logf("user.ExposeId: %s", user.ExposeId)
	t.Logf("user.ExposeEmailId: %s", user.ExposeEmailId)
	t.Logf("user.Name: %s", user.Name)
	t.Logf("user.BotFlag: %t", user.BotFlag)
}
