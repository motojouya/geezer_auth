package companyUser_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAcceptGetCompanyIdentifier(t *testing.T) {
	var identifier = "CP-TESTES"
	var token = "token-test"
	var companyUserAcceptRequest = companyUser.CompanyUserAcceptRequest{
		CompanyGetRequest: company.CompanyGetRequest{
			CompanyGet: company.CompanyGet{
				Identifier: identifier,
			},
		},
		CompanyAccept: companyUser.CompanyAccept{
			Token: token,
		},
	}

	var companyIdentifier, idErr = companyUserAcceptRequest.GetCompanyIdentifier()

	assert.Nil(t, idErr)
	assert.Equal(t, identifier, string(companyIdentifier))

	t.Logf("companyIdentifier: %s", companyIdentifier)

	var tokenResult, tokenErr = companyUserAcceptRequest.GetToken()
	assert.Nil(t, tokenErr)
	assert.Equal(t, token, string(tokenResult))
}

func TestAcceptGetCompanyIdentifierErr(t *testing.T) {
	var identifier = "cp-testes"
	var token = "token-test"
	var companyUserAcceptRequest = companyUser.CompanyUserAcceptRequest{
		CompanyGetRequest: company.CompanyGetRequest{
			CompanyGet: company.CompanyGet{
				Identifier: identifier,
			},
		},
		CompanyAccept: companyUser.CompanyAccept{
			Token: token,
		},
	}

	var _, err = companyUserAcceptRequest.GetCompanyIdentifier()

	assert.Error(t, err)

	t.Logf("error: %s", err)
}
