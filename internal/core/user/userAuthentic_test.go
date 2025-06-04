package user_test

import (
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	pkgUser "github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getCompany(exposeId pkgText.ExposeId) company.Company {
	var companyId uint = 1
	var name, _ = pkgText.NewName("TestCompany")
	var registeredDate = time.Now()
	return company.NewCompany(companyId, exposeId, name, registeredDate)
}

func getRoles(label pkgText.Label) []role.Role {
	var roleName, _ = pkgText.NewName("TestRole")
	var description, _ = text.NewText("Role for testing")
	var registeredDate = time.Now()
	return []role.Role{role.NewRole(roleId, roleName, label, description, registeredDate)}
}

func TestNewCompanyRole(t *testing.T) {
	var exposeId, _ = pkgText.NewExposeId("CP-TESTES")
	var company = getCompany(exposeId)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roles = getRoles(label)

	var companyRole = user.NewCompanyRole(company, roles)

	assert.Equal(t, string(exposeId), string(companyRole.Company.ExposeId))
	assert.Equal(t, len(roles), len(companyRole.Roles))
	assert.Equal(t, string(label), companyRole.Roles[0].Label)

	t.Logf("companyRole: %+v", companyRole)
	t.Logf("company: %+v", companyRole.Company)
	t.Logf("company.ExposeId: %s", companyRole.Company.ExposeId)
	t.Logf("role: %+v", companyRole.Roles[0])
	t.Logf("role.Label: %s", companyRole.Roles[0].Label)
}

func TestNewUserAuthentic(t *testing.T) {
	var userId uint = 1
	var userExposeId = text.NewExposeId("TestExposeId")
	var emailId = text.NewEmail("test@gmail.com")
	var email = text.NewEmail("test_2@gmail.com")
	var userName = text.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()

	var companyExposeId, _ = pkgText.NewExposeId("CP-TESTES")
	var company = getCompany(companyExposeId)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roles = getRoles(label)

	var companyRole = user.NewCompanyRole(company, roles)

	var user = model.NewUser(userId, userExposeId, userName, emailId, &email, botFlag, userRegisteredDate, updateDate, &companyRole)

	assert.Equal(t, userId, user.UserId)
	assert.Equal(t, string(userExposeId), string(user.ExposeId))
	assert.Equal(t, string(emailId), string(user.ExposeEmailId))
	assert.Equal(t, string(email), string(*user.Email))
	assert.Equal(t, string(userName), string(user.Name))
	assert.Equal(t, botFlag, user.BotFlag)
	assert.Equal(t, userRegisteredDate, user.RegisteredDate)
	assert.Equal(t, updateDate, user.UpdateDate)
	assert.Equal(t, string(companyExposeId), string(user.CompanyRole.Company.ExposeId))
	assert.Equal(t, 1, len(user.CompanyRole.Roles))
	assert.Equal(t, string(label), string(user.CompanyRole.Roles[0].Label))

	t.Logf("user: %+v", user)
	t.Logf("user.UserId: %d", user.UserId)
	t.Logf("user.ExposeId: %s", user.ExposeId)
	t.Logf("user.ExposeEmailId: %s", user.ExposeEmailId)
	t.Logf("user.Email: %s", *user.Email)
	t.Logf("user.Name: %s", user.Name)
	t.Logf("user.BotFlag: %t", user.BotFlag)
	t.Logf("user.RegisteredDate: %t", user.RegisteredDate)
	t.Logf("user.UpdateDate: %t", user.UpdateDate)

	t.Logf("companyRole: %+v", user.CompanyRole)
	t.Logf("company: %+v", user.CompanyRole.Company)
	t.Logf("company.ExposeId: %d", user.CompanyRole.Company.ExposeId)
	t.Logf("role: %+v", user.CompanyRole.Roles[0])
	t.Logf("role.Label: %d", user.CompanyRole.Roles[0].Label)
}

// TODO working

func TestModelToAccessTokenUser(t *testing.T) {
	var userId uint = 1
	var userExposeId = "TestExposeId"
	var emailId = "test@gmail.com"
	var email = "test_2@gmail.com"
	var userName = "TestName"
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()

	var companyId uint = 1
	var companyExposeId = "CP-TESTES"
	var companyName = "TestRole"
	var companyRegisteredDate = time.Now()
	var companyRoles = []model.Role{}
	var company = model.NewCompany(companyId, companyExposeId, companyName, companyRegisteredDate, companyRoles)

	var roleId uint = 1
	var roleName = "TestRole"
	var roleLabel = "TEST_ROLE"
	var description = "Role for testing"
	var roleRegisteredDate = time.Now()
	var role = model.NewRole(roleId, roleName, roleLabel, description, roleRegisteredDate)

	var companyRole = model.NewCompanyRole(company, role)
	var modelUser = model.NewUser(userId, userExposeId, userName, emailId, &email, botFlag, userRegisteredDate, updateDate, &companyRole)

	var accessTokenUser = utility.ModelToAccessTokenUser(modelUser)

	assert.Equal(t, userExposeId, accessTokenUser.ExposeId)
	assert.Equal(t, emailId, accessTokenUser.ExposeEmailId)
	assert.Equal(t, email, *accessTokenUser.Email)
	assert.Equal(t, userName, accessTokenUser.Name)
	assert.Equal(t, botFlag, accessTokenUser.BotFlag)
	assert.Equal(t, updateDate, accessTokenUser.UpdateDate)
	assert.Equal(t, companyExposeId, accessTokenUser.Company.ExposeId)
	assert.Equal(t, companyName, accessTokenUser.Company.Name)
	assert.Equal(t, roleLabel, accessTokenUser.Company.Role)
	assert.Equal(t, roleName, accessTokenUser.Company.RoleName)

	t.Logf("user: %+v", accessTokenUser)
	t.Logf("user.ExposeId: %s", accessTokenUser.ExposeId)
	t.Logf("user.ExposeEmailId: %s", accessTokenUser.ExposeEmailId)
	t.Logf("user.Email: %s", *accessTokenUser.Email)
	t.Logf("user.Name: %s", accessTokenUser.Name)
	t.Logf("user.BotFlag: %t", accessTokenUser.BotFlag)
	t.Logf("user.UpdateDate: %t", accessTokenUser.UpdateDate)
	t.Logf("company: %+v", accessTokenUser.Company)
	t.Logf("company.ExposeId: %s", accessTokenUser.Company.ExposeId)
	t.Logf("company.Name: %s", accessTokenUser.Company.Name)
	t.Logf("company.Role: %s", accessTokenUser.Company.Role)
	t.Logf("company.RoleName: %s", accessTokenUser.Company.RoleName)
}
