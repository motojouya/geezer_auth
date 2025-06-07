package authorization_test

import (
	"github.com/motojouya/geezer_auth/internal/authorization"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewAuthorizationError(t *testing.T) {
	var role = "TestRole"
	var action = "TestAction"
	var message = "This is a test authorization error"

	var err = authorization.NewAuthorizationError(role, action, message)

	assert.Equal(t, role, err.Role)
	assert.Equal(t, action, err.Action)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+" (role: "+role+", action: "+action+")", err.Error())
	assert.Equal(t, 403, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Role: %s", err.Role)
	t.Logf("error.Action: %s", err.Action)
}

func TestNewTokenExpiredError(t *testing.T) {
	var expiresAt = time.Now()
	var message = "This is a test authorization error"

	var err = authorization.NewTokenExpiredError(expiresAt, message)

	assert.Equal(t, expiresAt, err.ExpiresAt)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+" (expires at: "+expiresAt+")", err.Error())
	assert.Equal(t, 403, err.HttpStatus())

	t.Logf("error: %s", err.Error())
}
