package utility_test

import (
	"github.com/motojouya/geezer_auth/internal/utility"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO working

func TestNewRangeError(t *testing.T) {
	var name = "TestRangeError"
	var value = "TestValue"
	var min = 10
	var max = 20
	var message = "This is a test range error"

	var err = utility.NewRangeError(name, value, min, max, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message + " (name: " + name + ")", err.Error())
	assert.Equal(t, 400, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
}

func TestNewAuthenticationError(t *testing.T) {
	var name = "TestSystemConfigError"
	var message = "This is a test system config error"

	var err = utility.NewSystemConfigError(name, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message + " (name: " + name + ")", err.Error())
	assert.Equal(t, 500, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
}
