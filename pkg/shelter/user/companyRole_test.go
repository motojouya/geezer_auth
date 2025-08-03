package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/motojouya/geezer_auth/pkg/shelter/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCompanyRole(t *testing.T) {
	var companyIdentifier, _ = text.NewIdentifier("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyIdentifier, companyName)

	var roleLabel, _ = text.NewLabel("TestRole")
	var roleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(roleLabel, roleName)
	var roles = []user.Role{role}

	var companyRole = user.NewCompanyRole(company, roles)

	assert.Equal(t, string(companyIdentifier), string(companyRole.Company.Identifier))
	assert.Equal(t, len(roles), len(companyRole.Roles))
	assert.Equal(t, string(roleLabel), string(companyRole.Roles[0].Label))

	t.Logf("companyRole: %+v", companyRole)
	t.Logf("companyRole.Company.Identifier: %s", string(companyRole.Company.Identifier))
	t.Logf("companyRole.Role[0].Label: %s", string(companyRole.Roles[0].Label))
}
