package utility_test

import (
	"github.com/motojouya/geezer_auth/pkg/utility"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNilError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"

	var err = utility.NewNilError(name, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message + " (name: " + name + ")", err.Error())
	assert.Equal(t, 400, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
}

func TestNewSystemConfigError(t *testing.T) {
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

//TODO working

func TestNewPropertyError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = utility.NewNilError(name, message)

	var prop = "TestPath"
	var httpStatus = 210

	var propertyError = utility.NewPropertyError(prop, httpStatus, err)

	assert.Equal(t, prop, propertyError.Property)
	assert.Equal(t, httpStatus, propertyError.HttpStatus)
	assert.Equal(t, message, propertyError.Unwrap().Error())
	assert.Equal(t, message + " (property: " + prop + ", httpStatus: " + httpStatus + ")", propertyError.Error())

	t.Logf("error: %s", propertyError.Error())
	t.Logf("error.Name: %s", propertyError.Name)
}

func TestCreatePropertyError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = utility.NewNilError(name, message)

	var prop = "TestPath"

	var propertyError = utility.CreatePropertyError(prop, err)

	assert.Equal(t, prop, propertyError.Property)
	assert.Equal(t, 400, propertyError.HttpStatus)
	assert.Equal(t, message, propertyError.Unwrap().Error())
	assert.Equal(t, message + " (property: " + prop + ", httpStatus: " + httpStatus + ")", propertyError.Error())

	t.Logf("error: %s", propertyError.Error())
	t.Logf("error.Name: %s", propertyError.Name)
}
