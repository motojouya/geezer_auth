package text_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLabel(t *testing.T) {
	var labelString = "TEST_LABEL"

	var label, err = text.NewLabel(" " + labelString + " ")
	if err != nil {
		t.Fatalf("failed to create label: %v", err)
	}

	assert.Equal(t, labelString, string(label))

	t.Logf("label: %s", string(label))
}

func TestNewLabelEmptyError(t *testing.T) {
	var labelString = ""

	var _, err = text.NewLabel(labelString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if _, ok := err.(*text.LengthError); !ok {
		t.Fatalf("expected ErrInvalidLabelFormat, got %v", err)
	}
}

// 長い方のパターンは320文字とかなのでやるのがめんどい
func TestNewLabelLengthError(t *testing.T) {
	var labelString = "T"

	var _, err = text.NewLabel(labelString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if _, ok := err.(*text.LengthError); !ok {
		t.Fatalf("expected ErrInvalidLabelFormat, got %v", err)
	}
}

func TestNewLabelFormatError(t *testing.T) {
	var labelSources = []string{"AbC", "A-B", "A.C", "A1B"}

	for _, labelString := range labelSources {
		var _, err = text.NewLabel(labelString)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if _, ok := err.(*text.FormatError); !ok {
			t.Fatalf("expected ErrInvalidLabelFormat, got %v", err)
		}
	}
}
