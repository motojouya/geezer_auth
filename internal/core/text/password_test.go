package text_test

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	"testing"
)

func TestPasswordSuccess(t *testing.T) {
	var passwordStr = "password"
	var password, createErr = text.NewPassword(passwordStr)
	if createErr != nil {
		t.Error("Failed to hash password")
	}
	var hashed, hashErr = text.HashPassword(password)
	if hashErr != nil {
		t.Error("Failed to hash password")
	}

	var err = text.VerifyPassword(hashed, password)
	if err == nil {
		t.Log("Password verification succeeded")
	} else {
		t.Error("Password verification failed")
	}
}

func TestPasswordFailure(t *testing.T) {
	var password, createErr = text.NewPassword("password")
	if createErr != nil {
		t.Error("Failed to hash password")
	}
	var hashed, hashErr = text.HashPassword(password)
	if hashErr != nil {
		t.Error("Failed to hash password")
	}

	var err = text.VerifyPassword(hashed, "passward")
	if err == nil {
		t.Error("Password verification succeeded")
	} else {
		t.Log("Password verification failed")
	}
}
