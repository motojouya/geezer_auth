package company_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
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

func TestFromShelterCompany(t *testing.T) {
	var identifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var shelterCompany = getCompany(identifier)

	var response = company.FromShelterCompany(shelterCompany)

	assert.Equal(t, string(shelterCompany.Identifier), response.Company.Identifier)

	t.Logf("Response: %+v", response)
}

func TestFromShelterUserAuthentic(t *testing.T) {
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

	var shelterUsers = []shelterUser.UserAuthentic{*shelterUserObj}

	var response = company.FromShelterUserAuthentic(shelterUsers)

	assert.Len(t, response.Users, 1)
	assert.Equal(t, string(shelterUsers[0].Identifier), response.Users[0].Identifier)

	t.Logf("Response: %+v", response)
}
