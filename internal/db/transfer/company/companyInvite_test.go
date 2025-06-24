package company_test

import (
	"github.com/google/uuid"
	core "github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromCoreCompanyInvite(t *testing.T) {
	var persistKey uint = 1
	var companyValue = getCompany(persistKey)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var role = getRole(label)

	var tokenUUID, _ = uuid.NewUUID()
	var token, _ = text.CreateToken(tokenUUID)
	var registerDate = time.Now()
	var expireDate = registerDate.Add(50 * time.Hour)

	var coreCompanyInvite = core.NewCompanyInvite(persistKey, companyValue, token, role, registerDate, expireDate)
	var companyInvite = company.FromCoreCompanyInvite(coreCompanyInvite)

	assert.Equal(t, uint(0), companyInvite.PersistKey)
	assert.Equal(t, string(label), companyInvite.RoleLabel)
	assert.Equal(t, persistKey, companyInvite.CompanyPersistKey)
	assert.Equal(t, string(token), companyInvite.Token)
	assert.Equal(t, registerDate, companyInvite.RegisterDate)
	assert.Equal(t, expireDate, companyInvite.ExpireDate)

	t.Logf("companyInvite: %+v", companyInvite)
}

func TestToCoreCompanyInvite(t *testing.T) {
	var companyInvitePersistKey uint = 1

	var companyPersistKey uint = 1
	var companyValue = getCompany(companyPersistKey)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var role = getRole(label)

	var tokenUUID, _ = uuid.NewUUID()
	var token, _ = text.CreateToken(tokenUUID)
	var registerDate = time.Now()
	var expireDate = registerDate.Add(50 * time.Hour)

	var companyInviteFull = company.CompanyInviteFull{
		CompanyInvite: company.CompanyInvite{
			PersistKey:        companyInvitePersistKey,
			CompanyPersistKey: companyPersistKey,
			Token:             string(token),
			RoleLabel:         string(label),
			RegisterDate:      registerDate,
			ExpireDate:        expireDate,
		},
		CompanyIdentifier:     string(companyValue.Identifier),
		CompanyName:           string(companyValue.Name),
		CompanyRegisteredDate: companyValue.RegisteredDate,
		RoleName:              string(role.Name),
		RoleDescription:       string(role.Description),
		RoleRegisteredDate:    role.RegisteredDate,
	}

	var coreCompanyInvite, err = companyInviteFull.ToCoreCompanyInvite()

	assert.Nil(t, err)
	assert.Equal(t, companyValue.Identifier, coreCompanyInvite.Company.Identifier)
	assert.Equal(t, label, coreCompanyInvite.Role.Label)
	assert.Equal(t, token, coreCompanyInvite.Token)
	assert.Equal(t, registerDate, coreCompanyInvite.RegisterDate)
	assert.Equal(t, expireDate, coreCompanyInvite.ExpireDate)

	t.Logf("coreCompanyInvite: %+v", coreCompanyInvite)
}

func TestToCoreCompanyInviteError(t *testing.T) {
	var companyInvitePersistKey uint = 1

	var companyPersistKey uint = 1
	var companyValue = getCompany(companyPersistKey)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var role = getRole(label)

	var tokenUUID, _ = uuid.NewUUID()
	var token, _ = text.CreateToken(tokenUUID)
	var registerDate = time.Now()
	var expireDate = registerDate.Add(50 * time.Hour)

	var companyInviteFull = company.CompanyInviteFull{
		CompanyInvite: company.CompanyInvite{
			PersistKey:        companyInvitePersistKey,
			CompanyPersistKey: companyPersistKey,
			Token:             string(token),
			RoleLabel:         "invalid_label",
			RegisterDate:      registerDate,
			ExpireDate:        expireDate,
		},
		CompanyIdentifier:     string(companyValue.Identifier),
		CompanyName:           string(companyValue.Name),
		CompanyRegisteredDate: companyValue.RegisteredDate,
		RoleName:              string(role.Name),
		RoleDescription:       string(role.Description),
		RoleRegisteredDate:    role.RegisteredDate,
	}

	var _, err = companyInviteFull.ToCoreCompanyInvite()

	assert.NotNil(t, err)

	t.Logf("coreCompanyInvite error: %v", err)
}

func getCompany(persistKey uint) core.Company {
	var identifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var name, _ = pkgText.NewName("TestRole")
	var registeredDate = time.Now()

	return core.NewCompany(persistKey, identifier, name, registeredDate)
}

func getRole(label pkgText.Label) role.Role {
	var name, _ = pkgText.NewName("TestRole")
	var registeredDate = time.Now()
	var description, _ = text.NewText("This is a test role")

	return role.NewRole(name, label, description, registeredDate)
}
