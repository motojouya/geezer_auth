package companyUser_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getCompany(identifier pkgText.Identifier) shelter.Company {
	var name, _ = pkgText.NewName("TestRole")
	var registeredDate = time.Now()

	return shelter.NewCompany(1, identifier, name, registeredDate)
}

func getRoles(label pkgText.Label) []shelterRole.Role {
	var roleName, _ = pkgText.NewName("TestRole")
	var description, _ = text.NewText("Role for testing")
	var registeredDate = time.Now()
	return []shelterRole.Role{shelterRole.NewRole(roleName, label, description, registeredDate)}
}

func TestFromShelterUserAuthenticToGetResponse(t *testing.T) {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	var userValue = shelterUser.NewUser(userId, userIdentifier, userName, emailId, botFlag, userRegisteredDate, updateDate)

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var companyObj = getCompany(companyIdentifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roles = getRoles(label)

	var companyRole = shelterUser.NewCompanyRole(companyObj, roles)

	var email, _ = pkgText.NewEmail("test_2@gmail.com")
	var shelterUserObj = shelterUser.NewUserAuthentic(userValue, companyRole, &email)

	var getResponse = companyUser.FromShelterUserAuthenticToGetResponse(shelterUserObj)

	assert.NotNil(t, getResponse)
	assert.Equal(t, string(userIdentifier), getResponse.User.Identifier)
}

func TestFromShelterRoles(t *testing.T) {
	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roles = getRoles(label)

	var response = companyUser.FromShelterRoles(roles)

	assert.NotNil(t, response)
	assert.Len(t, response.Roles, 1)
	assert.Equal(t, string(label), response.Roles[0].Label)
}

func TestFromToken(t *testing.T) {
	var token, _ = text.NewToken("test-token")
	var response = companyUser.FromToken(token)
	assert.Equal(t, string(token), response.Token)
}
