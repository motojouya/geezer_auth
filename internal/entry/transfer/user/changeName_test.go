package user_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChangeNameToCoreUser(t *testing.T) {
	var name = "test_name"
	var userChangeNameRequest = user.UserChangeNameRequest{
		UserChangeName: user.UserChangeName{
			Name: name,
		},
	}

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var userName, _ = pkgText.NewName("different_name")
	var shelterUser = getUserForChangeName(userIdentifier, userName)

	var shelterUserNameChanged, err = userChangeNameRequest.ToCoreUser(shelterUser, time.Now())

	assert.Nil(t, err)
	assert.NotNil(t, shelterUserNameChanged)
	assert.Equal(t, shelterUser.Identifier, shelterUserNameChanged.Identifier)
	assert.Equal(t, name, string(shelterUserNameChanged.Name))
}

func getUserForChangeName(identifier pkgText.Identifier, name pkgText.Name) shelter.User {
	var userId uint = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	return shelter.NewUser(userId, identifier, name, emailId, botFlag, userRegisteredDate, updateDate)
}
