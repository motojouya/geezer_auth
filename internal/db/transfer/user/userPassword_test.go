package user_test

import (
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForPassword(persistKey uint) shelter.User {
	var identifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return shelter.NewUser(persistKey, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestFromCoreUserPassword(t *testing.T) {
	var persistKey uint = 1
	var userValue = getUserForPassword(persistKey)

	var password = text.NewHashedPassword("TestPassword")
	var registerDate = time.Now()

	var shelterUserPassword = shelter.CreateUserPassword(userValue, password, registerDate)

	var userPassword = user.FromCoreUserPassword(shelterUserPassword)

	assert.Equal(t, uint(0), userPassword.PersistKey)
	assert.Equal(t, persistKey, userPassword.UserPersistKey)
	assert.Equal(t, string(password), string(userPassword.Password))
	assert.Equal(t, registerDate, userPassword.RegisteredDate)
	assert.Nil(t, userPassword.ExpireDate)

	t.Logf("userPassword: %+v", userPassword)
}

func TestToCoreUserPassword(t *testing.T) {
	var now = time.Now()
	var expireDate = now.Add(1 * time.Hour)
	var userPasswordFull = user.UserPasswordFull{
		UserPassword: user.UserPassword{
			PersistKey:     1,
			UserPersistKey: 2,
			Password:       "password123",
			RegisteredDate: now.Add(3 * time.Hour),
			ExpireDate:     &expireDate,
		},
		UserIdentifier:     "US-TESTES",
		UserExposeEmailId:  "test2@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(2 * time.Hour),
	}

	var shelterUserPassword, err = userPasswordFull.ToCoreUserPassword()
	assert.NoError(t, err)

	assert.Equal(t, userPasswordFull.PersistKey, shelterUserPassword.PersistKey)
	assert.Equal(t, userPasswordFull.UserPersistKey, shelterUserPassword.User.PersistKey)
	assert.Equal(t, string(userPasswordFull.Password), string(shelterUserPassword.Password))
	assert.Equal(t, userPasswordFull.RegisteredDate, shelterUserPassword.RegisteredDate)
	assert.Equal(t, *userPasswordFull.ExpireDate, *shelterUserPassword.ExpireDate)

	t.Logf("shelterUserPassword: %+v", shelterUserPassword)
	t.Logf("shelterUserPassword.ExpireDate: %s", *shelterUserPassword.ExpireDate)
}

func TestToCoreUserPasswordError(t *testing.T) {
	var now = time.Now()
	var expireDate = now.Add(1 * time.Hour)
	var userPasswordFull = user.UserPasswordFull{
		UserPassword: user.UserPassword{
			PersistKey:     1,
			UserPersistKey: 2,
			Password:       "password123",
			RegisteredDate: now.Add(3 * time.Hour),
			ExpireDate:     &expireDate,
		},
		UserIdentifier:     "invalid-identifier",
		UserExposeEmailId:  "test2@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(2 * time.Hour),
	}

	var _, err = userPasswordFull.ToCoreUserPassword()

	assert.Error(t, err, "Expected error due to invalid user data")

	t.Logf("Error: %v", err)
}
