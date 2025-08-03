package user_test

import (
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
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
	var shelterUser = getUserForVerifyEmail(userIdentifier)

	var shelterUserEmail, err = userVerifyEmailRequest.ToCoreUserEmail(shelterUser, time.Now())

	assert.Nil(t, err)
	assert.NotNil(t, shelterUserEmail)
	assert.Equal(t, shelterUser.Identifier, shelterUserEmail.User.Identifier)
	assert.Equal(t, email, string(shelterUserEmail.Email))
	assert.Equal(t, verifyToken, string(shelterUserEmail.VerifyToken))
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
	var shelterUser = getUserForVerifyEmail(userIdentifier)

	var _, err = userVerifyEmailRequest.ToCoreUserEmail(shelterUser, time.Now())

	assert.Error(t, err)
}

func getUserForVerifyEmail(identifier pkgText.Identifier) shelter.User {
	var userId uint = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	return shelter.NewUser(userId, identifier, userName, emailId, botFlag, userRegisteredDate, updateDate)
}
