package user_test

import (
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForUserEmail(persistKey uint) shelter.User {
	var identifier, _ = pkgText.NewIdentifier("US-TESTES")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return shelter.NewUser(persistKey, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestFromShelterUserEmail(t *testing.T) {
	var persistKey uint = 1
	var userValue = getUserForUserEmail(persistKey)

	var email, _ = pkgText.NewEmail("test@google.com")
	var verifyToken, _ = text.NewToken("TestVerifyToken")
	var registerDate = time.Now()

	var shelterUserEmail = shelter.CreateUserEmail(userValue, email, verifyToken, registerDate)
	var userEmail = user.FromShelterUserEmail(shelterUserEmail)

	assert.Equal(t, uint(0), userEmail.PersistKey)
	assert.Equal(t, persistKey, userEmail.UserPersistKey)
	assert.Equal(t, string(email), string(userEmail.Email))
	assert.Equal(t, string(verifyToken), string(userEmail.VerifyToken))
	assert.Equal(t, registerDate, userEmail.RegisterDate)
	assert.Nil(t, userEmail.VerifyDate)
	assert.Nil(t, userEmail.ExpireDate)
}

func TestToShelterUserEmail(t *testing.T) {
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

	var shelterUserEmail, err = userEmail.ToShelterUserEmail()

	assert.NoError(t, err)
	assert.Equal(t, userEmail.PersistKey, shelterUserEmail.PersistKey)
	assert.Equal(t, userEmail.UserPersistKey, shelterUserEmail.User.PersistKey)
	assert.Equal(t, string(userEmail.Email), string(shelterUserEmail.Email))
	assert.Equal(t, string(userEmail.VerifyToken), string(shelterUserEmail.VerifyToken))
	assert.Equal(t, userEmail.RegisterDate, shelterUserEmail.RegisterDate)
	assert.Equal(t, userEmail.VerifyDate, shelterUserEmail.VerifyDate)
	assert.Equal(t, userEmail.ExpireDate, shelterUserEmail.ExpireDate)
	assert.Equal(t, userEmail.UserIdentifier, string(shelterUserEmail.User.Identifier))
	assert.Equal(t, userEmail.UserExposeEmailId, string(shelterUserEmail.User.ExposeEmailId))
	assert.Equal(t, userEmail.UserName, string(shelterUserEmail.User.Name))
	assert.Equal(t, userEmail.UserBotFlag, shelterUserEmail.User.BotFlag)
	assert.Equal(t, userEmail.UserRegisteredDate, shelterUserEmail.User.RegisteredDate)
	assert.Equal(t, userEmail.UserUpdateDate, shelterUserEmail.User.UpdateDate)

	t.Logf("shelterUserEmail: %+v", shelterUserEmail)
	t.Logf("shelterUserEmail.VerifyDate: %s", *shelterUserEmail.VerifyDate)
	t.Logf("shelterUserEmail.ExpireDate: %s", *shelterUserEmail.ExpireDate)
}

func TestToShelterUserEmailError(t *testing.T) {
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

	var _, err = userEmail.ToShelterUserEmail()

	assert.Error(t, err)

	t.Logf("Expected error: %v", err)
}
