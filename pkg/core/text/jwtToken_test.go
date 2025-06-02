package text_test

import (
	"github.com/motojouya/geezer_auth/internal/model/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJwtToken(t *testing.T) {
	var token = "test_token"

	var jwtToken = text.NewJwtToken(token)

	assert.Equal(t, token, string(jwtToken))

	t.Logf("jwtToken: %s", string(jwtToken))
}
