package company_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateCompanyExposeId(t *testing.T) {
	var exposeId, err = company.CreateCompanyExposeId("TESTES")

	assert.Nil(t, err)
	assert.NotEmpty(t, exposeId)
	assert.Equal(t, "CP-TESTES", string(exposeId))

	t.Logf("exposeId: %s", exposeId)
}

func TestCreateCompany(t *testing.T) {
	var exposeId, _ = text.NewExposeId("CP-TESTES")
	var name, _ = text.NewExposeId("TestRole")
	var registeredDate = time.Now()

	var company = company.CreateCompany(exposeId, name, registeredDate)

	assert.Equal(t, string(name), string(company.Name))
	assert.Equal(t, string(exposeId), string(company.ExposeId))
	assert.Equal(t, registeredDate, company.RegisteredDate)

	t.Logf("company: %+v", company)
	t.Logf("company.ExposeId: %s", company.ExposeId)
	t.Logf("company.Name: %s", company.Name)
	t.Logf("company.RegisteredDate: %s", company.RegisteredDate)
}

func TestNewCompany(t *testing.T) {
	var companyId uint = 1
	var exposeId, _ = text.NewExposeId("CP-TESTES")
	var name, _ = text.NewName("TestRole")
	var registeredDate = time.Now()

	var company = company.NewCompany(companyId, exposeId, name, registeredDate)

	assert.Equal(t, companyId, company.CompanyId)
	assert.Equal(t, string(name), string(company.Name))
	assert.Equal(t, string(exposeId), string(company.ExposeId))
	assert.Equal(t, registeredDate, company.RegisteredDate)

	t.Logf("company: %+v", company)
	t.Logf("company.CompanyId: %d", company.CompanyId)
	t.Logf("company.ExposeId: %s", company.ExposeId)
	t.Logf("company.Name: %s", company.Name)
	t.Logf("company.RegisteredDate: %s", company.RegisteredDate)
}
