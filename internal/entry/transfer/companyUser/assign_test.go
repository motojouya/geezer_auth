package companyUser_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssignGetCompanyIdentifier(t *testing.T) {
	var companyIdentifier = "CP-TESTES"
	var label = "LABEL_TEST"
	var userIdentifier = "US-TESTES"
	var companyUserAssignRequest = companyUser.CompanyUserAssignRequest{
		CompanyGetRequest: company.CompanyGetRequest{
			CompanyGet: company.CompanyGet{
				Identifier: companyIdentifier,
			},
		},
		RoleAssign: companyUser.RoleAssign{
			RoleInvite: companyUser.RoleInvite{
				RoleLabel: label,
			},
			UserIdentifier: userIdentifier,
		},
	}

	var companyIdentifierResult, cpIdErr = companyUserAssignRequest.GetCompanyIdentifier()

	assert.Nil(t, cpIdErr)
	assert.Equal(t, companyIdentifier, string(companyIdentifierResult))

	t.Logf("companyIdentifier: %s", string(companyIdentifierResult))

	var userIdentifierResult, usIdErr = companyUserAssignRequest.GetUserIdentifier()

	assert.Nil(t, usIdErr)
	assert.Equal(t, userIdentifier, string(userIdentifierResult))

	t.Logf("userIdentifier: %s", string(userIdentifierResult))

	var roleLabel, labelErr = companyUserAssignRequest.GetRoleLabel()

	assert.Nil(t, labelErr)
	assert.Equal(t, label, string(roleLabel))

	t.Logf("roleLabel: %s", string(roleLabel))
}

func TestAssignGetCompanyIdentifierErr(t *testing.T) {
	var companyIdentifier = "cp-testes"
	var label = "label-test"
	var userIdentifier = "us-testes"
	var companyUserAssignRequest = companyUser.CompanyUserAssignRequest{
		CompanyGetRequest: company.CompanyGetRequest{
			CompanyGet: company.CompanyGet{
				Identifier: companyIdentifier,
			},
		},
		RoleAssign: companyUser.RoleAssign{
			RoleInvite: companyUser.RoleInvite{
				RoleLabel: label,
			},
			UserIdentifier: userIdentifier,
		},
	}

	var _, cpIdErr = companyUserAssignRequest.GetCompanyIdentifier()

	assert.Error(t, cpIdErr)

	t.Logf("error: %s", cpIdErr)

	var _, usIdErr = companyUserAssignRequest.GetUserIdentifier()

	assert.Error(t, usIdErr)

	t.Logf("error: %s", usIdErr)

	var _, labelErr = companyUserAssignRequest.GetRoleLabel()

	assert.Error(t, labelErr)

	t.Logf("error: %s", labelErr)
}
