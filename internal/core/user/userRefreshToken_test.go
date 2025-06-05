package user_test

import (
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUser(exposeId pkgText.ExposeId) user.User {
	var userId = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return user.NewUser(userId, exposeId, emailId, name, botFlag, registeredDate, updateDate)
}

func TestCreateUserRefreshToken(t *testing.T) {
	var userExposeId, _ = pkgText.NewExposeId("TestExposeId")
	var user = getUser(userExposeId)

	var token, _ = text.NewToken("TestPassword")
	var registerDate = time.Now()
	var expireDate = registerDate.Add(50 * time.Hour * 24) // 50 days

	var userRefreshToken = user.CreateUserRefreshToken(user, token, registerDate)

	assert.Equal(t, string(userExposeId), string(userRefreshToken.User.ExposeId))
	assert.Equal(t, string(token), string(userRefreshToken.RefreshToken))
	assert.Equal(t, registerDate, userRefreshToken.RegisteredDate)
	assert.Equal(t, expireDate, userRefreshToken.ExpireDate)

	t.Logf("userRefreshToken: %+v", userRefreshToken)
	t.Logf("userRefreshToken.User.ExposeId: %s", userRefreshToken.User.ExposeId)
	t.Logf("userRefreshToken.RefreshToken: %s", userRefreshToken.RefreshToken)
	t.Logf("userRefreshToken.RegisteredDate: %s", userRefreshToken.RegisteredDate)
	t.Logf("userRefreshToken.ExpireDate: %s", *userRefreshToken.ExpireDate)
}

func TestNewUserRefreshToken(t *testing.T) {
	var userExposeId, _ = pkgText.NewExposeId("TestExposeId")
	var user = getUser(userExposeId)

	var token, _ = text.NewToken("TestPassword")
	var registerDate = time.Now()
	var expireDate = registerDate.Add(30 * time.Hour * 24) // 30 days

	var userRefreshToken = user.NewUserRefreshToken(1, user, token, registerDate, expireDate)

	assert.Equal(t, 1, userRefreshToken.UserRefreshTokenId)
	assert.Equal(t, string(userExposeId), string(userRefreshToken.User.ExposeId))
	assert.Equal(t, string(token), string(userRefreshToken.RefreshToken))
	assert.Equal(t, registerDate, userRefreshToken.RegisteredDate)
	assert.Equal(t, expireDate, userRefreshToken.ExpireDate)

	t.Logf("userRefreshToken: %+v", userRefreshToken)
	t.Logf("userRefreshToken.User.ExposeId: %s", userRefreshToken.User.ExposeId)
	t.Logf("userRefreshToken.RefreshToken: %s", userRefreshToken.RefreshToken)
	t.Logf("userRefreshToken.RegisteredDate: %s", userRefreshToken.RegisteredDate)
	t.Logf("userRefreshToken.ExpireDate: %s", *userRefreshToken.ExpireDate)
}
