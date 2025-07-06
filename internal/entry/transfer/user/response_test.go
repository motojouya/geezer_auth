package user_test

import (
	core "github.com/motojouya/geezer_auth/internal/core/company"
	coreRole "github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	coreUser "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getCompany(identifier pkgText.Identifier) core.Company {
	var name, _ = pkgText.NewName("TestRole")
	var registeredDate = time.Now()

	return core.NewCompany(1, identifier, name, registeredDate)
}

func getRoles(label pkgText.Label) []coreRole.Role {
	var roleName, _ = pkgText.NewName("TestRole")
	var description, _ = text.NewText("Role for testing")
	var registeredDate = time.Now()
	return []coreRole.Role{coreRole.NewRole(roleName, label, description, registeredDate)}
}

func TestFromCoreUserAuthenticToGetResponse(t *testing.T) {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	var userValue = coreUser.NewUser(userId, userIdentifier, userName, emailId, botFlag, userRegisteredDate, updateDate)

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var companyObj = getCompany(companyIdentifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roles = getRoles(label)

	var companyRole = coreUser.NewCompanyRole(companyObj, roles)

	var email, _ = pkgText.NewEmail("test_2@gmail.com")
	var coreUserObj = coreUser.NewUserAuthentic(userValue, companyRole, &email)

	var getResponse = user.FromCoreUserAuthenticToGetResponse(coreUserObj)

	assert.NotNil(t, getResponse)
	assert.Equal(t, string(userIdentifier), getResponse.User.Identifier)
}

func TestFromCoreUserAuthenticToUpdateResponse(t *testing.T) {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	var userValue = coreUser.NewUser(userId, userIdentifier, userName, emailId, botFlag, userRegisteredDate, updateDate)

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var companyObj = getCompany(companyIdentifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roles = getRoles(label)

	var companyRole = coreUser.NewCompanyRole(companyObj, roles)

	var email, _ = pkgText.NewEmail("test_2@gmail.com")
	var coreUserObj = coreUser.NewUserAuthentic(userValue, companyRole, &email)

	var accessToken = pkgText.NewJwtToken("access_token")

	var updateResponse = user.FromCoreUserAuthenticToUpdateResponse(coreUserObj, accessToken)

	assert.NotNil(t, updateResponse)
	assert.Equal(t, string(userIdentifier), updateResponse.User.Identifier)
	assert.Equal(t, string(accessToken), updateResponse.AccessToken)
}

func TestFromCoreUserAuthenticToRegisterResponse(t *testing.T) {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	var userValue = coreUser.NewUser(userId, userIdentifier, userName, emailId, botFlag, userRegisteredDate, updateDate)

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var companyObj = getCompany(companyIdentifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roles = getRoles(label)

	var companyRole = coreUser.NewCompanyRole(companyObj, roles)

	var email, _ = pkgText.NewEmail("test_2@gmail.com")
	var coreUserObj = coreUser.NewUserAuthentic(userValue, companyRole, &email)

	var accessToken = pkgText.NewJwtToken("access_token")
	var refreshToken, _ = text.NewToken("refresh_token")

	var registerResponse = user.FromCoreUserAuthenticToRegisterResponse(coreUserObj, refreshToken, accessToken)

	assert.NotNil(t, registerResponse)
	assert.Equal(t, string(userIdentifier), registerResponse.User.Identifier)
	assert.Equal(t, string(refreshToken), registerResponse.RefreshToken)
	assert.Equal(t, string(accessToken), registerResponse.AccessToken)
}
