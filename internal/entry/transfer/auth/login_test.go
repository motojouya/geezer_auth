package auth_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthLoginRequest(t *testing.T) {
	var identifier = "US-TESTES"
	var email = "test@example.com"
	var password = "testpassword"
	var authLoginRequest = auth.AuthLoginRequest{
		AuthLogin: auth.AuthLogin{
			AuthIdentifier: auth.AuthIdentifier{
				Identifier:      &identifier,
				EmailIdentifier: &email,
			},
			Password: password,
		},
	}

	var identifierResult, idErr = authLoginRequest.GetIdentifier()
	assert.Nil(t, idErr)
	assert.Equal(t, identifier, string(*identifierResult))

	var emailResult, emailErr = authLoginRequest.GetEmailIdentifier()
	assert.Nil(t, emailErr)
	assert.Equal(t, email, string(*emailResult))

	var passwordResult, passwordErr = authLoginRequest.GetPassword()
	assert.Nil(t, passwordErr)
	assert.Equal(t, password, string(passwordResult))
}

func TestAuthLoginRequestError(t *testing.T) {
	var identifier = "USB-TESTES"
	var email = "testexample.com"
	var password = "マルチバイト"
	var authLoginRequest = auth.AuthLoginRequest{
		AuthLogin: auth.AuthLogin{
			AuthIdentifier: auth.AuthIdentifier{
				Identifier:      &identifier,
				EmailIdentifier: &email,
			},
			Password: password,
		},
	}

	var _, idErr = authLoginRequest.GetIdentifier()
	assert.Error(t, idErr)

	var _, emailErr = authLoginRequest.GetEmailIdentifier()
	assert.Error(t, emailErr)

	var _, passwordErr = authLoginRequest.GetPassword()
	assert.Error(t, passwordErr)
}
