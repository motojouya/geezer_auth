package user_test

import (
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
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

	return core.NewUser(persistKey, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestFromCoreUserAccessToken(t *testing.T) {
	var persistKey uint = 1
	var coreUser = getUserForAccessToken(persistKey)

	var accessToken = pkgText.NewJwtToken("test.jwt.token")
	var registerDate = time.Now()
	var expireDate = time.Now().Add(24 * time.Hour)

	var coreUserAccessToken = core.CreateUserAccessToken(coreUser, accessToken, registerDate, expireDate)

	var userAccessToken = user.FromCoreUserAccessToken(coreUserAccessToken)

	assert.Equal(t, uint(0), userAccessToken.PersistKey)
	assert.Equal(t, persistKey, userAccessToken.UserPersistKey)
	assert.Equal(t, string(accessToken), string(userAccessToken.AccessToken))
	assert.WithinDuration(t, registerDate, userAccessToken.SourceUpdateDate, time.Second)
	assert.WithinDuration(t, registerDate, userAccessToken.RegisterDate, time.Second)
	assert.WithinDuration(t, expireDate, userAccessToken.ExpireDate, time.Second)

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
		PersistKey:         persistKey,
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

	var coreUserAccessToken, err = userAccessToken.ToCoreUserAccessToken()

	assert.NoError(t, err)
	assert.Equal(t, persistKey, coreUserAccessToken.PersistKey)
	assert.Equal(t, string(accessToken), string(coreUserAccessToken.AccessToken))
	assert.Equal(t, sourceUpdateDate, coreUserAccessToken.SourceUpdateDate)
	assert.Equal(t, registerDate, coreUserAccessToken.RegisterDate)
	assert.Equal(t, expireDate, coreUserAccessToken.ExpireDate)
	assert.Equal(t, userPersistKey, coreUserAccessToken.User.PersistKey)
	assert.Equal(t, userIdentifier, string(coreUserAccessToken.User.Identifier))
	assert.Equal(t, userEmailId, string(coreUserAccessToken.User.ExposeEmailId))
	assert.Equal(t, userName, string(coreUserAccessToken.User.Name))
	assert.Equal(t, userBotFlag, coreUserAccessToken.User.BotFlag)
	assert.Equal(t, userRegisteredDate, coreUserAccessToken.User.RegisteredDate)
	assert.Equal(t, userUpdateDate, coreUserAccessToken.User.UpdateDate)

	t.Logf("coreUserAccessToken: %+v", coreUserAccessToken)
}
