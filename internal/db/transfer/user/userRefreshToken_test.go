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

func getUserForRefreshToken(persistKey uint) core.User {
	var identifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return core.NewUser(persistKey, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func TestFromCoreUserRefreshToken(t *testing.T) {
	var persistKey uint = 1
	var userValue = getUserForRefreshToken(persistKey)

	var token, _ = text.NewToken("TestPassword")
	var registerDate = time.Now()
	var expireDate = registerDate.AddDate(0, 0, core.TokenValidityPeriodDays)

	var coreUserRefreshToken = core.CreateUserRefreshToken(userValue, token, registerDate)

	var userRefreshToken = user.FromCoreUserRefreshToken(coreUserRefreshToken)

	assert.Equal(t, uint(0), userRefreshToken.PersistKey)
	assert.Equal(t, persistKey, userRefreshToken.UserPersistKey)
	assert.Equal(t, string(token), string(userRefreshToken.RefreshToken))
	assert.WithinDuration(t, registerDate, userRefreshToken.RegisterDate, time.Second)
	assert.WithinDuration(t, expireDate, userRefreshToken.ExpireDate, time.Second)

	t.Logf("userRefreshToken: %+v", userRefreshToken)
}

func TestToCoreUserRefreshToken(t *testing.T) {
	var now = time.Now()
	var userRefreshTokenFull = user.UserRefreshTokenFull{
		UserRefreshToken: user.UserRefreshToken{
			PersistKey:     1,
			UserPersistKey: 2,
			RefreshToken:   "password123",
			RegisterDate:   now.Add(3 * time.Hour),
			ExpireDate:     now.Add(1 * time.Hour),
		},
		UserIdentifier:     "US-TESTES",
		UserExposeEmailId:  "test2@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(2 * time.Hour),
	}

	var coreUserRefreshToken, err = userRefreshTokenFull.ToCoreUserRefreshToken()
	assert.NoError(t, err)

	assert.Equal(t, userRefreshTokenFull.PersistKey, coreUserRefreshToken.PersistKey)
	assert.Equal(t, userRefreshTokenFull.UserPersistKey, coreUserRefreshToken.User.PersistKey)
	assert.Equal(t, userRefreshTokenFull.RefreshToken, string(coreUserRefreshToken.RefreshToken))
	assert.Equal(t, userRefreshTokenFull.RegisterDate, coreUserRefreshToken.RegisterDate)
	assert.Equal(t, userRefreshTokenFull.ExpireDate, coreUserRefreshToken.ExpireDate)

	t.Logf("coreUserRefreshToken: %+v", coreUserRefreshToken)
}

func TestToCoreUserRefreshTokenError(t *testing.T) {
	var now = time.Now()
	var userRefreshTokenFull = user.UserRefreshTokenFull{
		UserRefreshToken: user.UserRefreshToken{
			PersistKey:     1,
			UserPersistKey: 2,
			RefreshToken:   "password123",
			RegisterDate:   now.Add(3 * time.Hour),
			ExpireDate:     now.Add(1 * time.Hour),
		},
		UserIdentifier:     "invalid-identifier",
		UserExposeEmailId:  "test2@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(2 * time.Hour),
	}

	var _, err = userRefreshTokenFull.ToCoreUserRefreshToken()

	assert.Error(t, err)

	t.Logf("Error: %v", err)
}
