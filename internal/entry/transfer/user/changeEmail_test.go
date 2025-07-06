package user_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/core/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChangeEmailToCoreUserEmail(t *testing.T) {
	var email = "test@example.com"
	var userChangeEmailRequest = user.UserChangeEmailRequest{
		UserChangeEmail: user.UserChangeEmail{
			Email:       email,
		},
	}

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var coreUser = getUserForChangeEmail(userIdentifier)
	var verifyToken, _ = text.NewToken("verify_token_example")

	var coreUserEmail, err = userChangeEmailRequest.ToCoreUserEmail(coreUser, verifyToken, time.Now())

	assert.Nil(t, err)
	assert.NotNil(t, coreUserEmail)
	assert.Equal(t, coreUser.Identifier, coreUserEmail.User.Identifier)
	assert.Equal(t, email, string(coreUserEmail.Email))
}

func TestChangeEmailToCoreUserEmailError(t *testing.T) {
	var email = "testexample.com"
	var userChangeEmailRequest = user.UserChangeEmailRequest{
		UserChangeEmail: user.UserChangeEmail{
			Email:       email,
		},
	}

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var coreUser = getUserForChangeEmail(userIdentifier)
	var verifyToken, _ = text.NewToken("verify_token_example")

	var _, err = userChangeEmailRequest.ToCoreUserEmail(coreUser, verifyToken, time.Now())

	assert.Error(t, err)
}

func getUserForChangeEmail(identifier pkgText.Identifier) core.User {
	var userId uint = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	return core.NewUser(userId, identifier, userName, emailId, botFlag, userRegisteredDate, updateDate)
}
