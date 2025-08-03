package company_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToCoreCompany(t *testing.T) {
	var name = "TestCompany"

	var request = company.CompanyCreateRequest{
		CompanyCreate: company.CompanyCreate{
			Name: name,
		},
	}

	var identifier, _ = text.NewIdentifier("CP-TESTES")
	var registerDate = time.Now()

	shelterCompany, err := request.ToCoreCompany(identifier, registerDate)
	assert.Nil(t, err)

	assert.Equal(t, identifier, shelterCompany.Identifier)
	assert.Equal(t, name, string(shelterCompany.Name))
	assert.Equal(t, registerDate, shelterCompany.RegisteredDate)
}
