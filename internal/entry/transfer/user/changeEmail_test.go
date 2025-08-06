package user_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChangeEmailToCoreUserEmail(t *testing.T) {
	var email = "test@example.com"
	var userChangeEmailRequest = user.UserChangeEmailRequest{
		UserChangeEmail: user.UserChangeEmail{
			Email: email,
		},
	}

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var shelterUser = getUserForChangeEmail(userIdentifier)
	var verifyToken, _ = text.NewToken("verify_token_example")

	var shelterUserEmail, err = userChangeEmailRequest.ToCoreUserEmail(shelterUser, verifyToken, time.Now())

	assert.Nil(t, err)
	assert.NotNil(t, shelterUserEmail)
	assert.Equal(t, shelterUser.Identifier, shelterUserEmail.User.Identifier)
	assert.Equal(t, email, string(shelterUserEmail.Email))
}

func TestChangeEmailToCoreUserEmailError(t *testing.T) {
	var email = "testexample.com"
	var userChangeEmailRequest = user.UserChangeEmailRequest{
		UserChangeEmail: user.UserChangeEmail{
			Email: email,
		},
	}

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var shelterUser = getUserForChangeEmail(userIdentifier)
	var verifyToken, _ = text.NewToken("verify_token_example")

	var _, err = userChangeEmailRequest.ToCoreUserEmail(shelterUser, verifyToken, time.Now())

	assert.Error(t, err)
}

func getUserForChangeEmail(identifier pkgText.Identifier) shelter.User {
	var userId uint = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	return shelter.NewUser(userId, identifier, userName, emailId, botFlag, userRegisteredDate, updateDate)
}
