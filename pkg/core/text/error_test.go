package text_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLengthError(t *testing.T) {
	var name = "TestLengthError"
	var value = "short"
	var min uint = 5
	var max uint= 10
	var message = "This is a test length error"

	var err = text.NewLengthError(name, value, min, max, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, value, err.Value)
	assert.Equal(t, min, err.Min)
	assert.Equal(t, max, err.Max)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+" (name: "+name+", value: "+value+", min: 5, max: 10)", err.Error())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
	t.Logf("error.Value: %s", err.Value)
	t.Logf("error.Min: %d", err.Min)
	t.Logf("error.Max: %d", err.Max)
}

func TestNewCharacterError(t *testing.T) {
	var name = "TestCharacterError"
	var chars = "abc"
	var value = "xyz"
	var message = "This is a test character error"

	var err = text.NewCharacterError(name, chars, value, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, chars, err.Chars)
	assert.Equal(t, value, err.Value)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+" (name: "+name+", chars: "+chars+", value: "+value+")", err.Error())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
	t.Logf("error.Value: %s", err.Value)
	t.Logf("error.Chars: %s", err.Chars)
}

func TestNewFormatError(t *testing.T) {
	var name = "TestFormatError"
	var value = "invalid_format"
	var format = "expected_format"
	var message = "This is a test format error"

	var err = text.NewFormatError(name, value, format, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, value, err.Value)
	assert.Equal(t, format, err.Format)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+" (name: "+name+", value: "+value+", format: "+format+")", err.Error())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
	t.Logf("error.Value: %s", err.Value)
	t.Logf("error.Format: %s", err.Format)
}
