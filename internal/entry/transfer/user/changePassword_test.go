package user_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChangePassword(t *testing.T) {
	var next = "testpassnext"
	var userChangePasswordRequest = user.UserChangePasswordRequest{
		UserChangePassword: user.UserChangePassword{
			Password: next,
		},
	}

	var nextPassword, nextErr = userChangePasswordRequest.GetPassword()

	assert.Nil(t, nextErr)
	assert.Equal(t, next, string(nextPassword))
}

func TestChangePasswordError(t *testing.T) {
	var next = "test_pass_next"
	var userChangePasswordRequest = user.UserChangePasswordRequest{
		UserChangePassword: user.UserChangePassword{
			Password: next,
		},
	}

	var _, nextErr = userChangePasswordRequest.GetPassword()

	assert.Error(t, nextErr)
}
