package text_test

import (
	"errors"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewName(t *testing.T) {
	var nameString = "TEST_LABEL"

	var name, err = text.NewName(" " + nameString + " ")
	if err != nil {
		t.Fatalf("failed to create name: %v", err)
	}

	assert.Equal(t, nameString, string(name))

	t.Logf("name: %s", string(name))
}

// 長い方のパターンは255文字とかなのでやるのがめんどい
func TestNewNameEmptyError(t *testing.T) {
	var nameString = ""

	var _, err = text.NewName(nameString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.As(err, &text.LengthError{}) {
		t.Fatalf("expected ErrInvalidNameFormat, got %v", err)
	}
}
