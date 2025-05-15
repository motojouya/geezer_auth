package model_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/motojouya/geezer_auth/internal/model"
)

func TestCreateCompany(t *testing.T) {
	var exposeId = "CP-TESTES"
	var name = "TestRole"

	var company = model.CreateCompany(exposeId, name)

	assert.Equal(t, name, company.Name)
	assert.Equal(t, exposeId, company.ExposeId)

	t.Logf("company: %+v", company)
	t.Logf("company.ExposeId: %s", company.ExposeId)
	t.Logf("company.Name: %s", company.Name)
}

func TestNewCompany(t *testing.T) {
	var companyId uint = 1
	var exposeId = "CP-TESTES"
	var name = "TestRole"
	var registeredDate = time.Now()

	var roleId uint = 1
	var role = model.NewRole(roleId, "TestRole", "TEST_ROLE", "Role for testing", registeredDate)
	var roles = []model.Role{role}

	var company = model.NewCompany(companyId, exposeId, name, registeredDate, roles)

	assert.Equal(t, companyId, company.CompanyId)
	assert.Equal(t, name, company.Name)
	assert.Equal(t, exposeId, company.ExposeId)
	assert.Equal(t, registeredDate, company.RegisteredDate)

	assert.Equal(t, len(roles), len(company.Roles))
	assert.Equal(t, roleId, company.Roles[0].RoleId)

	t.Logf("company: %+v", company)
	t.Logf("company.CompanyId: %d", company.CompanyId)
	t.Logf("company.ExposeId: %s", company.ExposeId)
	t.Logf("company.Name: %s", company.Name)
	t.Logf("company.RegisteredDate: %s", company.RegisteredDate)

	t.Logf("company.Roles.length: %d", len(company.Roles))
	t.Logf("company.Roles[0].RoleId: %d", company.Roles[0].RoleId)
}
