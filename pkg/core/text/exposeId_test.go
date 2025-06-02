package text_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"errors"
)

func TestNewExposeId(t *testing.T) {
	var exposeIdString = "US-ABCDEF"

	var exposeId, err = text.NewEmail(" " + exposeIdString + " ")
	if err != nil {
		t.Fatalf("failed to create exposeId: %v", err)
	}

	assert.Equal(t, exposeIdString, string(exposeId))

	t.Logf("exposeId: %s", string(exposeId))
}

func TestNewExposeIdEmptyError(t *testing.T) {
	var exposeIdString = ""

	var exposeId, err = text.NewEmail(exposeIdString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.As(err, &text.LengthError{}) {
		t.Fatalf("expected ErrInvalidEmailFormat, got %v", err)
	}
}

func TestNewExposeIdLengthError(t *testing.T) {
	var exposeIdSources = []string{"US-ABCDE", "US-ABCDEFG"}

	for _, exposeIdString := range exposeIdSources {
		var exposeId, err = text.NewEmail(exposeIdString)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.As(err, &text.LengthError{}) {
			t.Fatalf("expected ErrInvalidEmailFormat, got %v", err)
		}
	}
}

func TestNewExposeIdFormatError(t *testing.T) {
	var exposeIdSources = []string{"USABCDEFG", "US.ABCDEF", "US_ABCDEF", "USA-BCDEF", "US-ABCDE1", "U1-ABCDEF"}

	for _, exposeIdString := range exposeIdSources {
		var exposeId, err = text.NewEmail(exposeIdString)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.As(err, &text.FormatError{}) {
			t.Fatalf("expected ErrInvalidEmailFormat, got %v", err)
		}
	}
}
