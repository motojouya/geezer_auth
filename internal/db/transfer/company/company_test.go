package company_test

import (
	core "github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromCoreCompany(t *testing.T) {
	var identifier, _ = text.NewIdentifier("CP-TESTES")
	var name, _ = text.NewName("TestRole")
	var registeredDate = time.Now()

	var coreCompany = core.CreateCompany(identifier, name, registeredDate)

	var company = company.FromCoreCompany(coreCompany)

	assert.Equal(t, uint(0), company.PersistKey)
	assert.Equal(t, string(name), company.Name)
	assert.Equal(t, string(identifier), company.Identifier)
	assert.Equal(t, registeredDate, company.RegisteredDate)

	t.Logf("company: %+v", company)
}

func TestToCoreCompany(t *testing.T) {
	var identifier = "CP-TESTES"
	var name = "TestRole"
	var registeredDate = time.Now()

	var company = company.Company{
		PersistKey:     1,
		Identifier:     identifier,
		Name:           name,
		RegisteredDate: registeredDate,
	}

	coreCompany, err := company.ToCoreCompany()
	assert.NoError(t, err)

	assert.Equal(t, uint(1), coreCompany.PersistKey)
	assert.Equal(t, name, string(coreCompany.Name))
	assert.Equal(t, identifier, string(coreCompany.Identifier))
	assert.Equal(t, registeredDate, coreCompany.RegisteredDate)

	t.Logf("coreCompany: %+v", coreCompany)
}

func TestToCoreCompanyError(t *testing.T) {
	var company = company.Company{
		PersistKey:     1,
		Identifier:     "cp-testes",
		Name:           "TestRole",
		RegisteredDate: time.Now(),
	}

	var _, err = company.ToCoreCompany()
	assert.Error(t, err)

	t.Logf("error: %v", err)
}
