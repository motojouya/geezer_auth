package user_test

import (
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUser(identifier text.Identifier, updateDate time.Time) user.User {
	var userId = 1
	var emailId, _ = text.NewEmail("test@gmail.com")
	var name, _ = text.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()

	return user.NewUser(userId, identifier, emailId, name, botFlag, registeredDate, updateDate)
}

func TestCreateUserAccessToken(t *testing.T) {
	var identifier, _ = text.NewIdentifier("TestIdentifier")
	var updateDate = time.Now()
	var user = getUser(identifier, updateDate)

	var accessToken, _ = text.NewJwtToken("test.jwt.token")
	var registeredDate = time.Now()
	var expireDate = time.Now().Add(24 * time.Hour)

	var userAccessToken = user.CreateUserAccessToken(user, accessToken, registeredDate, &expireDate)

	assert.Equal(t, string(identifier), string(userAccessToken.User.Identifier))
	assert.Equal(t, string(accessToken), userAccessToken.User.BotFlag)
	assert.Equal(t, updateDate, userAccessToken.sourceUpdateDate)
	assert.Equal(t, registeredDate, userAccessToken.RegisteredDate)
	assert.Equal(t, expireDate, userAccessToken.ExpireDate)

	t.Logf("userAccessToken: %+v", userAccessToken)
	t.Logf("userAccessToken.User.Identifier: %s", userAccessToken.User.Identifier)
	t.Logf("userAccessToken.AccessToken: %s", userAccessToken.AccessToken)
	t.Logf("userAccessToken.SourceUpdateDate: %s", userAccessToken.SourceUpdateDate)
	t.Logf("userAccessToken.RegisteredDate: %s", userAccessToken.RegisteredDate)
	t.Logf("userAccessToken.ExpireDate: %s", userAccessToken.ExpireDate)
}

func TestNewUserAccessToken(t *testing.T) {
	var identifier, _ = text.NewIdentifier("TestIdentifier")
	var updateDate = time.Now()
	var user = getUser(identifier, updateDate)

	var accessToken, _ = text.NewJwtToken("test.jwt.token")
	var sourceUpdateDate = time.Now()
	var registeredDate = time.Now()
	var expireDate = time.Now().Add(24 * time.Hour)

	var userAccessToken = user.NewUserAccessToken(1, user, accessToken, sourceUpdateDate, registeredDate, &expireDate)

	assert.Equal(t, 1, userAccessToken.UserAccessTokenId)
	assert.Equal(t, string(identifier), string(userAccessToken.User.Identifier))
	assert.Equal(t, string(accessToken), string(userAccessToken.AccessToken))
	assert.Equal(t, sourceUpdateDate, userAccessToken.SourceUpdateDate)
	assert.Equal(t, registeredDate, userAccessToken.RegisteredDate)
	assert.Equal(t, expireDate, *userAccessToken.ExpireDate)

	t.Logf("userAccessToken: %+v", userAccessToken)
	t.Logf("userAccessToken.UserAccessTokenId: %d", userAccessToken.UserAccessTokenId)
	t.Logf("userAccessToken.User.Identifier: %s", userAccessToken.User.Identifier)
	t.Logf("userAccessToken.AccessToken: %s", userAccessToken.AccessToken)
	t.Logf("userAccessToken.SourceUpdateDate: %s", userAccessToken.SourceUpdateDate)
	t.Logf("userAccessToken.RegisteredDate: %s", userAccessToken.RegisteredDate)
	t.Logf("userAccessToken.ExpireDate: %s", *userAccessToken.ExpireDate)
}
