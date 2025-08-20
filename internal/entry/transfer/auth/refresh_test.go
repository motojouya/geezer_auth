package auth_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthRefreshRequest(t *testing.T) {
	var refreshToken = "testtoken"
	var authRefreshRequest = auth.AuthRefreshRequest{
		AuthRefresh: auth.AuthRefresh{
			RefreshToken: refreshToken,
		},
	}

	var refreshResult, refreshErr = authRefreshRequest.GetRefreshToken()
	assert.Nil(t, refreshErr)
	assert.Equal(t, refreshToken, string(refreshResult))
}

func TestAuthRefreshRequestError(t *testing.T) {
	var refreshToken = "マルチバイト"
	var authRefreshRequest = auth.AuthRefreshRequest{
		AuthRefresh: auth.AuthRefresh{
			RefreshToken: refreshToken,
		},
	}

	var _, tokenErr = authRefreshRequest.GetRefreshToken()
	assert.Nil(t, tokenErr)
}
