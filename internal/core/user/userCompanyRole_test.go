package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUser(identifier pkgText.Identifier) user.User {
	var userId = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return user.NewUser(userId, identifier, emailId, name, botFlag, registeredDate, updateDate)
}

func getCompany(identifier pkgText.Identifier) company.Company {
	var companyId = 1
	var name, _ = pkgText.NewName("TestCompany")
	var registeredDate = time.Now()

	return company.NewCompany(companyId, identifier, name, registeredDate)
}

func getRole(label pkgText.Label) user.Role {
	var name, _ = pkgText.NewName("TestRole")
	var description, _ = text.NewText("TestRoleDescription")
	var registeredDate = time.Now()

	return user.NewRole(roleId, identifier, name, registeredDate)
}

func TestCreateUserCompanyRole(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var user = getUser(userIdentifier)

	var companyIdentifier, _ = pkgText.NewIdentifier("TestCompanyIdentifier")
	var company = getCompany(companyIdentifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role = user.NewRole(label)

	var registerDate = time.Now()

	var userCompanyRole = user.CreateUserCompanyRole(user, company, role, registerDate)

	assert.Equal(t, string(identifier), string(userCompanyRole.User.Identifier))
	assert.Equal(t, string(companyIdentifier), string(userCompanyRole.Company.Identifier))
	assert.Equal(t, string(label), string(userCompanyRole.Role.Label))
	assert.Equal(t, registerDate, userCompanyRole.RegisteredDate)
	assert.Nil(t, *userCompanyRole.ExpireDate)

	t.Logf("userCompanyRole: %+v", userCompanyRole)
	t.Logf("userCompanyRole.User.Identifier: %s", userCompanyRole.User.Identifier)
	t.Logf("userCompanyRole.Company.Identifier: %s", userCompanyRole.Company.Identifier)
	t.Logf("userCompanyRole.Role.Label: %s", userCompanyRole.Role.Label)
	t.Logf("userCompanyRole.RegisteredDate: %s", userCompanyRole.RegisteredDate)
	t.Logf("userCompanyRole.ExpireDate: %s", *userCompanyRole.ExpireDate)
}

func TestNewUserCompanyRole(t *testing.T) {
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var user = getUser(userIdentifier)

	var companyIdentifier, _ = pkgText.NewIdentifier("TestCompanyIdentifier")
	var company = getCompany(companyIdentifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role = user.NewRole(label)

	var registerDate = time.Now()
	var expireDate = registerDate.Add(24 * time.Hour)

	var userCompanyRole = user.NewUserCompanyRole(1, user, company, role, registerDate, &expireDate)

	assert.Equal(t, 1, userCompanyRole.UserCompanyRoleID)
	assert.Equal(t, string(userIdentifier), string(userCompanyRole.User.Identifier))
	assert.Equal(t, string(companyIdentifier), string(userCompanyRole.Company.Identifier))
	assert.Equal(t, string(label), string(userCompanyRole.Role.Label))
	assert.Equal(t, registerDate, userCompanyRole.RegisteredDate)
	assert.Equal(t, expireDate, *userCompanyRole.ExpireDate)

	t.Logf("userCompanyRole: %+v", userCompanyRole)
	t.Logf("userCompanyRole.UserCompanyRoleID: %d", userCompanyRole.UserCompanyRoleID)
	t.Logf("userCompanyRole.User.Identifier: %s", userCompanyRole.User.Identifier)
	t.Logf("userCompanyRole.Company.Identifier: %s", userCompanyRole.Company.Identifier)
	t.Logf("userCompanyRole.Role.Label: %s", userCompanyRole.Role.Label)
	t.Logf("userCompanyRole.RegisteredDate: %s", userCompanyRole.RegisteredDate)
	t.Logf("userCompanyRole.ExpireDate: %s", *userCompanyRole.ExpireDate)
}
