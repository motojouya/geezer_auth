package user_test

import (
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForCompanyRole(identifier pkgText.Identifier) user.User {
	var userId uint = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return user.NewUser(userId, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func getCompanyForCompanyRole(identifier pkgText.Identifier) company.Company {
	var companyId uint = 1
	var name, _ = pkgText.NewName("TestCompany")
	var registeredDate = time.Now()

	return company.NewCompany(companyId, identifier, name, registeredDate)
}

func getRole(label pkgText.Label) role.Role {
	var name, _ = pkgText.NewName("TestRole")
	var description, _ = text.NewText("TestRoleDescription")
	var registeredDate = time.Now()

	return role.NewRole(name, label, description, registeredDate)
}

func TestCreateUserCompanyRole(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var userValue = getUserForCompanyRole(userIdentifier)

	var companyIdentifier, _ = pkgText.NewIdentifier("TestCompanyIdentifier")
	var company = getCompanyForCompanyRole(companyIdentifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role = getRole(label)

	var registerDate = time.Now()

	var userCompanyRole = user.CreateUserCompanyRole(userValue, company, role, registerDate)

	assert.Equal(t, string(userIdentifier), string(userCompanyRole.User.Identifier))
	assert.Equal(t, string(companyIdentifier), string(userCompanyRole.Company.Identifier))
	assert.Equal(t, string(label), string(userCompanyRole.Role.Label))
	assert.Equal(t, registerDate, userCompanyRole.RegisterDate)
	assert.Nil(t, userCompanyRole.ExpireDate)

	t.Logf("userCompanyRole: %+v", userCompanyRole)
	t.Logf("userCompanyRole.User.Identifier: %s", userCompanyRole.User.Identifier)
	t.Logf("userCompanyRole.Company.Identifier: %s", userCompanyRole.Company.Identifier)
	t.Logf("userCompanyRole.Role.Label: %s", userCompanyRole.Role.Label)
	t.Logf("userCompanyRole.RegisterDate: %s", userCompanyRole.RegisterDate)
}

func TestNewUserCompanyRole(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var userValue = getUserForCompanyRole(userIdentifier)

	var companyIdentifier, _ = pkgText.NewIdentifier("TestCompanyIdentifier")
	var company = getCompanyForCompanyRole(companyIdentifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role = getRole(label)

	var registerDate = time.Now()
	var expireDate = registerDate.Add(24 * time.Hour)

	var userCompanyRole = user.NewUserCompanyRole(1, userValue, company, role, registerDate, &expireDate)

	assert.Equal(t, uint(1), userCompanyRole.PersistKey)
	assert.Equal(t, string(userIdentifier), string(userCompanyRole.User.Identifier))
	assert.Equal(t, string(companyIdentifier), string(userCompanyRole.Company.Identifier))
	assert.Equal(t, string(label), string(userCompanyRole.Role.Label))
	assert.Equal(t, registerDate, userCompanyRole.RegisterDate)
	assert.Equal(t, expireDate, *userCompanyRole.ExpireDate)

	t.Logf("userCompanyRole: %+v", userCompanyRole)
	t.Logf("userCompanyRole.PersistKey: %d", userCompanyRole.PersistKey)
	t.Logf("userCompanyRole.User.Identifier: %s", userCompanyRole.User.Identifier)
	t.Logf("userCompanyRole.Company.Identifier: %s", userCompanyRole.Company.Identifier)
	t.Logf("userCompanyRole.Role.Label: %s", userCompanyRole.Role.Label)
	t.Logf("userCompanyRole.RegisterDate: %s", userCompanyRole.RegisterDate)
	t.Logf("userCompanyRole.ExpireDate: %s", *userCompanyRole.ExpireDate)
}

func TestIsUserUCR(t *testing.T) {
	var diffUserIdentifier, _ = pkgText.NewIdentifier("US-TOSTOS")
	var diffUserValue = getUserForCompanyRole(diffUserIdentifier)

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var userValue = getUserForCompanyRole(userIdentifier)

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var company = getCompanyForCompanyRole(companyIdentifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role = getRole(label)

	var registerDate = time.Now()
	var expireDate = registerDate.Add(24 * time.Hour)

	var userCompanyRole = user.NewUserCompanyRole(1, userValue, company, role, registerDate, &expireDate)

	var same = user.IsUserUCR(userValue)(userCompanyRole)
	assert.True(t, same)

	var diff = user.IsUserUCR(diffUserValue)(userCompanyRole)
	assert.False(t, diff)

	t.Logf("userCompanyRole: %v", userCompanyRole)
}

func TestSameCompanyUCR(t *testing.T) {
	var diffCompanyIdentifier, _ = pkgText.NewIdentifier("CP-TOSTOS")
	var diffCompany = getCompanyForCompanyRole(diffCompanyIdentifier)

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var userValue = getUserForCompanyRole(userIdentifier)

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var company = getCompanyForCompanyRole(companyIdentifier)

	var label01, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role01 = getRole(label01)

	var label02, _ = pkgText.NewLabel("TOST_ROLE_LABEL")
	var role02 = getRole(label02)

	var registerDate = time.Now()
	var expireDate = registerDate.Add(24 * time.Hour)

	var userCompanyRole1 = user.NewUserCompanyRole(1, userValue, company, role01, registerDate, &expireDate)
	var userCompanyRole2 = user.NewUserCompanyRole(2, userValue, company, role02, registerDate, &expireDate)
	var userCompanyRole3 = user.NewUserCompanyRole(3, userValue, diffCompany, role01, registerDate, &expireDate)

	assert.True(t, user.SameCompanyUCR(userCompanyRole1, userCompanyRole2))
	assert.False(t, user.SameCompanyUCR(userCompanyRole1, userCompanyRole3))

	t.Logf("userCompanyRole1: %v", userCompanyRole1)
	t.Logf("userCompanyRole2: %v", userCompanyRole2)
	t.Logf("userCompanyRole3: %v", userCompanyRole3)
}

func TestGetRoleUCR(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var userValue = getUserForCompanyRole(userIdentifier)

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var company = getCompanyForCompanyRole(companyIdentifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role = getRole(label)

	var registerDate = time.Now()
	var expireDate = registerDate.Add(24 * time.Hour)

	var userCompanyRole = user.NewUserCompanyRole(1, userValue, company, role, registerDate, &expireDate)

	assert.Equal(t, role.Label, user.GetRoleUCR(userCompanyRole).Label)

	t.Logf("userCompanyRole: %v", userCompanyRole)
}

func TestListToCompanyRole(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var userValue = getUserForCompanyRole(userIdentifier)

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var company = getCompanyForCompanyRole(companyIdentifier)

	var label1, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role1 = getRole(label1)

	var label2, _ = pkgText.NewLabel("TOST_ROLE_LABEL")
	var role2 = getRole(label2)

	var registerDate = time.Now()
	var expireDate = registerDate.Add(24 * time.Hour)

	var userCompanyRole1 = user.NewUserCompanyRole(1, userValue, company, role1, registerDate, &expireDate)
	var userCompanyRole2 = user.NewUserCompanyRole(2, userValue, company, role2, registerDate, &expireDate)

	var userCompanyRoles = []user.UserCompanyRole{*userCompanyRole1, *userCompanyRole2}
	var companyRole, err = user.ListToCompanyRole(userValue, userCompanyRoles)

	assert.Nil(t, err)
	assert.Equal(t, string(companyIdentifier), string(companyRole.Company.Identifier))
	assert.Equal(t, 2, len(companyRole.Roles))
	assert.Equal(t, string(label1), string(companyRole.Roles[0].Label))
	assert.Equal(t, string(label2), string(companyRole.Roles[1].Label))

	t.Logf("companyRole: %v", companyRole)
}

func TestListToCompanyRoleErrUser(t *testing.T) {
	var userIdentifier1, _ = pkgText.NewIdentifier("US-TESTES")
	var userValue1 = getUserForCompanyRole(userIdentifier1)

	var userIdentifier2, _ = pkgText.NewIdentifier("US-TOSTOS")
	var userValue2 = getUserForCompanyRole(userIdentifier2)

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var company = getCompanyForCompanyRole(companyIdentifier)

	var label1, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role1 = getRole(label1)

	var label2, _ = pkgText.NewLabel("TOST_ROLE_LABEL")
	var role2 = getRole(label2)

	var registerDate = time.Now()
	var expireDate = registerDate.Add(24 * time.Hour)

	var userCompanyRole1 = user.NewUserCompanyRole(1, userValue1, company, role1, registerDate, &expireDate)
	var userCompanyRole2 = user.NewUserCompanyRole(2, userValue2, company, role2, registerDate, &expireDate)

	var userCompanyRoles = []user.UserCompanyRole{*userCompanyRole1, *userCompanyRole2}
	var _, err = user.ListToCompanyRole(userValue1, userCompanyRoles)

	assert.Error(t, err)

	t.Logf("Error: %v", err)
}

func TestListToCompanyRoleErrCompany(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var userValue = getUserForCompanyRole(userIdentifier)

	var companyIdentifier1, _ = pkgText.NewIdentifier("CP-TESTES")
	var company1 = getCompanyForCompanyRole(companyIdentifier1)

	var companyIdentifier2, _ = pkgText.NewIdentifier("CP-TOSTOS")
	var company2 = getCompanyForCompanyRole(companyIdentifier2)

	var label1, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role1 = getRole(label1)

	var label2, _ = pkgText.NewLabel("TOST_ROLE_LABEL")
	var role2 = getRole(label2)

	var registerDate = time.Now()
	var expireDate = registerDate.Add(24 * time.Hour)

	var userCompanyRole1 = user.NewUserCompanyRole(1, userValue, company1, role1, registerDate, &expireDate)
	var userCompanyRole2 = user.NewUserCompanyRole(2, userValue, company2, role2, registerDate, &expireDate)

	var userCompanyRoles = []user.UserCompanyRole{*userCompanyRole1, *userCompanyRole2}
	var _, err = user.ListToCompanyRole(userValue, userCompanyRoles)

	assert.Error(t, err)

	t.Logf("Error: %v", err)
}
