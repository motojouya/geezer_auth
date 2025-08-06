package company_test

import (
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/company"
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromShelterCompany(t *testing.T) {
	var identifier, _ = text.NewIdentifier("CP-TESTES")
	var name, _ = text.NewName("TestRole")
	var registeredDate = time.Now()

	var shelterCompany = shelter.CreateCompany(identifier, name, registeredDate)

	var company = company.FromShelterCompany(shelterCompany)

	assert.Equal(t, uint(0), company.PersistKey)
	assert.Equal(t, string(name), company.Name)
	assert.Equal(t, string(identifier), company.Identifier)
	assert.Equal(t, registeredDate, company.RegisteredDate)

	t.Logf("company: %+v", company)
}

func TestToShelterCompany(t *testing.T) {
	var identifier = "CP-TESTES"
	var name = "TestRole"
	var registeredDate = time.Now()

	var company = company.Company{
		PersistKey:     1,
		Identifier:     identifier,
		Name:           name,
		RegisteredDate: registeredDate,
	}

	shelterCompany, err := company.ToShelterCompany()
	assert.NoError(t, err)

	assert.Equal(t, uint(1), shelterCompany.PersistKey)
	assert.Equal(t, name, string(shelterCompany.Name))
	assert.Equal(t, identifier, string(shelterCompany.Identifier))
	assert.Equal(t, registeredDate, shelterCompany.RegisteredDate)

	t.Logf("shelterCompany: %+v", shelterCompany)
}

func TestToShelterCompanyError(t *testing.T) {
	var company = company.Company{
		PersistKey:     1,
		Identifier:     "cp-testes",
		Name:           "TestRole",
		RegisteredDate: time.Now(),
	}

	var _, err = company.ToShelterCompany()
	assert.Error(t, err)

	t.Logf("error: %v", err)
}
