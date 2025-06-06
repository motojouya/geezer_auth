package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateUserIdentifier(t *testing.T) {
	var identifier, err = company.CreateUserIdentifier("TESTES")

	assert.Nil(t, err)
	assert.NotEmpty(t, identifier)
	assert.Equal(t, "US-TESTES", string(identifier))

	t.Logf("identifier: %s", identifier)
}

func TestCreateUser(t *testing.T) {
	var identifier, _ = text.NewIdentifier("TestIdentifier")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var name, _ = text.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()

	var userObj = user.CreateUser(identifier, emailId, name, botFlag, registeredDate)

	assert.Equal(t, string(identifier), string(userObj.Identifier))
	assert.Equal(t, string(emailId), string(userObj.ExposeEmailId))
	assert.Equal(t, string(name), string(userObj.Name))
	assert.Equal(t, botFlag, userObj.BotFlag)
	assert.Equal(t, registeredDate, userObj.RegisteredDate)
	assert.Equal(t, registeredDate, userObj.UpdateDate)

	t.Logf("user: %+v", userObj)
	t.Logf("user.Identifier: %s", userObj.Identifier)
	t.Logf("user.ExposeEmailId: %s", userObj.ExposeEmailId)
	t.Logf("user.Name: %s", userObj.Name)
	t.Logf("user.BotFlag: %t", userObj.BotFlag)
	t.Logf("user.RegisteredDate: %s", userObj.RegisteredDate)
	t.Logf("user.UpdateDate: %s", userObj.UpdateDate)
}

func TestNewUser(t *testing.T) {
	var userId = 1
	var identifier, _ = text.NewIdentifier("TestIdentifier")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var name, _ = text.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	var userObj = user.NewUser(userId, identifier, emailId, name, botFlag, registeredDate, updateDate)

	assert.Equal(t, userId, userObj.UserId)
	assert.Equal(t, string(identifier), string(userObj.Identifier))
	assert.Equal(t, string(emailId), string(userObj.ExposeEmailId))
	assert.Equal(t, string(name), string(userObj.Name))
	assert.Equal(t, botFlag, userObj.BotFlag)
	assert.Equal(t, registeredDate, userObj.RegisteredDate)
	assert.Equal(t, updateDate, userObj.UpdateDate)

	t.Logf("user: %+v", userObj)
	t.Logf("user.UserId: %d", userObj.UserId)
	t.Logf("user.Identifier: %s", userObj.Identifier)
	t.Logf("user.ExposeEmailId: %s", userObj.ExposeEmailId)
	t.Logf("user.Name: %s", userObj.Name)
	t.Logf("user.BotFlag: %t", userObj.BotFlag)
	t.Logf("user.RegisteredDate: %s", userObj.RegisteredDate)
	t.Logf("user.UpdateDate: %s", userObj.UpdateDate)
}
