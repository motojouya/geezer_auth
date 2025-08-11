package user_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyEmailGetEmail(t *testing.T) {
	var email = "test@example.com"
	var verifyToken = "verify_token_example"
	var userVerifyEmailRequest = user.UserVerifyEmailRequest{
		UserVerifyEmail: user.UserVerifyEmail{
			Email:       email,
			VerifyToken: verifyToken,
		},
	}

	var emailResult, emailErr = userVerifyEmailRequest.GetEmail()
	var tokenResult, tokenErr = userVerifyEmailRequest.GetVerifyToken()

	assert.Nil(t, emailErr)
	assert.Nil(t, tokenErr)
	assert.Equal(t, email, string(emailResult))
	assert.Equal(t, verifyToken, string(tokenResult))
}

func TestVerifyEmailGetEmailError(t *testing.T) {
	var email = "testexample.com"
	var verifyToken = "verify_token_example"
	var userVerifyEmailRequest = user.UserVerifyEmailRequest{
		UserVerifyEmail: user.UserVerifyEmail{
			Email:       email,
			VerifyToken: verifyToken,
		},
	}

	var _, err = userVerifyEmailRequest.GetEmail()

	assert.Error(t, err)
}
