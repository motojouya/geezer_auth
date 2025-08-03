package text_test

import (
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewToken(t *testing.T) {
	var tokenStr = "valid_token_string"
	var token, err = text.NewToken(tokenStr)
	if err != nil {
		t.Error("Failed to create new token")
	}

	assert.Equal(t, tokenStr, string(token))
}
