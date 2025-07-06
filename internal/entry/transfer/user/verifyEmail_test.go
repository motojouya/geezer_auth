package user_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVerifyEmailToCoreUserEmail(t *testing.T) {
	var email = "test@example.com"
	var verifyToken = "verify_token_example"
	var userVerifyEmailRequest = user.UserVerifyEmailRequest{
		UserVerifyEmail: user.UserVerifyEmail{
			Email:       email,
			VerifyToken: verifyToken,
		},
	}

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var coreUser = getUserForVerifyEmail(userIdentifier)

	var coreUserEmail, err = userVerifyEmailRequest.ToCoreUserEmail(coreUser, time.Now())

	assert.Nil(t, err)
	assert.NotNil(t, coreUserEmail)
	assert.Equal(t, coreUser.Identifier, coreUserEmail.User.Identifier)
	assert.Equal(t, email, string(coreUserEmail.Email))
	assert.Equal(t, verifyToken, string(coreUserEmail.VerifyToken))
}

func TestVerifyEmailToCoreUserEmailError(t *testing.T) {
	var email = "testexample.com"
	var verifyToken = "verify_token_example"
	var userVerifyEmailRequest = user.UserVerifyEmailRequest{
		UserVerifyEmail: user.UserVerifyEmail{
			Email:       email,
			VerifyToken: verifyToken,
		},
	}

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var coreUser = getUserForVerifyEmail(userIdentifier)

	var _, err = userVerifyEmailRequest.ToCoreUserEmail(coreUser, time.Now())

	assert.Error(t, err)
}

func getUserForVerifyEmail(identifier pkgText.Identifier) core.User {
	var userId uint = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	return core.NewUser(userId, identifier, userName, emailId, botFlag, userRegisteredDate, updateDate)
}
