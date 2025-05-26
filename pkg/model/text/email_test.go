package text_test

import (
	"github.com/motojouya/geezer_auth/internal/model/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"errors"
)

func TestNewEmail(t *testing.T) {
	var emailString = "test@google.com"

	var email, err = text.NewEmail(" " + emailString + " ")
	if err != nil {
		t.Fatalf("failed to create email: %v", err)
	}

	assert.Equal(t, emailString, string(email))

	t.Logf("email: %s", string(email))
}

func TestNewEmailEmptyError(t *testing.T) {
	var emailString = ""

	var email, err = text.NewEmail(emailString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.As(err, &text.LengthError{}) {
		t.Fatalf("expected ErrInvalidEmailFormat, got %v", err)
	}
}

// 長い方のパターンは320文字とかなのでやるのがめんどい
func TestNewEmailLengthError(t *testing.T) {
	var emailString = "ts"

	var email, err = text.NewEmail(emailString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.As(err, &text.LengthError{}) {
		t.Fatalf("expected ErrInvalidEmailFormat, got %v", err)
	}
}

func TestNewEmailFormatError(t *testing.T) {
	var emailString = "test_token"

	var email, err = text.NewEmail(emailString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.As(err, &text.FormatError{}) {
		t.Fatalf("expected ErrInvalidEmailFormat, got %v", err)
	}
}
