package user_test

import (
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromCoreUser(t *testing.T) {
	var identifier, _ = text.NewIdentifier("US-TESTES")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var name, _ = text.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()

	var shelterUser = shelter.CreateUser(identifier, emailId, name, botFlag, registeredDate)
	var transferUser = user.FromCoreUser(shelterUser)

	assert.Equal(t, uint(0), transferUser.PersistKey)
	assert.Equal(t, string(identifier), string(transferUser.Identifier))
	assert.Equal(t, string(emailId), string(transferUser.ExposeEmailId))
	assert.Equal(t, string(name), string(transferUser.Name))
	assert.Equal(t, botFlag, transferUser.BotFlag)
	assert.Equal(t, registeredDate, transferUser.RegisteredDate)
	assert.Equal(t, registeredDate, transferUser.UpdateDate)

	t.Logf("transferUser: %+v", transferUser)
}

func TestToCoreUser(t *testing.T) {
	var persistKey uint = 1
	var identifier = "US-TESTES"
	var emailId = "test@gmail.com"
	var name = "TestName"
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	var user = user.User{
		PersistKey:     persistKey,
		Identifier:     identifier,
		ExposeEmailId:  emailId,
		Name:           name,
		BotFlag:        botFlag,
		RegisteredDate: registeredDate,
		UpdateDate:     updateDate,
	}

	var shelterUser, err = user.ToCoreUser()

	assert.NoError(t, err)
	assert.Equal(t, persistKey, shelterUser.PersistKey)
	assert.Equal(t, identifier, string(shelterUser.Identifier))
	assert.Equal(t, emailId, string(shelterUser.ExposeEmailId))
	assert.Equal(t, name, string(shelterUser.Name))
	assert.Equal(t, botFlag, shelterUser.BotFlag)
	assert.Equal(t, registeredDate, shelterUser.RegisteredDate)
	assert.Equal(t, updateDate, shelterUser.UpdateDate)

	t.Logf("shelterUser: %+v", shelterUser)
}

func TestToCoreUserErr(t *testing.T) {
	var persistKey uint = 1
	var identifier = "invalid-identifier"
	var emailId = "test@gmail.com"
	var name = "TestName"
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	var user = user.User{
		PersistKey:     persistKey,
		Identifier:     identifier,
		ExposeEmailId:  emailId,
		Name:           name,
		BotFlag:        botFlag,
		RegisteredDate: registeredDate,
		UpdateDate:     updateDate,
	}

	var _, err = user.ToCoreUser()
	assert.Error(t, err)

	t.Logf("error: %v", err)
}
