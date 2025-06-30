package user_test

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForUserEmail(persistKey uint) core.User {
	var identifier, _ = pkgText.NewIdentifier("US-TESTES")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return core.NewUser(persistKey, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestFromCoreUserEmail(t *testing.T) {
	var persistKey uint = 1
	var userValue = getUserForUserEmail(persistKey)

	var email, _ = pkgText.NewEmail("test@google.com")
	var verifyToken, _ = text.NewToken("TestVerifyToken")
	var registerDate = time.Now()

	var coreUserEmail = core.CreateUserEmail(userValue, email, verifyToken, registerDate)
	var userEmail = user.FromCoreUserEmail(coreUserEmail)

	assert.Equal(t, uint(0), userEmail.PersistKey)
	assert.Equal(t, persistKey, userEmail.UserPersistKey)
	assert.Equal(t, string(email), string(userEmail.Email))
	assert.Equal(t, string(verifyToken), string(userEmail.VerifyToken))
	assert.Equal(t, registerDate, userEmail.RegisterDate)
	assert.Nil(t, userEmail.VerifyDate)
	assert.Nil(t, userEmail.ExpireDate)
}

func TestToCoreUserEmail(t *testing.T) {
	var now = time.Now()
	var verifiyDate = now.Add(1 * time.Hour)
	var expireDate = now.Add(1 * time.Hour)
	var userEmail = &user.UserEmailFull{
		UserEmail: user.UserEmail{
			PersistKey:     1,
			UserPersistKey: 2,
			Email:          "test01@example.com",
			VerifyToken:    "TestVerifyToken",
			RegisterDate:   now,
			VerifyDate:     &verifiyDate,
			ExpireDate:     &expireDate,
		},
		UserIdentifier:     "US-TESTES",
		UserExposeEmailId:  "test02@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now.Add(2 * time.Hour),
		UserUpdateDate:     now.Add(3 * time.Hour),
	}

	var coreUserEmail, err = userEmail.ToCoreUserEmail()

	assert.NoError(t, err)
	assert.Equal(t, userEmail.PersistKey, coreUserEmail.PersistKey)
	assert.Equal(t, userEmail.UserPersistKey, coreUserEmail.User.PersistKey)
	assert.Equal(t, string(userEmail.Email), string(coreUserEmail.Email))
	assert.Equal(t, string(userEmail.VerifyToken), string(coreUserEmail.VerifyToken))
	assert.Equal(t, userEmail.RegisterDate, coreUserEmail.RegisterDate)
	assert.Equal(t, userEmail.VerifyDate, coreUserEmail.VerifyDate)
	assert.Equal(t, userEmail.ExpireDate, coreUserEmail.ExpireDate)
	assert.Equal(t, userEmail.UserIdentifier, string(coreUserEmail.User.Identifier))
	assert.Equal(t, userEmail.UserExposeEmailId, string(coreUserEmail.User.ExposeEmailId))
	assert.Equal(t, userEmail.UserName, string(coreUserEmail.User.Name))
	assert.Equal(t, userEmail.UserBotFlag, coreUserEmail.User.BotFlag)
	assert.Equal(t, userEmail.UserRegisteredDate, coreUserEmail.User.RegisteredDate)
	assert.Equal(t, userEmail.UserUpdateDate, coreUserEmail.User.UpdateDate)

	t.Logf("coreUserEmail: %+v", coreUserEmail)
	t.Logf("coreUserEmail.VerifyDate: %s", *coreUserEmail.VerifyDate)
	t.Logf("coreUserEmail.ExpireDate: %s", *coreUserEmail.ExpireDate)
}

func TestToCoreUserEmailError(t *testing.T) {
	var now = time.Now()
	var verifiyDate = now.Add(1 * time.Hour)
	var expireDate = now.Add(1 * time.Hour)
	var userEmail = &user.UserEmailFull{
		UserEmail: user.UserEmail{
			PersistKey:     1,
			UserPersistKey: 2,
			Email:          "test01@example.com",
			VerifyToken:    "TestVerifyToken",
			RegisterDate:   now,
			VerifyDate:     &verifiyDate,
			ExpireDate:     &expireDate,
		},
		UserIdentifier:     "invalid-identifier",
		UserExposeEmailId:  "test02@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now.Add(2 * time.Hour),
		UserUpdateDate:     now.Add(3 * time.Hour),
	}

	var _, err = userEmail.ToCoreUserEmail()

	assert.Error(t, err)

	t.Logf("Expected error: %v", err)
}
