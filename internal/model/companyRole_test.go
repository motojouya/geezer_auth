package model_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/motojouya/geezer_auth/internal/model"
)

func TestNewCompanyRole(t *testing.T) {
	var companyId uint = 1
	var exposeId = "CP-TESTES"
	var companyName = "TestRole"
	var companyRegisteredDate = time.Now()
	var companyRoles = []model.Role{}
	var company = model.NewCompany(companyId, exposeId, companyName, companyRegisteredDate, companyRoles)

	var roleId uint = 1
	var roleName = "TestRole"
	var label = "TEST_ROLE"
	var description = "Role for testing"
	var roleRegisteredDate = time.Now()
	var role = model.NewRole(roleId, roleName, label, description, roleRegisteredDate)

	var companyRole = model.NewCompanyRole(company, role)

	assert.Equal(t, companyId, companyRole.Company.CompanyId)
	assert.Equal(t, companyName, companyRole.Company.Name)
	assert.Equal(t, exposeId, companyRole.Company.ExposeId)
	assert.Equal(t, companyRegisteredDate, companyRole.Company.RegisteredDate)
	assert.Equal(t, len(companyRoles), len(companyRole.Company.Roles))

	assert.Equal(t, roleId, companyRole.Role.RoleId)
	assert.Equal(t, roleName, companyRole.Role.Name)
	assert.Equal(t, label, companyRole.Role.Label)
	assert.Equal(t, description, companyRole.Role.Description)
	assert.Equal(t, roleRegisteredDate, companyRole.Role.RegisteredDate)

	t.Logf("companyRole: %+v", companyRole)

	t.Logf("company: %+v", companyRole.Company)
	t.Logf("company.CompanyId: %d", companyRole.Company.CompanyId)
	t.Logf("company.ExposeId: %s", companyRole.Company.ExposeId)
	t.Logf("company.Name: %s", companyRole.Company.Name)
	t.Logf("company.RegisteredDate: %s", companyRole.Company.RegisteredDate)
	t.Logf("company.Roles.length: %d", len(companyRole.Company.Roles))

	t.Logf("role: %+v", companyRole.Role)
	t.Logf("role.RoleId: %d", companyRole.Role.RoleId)
	t.Logf("role.Name: %s", companyRole.Role.Name)
	t.Logf("role.Label: %s", companyRole.Role.Label)
	t.Logf("role.Description: %s", companyRole.Role.Description)
	t.Logf("role.RegisteredDate: %s", companyRole.Role.RegisteredDate)
}
