package company_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	"github.com/motojouya/geezer_auth/pkg/core/text"
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

	coreCompany, err := request.ToCoreCompany(identifier, registerDate)
	assert.Nil(t, err)

	assert.Equal(t, identifier, coreCompany.Identifier)
	assert.Equal(t, name, string(coreCompany.Name))
	assert.Equal(t, registerDate, coreCompany.RegisteredDate)
}
