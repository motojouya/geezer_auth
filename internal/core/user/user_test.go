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

	var user = model.CreateUser(exposeId, emailId, name, botFlag, registeredDate)

	assert.Equal(t, string(exposeId), string(user.ExposeId))
	assert.Equal(t, string(emailId), string(user.ExposeEmailId))
	assert.Equal(t, string(name), string(user.Name))
	assert.Equal(t, botFlag, user.BotFlag)
	assert.Equal(t, registeredDate, user.RegisteredDate)
	assert.Equal(t, registeredDate, user.UpdateDate)

	t.Logf("user: %+v", user)
	t.Logf("user.ExposeId: %s", user.ExposeId)
	t.Logf("user.ExposeEmailId: %s", user.ExposeEmailId)
	t.Logf("user.Name: %s", user.Name)
	t.Logf("user.BotFlag: %t", user.BotFlag)
	t.Logf("user.RegisteredDate: %s", user.RegisteredDate)
	t.Logf("user.UpdateDate: %s", user.UpdateDate)
}

func TestCreateUser(t *testing.T) {
	var userId = 1
	var exposeId, _ = text.NewExposeId("TestExposeId")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var name, _ = text.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	var user = model.NewUser(userId, exposeId, emailId, name, botFlag, registeredDate)

	assert.Equal(t, userId, user.UserId)
	assert.Equal(t, string(exposeId), string(user.ExposeId))
	assert.Equal(t, string(emailId), string(user.ExposeEmailId))
	assert.Equal(t, string(name), string(user.Name))
	assert.Equal(t, botFlag, user.BotFlag)
	assert.Equal(t, registeredDate, user.RegisteredDate)
	assert.Equal(t, updateDate, user.UpdateDate)

	t.Logf("user: %+v", user)
	t.Logf("user.UserId: %d", user.UserId)
	t.Logf("user.ExposeId: %s", user.ExposeId)
	t.Logf("user.ExposeEmailId: %s", user.ExposeEmailId)
	t.Logf("user.Name: %s", user.Name)
	t.Logf("user.BotFlag: %t", user.BotFlag)
	t.Logf("user.RegisteredDate: %s", user.RegisteredDate)
	t.Logf("user.UpdateDate: %s", user.UpdateDate)
}
