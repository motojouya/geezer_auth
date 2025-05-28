package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/model/text"
	"github.com/motojouya/geezer_auth/pkg/model/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getCompanyRole(companyExposeIdStr string, roleLabelStr string) user.CompanyRole {
	var companyExposeId, _ = text.NewExposeId(companyExposeIdStr)
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyExposeId, companyName)

	var roleLabel, _ = text.NewLabel(roleLabelStr)
	var roleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(roleLabel, roleName)
	var roles = []user.Roles{role}

	return user.NewCompanyRole(company, roles)
}

func TestNewUser(t *testing.T) {
	var companyExposeId = "CP-TESTES"
	var roleLabel = "TestRole"
	var companyRole = getCompanyRole(companyExposeId, roleLabel)

	var userExposeId = text.NewExposeId("TestExposeId")
	var emailId = text.NewEmail("test@gmail.com")
	var email = text.NewEmail("test_2@gmail.com")
	var userName = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var user = user.NewUser(userExposeId, emailId, email, userName, botFlag, companyRole, updateDate)

	assert.Equal(t, string(userExposeId), string(user.ExposeId))
	assert.Equal(t, string(emailId), string(user.ExposeEmailId))
	assert.Equal(t, string(email), string(*user.Email))
	assert.Equal(t, string(userName), string(user.Name))
	assert.Equal(t, botFlag, user.BotFlag)
	assert.Equal(t, updateDate, user.UpdateDate)

	assert.Equal(t, companyExposeId, string(companyRole.Company.ExposeId))
	assert.Equal(t, len(roles), len(companyRole.Roles))
	assert.Equal(t, roleLabel, string(companyRole.Roles[0].Label))

	t.Logf("user: %+v", user)
	t.Logf("user.ExposeId: %s", user.ExposeId)
	t.Logf("user.ExposeEmailId: %s", user.ExposeEmailId)
	t.Logf("user.Email: %s", *user.Email)
	t.Logf("user.Name: %s", user.Name)
	t.Logf("user.BotFlag: %t", user.BotFlag)
	t.Logf("user.UpdateDate: %t", user.UpdateDate)

	t.Logf("user.CompanyRole: %+v", user.CompanyRole)
	t.Logf("user.CompanyRole.Company.ExposeId: %s", string(user.CompanyRole.Company.ExposeId))
	t.Logf("user.CompanyRole.Role[0].Label: %s", string(user.CompanyRole.Roles[0].Label))
}
