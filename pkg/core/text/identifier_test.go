package text_test

import (
	"errors"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewIdentifier(t *testing.T) {
	var identifierString = "US-ABCDEF"

	var identifier, err = text.NewEmail(" " + identifierString + " ")
	if err != nil {
		t.Fatalf("failed to create identifier: %v", err)
	}

	assert.Equal(t, identifierString, string(identifier))

	t.Logf("identifier: %s", string(identifier))
}

func TestNewIdentifierEmptyError(t *testing.T) {
	var identifierString = ""

	var _, err = text.NewEmail(identifierString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.As(err, &text.LengthError{}) {
		t.Fatalf("expected ErrInvalidEmailFormat, got %v", err)
	}
}

func TestNewIdentifierLengthError(t *testing.T) {
	var identifierSources = []string{"US-ABCDE", "US-ABCDEFG"}

	for _, identifierString := range identifierSources {
		var _, err = text.NewEmail(identifierString)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.As(err, &text.LengthError{}) {
			t.Fatalf("expected ErrInvalidEmailFormat, got %v", err)
		}
	}
}

func TestNewIdentifierFormatError(t *testing.T) {
	var identifierSources = []string{"USABCDEFG", "US.ABCDEF", "US_ABCDEF", "USA-BCDEF", "US-ABCDE1", "U1-ABCDEF"}

	for _, identifierString := range identifierSources {
		var _, err = text.NewEmail(identifierString)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.As(err, &text.FormatError{}) {
			t.Fatalf("expected ErrInvalidEmailFormat, got %v", err)
		}
	}
}
