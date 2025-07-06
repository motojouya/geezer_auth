package user_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
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
	var coreUser = getUserForChangeName(userIdentifier, userName)

	var coreUserNameChanged, err = userChangeNameRequest.ToCoreUser(coreUser, time.Now())

	assert.Nil(t, err)
	assert.NotNil(t, coreUserNameChanged)
	assert.Equal(t, coreUser.Identifier, coreUserNameChanged.Identifier)
	assert.Equal(t, name, string(coreUserNameChanged.Name))
}

func getUserForChangeName(identifier pkgText.Identifier, name pkgText.Name) core.User {
	var userId uint = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	return core.NewUser(userId, identifier, name, emailId, botFlag, userRegisteredDate, updateDate)
}
