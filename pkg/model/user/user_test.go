package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/accessToken"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewCompany(t *testing.T) {
	var exposeId = "CP-TESTES"
	var name = "TestCompany"
	var role = "TestRole"
	var roleName = "TestRoleName"

	var company = accessToken.CreateCompany(exposeId, name, role, roleName)

	assert.Equal(t, name, company.Name)
	assert.Equal(t, exposeId, company.ExposeId)
	assert.Equal(t, role, company.Role)
	assert.Equal(t, roleName, company.RoleName)

	t.Logf("company: %+v", company)
	t.Logf("company.ExposeId: %s", company.ExposeId)
	t.Logf("company.Name: %s", company.Name)
	t.Logf("company.Role: %s", company.Role)
	t.Logf("company.RoleName: %s", company.RoleName)
}

func TestNewUser(t *testing.T) {
	var companyExposeId = "CP-TESTES"
	var companyName = "TestCompany"
	var companyRole = "TestRole"
	var companyRoleName = "TestRoleName"

	var company = accessToken.CreateCompany(exposeId, name, role, roleName)

	var userExposeId = "TestExposeId"
	var emailId = "test@gmail.com"
	var email = "test_2@gmail.com"
	var userName = "TestName"
	var botFlag = false
	var updateDate = time.Now()

	var user = accessToken.NewUser(userExposeId, emailId, email, userName, botFlag, company, updateDate)

	assert.Equal(t, userExposeId, user.ExposeId)
	assert.Equal(t, emailId, user.ExposeEmailId)
	assert.Equal(t, email, *user.Email)
	assert.Equal(t, userName, user.Name)
	assert.Equal(t, botFlag, user.BotFlag)
	assert.Equal(t, updateDate, user.UpdateDate)
	assert.Equal(t, companyExposeId, user.Company.ExposeId)

	t.Logf("user: %+v", user)
	t.Logf("user.ExposeId: %s", user.ExposeId)
	t.Logf("user.ExposeEmailId: %s", user.ExposeEmailId)
	t.Logf("user.Email: %s", *user.Email)
	t.Logf("user.Name: %s", user.Name)
	t.Logf("user.BotFlag: %t", user.BotFlag)
	t.Logf("user.UpdateDate: %t", user.UpdateDate)
	t.Logf("company: %+v", user.Company)
	t.Logf("company.ExposeId: %s", user.Company.ExposeId)
}
