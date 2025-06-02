package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewCompanyRole(t *testing.T) {
	var companyExposeId, _ = text.NewExposeId("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyExposeId, companyName)

	var roleLabel, _ = text.NewLabel("TestRole")
	var roleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(roleLabel, roleName)
	var roles = []user.Roles{role}

	var companyRole = user.NewCompanyRole(company, roles)

	assert.Equal(t, string(companyExposeId), string(companyRole.Company.ExposeId))
	assert.Equal(t, len(roles), len(companyRole.Roles))
	assert.Equal(t, string(roleLabel), string(companyRole.Roles[0].Label))

	t.Logf("companyRole: %+v", companyRole)
	t.Logf("companyRole.Company.ExposeId: %s", string(companyRole.Company.ExposeId))
	t.Logf("companyRole.Role[0].Label: %s", string(companyRole.Roles[0].Label))
}
