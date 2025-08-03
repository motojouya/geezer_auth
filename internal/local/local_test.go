package io_test

import (
	"github.com/motojouya/geezer_auth/internal/io"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateRamdomString(t *testing.T) {
	l := io.CreateLocal()
	randomString := l.GenerateRamdomString(10, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	assert.Len(t, randomString, 10, "Random string should be of length 10")
	assert.Regexp(t, "^[a-zA-Z0-9]+$", randomString, "Random string should only contain alphanumeric characters")
}
