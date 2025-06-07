package user_test

import (
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	pkgUser "github.com/motojouya/geezer_auth/pkg/core/user"
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

func getRoles(label pkgText.Label) []role.Role {
	var roleName, _ = pkgText.NewName("TestRole")
	var description, _ = text.NewText("Role for testing")
	var registeredDate = time.Now()
	return []role.Role{role.NewRole(roleId, roleName, label, description, registeredDate)}
}

func TestNewCompanyRole(t *testing.T) {
	var identifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var company = getCompany(identifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roles = getRoles(label)

	var companyRole = user.NewCompanyRole(company, roles)

	assert.Equal(t, string(identifier), string(companyRole.Company.Identifier))
	assert.Equal(t, len(roles), len(companyRole.Roles))
	assert.Equal(t, string(label), companyRole.Roles[0].Label)

	t.Logf("companyRole: %+v", companyRole)
	t.Logf("company: %+v", companyRole.Company)
	t.Logf("company.Identifier: %s", companyRole.Company.Identifier)
	t.Logf("role: %+v", companyRole.Roles[0])
	t.Logf("role.Label: %s", companyRole.Roles[0].Label)
}

func TestNewUserAuthentic(t *testing.T) {
	var userId uint = 1
	var userIdentifier = text.NewIdentifier("TestIdentifier")
	var emailId = text.NewEmail("test@gmail.com")
	var email = text.NewEmail("test_2@gmail.com")
	var userName = text.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var company = getCompany(companyIdentifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roles = getRoles(label)

	var companyRole = user.NewCompanyRole(company, roles)

	var userObj = user.NewUserAuthentic(userId, userIdentifier, userName, emailId, &email, botFlag, userRegisteredDate, updateDate, &companyRole)

	assert.Equal(t, userId, userObj.UserId)
	assert.Equal(t, string(userIdentifier), string(userObj.Identifier))
	assert.Equal(t, string(emailId), string(userObj.ExposeEmailId))
	assert.Equal(t, string(email), string(*userObj.Email))
	assert.Equal(t, string(userName), string(userObj.Name))
	assert.Equal(t, botFlag, userObj.BotFlag)
	assert.Equal(t, userRegisteredDate, userObj.RegisteredDate)
	assert.Equal(t, updateDate, userObj.UpdateDate)
	assert.Equal(t, string(companyIdentifier), string(userObj.CompanyRole.Company.Identifier))
	assert.Equal(t, 1, len(userObj.CompanyRole.Roles))
	assert.Equal(t, string(label), string(userObj.CompanyRole.Roles[0].Label))

	t.Logf("user: %+v", userObj)
	t.Logf("user.UserId: %d", userObj.UserId)
	t.Logf("user.Identifier: %s", userObj.Identifier)
	t.Logf("user.ExposeEmailId: %s", userObj.ExposeEmailId)
	t.Logf("user.Email: %s", *userObj.Email)
	t.Logf("user.Name: %s", userObj.Name)
	t.Logf("user.BotFlag: %t", userObj.BotFlag)
	t.Logf("user.RegisteredDate: %t", userObj.RegisteredDate)
	t.Logf("user.UpdateDate: %t", userObj.UpdateDate)

	t.Logf("companyRole: %+v", userObj.CompanyRole)
	t.Logf("company: %+v", userObj.CompanyRole.Company)
	t.Logf("company.Identifier: %d", userObj.CompanyRole.Company.Identifier)
	t.Logf("role: %+v", userObj.CompanyRole.Roles[0])
	t.Logf("role.Label: %d", userObj.CompanyRole.Roles[0].Label)
}

func TestModelToAccessTokenUser(t *testing.T) {
	var userId uint = 1
	var userIdentifier = text.NewIdentifier("TestIdentifier")
	var emailId = text.NewEmail("test@gmail.com")
	var email = text.NewEmail("test_2@gmail.com")
	var userName = text.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()

	var companyId uint = 1
	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()
	var company = company.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roleName, _ = pkgText.NewName("TestRole")
	var description, _ = text.NewText("Role for testing")
	var roleRegisteredDate = time.Now()
	var roles = []role.Role{role.NewRole(roleId, roleName, label, description, roleRegisteredDate)}

	var companyRole = user.NewCompanyRole(company, roles)

	var userObj = user.NewUserAuthentic(userId, userIdentifier, userName, emailId, &email, botFlag, userRegisteredDate, updateDate, &companyRole)

	var jwtUser = userObj.ToJwtUser()

	assert.Equal(t, string(userIdentifier), string(jwtUser.Identifier))
	assert.Equal(t, string(emailId), string(jwtUser.ExposeEmailId))
	assert.Equal(t, string(email), string(*jwtUser.Email))
	assert.Equal(t, string(userName), string(jwtUser.Name))
	assert.Equal(t, botFlag, jwtUser.BotFlag)
	assert.Equal(t, updateDate, jwtUser.UpdateDate)
	assert.Equal(t, companyIdentifier, jwtUser.CompanyRole.Company.Identifier)
	assert.Equal(t, companyName, jwtUser.CompanyRole.Company.Name)
	assert.Equal(t, len(roles), len(jwtUser.CompanyRole.Roles))
	assert.Equal(t, label, jwtUser.CompanyRole.Roles[0].Label)
	assert.Equal(t, roleName, jwtUser.CompanyRole.Roles[0].Name)

	t.Logf("user: %+v", jwtUser)
	t.Logf("user.Identifier: %s", jwtUser.Identifier)
	t.Logf("user.ExposeEmailId: %s", jwtUser.ExposeEmailId)
	t.Logf("user.Email: %s", *jwtUser.Email)
	t.Logf("user.Name: %s", jwtUser.Name)
	t.Logf("user.BotFlag: %t", jwtUser.BotFlag)
	t.Logf("user.UpdateDate: %t", jwtUser.UpdateDate)
	t.Logf("company: %+v", jwtUser.Company)
	t.Logf("company.Identifier: %s", jwtUser.CompanyRole.Company.Identifier)
	t.Logf("company.Name: %s", jwtUser.CompanyRole.Company.Name)
	t.Logf("Role.Label: %s", jwtUser.CompanyRole.Roles[0].Label)
	t.Logf("Role.Name: %s", jwtUser.CompanyRole.Roles[0].Name)
}
