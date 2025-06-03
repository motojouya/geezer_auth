package text_test

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordSuccess(t *testing.T) {
	var passwordStr = "password"
	var password, err = text.NewPassword(passwordStr)
	if err != nil {
		t.Error("Failed to hash password")
	}
	var hashed, err = text.HashPassword(password)
	if err != nil {
		t.Error("Failed to hash password")
	}

	var result = text.VerifyPassword(hashed, password)
	if result {
		t.Log("Password verification succeeded")
	} else {
		t.Error("Password verification failed")
	}
}

func TestPasswordFailure(t *testing.T) {
	var password, err = text.NewPassword("password")
	if err != nil {
		t.Error("Failed to hash password")
	}
	var hashed, err = text.HashPassword(password)
	if err != nil {
		t.Error("Failed to hash password")
	}

	var result = text.VerifyPassword(hashed, "passward")
	if result {
		t.Error("Password verification succeeded")
	} else {
		t.Log("Password verification failed")
	}
}
