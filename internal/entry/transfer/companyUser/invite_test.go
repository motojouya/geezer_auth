package companyUser_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInviteGetCompanyIdentifier(t *testing.T) {
	var identifier = "CP-TESTES"
	var label = "LABEL_TEST"
	var companyUserInviteRequest = companyUser.CompanyUserInviteRequest{
		CompanyGetRequest: company.CompanyGetRequest{
			CompanyGet: company.CompanyGet{
				Identifier: identifier,
			},
		},
		RoleInvite: companyUser.RoleInvite{
			RoleLabel: label,
		},
	}

	var companyIdentifier, idErr = companyUserInviteRequest.GetCompanyIdentifier()

	assert.Nil(t, idErr)
	assert.Equal(t, identifier, string(companyIdentifier))

	t.Logf("companyIdentifier: %s", companyIdentifier)

	var roleLabel, labelErr = companyUserInviteRequest.GetRoleLabel()
	assert.Nil(t, labelErr)
	assert.Equal(t, label, string(roleLabel))
}

func TestInviteGetCompanyIdentifierErr(t *testing.T) {
	var identifier = "cp-testes"
	var label = "label-test"
	var companyUserInviteRequest = companyUser.CompanyUserInviteRequest{
		CompanyGetRequest: company.CompanyGetRequest{
			CompanyGet: company.CompanyGet{
				Identifier: identifier,
			},
		},
		RoleInvite: companyUser.RoleInvite{
			RoleLabel: label,
		},
	}

	var _, idErr = companyUserInviteRequest.GetCompanyIdentifier()

	assert.Error(t, idErr)

	t.Logf("error: %s", idErr)

	var _, labelErr = companyUserInviteRequest.GetRoleLabel()

	assert.Error(t, labelErr)

	t.Logf("error: %s", labelErr)
}
