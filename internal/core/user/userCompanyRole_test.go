package user_test

import (
	"github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
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
	assert.Nil(t, *userCompanyRole.ExpireDate)

	t.Logf("userCompanyRole: %+v", userCompanyRole)
	t.Logf("userCompanyRole.User.Identifier: %s", userCompanyRole.User.Identifier)
	t.Logf("userCompanyRole.Company.Identifier: %s", userCompanyRole.Company.Identifier)
	t.Logf("userCompanyRole.Role.Label: %s", userCompanyRole.Role.Label)
	t.Logf("userCompanyRole.RegisterDate: %s", userCompanyRole.RegisterDate)
	t.Logf("userCompanyRole.ExpireDate: %s", *userCompanyRole.ExpireDate)
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

	var userCompanyRole = user.NewUserCompanyRole(1, userValue, company, role, registerDate, expireDate)

	assert.Equal(t, 1, userCompanyRole.PersistKey)
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
