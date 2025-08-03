package text_test

import (
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewText(t *testing.T) {
	var textStr = "This is a test text."
	var textValue, err = text.NewText(" " + textStr + " ")
	if err != nil {
		t.Error("Failed to create new text:", err)
	}

	assert.Equal(t, textStr, string(textValue))
}
