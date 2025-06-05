package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateUserExposeId(t *testing.T) {
	var exposeId, err = company.CreateUserExposeId("TESTES")

	assert.Nil(t, err)
	assert.NotEmpty(t, exposeId)
	assert.Equal(t, "US-TESTES", string(exposeId))

	t.Logf("exposeId: %s", exposeId)
}

func TestCreateUser(t *testing.T) {
	var exposeId, _ = text.NewExposeId("TestExposeId")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var name, _ = text.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()

	var userObj = user.CreateUser(exposeId, emailId, name, botFlag, registeredDate)

	assert.Equal(t, string(exposeId), string(userObj.ExposeId))
	assert.Equal(t, string(emailId), string(userObj.ExposeEmailId))
	assert.Equal(t, string(name), string(userObj.Name))
	assert.Equal(t, botFlag, userObj.BotFlag)
	assert.Equal(t, registeredDate, userObj.RegisteredDate)
	assert.Equal(t, registeredDate, userObj.UpdateDate)

	t.Logf("user: %+v", userObj)
	t.Logf("user.ExposeId: %s", userObj.ExposeId)
	t.Logf("user.ExposeEmailId: %s", userObj.ExposeEmailId)
	t.Logf("user.Name: %s", userObj.Name)
	t.Logf("user.BotFlag: %t", userObj.BotFlag)
	t.Logf("user.RegisteredDate: %s", userObj.RegisteredDate)
	t.Logf("user.UpdateDate: %s", userObj.UpdateDate)
}

func TestNewUser(t *testing.T) {
	var userId = 1
	var exposeId, _ = text.NewExposeId("TestExposeId")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var name, _ = text.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	var userObj = user.NewUser(userId, exposeId, emailId, name, botFlag, registeredDate, updateDate)

	assert.Equal(t, userId, userObj.UserId)
	assert.Equal(t, string(exposeId), string(userObj.ExposeId))
	assert.Equal(t, string(emailId), string(userObj.ExposeEmailId))
	assert.Equal(t, string(name), string(userObj.Name))
	assert.Equal(t, botFlag, userObj.BotFlag)
	assert.Equal(t, registeredDate, userObj.RegisteredDate)
	assert.Equal(t, updateDate, userObj.UpdateDate)

	t.Logf("user: %+v", userObj)
	t.Logf("user.UserId: %d", userObj.UserId)
	t.Logf("user.ExposeId: %s", userObj.ExposeId)
	t.Logf("user.ExposeEmailId: %s", userObj.ExposeEmailId)
	t.Logf("user.Name: %s", userObj.Name)
	t.Logf("user.BotFlag: %t", userObj.BotFlag)
	t.Logf("user.RegisteredDate: %s", userObj.RegisteredDate)
	t.Logf("user.UpdateDate: %s", userObj.UpdateDate)
}
