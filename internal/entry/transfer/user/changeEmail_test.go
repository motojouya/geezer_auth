package user_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChangeEmailGetEmail(t *testing.T) {
	var email = "test@example.com"
	var userChangeEmailRequest = user.UserChangeEmailRequest{
		UserChangeEmail: user.UserChangeEmail{
			Email: email,
		},
	}

	var emailResult, err = userChangeEmailRequest.GetEmail()

	assert.Nil(t, err)
	assert.Equal(t, email, string(emailResult))
}

func TestChangeEmailToShelterUserEmailError(t *testing.T) {
	var email = "testexample.com"
	var userChangeEmailRequest = user.UserChangeEmailRequest{
		UserChangeEmail: user.UserChangeEmail{
			Email: email,
		},
	}

	var _, err = userChangeEmailRequest.GetEmail()

	assert.Error(t, err)
}
