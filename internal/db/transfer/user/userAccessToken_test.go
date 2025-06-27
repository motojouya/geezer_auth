package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForAccessToken(persistKey uint) core.User {
	var identifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return user.NewUser(persistKey, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestFromCoreUserAccessToken(t *testing.T) {
	var persistKey uint = 1
	var coreUser = getUserForAccessToken(persistKey)

	var accessToken = pkgText.NewJwtToken("test.jwt.token")
	var registerDate = time.Now()
	var expireDate = time.Now().Add(24 * time.Hour)

	var coreUserAccessToken = core.CreateUserAccessToken(userValue, accessToken, registerDate, expireDate)

	var userAccessToken = user.FromCoreUserAccessToken(coreUserAccessToken)

	assert.Equal(t, uint(0), userAccessToken.PersistKey)
	assert.Equal(t, persistKey, userAccessToken.User.PersistKey)
	assert.Equal(t, string(accessToken), string(userAccessToken.AccessToken))
	assert.Equal(t, registerDate, userAccessToken.SourceUpdateDate)
	assert.Equal(t, registerDate, userAccessToken.RegisterDate)
	assert.Equal(t, expireDate, userAccessToken.ExpireDate)

	t.Logf("userAccessToken: %+v", userAccessToken)
}

func TestToCoreUserAccessToken(t *testing.T) {
	var userPersistKey uint = 1
	var userIdentifier = "US-TESTES"
	var userEmailId = "test@gmail.com"
	var userName = "TestName"
	var userBotFlag = false
	var userRegisteredDate = time.Now()
	var userUpdateDate = time.Now()

	var persistKey uint = 1
	var accessToken = pkgText.NewJwtToken("test.jwt.token")
	var sourceUpdateDate = time.Now()
	var registerDate = time.Now()
	var expireDate = time.Now().Add(24 * time.Hour)

	var userAccessToken = user.UserAccessTokenFull{
		PersistKey:     persistKey,
		UserPersistKey:     userPersistKey,
		UserIdentifier:     userIdentifier,
		UserExposeEmailId:  userEmailId,
		UserName:           userName,
		UserBotFlag:        userBotFlag,
		UserRegisteredDate: userRegisteredDate,
		UserUpdateDate:     userUpdateDate,
		AccessToken:        string(accessToken),
		SourceUpdateDate:   sourceUpdateDate,
		RegisterDate:       registerDate,
		ExpireDate:         expireDate,
	}

	var userAccessToken, err = userAccessToken.ToCoreUserAccessToken()

	var coreUserAccessToken, err = userAccessToken.ToCoreUserAccessToken()

	assert.NoError(t, err)
	assert.Equal(t, persistKey, coreUserAccessToken.User.PersistKey)
	assert.Equal(t, string(accessToken), string(coreUserAccessToken.AccessToken))
	assert.Equal(t, sourceUpdateDate, coreUserAccessToken.SourceUpdateDate)
	assert.Equal(t, registerDate, coreUserAccessToken.RegisterDate)
	assert.Equal(t, expireDate, coreUserAccessToken.ExpireDate)

	t.Logf("coreUserAccessToken: %+v", coreUserAccessToken)
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
