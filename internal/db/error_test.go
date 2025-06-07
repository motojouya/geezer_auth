package db_test

import (
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNotFoundError(t *testing.T) {
	var table = "TestUser"
	var key = "TestUserId"
	var value = "TestValue"
	var keys = map[string]string{key: value}
	var message = "This is a test range error"

	var err = db.NewNotFoundError(table, keys, message)

	assert.Equal(t, name, err.Name)
	var val, exist = err.keys[key]
	assert.True(t, exist)
	assert.Equal(t, value, val)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+" (table: "+table+", keys: {"+key+": "+value+",})", err.Error())
	assert.Equal(t, 400, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
}

func TestNewDuplicateError(t *testing.T) {
	var table = "TestUser"
	var key = "TestUserId"
	var value = "TestValue"
	var keys = map[string]string{key: value}
	var message = "This is a test range error"

	var err = db.NewDuplicateError(table, keys, message)

	assert.Equal(t, name, err.Name)
	var val, exist = err.keys[key]
	assert.True(t, exist)
	assert.Equal(t, value, val)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+" (table: "+table+", keys: {"+key+": "+value+",})", err.Error())
	assert.Equal(t, 400, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
}
