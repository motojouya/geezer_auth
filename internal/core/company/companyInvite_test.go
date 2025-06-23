package company_test

import (
	"github.com/google/uuid"
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getCompany(identifier pkgText.Identifier) company.Company {
	var name, _ = pkgText.NewName("TestRole")
	var registeredDate = time.Now()

	return company.NewCompany(1, identifier, name, registeredDate)
}

func getRole(label pkgText.Label) role.Role {
	var name, _ = pkgText.NewName("TestRole")
	var registeredDate = time.Now()
	var description, _ = text.NewText("This is a test role")

	return role.NewRole(name, label, description, registeredDate)
}

func TestCreateCompanyInvite(t *testing.T) {
	var identifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var companyValue = getCompany(identifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var role = getRole(label)

	var tokenUUID, _ = uuid.NewUUID()
	var token, _ = text.CreateToken(tokenUUID)
	var registerDate = time.Now()
	var expireDate = registerDate.Add(50 * time.Hour)

	var companyInvite = company.CreateCompanyInvite(companyValue, token, role, registerDate)

	assert.Equal(t, string(label), string(companyInvite.Role.Label))
	assert.Equal(t, string(identifier), string(companyInvite.Company.Identifier))
	assert.Equal(t, registerDate, companyInvite.RegisterDate)
	assert.Equal(t, expireDate, companyInvite.ExpireDate)

	t.Logf("companyInvite: %+v", companyInvite)
	t.Logf("companyInvite.Company.Identifier: %s", companyInvite.Company.Identifier)
	t.Logf("companyInvite.Role.Label: %s", companyInvite.Role.Label)
	t.Logf("companyInvite.RegisteredDate: %s", companyInvite.RegisterDate)
	t.Logf("companyInvite.ExpireDate: %s", companyInvite.ExpireDate)
}

func TestNewCompanyInvite(t *testing.T) {
	var companyInviteId uint = 1

	var identifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var companyValue = getCompany(identifier)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var role = getRole(label)

	var tokenUUID, _ = uuid.NewUUID()
	var token, _ = text.CreateToken(tokenUUID)
	var registeredDate = time.Now()
	var expireDate = registeredDate.Add(50 * time.Hour)

	var companyInvite = company.NewCompanyInvite(companyInviteId, companyValue, token, role, registeredDate, expireDate)

	assert.Equal(t, companyInviteId, companyInvite.PersistKey)
	assert.Equal(t, string(label), string(companyInvite.Role.Label))
	assert.Equal(t, string(identifier), string(companyInvite.Company.Identifier))
	assert.Equal(t, registeredDate, companyInvite.RegisterDate)
	assert.Equal(t, expireDate, companyInvite.ExpireDate)

	t.Logf("companyInvite: %+v", companyInvite)
	t.Logf("companyInvite.Company.CompanyId: %d", companyInvite.Company.PersistKey)
	t.Logf("companyInvite.Company.Identifier: %s", companyInvite.Company.Identifier)
	t.Logf("companyInvite.Role.Label: %s", companyInvite.Role.Label)
	t.Logf("companyInvite.RegisteredDate: %s", companyInvite.RegisterDate)
	t.Logf("companyInvite.ExpireDate: %s", companyInvite.ExpireDate)
}
