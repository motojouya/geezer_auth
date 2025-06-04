package company_test

import (
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getCompany(exposeId pkgText.ExposeId) *company.Company {
	var name, _ = pkgText.NewName("TestRole")
	var registeredDate = time.Now()

	return company.NewCompany(1, exposeId, name, registeredDate)
}

func getRole(label pkgText.Label) *role.Role {
	var name, _ = pkgText.NewName("TestRole")
	var registeredDate = time.Now()

	return role.NewRole(label, name, registeredDate)
}

func TestCreateCompanyInvite(t *testing.T) {
	var exposeId, _ = text.NewExposeId("CP-TESTES")
	var company = getCompany(exposeId)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var role = getRole(label)

	var token, _ = uuid.NewUUID()
	var registeredDate = time.Now()

	var companyInvite = company.CreateCompanyInvite(company, token, role, registeredDate)

	assert.Equal(t, string(label), string(companyInvite.Role.Label))
	assert.Equal(t, string(exposeId), string(companyInvite.Company.ExposeId))
	assert.Equal(t, registeredDate, companyInvite.RegisteredDate)
	assert.Equal(t, expireDate, companyInvite.ExpireDate)

	t.Logf("companyInvite: %+v", companyInvite)
	t.Logf("companyInvite.Company.ExposeId: %s", companyInvite.Company.ExposeId)
	t.Logf("companyInvite.Role.Label: %s", companyInvite.Role.Label)
	t.Logf("companyInvite.RegisteredDate: %s", companyInvite.RegisteredDate)
	t.Logf("companyInvite.ExpireDate: %s", companyInvite.ExpireDate)
}

func TestNewCompanyInvite(t *testing.T) {
	var companyInviteId uint = 1

	var exposeId, _ = text.NewExposeId("CP-TESTES")
	var company = getCompany(exposeId)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var role = getRole(label)

	var token, _ = uuid.NewUUID()
	var registeredDate = time.Now()
	var expireDate = registeredDate.Add(company.TokenValidityPeriodHours * time.Hour)

	var companyInvite = company.NewCompanyInvite(companyInviteId, company, token, role, registeredDate, expireDate)

	assert.Equal(t, companyInviteId, companyInvite.CompanyInviteId)
	assert.Equal(t, string(label), string(companyInvite.Role.Label))
	assert.Equal(t, string(exposeId), string(companyInvite.Company.ExposeId))
	assert.Equal(t, registeredDate, companyInvite.RegisteredDate)
	assert.Equal(t, expireDate, companyInvite.ExpireDate)

	t.Logf("companyInvite: %+v", companyInvite)
	t.Logf("companyInvite.Company.CompanyId: %d", companyInvite.Company.CompanyId)
	t.Logf("companyInvite.Company.ExposeId: %s", companyInvite.Company.ExposeId)
	t.Logf("companyInvite.Role.Label: %s", companyInvite.Role.Label)
	t.Logf("companyInvite.RegisteredDate: %s", companyInvite.RegisteredDate)
	t.Logf("companyInvite.ExpireDate: %s", companyInvite.ExpireDate)
}
