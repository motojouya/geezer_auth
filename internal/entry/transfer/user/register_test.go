package user_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRegisterToShelterUser(t *testing.T) {
	var email = "test@example.com"
	var name = "Test_User"
	var bot = false
	var password = "password"
	var userRegisterRequest = user.UserRegisterRequest{
		UserRegister: user.UserRegister{
			Email:    email,
			Name:     name,
			Bot:      bot,
			Password: password,
		},
	}

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var userRegisterDate = time.Now()
	var unsavedUser, err = userRegisterRequest.ToShelterUser(userIdentifier, userRegisterDate)

	assert.Nil(t, err)
	assert.NotNil(t, unsavedUser)
	assert.Equal(t, userIdentifier, unsavedUser.Identifier)
	assert.Equal(t, email, string(unsavedUser.ExposeEmailId))
	assert.Equal(t, name, string(unsavedUser.Name))
	assert.Equal(t, bot, unsavedUser.BotFlag)
	assert.Equal(t, userRegisterDate, unsavedUser.RegisteredDate)
	assert.Equal(t, userRegisterDate, unsavedUser.UpdateDate)

	t.Logf("Unsaved User: %+v", unsavedUser)

	var passwordText, passwordErr = userRegisterRequest.GetPassword()
	assert.Nil(t, passwordErr)
	assert.NotNil(t, passwordText)
	assert.Equal(t, password, string(passwordText))

	t.Logf("password: %+s", passwordText)
}

func TestRegisterToShelterUserError(t *testing.T) {
	var email = "testexample.com"
	var name = "Test_User"
	var bot = false
	var password = "password_123"
	var userRegisterRequest = user.UserRegisterRequest{
		UserRegister: user.UserRegister{
			Email:    email,
			Name:     name,
			Bot:      bot,
			Password: password,
		},
	}

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var userRegisterDate = time.Now()
	var _, err = userRegisterRequest.ToShelterUser(userIdentifier, userRegisterDate)

	assert.Error(t, err)
	t.Logf("Error: %v", err)

	var _, passwordErr = userRegisterRequest.GetPassword()
	assert.Error(t, passwordErr)

	t.Logf("Error: %v", passwordErr)
}
