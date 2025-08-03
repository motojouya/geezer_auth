package text_test

import (
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPassword(t *testing.T) {
	var passwordStr = "PassWord123"
	var password, unexpectErr = text.NewPassword(passwordStr)
	assert.Nil(t, unexpectErr)
	assert.Equal(t, passwordStr, string(password))

	var errorStr = "Pass_Word123"
	var _, expectErr = text.NewPassword(errorStr)
	assert.Error(t, expectErr)
}

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
