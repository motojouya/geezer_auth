package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUser(exposeId pkgText.ExposeId) user.User {
	var userId = 1
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return user.NewUser(userId, exposeId, emailId, name, botFlag, registeredDate, updateDate)
}

func getCompany(exposeId pkgText.ExposeId) company.Company {
	var companyId = 1
	var name, _ = pkgText.NewName("TestCompany")
	var registeredDate = time.Now()

	return company.NewCompany(companyId, exposeId, name, registeredDate)
}

func getRole(label pkgText.Label) user.Role {
	var name, _ = pkgText.NewName("TestRole")
	var description, _ = text.NewText("TestRoleDescription")
	var registeredDate = time.Now()

	return user.NewRole(roleId, exposeId, name, registeredDate)
}

func TestCreateUserCompanyRole(t *testing.T) {
	var userExposeId, _ = pkgText.NewExposeId("TestExposeId")
	var user = getUser(userExposeId)

	var companyExposeId, _ = pkgText.NewExposeId("TestCompanyExposeId")
	var company = getCompany(companyExposeId)

	var label, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role = user.NewRole(label)

	var registerDate = time.Now()

	var userCompanyRole = user.CreateUserCompanyRole(user, company, role, registerDate)

	assert.Equal(t, string(exposeId), string(userCompanyRole.User.ExposeId))
	assert.Equal(t, string(companyExposeId), string(userCompanyRole.Company.ExposeId))
	assert.Equal(t, string(label), string(userCompanyRole.Role.Label))
	assert.Equal(t, registerDate, userCompanyRole.RegisteredDate)
	assert.Nil(t, *userCompanyRole.ExpireDate)

	t.Logf("userCompanyRole: %+v", userCompanyRole)
	t.Logf("userCompanyRole.User.ExposeId: %s", userCompanyRole.User.ExposeId)
	t.Logf("userCompanyRole.Company.ExposeId: %s", userCompanyRole.Company.ExposeId)
	t.Logf("userCompanyRole.Role.Label: %s", userCompanyRole.Role.Label)
	t.Logf("userCompanyRole.RegisteredDate: %s", userCompanyRole.RegisteredDate)
	t.Logf("userCompanyRole.ExpireDate: %s", *userCompanyRole.ExpireDate)
}

func TestNewUserCompanyRole(t *testing.T) {
	var userExposeId, _ = pkgText.NewExposeId("TestExposeId")
	var user = getUser(userExposeId)

	var companyExposeId, _ = pkgText.NewExposeId("TestCompanyExposeId")
	var company = getCompany(companyExposeId)

	var label, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role = user.NewRole(label)

	var registerDate = time.Now()
	var expireDate = registerDate.Add(24 * time.Hour)

	var userCompanyRole = user.NewUserCompanyRole(1, user, company, role, registerDate, &expireDate)

	assert.Equal(t, 1, userCompanyRole.UserCompanyRoleID)
	assert.Equal(t, string(userExposeId), string(userCompanyRole.User.ExposeId))
	assert.Equal(t, string(companyExposeId), string(userCompanyRole.Company.ExposeId))
	assert.Equal(t, string(label), string(userCompanyRole.Role.Label))
	assert.Equal(t, registerDate, userCompanyRole.RegisteredDate)
	assert.Equal(t, expireDate, *userCompanyRole.ExpireDate)

	t.Logf("userCompanyRole: %+v", userCompanyRole)
	t.Logf("userCompanyRole.UserCompanyRoleID: %d", userCompanyRole.UserCompanyRoleID)
	t.Logf("userCompanyRole.User.ExposeId: %s", userCompanyRole.User.ExposeId)
	t.Logf("userCompanyRole.Company.ExposeId: %s", userCompanyRole.Company.ExposeId)
	t.Logf("userCompanyRole.Role.Label: %s", userCompanyRole.Role.Label)
	t.Logf("userCompanyRole.RegisteredDate: %s", userCompanyRole.RegisteredDate)
	t.Logf("userCompanyRole.ExpireDate: %s", *userCompanyRole.ExpireDate)
}
