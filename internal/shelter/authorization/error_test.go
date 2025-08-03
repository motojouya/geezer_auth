package authorization_test

import (
	"github.com/motojouya/geezer_auth/internal/shelter/authorization"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewAuthorizationError(t *testing.T) {
	var role = "TestRole"
	var action = "TestAction"
	var message = "This is a test authorization error"
	var httpStatus uint = 403

	var err = authorization.NewAuthorizationError(role, action, message)

	assert.Equal(t, role, err.Role)
	assert.Equal(t, action, err.Action)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+" (role: "+role+", action: "+action+")", err.Error())
	assert.Equal(t, httpStatus, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Role: %s", err.Role)
	t.Logf("error.Action: %s", err.Action)
}

func TestNewTokenExpiredError(t *testing.T) {
	var expiresAt = time.Now()
	var message = "This is a test authorization error"
	var httpStatus uint = 403

	var err = authorization.NewTokenExpiredError(expiresAt, message)

	assert.Equal(t, expiresAt, err.ExpiresAt)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+" (expires at: "+expiresAt.Format("2006-01-02T15:04:05Z")+")", err.Error())
	assert.Equal(t, httpStatus, err.HttpStatus())

	t.Logf("error: %s", err.Error())
}
