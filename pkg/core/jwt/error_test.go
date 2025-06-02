package jwt_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJwtError(t *testing.T) {
	var claims = "TestClaims"
	var value = "TestValue"
	var message = "This is a test nil error"

	var err = jwt.NewJwtError(claims, value, message)

	assert.Equal(t, claims, err.Claims)
	assert.Equal(t, value, err.Value)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message + " (claims: " + claims + ", value: " + value + ")", err.Error())
	assert.Equal(t, 400, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Claims: %s", err.Claims)
	t.Logf("error.Value: %s", err.Value)
}
