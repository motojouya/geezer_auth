package model_test

import (
	"github.com/motojouya/geezer_auth/pkg/model"
	"github.com/motojouya/geezer_auth/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	var exposeId = "TestExposeId"
	var emailId = "test@gmail.com"
	var name = "TestName"
	var botFlag = false

	var user = model.CreateUser(exposeId, emailId, name, botFlag)

	assert.Equal(t, exposeId, user.ExposeId)
	assert.Equal(t, emailId, user.ExposeEmailId)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, botFlag, user.BotFlag)

	t.Logf("user: %+v", user)
	t.Logf("user.ExposeId: %s", user.ExposeId)
	t.Logf("user.ExposeEmailId: %s", user.ExposeEmailId)
	t.Logf("user.Name: %s", user.Name)
	t.Logf("user.BotFlag: %t", user.BotFlag)
}

func TestNewUser(t *testing.T) {
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
	var label = "TEST_ROLE"
	var description = "Role for testing"
	var roleRegisteredDate = time.Now()
	var role = model.NewRole(roleId, roleName, label, description, roleRegisteredDate)

	var companyRole = model.NewCompanyRole(company, role)

	var user = model.NewUser(userId, userExposeId, userName, emailId, &email, botFlag, userRegisteredDate, updateDate, &companyRole)

	assert.Equal(t, userId, user.UserId)
	assert.Equal(t, userExposeId, user.ExposeId)
	assert.Equal(t, emailId, user.ExposeEmailId)
	assert.Equal(t, email, *user.Email)
	assert.Equal(t, userName, user.Name)
	assert.Equal(t, botFlag, user.BotFlag)
	assert.Equal(t, userRegisteredDate, user.RegisteredDate)
	assert.Equal(t, updateDate, user.UpdateDate)
	assert.Equal(t, companyId, user.CompanyRole.Company.CompanyId)
	assert.Equal(t, roleId, user.CompanyRole.Role.RoleId)

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
	t.Logf("company.CompanyId: %d", user.CompanyRole.Company.CompanyId)
	t.Logf("role: %+v", user.CompanyRole.Role)
	t.Logf("role.RoleId: %d", user.CompanyRole.Role.RoleId)
}

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
