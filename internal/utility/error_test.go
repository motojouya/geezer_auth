package utility_test

import (
	"github.com/motojouya/geezer_auth/internal/utility"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRangeError(t *testing.T) {
	var name = "TestRangeError"
	var value = "TestValue"
	var min = 10
	var max = 20
	var message = "This is a test range error"

	var err = utility.NewRangeError(name, value, min, max, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, value, err.Value)
	assert.Equal(t, min, err.Min)
	assert.Equal(t, max, err.Max)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message + " (name: " + name + ", value: " + valu + ", min: " + min + ", max: " + max + ")", err.Error())
	assert.Equal(t, 400, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
}

func TestNewAuthenticationError(t *testing.T) {
	var userIdentifier = "TestUserIdentifier"
	var message = "This is a test system config error"

	var err = utility.NewAuthenticationError(userIdentifier, message)

	assert.Equal(t, userIdentifier, err.UserIdentifier)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message + " (name: " + name + ")", err.Error())
	assert.Equal(t, 500, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.UserIdentifier)
}
