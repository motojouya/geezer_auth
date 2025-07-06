package user_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChangePassword(t *testing.T) {
	var next = "testpassnext"
	var now = "testpassnow"
	var userChangePasswordRequest = user.UserChangePasswordRequest{
		UserChangePassword: user.UserChangePassword{
			NextPassword: next,
			NowPassword:  now,
		},
	}

	var nextPassword, nextErr = userChangePasswordRequest.GetNextPassword()
	var nowPassword, nowErr = userChangePasswordRequest.GetNowPassword()

	assert.Nil(t, nextErr)
	assert.Nil(t, nowErr)
	assert.Equal(t, next, string(nextPassword))
	assert.Equal(t, now, string(nowPassword))
}

func TestChangePasswordError(t *testing.T) {
	var next = "test_pass_next"
	var now = "test_pass_now"
	var userChangePasswordRequest = user.UserChangePasswordRequest{
		UserChangePassword: user.UserChangePassword{
			NextPassword: next,
			NowPassword:  now,
		},
	}

	var _, nextErr = userChangePasswordRequest.GetNextPassword()
	var _, nowErr = userChangePasswordRequest.GetNowPassword()

	assert.Error(t, nextErr)
	assert.Error(t, nowErr)
}
