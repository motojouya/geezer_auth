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

func TestCreateUserEmail(t *testing.T) {
	var userExposeId, _ = pkgText.NewExposeId("TestExposeId")
	var user = getUser(userExposeId)

	var email = pkgText.NewEmail("test@google.com")
	var verifyToken, _ = text.NewToken("TestVerifyToken")
	var registerDate = time.Now()

	var userEmail = user.CreateUserEmail(user, email, verifyToken, registerDate)

	assert.Equal(t, string(userExposeId), string(userEmail.User.ExposeId))
	assert.Equal(t, string(email), string(userEmail.Email))
	assert.Equal(t, string(verifyToken), string(userEmail.VerifyToken))
	assert.Equal(t, registerDate, userEmail.RegisteredDate)
	assert.Nil(t, *userEmail.VerifyDate)
	assert.Nil(t, *userEmail.ExpireDate)

	t.Logf("userEmail: %+v", userEmail)
	t.Logf("userEmail.User.ExposeId: %s", userEmail.User.ExposeId)
	t.Logf("userEmail.Email: %s", userEmail.Email)
	t.Logf("userEmail.VerifyToken: %s", userEmail.VerifyToken)
	t.Logf("userEmail.RegisteredDate: %s", userEmail.RegisteredDate)
	t.Logf("userEmail.VerifyDate: %s", *userEmail.VerifyDate)
	t.Logf("userEmail.ExpireDate: %s", *userEmail.ExpireDate)
}

func TestNewUserEmail(t *testing.T) {
	var userExposeId, _ = pkgText.NewExposeId("TestExposeId")
	var user = getUser(userExposeId)

	var email = pkgText.NewEmail("test@google.com")
	var verifyToken, _ = text.NewToken("TestVerifyToken")
	var registerDate = time.Now()
	var verifyDate = registerDate.Add(24 * time.Hour)
	var expireDate = registerDate.Add(48 * time.Hour)

	var userEmail = user.NewUserEmail(1, user, email, verifyToken, registerDate, &verifyDate, &expireDate)

	assert.Equal(t, 1, userEmail.UserEmailID)
	assert.Equal(t, string(userExposeId), string(userEmail.User.ExposeId))
	assert.Equal(t, string(email), string(userEmail.Email))
	assert.Equal(t, string(verifyToken), string(userEmail.VerifyToken))
	assert.Equal(t, registerDate, userEmail.RegisteredDate)
	assert.Equal(t, verifyDate, *userEmail.VerifyDate)
	assert.Equal(t, expireDate, *userEmail.ExpireDate)

	t.Logf("userEmail: %+v", userEmail)
	t.Logf("userEmail.User.ExposeId: %s", userEmail.User.ExposeId)
	t.Logf("userEmail.Email: %s", userEmail.Email)
	t.Logf("userEmail.VerifyToken: %s", userEmail.VerifyToken)
	t.Logf("userEmail.RegisteredDate: %s", userEmail.RegisteredDate)
	t.Logf("userEmail.VerifyDate: %s", *userEmail.VerifyDate)
	t.Logf("userEmail.ExpireDate: %s", *userEmail.ExpireDate)
}
