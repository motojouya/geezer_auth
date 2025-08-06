package common_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	"github.com/motojouya/geezer_auth/internal/shelter/company"
	"github.com/motojouya/geezer_auth/internal/shelter/role"
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getCompany(identifier pkgText.Identifier) company.Company {
	var companyId uint = 1
	var name, _ = pkgText.NewName("TestCompany")
	var registeredDate = time.Now()
	return company.NewCompany(companyId, identifier, name, registeredDate)
}

func getRole(label pkgText.Label) role.Role {
	var roleName, _ = pkgText.NewName("TestRole")
	var description, _ = text.NewText("Role for testing")
	var registeredDate = time.Now()
	return role.NewRole(roleName, label, description, registeredDate)
}

func TestFromShelterUser(t *testing.T) {
	var identifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var company = getCompany(identifier)

	var label1, _ = pkgText.NewLabel("TEST_ROLE")
	var role1 = getRole(label1)

	var label2, _ = pkgText.NewLabel("TOST_ROLE")
	var role2 = getRole(label2)

	var roles = []role.Role{role1, role2}

	var companyRole = user.NewCompanyRole(company, roles)

	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	var shelterUser = user.NewUser(userId, userIdentifier, userName, emailId, botFlag, userRegisteredDate, updateDate)

	var email, _ = pkgText.NewEmail("test2@gmail.com")

	var shelterUserAuthentic = user.NewUserAuthentic(shelterUser, companyRole, &email)

	var transferUser = common.FromShelterUser(shelterUserAuthentic)

	assert.Equal(t, string(userIdentifier), transferUser.Identifier)
	assert.Equal(t, string(emailId), transferUser.IdentifierEmail)
	assert.Equal(t, string(userName), transferUser.Name)
	assert.Equal(t, botFlag, transferUser.BotFlag)
	assert.Equal(t, updateDate, transferUser.UpdateDate)
	assert.Equal(t, string(identifier), transferUser.CompanyRole.Company.Identifier)
	assert.Equal(t, len(roles), len(transferUser.CompanyRole.Roles))
	assert.Equal(t, string(label1), transferUser.CompanyRole.Roles[0].Label)
	assert.Equal(t, string(label2), transferUser.CompanyRole.Roles[1].Label)

	t.Logf("transferUser: %+v", transferUser)
	t.Logf("transferUser.CompanyRole: %s", *transferUser.CompanyRole)
	t.Logf("transferUser.Email: %s", *transferUser.Email)
}
