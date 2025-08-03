package essence_test

import (
	"github.com/motojouya/geezer_auth/internal/core/essence"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestInvalidArgumentError(t *testing.T) {
	var name = "TestInvalidArgumentError"
	var value = "invalid_value"
	var message = "This is a test invalid argument error"
	var httpStatus uint = 400

	var err = essence.NewInvalidArgumentError(name, value, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, value, err.Value)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+", name: "+name+", value: "+value, err.Error())
	assert.Equal(t, httpStatus, err.HttpStatus())

	t.Logf("error: %s", err)
	t.Logf("error message: %s", err.Error())
}

func TestNewRangeError(t *testing.T) {
	var name = "TestRangeError"
	var value = 100
	var min = 10
	var max = 20
	var message = "This is a test range error"
	var httpStatus uint = 400

	var err = essence.NewRangeError(name, value, min, max, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, value, err.Value)
	assert.Equal(t, min, err.Min)
	assert.Equal(t, max, err.Max)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+", name: "+name+", value: "+strconv.Itoa(value)+", min: "+strconv.Itoa(min)+", max: "+strconv.Itoa(max), err.Error())
	assert.Equal(t, httpStatus, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
}

func TestNewAuthenticationError(t *testing.T) {
	var userIdentifier = "TestUserIdentifier"
	var message = "This is a test system config error"
	var httpStatus uint = 401

	var err = essence.NewAuthenticationError(userIdentifier, message)

	assert.Equal(t, userIdentifier, err.UserIdentifier)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+", userIdentifier: "+userIdentifier, err.Error())
	assert.Equal(t, httpStatus, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.UserIdentifier)
}

func TestNewNotFoundError(t *testing.T) {
	var table = "TestUser"
	var key = "TestUserId"
	var value = "TestValue"
	var keys = map[string]string{key: value}
	var message = "This is a test range error"
	var httpStatus uint = 400

	var err = essence.NewNotFoundError(table, keys, message)

	assert.Equal(t, table, err.Table)
	var val, exist = err.Keys[key]
	assert.True(t, exist)
	assert.Equal(t, value, val)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+", table: "+table+", keys: {"+key+": "+value+", }", err.Error())
	assert.Equal(t, httpStatus, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Table: %s", err.Table)
}

func TestNewDuplicateError(t *testing.T) {
	var table = "TestUser"
	var key = "TestUserId"
	var value = "TestValue"
	var keys = map[string]string{key: value}
	var message = "This is a test range error"
	var httpStatus uint = 400

	var err = essence.NewDuplicateError(table, keys, message)

	assert.Equal(t, table, err.Table)
	var val, exist = err.Keys[key]
	assert.True(t, exist)
	assert.Equal(t, value, val)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+", table: "+table+", keys: {"+key+": "+value+", }", err.Error())
	assert.Equal(t, httpStatus, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Table: %s", err.Table)
}
