package company_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateCompanyIdentifier(t *testing.T) {
	var identifier, err = company.CreateCompanyIdentifier("TESTES")

	assert.Nil(t, err)
	assert.NotEmpty(t, identifier)
	assert.Equal(t, "CP-TESTES", string(identifier))

	t.Logf("identifier: %s", identifier)
}

func TestCreateCompany(t *testing.T) {
	var identifier, _ = text.NewIdentifier("CP-TESTES")
	var name, _ = text.NewIdentifier("TestRole")
	var registeredDate = time.Now()

	var company = company.CreateCompany(identifier, name, registeredDate)

	assert.Equal(t, string(name), string(company.Name))
	assert.Equal(t, string(identifier), string(company.Identifier))
	assert.Equal(t, registeredDate, company.RegisteredDate)

	t.Logf("company: %+v", company)
	t.Logf("company.Identifier: %s", company.Identifier)
	t.Logf("company.Name: %s", company.Name)
	t.Logf("company.RegisteredDate: %s", company.RegisteredDate)
}

func TestNewCompany(t *testing.T) {
	var companyId uint = 1
	var identifier, _ = text.NewIdentifier("CP-TESTES")
	var name, _ = text.NewName("TestRole")
	var registeredDate = time.Now()

	var company = company.NewCompany(companyId, identifier, name, registeredDate)

	assert.Equal(t, companyId, company.CompanyId)
	assert.Equal(t, string(name), string(company.Name))
	assert.Equal(t, string(identifier), string(company.Identifier))
	assert.Equal(t, registeredDate, company.RegisteredDate)

	t.Logf("company: %+v", company)
	t.Logf("company.CompanyId: %d", company.CompanyId)
	t.Logf("company.Identifier: %s", company.Identifier)
	t.Logf("company.Name: %s", company.Name)
	t.Logf("company.RegisteredDate: %s", company.RegisteredDate)
}
