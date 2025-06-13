package user_test

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForUserEmail(identifier pkgText.Identifier) user.User {
	var userId uint = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return user.NewUser(userId, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestCreateUserEmail(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var userValue = getUserForUserEmail(userIdentifier)

	var email, _ = pkgText.NewEmail("test@google.com")
	var verifyToken, _ = text.NewToken("TestVerifyToken")
	var registerDate = time.Now()

	var userEmail = user.CreateUserEmail(userValue, email, verifyToken, registerDate)

	assert.Equal(t, string(userIdentifier), string(userEmail.User.Identifier))
	assert.Equal(t, string(email), string(userEmail.Email))
	assert.Equal(t, string(verifyToken), string(userEmail.VerifyToken))
	assert.Equal(t, registerDate, userEmail.RegisterDate)
	assert.Nil(t, *userEmail.VerifyDate)
	assert.Nil(t, *userEmail.ExpireDate)

	t.Logf("userEmail: %+v", userEmail)
	t.Logf("userEmail.User.Identifier: %s", userEmail.User.Identifier)
	t.Logf("userEmail.Email: %s", userEmail.Email)
	t.Logf("userEmail.VerifyToken: %s", userEmail.VerifyToken)
	t.Logf("userEmail.RegisteredDate: %s", userEmail.RegisterDate)
	t.Logf("userEmail.VerifyDate: %s", *userEmail.VerifyDate)
	t.Logf("userEmail.ExpireDate: %s", *userEmail.ExpireDate)
}

func TestNewUserEmail(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var userValue = getUserForUserEmail(userIdentifier)

	var email, _ = pkgText.NewEmail("test@google.com")
	var verifyToken, _ = text.NewToken("TestVerifyToken")
	var registerDate = time.Now()
	var verifyDate = registerDate.Add(24 * time.Hour)
	var expireDate = registerDate.Add(48 * time.Hour)

	var userEmail = user.NewUserEmail(1, userValue, email, verifyToken, registerDate, verifyDate, expireDate)

	assert.Equal(t, 1, userEmail.PersistKey)
	assert.Equal(t, string(userIdentifier), string(userEmail.User.Identifier))
	assert.Equal(t, string(email), string(userEmail.Email))
	assert.Equal(t, string(verifyToken), string(userEmail.VerifyToken))
	assert.Equal(t, registerDate, userEmail.RegisterDate)
	assert.Equal(t, verifyDate, *userEmail.VerifyDate)
	assert.Equal(t, expireDate, *userEmail.ExpireDate)

	t.Logf("userEmail: %+v", userEmail)
	t.Logf("userEmail.User.Identifier: %s", userEmail.User.Identifier)
	t.Logf("userEmail.Email: %s", userEmail.Email)
	t.Logf("userEmail.VerifyToken: %s", userEmail.VerifyToken)
	t.Logf("userEmail.RegisteredDate: %s", userEmail.RegisterDate)
	t.Logf("userEmail.VerifyDate: %s", *userEmail.VerifyDate)
	t.Logf("userEmail.ExpireDate: %s", *userEmail.ExpireDate)
}
