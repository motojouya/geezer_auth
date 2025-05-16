package model_test

import (
	"github.com/motojouya/geezer_auth/internal/model"
	"testing"
)

func TestPassword(t *testing.T) {
	var password = "password"
	var hashed, err = model.GetPassword("password")
	if err != nil {
		t.Error("Failed to hash password")
	}

	var result = model.VerifyPassword(hashed, password)
	if result {
		t.Log("Password verification succeeded")
	} else {
		t.Error("Password verification failed")
	}
}
