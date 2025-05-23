package model_test

import (
	"github.com/motojouya/geezer_auth/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJwtToken(t *testing.T) {
	var token = "test_token"

	var jwtToken = model.NewJwtToken(token)

	assert.Equal(t, token, string(jwtToken))

	t.Logf("jwtToken: %s", string(jwtToken))
}
