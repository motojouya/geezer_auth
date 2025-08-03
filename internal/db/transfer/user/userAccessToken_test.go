package user_test

import (
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForAccessToken(persistKey uint) shelter.User {
	var identifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return shelter.NewUser(persistKey, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestFromCoreUserAccessToken(t *testing.T) {
	var persistKey uint = 1
	var shelterUser = getUserForAccessToken(persistKey)

	var accessToken = pkgText.NewJwtToken("test.jwt.token")
	var registerDate = time.Now()
	var expireDate = time.Now().Add(24 * time.Hour)

	var shelterUserAccessToken = shelter.CreateUserAccessToken(shelterUser, accessToken, registerDate, expireDate)

	var userAccessToken = user.FromCoreUserAccessToken(shelterUserAccessToken)

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
		UserAccessToken: user.UserAccessToken{
			PersistKey:       persistKey,
			UserPersistKey:   userPersistKey,
			AccessToken:      string(accessToken),
			SourceUpdateDate: sourceUpdateDate,
			RegisterDate:     registerDate,
			ExpireDate:       expireDate,
		},
		UserIdentifier:     userIdentifier,
		UserExposeEmailId:  userEmailId,
		UserName:           userName,
		UserBotFlag:        userBotFlag,
		UserRegisteredDate: userRegisteredDate,
		UserUpdateDate:     userUpdateDate,
	}

	var shelterUserAccessToken, err = userAccessToken.ToCoreUserAccessToken()

	assert.NoError(t, err)
	assert.Equal(t, persistKey, shelterUserAccessToken.PersistKey)
	assert.Equal(t, string(accessToken), string(shelterUserAccessToken.AccessToken))
	assert.Equal(t, sourceUpdateDate, shelterUserAccessToken.SourceUpdateDate)
	assert.Equal(t, registerDate, shelterUserAccessToken.RegisterDate)
	assert.Equal(t, expireDate, shelterUserAccessToken.ExpireDate)
	assert.Equal(t, userPersistKey, shelterUserAccessToken.User.PersistKey)
	assert.Equal(t, userIdentifier, string(shelterUserAccessToken.User.Identifier))
	assert.Equal(t, userEmailId, string(shelterUserAccessToken.User.ExposeEmailId))
	assert.Equal(t, userName, string(shelterUserAccessToken.User.Name))
	assert.Equal(t, userBotFlag, shelterUserAccessToken.User.BotFlag)
	assert.Equal(t, userRegisteredDate, shelterUserAccessToken.User.RegisteredDate)
	assert.Equal(t, userUpdateDate, shelterUserAccessToken.User.UpdateDate)

	t.Logf("shelterUserAccessToken: %+v", shelterUserAccessToken)
}
