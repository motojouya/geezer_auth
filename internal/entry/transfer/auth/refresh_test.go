package auth_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthRefreshRequest(t *testing.T) {
	var identifier = "US-TESTES"
	var email = "test@example.com"
	var refreshToken = "testtoken"
	var authRefreshRequest = auth.AuthRefreshRequest{
		AuthRefresh: auth.AuthRefresh{
			AuthIdentifier: auth.AuthIdentifier{
				Identifier:      &identifier,
				EmailIdentifier: &email,
			},
			RefreshToken: refreshToken,
		},
	}

	var identifierResult, idErr = authRefreshRequest.GetIdentifier()
	assert.Nil(t, idErr)
	assert.Equal(t, identifier, string(*identifierResult))

	var emailResult, emailErr = authRefreshRequest.GetEmailIdentifier()
	assert.Nil(t, emailErr)
	assert.Equal(t, email, string(*emailResult))

	var refreshResult, refreshErr = authRefreshRequest.GetRefreshToken()
	assert.Nil(t, refreshErr)
	assert.Equal(t, refreshToken, string(refreshResult))
}

func TestAuthRefreshRequestError(t *testing.T) {
	var identifier = "USB-TESTES"
	var email = "testexample.com"
	var refreshToken = "マルチバイト"
	var authRefreshRequest = auth.AuthRefreshRequest{
		AuthRefresh: auth.AuthRefresh{
			AuthIdentifier: auth.AuthIdentifier{
				Identifier:      &identifier,
				EmailIdentifier: &email,
			},
			RefreshToken: refreshToken,
		},
	}

	var _, idErr = authRefreshRequest.GetIdentifier()
	assert.Error(t, idErr)

	var _, emailErr = authRefreshRequest.GetEmailIdentifier()
	assert.Error(t, emailErr)

	var _, tokenErr = authRefreshRequest.GetRefreshToken()
	assert.Nil(t, tokenErr)
}
