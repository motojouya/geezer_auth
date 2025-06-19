package user_test

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForPassword(identifier pkgText.Identifier) user.User {
	var userId uint = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return user.NewUser(userId, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestCreateUserPassword(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var userValue = getUserForPassword(userIdentifier)

	var password = text.NewHashedPassword("TestPassword")
	var registerDate = time.Now()

	var userPassword = user.CreateUserPassword(userValue, password, registerDate)

	assert.Equal(t, string(userIdentifier), string(userPassword.User.Identifier))
	assert.Equal(t, string(password), string(userPassword.Password))
	assert.Equal(t, registerDate, userPassword.RegisteredDate)
	assert.Nil(t, userPassword.ExpireDate)

	t.Logf("userPassword: %+v", userPassword)
	t.Logf("userPassword.User.Identifier: %s", userPassword.User.Identifier)
	t.Logf("userPassword.Password: %s", userPassword.Password)
	t.Logf("userPassword.RegisteredDate: %s", userPassword.RegisteredDate)
}

func TestNewUserPassword(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var userValue = getUserForPassword(userIdentifier)

	var password = text.NewHashedPassword("TestPassword")
	var registerDate = time.Now()
	var expireDate = registerDate.Add(24 * time.Hour)

	var userPassword = user.NewUserPassword(1, userValue, password, registerDate, expireDate)

	assert.Equal(t, uint(1), userPassword.PersistKey)
	assert.Equal(t, string(userIdentifier), string(userPassword.User.Identifier))
	assert.Equal(t, string(password), string(userPassword.Password))
	assert.Equal(t, registerDate, userPassword.RegisteredDate)
	assert.Equal(t, expireDate, *userPassword.ExpireDate)

	t.Logf("userPassword: %+v", userPassword)
	t.Logf("userPassword.User.Identifier: %s", userPassword.User.Identifier)
	t.Logf("userPassword.Password: %s", userPassword.Password)
	t.Logf("userPassword.RegisteredDate: %s", userPassword.RegisteredDate)
	t.Logf("userPassword.ExpireDate: %s", *userPassword.ExpireDate)
}
