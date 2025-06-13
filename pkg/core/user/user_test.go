package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getCompanyRole(companyIdentifierStr string, roleLabelStr string) *user.CompanyRole {
	var companyIdentifier, _ = text.NewIdentifier(companyIdentifierStr)
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyIdentifier, companyName)

	var roleLabel, _ = text.NewLabel(roleLabelStr)
	var roleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(roleLabel, roleName)
	var roles = []user.Role{role}

	return user.NewCompanyRole(company, roles)
}

func TestNewUser(t *testing.T) {
	var companyIdentifier = "CP-TESTES"
	var roleLabel = "TestRole"
	var companyRole = getCompanyRole(companyIdentifier, roleLabel)

	var userIdentifier, _ = text.NewIdentifier("TestIdentifier")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var email, _ = text.NewEmail("test_2@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var user = user.NewUser(userIdentifier, emailId, &email, userName, botFlag, companyRole, updateDate)

	assert.Equal(t, string(userIdentifier), string(user.Identifier))
	assert.Equal(t, string(emailId), string(user.EmailId))
	assert.Equal(t, string(email), string(*user.Email))
	assert.Equal(t, string(userName), string(user.Name))
	assert.Equal(t, botFlag, user.BotFlag)
	assert.Equal(t, updateDate, user.UpdateDate)

	assert.Equal(t, companyIdentifier, string(companyRole.Company.Identifier))
	assert.Equal(t, 1, len(companyRole.Roles))
	assert.Equal(t, roleLabel, string(companyRole.Roles[0].Label))

	t.Logf("user: %+v", user)
	t.Logf("user.Identifier: %s", user.Identifier)
	t.Logf("user.ExposeEmailId: %s", user.EmailId)
	t.Logf("user.Email: %s", *user.Email)
	t.Logf("user.Name: %s", user.Name)
	t.Logf("user.BotFlag: %t", user.BotFlag)
	t.Logf("user.UpdateDate: %s", user.UpdateDate)

	t.Logf("user.CompanyRole: %+v", user.CompanyRole)
	t.Logf("user.CompanyRole.Company.Identifier: %s", string(user.CompanyRole.Company.Identifier))
	t.Logf("user.CompanyRole.Role[0].Label: %s", string(user.CompanyRole.Roles[0].Label))
}
