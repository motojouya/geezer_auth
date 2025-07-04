package company_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCompanyIdentifier(t *testing.T) {
	var identifier = "CP-TESTES"
	var request = company.CompanyGetRequest{
		CompanyGet: company.CompanyGet{
			Identifier: identifier,
		},
	}

	var companyIdentifier, err = request.GetCompanyIdentifier()

	assert.Nil(t, err)
	assert.Equal(t, identifier, string(companyIdentifier))

	t.Logf("companyIdentifier: %s", companyIdentifier)
}

func TestGetCompanyIdentifierErr(t *testing.T) {
	var identifier = "cp-testes"
	var request = company.CompanyGetRequest{
		CompanyGet: company.CompanyGet{
			Identifier: identifier,
		},
	}

	var _, err = request.GetCompanyIdentifier()

	assert.Error(t, err)

	t.Logf("error: %s", err)
}
