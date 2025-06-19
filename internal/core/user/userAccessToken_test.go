package user_test

import (
	"github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForAccessToken(identifier pkgText.Identifier) user.User {
	var userId uint = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return user.NewUser(userId, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestCreateUserAccessToken(t *testing.T) {
	var identifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var userValue = getUserForAccessToken(identifier)

	var accessToken = pkgText.NewJwtToken("test.jwt.token")
	var registerDate = time.Now()
	var expireDate = time.Now().Add(24 * time.Hour)

	var userAccessToken = user.CreateUserAccessToken(userValue, accessToken, registerDate, expireDate)

	assert.Equal(t, string(identifier), string(userAccessToken.User.Identifier))
	assert.Equal(t, string(accessToken), string(userAccessToken.AccessToken))
	assert.Equal(t, registerDate, userAccessToken.RegisterDate)
	assert.Equal(t, expireDate, userAccessToken.ExpireDate)

	t.Logf("userAccessToken: %+v", userAccessToken)
	t.Logf("userAccessToken.User.Identifier: %s", userAccessToken.User.Identifier)
	t.Logf("userAccessToken.AccessToken: %s", userAccessToken.AccessToken)
	t.Logf("userAccessToken.SourceUpdateDate: %s", userAccessToken.SourceUpdateDate)
	t.Logf("userAccessToken.RegisterDate: %s", userAccessToken.RegisterDate)
	t.Logf("userAccessToken.ExpireDate: %s", userAccessToken.ExpireDate)
}

func TestNewUserAccessToken(t *testing.T) {
	var identifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var userValue = getUserForAccessToken(identifier)

	var accessToken = pkgText.NewJwtToken("test.jwt.token")
	var sourceUpdateDate = time.Now()
	var registerDate = time.Now()
	var expireDate = time.Now().Add(24 * time.Hour)

	var userAccessToken = user.NewUserAccessToken(1, userValue, accessToken, sourceUpdateDate, registerDate, expireDate)

	assert.Equal(t, uint(1), userAccessToken.PersistKey)
	assert.Equal(t, string(identifier), string(userAccessToken.User.Identifier))
	assert.Equal(t, string(accessToken), string(userAccessToken.AccessToken))
	assert.Equal(t, sourceUpdateDate, userAccessToken.SourceUpdateDate)
	assert.Equal(t, registerDate, userAccessToken.RegisterDate)
	assert.Equal(t, expireDate, userAccessToken.ExpireDate)

	t.Logf("userAccessToken: %+v", userAccessToken)
	t.Logf("userAccessToken.PersistKey: %d", userAccessToken.PersistKey)
	t.Logf("userAccessToken.User.Identifier: %s", userAccessToken.User.Identifier)
	t.Logf("userAccessToken.AccessToken: %s", userAccessToken.AccessToken)
	t.Logf("userAccessToken.SourceUpdateDate: %s", userAccessToken.SourceUpdateDate)
	t.Logf("userAccessToken.RegisterDate: %s", userAccessToken.RegisterDate)
	t.Logf("userAccessToken.ExpireDate: %s", userAccessToken.ExpireDate)
}
