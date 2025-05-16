package model_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/motojouya/geezer_auth/internal/model"
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

	var user = model.NewUser(userId, userExposeId, userName, emailId, &email, botFlag, userRegisteredDate, &companyRole)

	assert.Equal(t, userId, user.UserId)
	assert.Equal(t, userExposeId, user.ExposeId)
	assert.Equal(t, emailId, user.ExposeEmailId)
	assert.Equal(t, email, *user.Email)
	assert.Equal(t, userName, user.Name)
	assert.Equal(t, botFlag, user.BotFlag)
	assert.Equal(t, companyId, user.CompanyRole.Company.CompanyId)
	assert.Equal(t, roleId, user.CompanyRole.Role.RoleId)

	t.Logf("user: %+v", user)
	t.Logf("user.UserId: %d", user.UserId)
	t.Logf("user.ExposeId: %s", user.ExposeId)
	t.Logf("user.ExposeEmailId: %s", user.ExposeEmailId)
	t.Logf("user.Email: %s", *user.Email)
	t.Logf("user.Name: %s", user.Name)
	t.Logf("user.BotFlag: %t", user.BotFlag)

	t.Logf("companyRole: %+v", user.CompanyRole)
	t.Logf("company: %+v", user.CompanyRole.Company)
	t.Logf("company.CompanyId: %d", user.CompanyRole.Company.CompanyId)
	t.Logf("role: %+v", user.CompanyRole.Role)
	t.Logf("role.RoleId: %d", user.CompanyRole.Role.RoleId)
}
