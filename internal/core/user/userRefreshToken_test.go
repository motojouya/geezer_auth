package user_test

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForRefreshToken(identifier pkgText.Identifier) user.User {
	var userId uint = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return user.NewUser(userId, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestCreateUserRefreshToken(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var userValue = getUserForRefreshToken(userIdentifier)

	var token, _ = text.NewToken("TestPassword")
	var registerDate = time.Now()
	var expireDate = registerDate.Add(50 * time.Hour * 24) // 50 days

	var userRefreshToken = user.CreateUserRefreshToken(userValue, token, registerDate)

	assert.Equal(t, string(userIdentifier), string(userRefreshToken.User.Identifier))
	assert.Equal(t, string(token), string(userRefreshToken.RefreshToken))
	assert.WithinDuration(t, registerDate, userRefreshToken.RegisterDate, time.Second)
	assert.WithinDuration(t, expireDate, userRefreshToken.ExpireDate, time.Second)

	t.Logf("userRefreshToken: %+v", userRefreshToken)
	t.Logf("userRefreshToken.User.Identifier: %s", userRefreshToken.User.Identifier)
	t.Logf("userRefreshToken.RefreshToken: %s", userRefreshToken.RefreshToken)
	t.Logf("userRefreshToken.RegisterDate: %s", userRefreshToken.RegisterDate)
	t.Logf("userRefreshToken.ExpireDate: %s", userRefreshToken.ExpireDate)
}

func TestNewUserRefreshToken(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var userValue = getUserForRefreshToken(userIdentifier)

	var token, _ = text.NewToken("TestPassword")
	var registerDate = time.Now()
	var expireDate = registerDate.Add(30 * time.Hour * 24) // 30 days

	var userRefreshToken = user.NewUserRefreshToken(1, userValue, token, registerDate, expireDate)

	assert.Equal(t, uint(1), userRefreshToken.PersistKey)
	assert.Equal(t, string(userIdentifier), string(userRefreshToken.User.Identifier))
	assert.Equal(t, string(token), string(userRefreshToken.RefreshToken))
	assert.Equal(t, registerDate, userRefreshToken.RegisterDate)
	assert.Equal(t, expireDate, userRefreshToken.ExpireDate)

	t.Logf("userRefreshToken: %+v", userRefreshToken)
	t.Logf("userRefreshToken.User.Identifier: %s", userRefreshToken.User.Identifier)
	t.Logf("userRefreshToken.RefreshToken: %s", userRefreshToken.RefreshToken)
	t.Logf("userRefreshToken.RegisterDate: %s", userRefreshToken.RegisterDate)
	t.Logf("userRefreshToken.ExpireDate: %s", userRefreshToken.ExpireDate)
}
