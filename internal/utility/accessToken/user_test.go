package accessToken_test

import (
	"github.com/motojouya/geezer_auth/pkg/accessToken"
	"github.com/motojouya/geezer_auth/internal/model"
	utility "github.com/motojouya/geezer_auth/internal/utility/accessToken"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestModelToAccessTokenUser(t *testing.T) {
	var userId uint = 1
	var userExposeId = "TestExposeId"
	var emailId = "test@gmail.com"
	var email = "test_2@gmail.com"
	var userName = "TestName"
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()

	var companyId uint = 1
	var companyExposeId = "CP-TESTES"
	var companyName = "TestRole"
	var companyRegisteredDate = time.Now()
	var companyRoles = []model.Role{}
	var company = model.NewCompany(companyId, companyExposeId, companyName, companyRegisteredDate, companyRoles)

	var roleId uint = 1
	var roleName = "TestRole"
	var roleLabel = "TEST_ROLE"
	var description = "Role for testing"
	var roleRegisteredDate = time.Now()
	var role = model.NewRole(roleId, roleName, roleLabel, description, roleRegisteredDate)

	var companyRole = model.NewCompanyRole(company, role)
	var modelUser = model.NewUser(userId, userExposeId, userName, emailId, &email, botFlag, userRegisteredDate, updateDate, &companyRole)

	var accessTokenUser = utility.ModelToAccessTokenUser(modelUser)

	assert.Equal(t, userExposeId, accessTokenUser.ExposeId)
	assert.Equal(t, emailId, accessTokenUser.ExposeEmailId)
	assert.Equal(t, email, *accessTokenUser.Email)
	assert.Equal(t, userName, accessTokenUser.Name)
	assert.Equal(t, botFlag, accessTokenUser.BotFlag)
	assert.Equal(t, updateDate, accessTokenUser.UpdateDate)
	assert.Equal(t, companyExposeId, accessTokenUser.Company.ExposeId)
	assert.Equal(t, companyName, accessTokenUser.Company.Name)
	assert.Equal(t, roleLabel, accessTokenUser.Company.Role)
	assert.Equal(t, roleName, accessTokenUser.Company.RoleName)

	t.Logf("user: %+v", accessTokenUser)
	t.Logf("user.ExposeId: %s", accessTokenUser.ExposeId)
	t.Logf("user.ExposeEmailId: %s", accessTokenUser.ExposeEmailId)
	t.Logf("user.Email: %s", *accessTokenUser.Email)
	t.Logf("user.Name: %s", accessTokenUser.Name)
	t.Logf("user.BotFlag: %t", accessTokenUser.BotFlag)
	t.Logf("user.UpdateDate: %t", accessTokenUser.UpdateDate)
	t.Logf("company: %+v", accessTokenUser.Company)
	t.Logf("company.ExposeId: %s", accessTokenUser.Company.ExposeId)
	t.Logf("company.Name: %s", accessTokenUser.Company.Name)
	t.Logf("company.Role: %s", accessTokenUser.Company.Role)
	t.Logf("company.RoleName: %s", accessTokenUser.Company.RoleName)
}
