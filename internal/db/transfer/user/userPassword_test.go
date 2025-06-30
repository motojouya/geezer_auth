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

func getUserForPassword(persistKey uint) core.User {
	var identifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return core.NewUser(persistKey, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestFromCoreUserPassword(t *testing.T) {
	var persistKey uint = 1
	var userValue = getUserForPassword(persistKey)

	var password = text.NewHashedPassword("TestPassword")
	var registerDate = time.Now()

	var coreUserPassword = core.CreateUserPassword(userValue, password, registerDate)

	var userPassword = user.FromCoreUserPassword(coreUserPassword)

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

	var coreUserPassword, err = userPasswordFull.ToCoreUserPassword()
	assert.NoError(t, err)

	assert.Equal(t, userPasswordFull.PersistKey, coreUserPassword.PersistKey)
	assert.Equal(t, userPasswordFull.UserPersistKey, coreUserPassword.User.PersistKey)
	assert.Equal(t, string(userPasswordFull.Password), string(coreUserPassword.Password))
	assert.Equal(t, userPasswordFull.RegisteredDate, coreUserPassword.RegisteredDate)
	assert.Equal(t, *userPasswordFull.ExpireDate, *coreUserPassword.ExpireDate)

	t.Logf("coreUserPassword: %+v", coreUserPassword)
	t.Logf("coreUserPassword.ExpireDate: %s", *coreUserPassword.ExpireDate)
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
