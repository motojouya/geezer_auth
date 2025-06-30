package user_test

import (
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForCompanyRole(persistKey uint) core.User {
	var identifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()
	var updateDate = time.Now()

	return core.NewUser(persistKey, identifier, name, emailId, botFlag, registeredDate, updateDate)
}

func getCompanyForCompanyRole(persistKey uint) company.Company {
	var identifier, _ = pkgText.NewIdentifier("TestCompanyIdentifier")
	var name, _ = pkgText.NewName("TestCompany")
	var registeredDate = time.Now()

	return company.NewCompany(persistKey, identifier, name, registeredDate)
}

func getRole(label pkgText.Label) role.Role {
	var name, _ = pkgText.NewName("TestRole")
	var description, _ = text.NewText("TestRoleDescription")
	var registeredDate = time.Now()

	return role.NewRole(name, label, description, registeredDate)
}

func TestFromCoreUserCompanyRole(t *testing.T) {
	var userPersistKey uint = 1
	var userValue = getUserForCompanyRole(userPersistKey)

	var companyPersistKey uint = 1
	var company = getCompanyForCompanyRole(companyPersistKey)

	var label, _ = pkgText.NewLabel("TEST_ROLE_LABEL")
	var role = getRole(label)

	var registerDate = time.Now()

	var coreUserCompanyRole = core.CreateUserCompanyRole(userValue, company, role, registerDate)
	var userCompanyRole = user.FromCoreUserCompanyRole(coreUserCompanyRole)

	assert.Equal(t, uint(0), userCompanyRole.PersistKey)
	assert.Equal(t, userPersistKey, userCompanyRole.UserPersistKey)
	assert.Equal(t, companyPersistKey, userCompanyRole.CompanyPersistKey)
	assert.Equal(t, string(role.Label), userCompanyRole.RoleLabel)
	assert.Equal(t, registerDate, userCompanyRole.RegisterDate)
	assert.Nil(t, userCompanyRole.ExpireDate)

	t.Logf("userCompanyRole: %+v", userCompanyRole)
}

func TestToCoreUserCompanyRole(t *testing.T) {
	var now = time.Now()
	var expireDate = now.Add(1 * time.Hour)
	var userCompanyRoleFull = &user.UserCompanyRoleFull{
		UserCompanyRole: user.UserCompanyRole{
			PersistKey:        1,
			UserPersistKey:    2,
			CompanyPersistKey: 3,
			RoleLabel:         "TEST_ROLE_LABEL",
			RegisterDate:      now,
			ExpireDate:        &expireDate,
		},
		UserIdentifier:        "US-TESTES",
		UserExposeEmailId:     "test02@example.com",
		UserName:              "TestUserName",
		UserBotFlag:           false,
		UserRegisteredDate:    now.Add(2 * time.Hour),
		UserUpdateDate:        now.Add(3 * time.Hour),
		CompanyIdentifier:     "CP-TESTES",
		CompanyName:           "TestCompanyName",
		CompanyRegisteredDate: now.Add(4 * time.Hour),
		RoleName:              "TestRoleName",
		RoleDescription:       "TestRoleDescription",
		RoleRegisteredDate:    now.Add(5 * time.Hour),
	}

	var coreUserCompanyRole, err = userCompanyRoleFull.ToCoreUserCompanyRole()

	assert.NoError(t, err)
	assert.Equal(t, userCompanyRoleFull.PersistKey, coreUserCompanyRole.PersistKey)
	assert.Equal(t, userCompanyRoleFull.UserPersistKey, coreUserCompanyRole.User.PersistKey)
	assert.Equal(t, userCompanyRoleFull.CompanyPersistKey, coreUserCompanyRole.Company.PersistKey)
	assert.Equal(t, userCompanyRoleFull.RoleLabel, string(coreUserCompanyRole.Role.Label))
	assert.Equal(t, userCompanyRoleFull.RegisterDate, coreUserCompanyRole.RegisterDate)
	assert.Equal(t, userCompanyRoleFull.ExpireDate, coreUserCompanyRole.ExpireDate)

	t.Logf("coreUserCompanyRole: %+v", coreUserCompanyRole)
	t.Logf("coreUserCompanyRole.ExpireDate: %s", *coreUserCompanyRole.ExpireDate)
}

func TestToCoreUserCompanyRoleError(t *testing.T) {
	var now = time.Now()
	var expireDate = now.Add(1 * time.Hour)
	var userCompanyRoleFull = &user.UserCompanyRoleFull{
		UserCompanyRole: user.UserCompanyRole{
			PersistKey:        1,
			UserPersistKey:    2,
			CompanyPersistKey: 3,
			RoleLabel:         "test_role_label",
			RegisterDate:      now,
			ExpireDate:        &expireDate,
		},
		UserIdentifier:        "US-TESTES",
		UserExposeEmailId:     "test02@example.com",
		UserName:              "TestUserName",
		UserBotFlag:           false,
		UserRegisteredDate:    now.Add(2 * time.Hour),
		UserUpdateDate:        now.Add(3 * time.Hour),
		CompanyIdentifier:     "CP-TESTES",
		CompanyName:           "TestCompanyName",
		CompanyRegisteredDate: now.Add(4 * time.Hour),
		RoleName:              "TestRoleName",
		RoleDescription:       "TestRoleDescription",
		RoleRegisteredDate:    now.Add(5 * time.Hour),
	}

	var _, err = userCompanyRoleFull.ToCoreUserCompanyRole()

	assert.Error(t, err)

	t.Logf("Error: %v", err)
}
